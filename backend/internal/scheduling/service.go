package scheduling

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"schedule/internal/shared"
)

const (
	cacheKeyPrefix = "availability"
	cacheTTL       = 30 * time.Second
)

// AvailabilityOptions agrupa os parâmetros necessários para calcular disponibilidade.
type AvailabilityOptions struct {
	EstablishmentID string
	ServiceID       string
	ProfessionalID  string // opcional — vazio = qualquer profissional
	DateFrom        time.Time
	DateTo          time.Time
	Timezone        string
}

// Service orquestra a busca de disponibilidade, cache Redis e paralelismo.
type Service struct {
	repo  Repository
	redis *redis.Client
}

// NewService cria um novo Service com injeção de dependências.
func NewService(repo Repository, rdb *redis.Client) *Service {
	return &Service{repo: repo, redis: rdb}
}

// GetAvailability retorna os slots disponíveis agrupados por profissional.
// Se opts.ProfessionalID for vazio, busca todos os profissionais ativos em paralelo.
func (s *Service) GetAvailability(ctx context.Context, opts AvailabilityOptions) (map[string][]Slot, error) {
	// Busca duração do serviço (obrigatório).
	duration, err := s.repo.GetServiceDuration(ctx, opts.ServiceID, opts.EstablishmentID)
	if err != nil {
		return nil, err
	}

	// Busca horários de funcionamento do estabelecimento (compartilhado por todos os profissionais).
	bizHours, err := s.repo.GetBusinessHours(ctx, opts.EstablishmentID)
	if err != nil {
		return nil, err
	}

	// Determina a lista de profissionais a processar.
	var professionals []ProfessionalRef

	if opts.ProfessionalID != "" {
		professionals = []ProfessionalRef{{ID: opts.ProfessionalID}}
	} else {
		professionals, err = s.repo.GetActiveProfessionals(ctx, opts.EstablishmentID, opts.ServiceID)
		if err != nil {
			return nil, err
		}
	}

	if len(professionals) == 0 {
		return map[string][]Slot{}, nil
	}

	// Executa a busca em paralelo para cada profissional.
	type result struct {
		professionalID string
		slots          []Slot
		err            error
	}

	results := make(chan result, len(professionals))
	var wg sync.WaitGroup

	for _, prof := range professionals {
		wg.Add(1)
		go func(profID string) {
			defer wg.Done()

			slots, err := s.getSlotsForProfessional(ctx, profID, bizHours, duration, opts)
			results <- result{professionalID: profID, slots: slots, err: err}
		}(prof.ID)
	}

	// Fecha o canal depois que todas as goroutines terminarem.
	go func() {
		wg.Wait()
		close(results)
	}()

	// Coleta resultados, acumulando o primeiro erro encontrado.
	out := make(map[string][]Slot, len(professionals))
	var firstErr error

	for r := range results {
		if r.err != nil {
			if firstErr == nil {
				firstErr = r.err
			}
			continue
		}
		out[r.professionalID] = r.slots
	}

	if firstErr != nil {
		return nil, firstErr
	}

	return out, nil
}

// getSlotsForProfessional calcula os slots para um único profissional,
// respeitando o cache Redis por dia.
func (s *Service) getSlotsForProfessional(
	ctx context.Context,
	professionalID string,
	bizHours []BusinessHourRow,
	durationMinutes int,
	opts AvailabilityOptions,
) ([]Slot, error) {
	// Busca dados do profissional (jornada, agendamentos, bloqueios) uma única vez.
	profHours, err := s.repo.GetProfessionalHours(ctx, professionalID)
	if err != nil {
		return nil, err
	}

	appointments, err := s.repo.GetConfirmedAppointments(ctx, professionalID, opts.DateFrom, opts.DateTo)
	if err != nil {
		return nil, err
	}

	blockedPeriods, err := s.repo.GetBlockedPeriods(ctx, professionalID, opts.DateFrom, opts.DateTo)
	if err != nil {
		return nil, err
	}

	// Tenta cache Redis dia a dia.
	// Se todos os dias derem hit, retorna do cache; caso contrário, recalcula tudo e popula.
	var allSlots []Slot

	// Itera dia a dia para granularidade de cache por data.
	loc, err := time.LoadLocation(opts.Timezone)
	if err != nil {
		loc = time.UTC
	}

	current := truncateToDay(opts.DateFrom, loc)
	last := truncateToDay(opts.DateTo, loc)

	for !current.After(last) {
		dateStr := current.In(loc).Format("2006-01-02")
		cacheKey := buildCacheKey(opts.EstablishmentID, professionalID, dateStr)

		cached, err := s.redis.Get(ctx, cacheKey).Bytes()
		if err == nil {
			// Cache hit: desserializa e adiciona ao resultado.
			var daySlots []Slot
			if jsonErr := json.Unmarshal(cached, &daySlots); jsonErr == nil {
				allSlots = append(allSlots, daySlots...)
				current = current.Add(24 * time.Hour)
				continue
			}
			// Se desserialização falhar, recalcula para este dia.
		}

		// Cache miss: calcula slots apenas para este dia.
		dayFrom := current
		dayTo := current.Add(24 * time.Hour)

		daySlots := CalculateSlots(
			bizHours,
			profHours,
			appointments,
			blockedPeriods,
			nil, // externalBusy: Google Calendar integrado na Fase seguinte
			durationMinutes,
			dayFrom,
			dayTo,
			opts.Timezone,
			professionalID,
		)

		// Serializa e armazena no cache (ignora erros de cache — não bloqueante).
		if data, jsonErr := json.Marshal(daySlots); jsonErr == nil {
			_ = s.redis.Set(ctx, cacheKey, data, cacheTTL).Err()
		}

		allSlots = append(allSlots, daySlots...)
		current = current.Add(24 * time.Hour)
	}

	return allSlots, nil
}

// InvalidateCache remove as chaves de disponibilidade afetadas no Redis para
// o profissional e data fornecidos. O timezone deve corresponder ao fuso do
// estabelecimento para que a chave gerada bata com as chaves armazenadas.
// Deve ser chamada ao criar ou cancelar agendamentos.
func (s *Service) InvalidateCache(ctx context.Context, establishmentID, professionalID, timezone string, date time.Time) error {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		loc = time.UTC
	}
	dateStr := date.In(loc).Format("2006-01-02")
	cacheKey := buildCacheKey(establishmentID, professionalID, dateStr)
	return s.redis.Del(ctx, cacheKey).Err()
}

// buildCacheKey constrói a chave Redis no formato:
// availability:{establishmentID}:{professionalID}:{date}
func buildCacheKey(establishmentID, professionalID, date string) string {
	return fmt.Sprintf("%s:%s:%s:%s", cacheKeyPrefix, establishmentID, professionalID, date)
}

func (s *Service) CreatePublicAppointment(ctx context.Context, input CreateAppointmentInput) (*PublicAppointmentResult, error) {
	if input.EstablishmentID == "" || input.ServiceID == "" || input.ProfessionalID == "" || input.IdempotencyKey == "" {
		return nil, shared.ErrInvalidInput
	}

	input.ClientName = strings.TrimSpace(input.ClientName)
	input.ClientEmail = strings.TrimSpace(input.ClientEmail)
	input.ClientPhone = strings.TrimSpace(input.ClientPhone)
	input.ClientBirthDate = strings.TrimSpace(input.ClientBirthDate)
	if input.ClientName == "" || input.ClientEmail == "" || input.ClientPhone == "" || input.ClientBirthDate == "" {
		return nil, shared.ErrInvalidInput
	}

	if existing, err := s.repo.FindAppointmentByIdempotencyKey(ctx, input.EstablishmentID, input.IdempotencyKey); err == nil {
		return toPublicAppointmentResult(existing), nil
	} else if err != shared.ErrNotFound {
		return nil, err
	}

	duration, err := s.repo.GetServiceDuration(ctx, input.ServiceID, input.EstablishmentID)
	if err != nil {
		return nil, err
	}

	if err := s.repo.ValidateProfessionalForService(ctx, input.EstablishmentID, input.ProfessionalID, input.ServiceID); err != nil {
		return nil, err
	}

	timezone, err := s.repo.GetEstablishmentTimezone(ctx, input.EstablishmentID)
	if err != nil {
		return nil, err
	}

	endsAt := input.StartsAt.UTC().Add(time.Duration(duration) * time.Minute)
	dayStart, dayEnd := dayBounds(input.StartsAt.UTC(), timezone)

	bizHours, err := s.repo.GetBusinessHours(ctx, input.EstablishmentID)
	if err != nil {
		return nil, err
	}
	profHours, err := s.repo.GetProfessionalHours(ctx, input.ProfessionalID)
	if err != nil {
		return nil, err
	}
	appointments, err := s.repo.GetConfirmedAppointments(ctx, input.ProfessionalID, dayStart, dayEnd)
	if err != nil {
		return nil, err
	}
	blockedPeriods, err := s.repo.GetBlockedPeriods(ctx, input.ProfessionalID, dayStart, dayEnd)
	if err != nil {
		return nil, err
	}

	slots := CalculateSlots(
		bizHours,
		profHours,
		appointments,
		blockedPeriods,
		nil,
		duration,
		dayStart,
		dayEnd,
		timezone,
		input.ProfessionalID,
	)

	var slotAvailable bool
	for _, slot := range slots {
		if slot.StartsAt.Equal(input.StartsAt.UTC()) && slot.EndsAt.Equal(endsAt.UTC()) {
			slotAvailable = true
			break
		}
	}
	if !slotAvailable {
		return nil, shared.ErrSlotUnavailable
	}

	appointment, err := s.repo.CreateAppointment(ctx, input, endsAt)
	if err != nil {
		return nil, err
	}

	_ = s.InvalidateCache(ctx, input.EstablishmentID, input.ProfessionalID, timezone, input.StartsAt)

	return toPublicAppointmentResult(appointment), nil
}

func toPublicAppointmentResult(appointment *Appointment) *PublicAppointmentResult {
	return &PublicAppointmentResult{
		ID:              appointment.ID,
		ServiceID:       appointment.ServiceID,
		ProfessionalID:  appointment.ProfessionalID,
		ClientName:      appointment.ClientName,
		ClientEmail:     appointment.ClientEmail,
		ClientPhone:     appointment.ClientPhone,
		ClientBirthDate: appointment.ClientBirthDate,
		StartsAt:        appointment.StartsAt,
		EndsAt:          appointment.EndsAt,
		Status:          appointment.Status,
	}
}

func dayBounds(t time.Time, timezone string) (time.Time, time.Time) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		loc = time.UTC
	}
	dayStart := truncateToDay(t, loc)
	return dayStart, dayStart.Add(24 * time.Hour)
}

func (s *Service) GetPublicAppointment(ctx context.Context, establishmentID, appointmentID, phone string) (*PublicAppointmentDetail, error) {
	phone = strings.TrimSpace(phone)
	if establishmentID == "" || appointmentID == "" || phone == "" {
		return nil, shared.ErrInvalidInput
	}

	appointment, err := s.repo.FindAppointmentByIDAndPhone(ctx, establishmentID, appointmentID, phone)
	if err != nil {
		return nil, err
	}

	minAdvanceCancelHours, err := s.repo.GetMinAdvanceCancelHours(ctx, establishmentID)
	if err != nil {
		return nil, err
	}

	return toPublicAppointmentDetail(appointment, minAdvanceCancelHours), nil
}

func (s *Service) CancelPublicAppointment(ctx context.Context, establishmentID, appointmentID, phone string) (*PublicAppointmentDetail, error) {
	phone = strings.TrimSpace(phone)
	if establishmentID == "" || appointmentID == "" || phone == "" {
		return nil, shared.ErrInvalidInput
	}

	appointment, err := s.repo.FindAppointmentByIDAndPhone(ctx, establishmentID, appointmentID, phone)
	if err != nil {
		return nil, err
	}

	minAdvanceCancelHours, err := s.repo.GetMinAdvanceCancelHours(ctx, establishmentID)
	if err != nil {
		return nil, err
	}

	if appointment.Status == "cancelled" {
		return toPublicAppointmentDetail(appointment, minAdvanceCancelHours), nil
	}

	if !canCancelAppointment(appointment, minAdvanceCancelHours) {
		return nil, shared.ErrCancelTooLate
	}

	cancelled, err := s.repo.CancelAppointment(ctx, establishmentID, appointmentID)
	if err != nil {
		return nil, err
	}

	timezone, tzErr := s.repo.GetEstablishmentTimezone(ctx, establishmentID)
	if tzErr == nil {
		_ = s.InvalidateCache(ctx, establishmentID, cancelled.ProfessionalID, timezone, cancelled.StartsAt)
	}

	return toPublicAppointmentDetail(cancelled, minAdvanceCancelHours), nil
}

func (s *Service) ReschedulePublicAppointment(ctx context.Context, establishmentID, appointmentID, phone string, newStartsAt time.Time) (*PublicAppointmentResult, error) {
	phone = strings.TrimSpace(phone)
	if establishmentID == "" || appointmentID == "" || phone == "" {
		return nil, shared.ErrInvalidInput
	}

	current, err := s.repo.FindAppointmentByIDAndPhone(ctx, establishmentID, appointmentID, phone)
	if err != nil {
		return nil, err
	}

	minAdvanceCancelHours, err := s.repo.GetMinAdvanceCancelHours(ctx, establishmentID)
	if err != nil {
		return nil, err
	}

	if !canCancelAppointment(current, minAdvanceCancelHours) {
		return nil, shared.ErrCancelTooLate
	}

	if current.Status == "cancelled" {
		return nil, shared.ErrInvalidInput
	}

	idempotencyKey := fmt.Sprintf("reschedule:%s:%d", current.ID, newStartsAt.UTC().UnixNano())
	input := CreateAppointmentInput{
		EstablishmentID: establishmentID,
		ServiceID:       current.ServiceID,
		ProfessionalID:  current.ProfessionalID,
		StartsAt:        newStartsAt.UTC(),
		ClientName:      current.ClientName,
		ClientEmail:     derefString(current.ClientEmail),
		ClientPhone:     current.ClientPhone,
		ClientBirthDate: derefString(current.ClientBirthDate),
		IdempotencyKey:  idempotencyKey,
	}

	duration, err := s.repo.GetServiceDuration(ctx, input.ServiceID, input.EstablishmentID)
	if err != nil {
		return nil, err
	}

	timezone, err := s.repo.GetEstablishmentTimezone(ctx, input.EstablishmentID)
	if err != nil {
		return nil, err
	}

	endsAt := input.StartsAt.UTC().Add(time.Duration(duration) * time.Minute)
	dayStart, dayEnd := dayBounds(input.StartsAt.UTC(), timezone)

	bizHours, err := s.repo.GetBusinessHours(ctx, input.EstablishmentID)
	if err != nil {
		return nil, err
	}
	profHours, err := s.repo.GetProfessionalHours(ctx, input.ProfessionalID)
	if err != nil {
		return nil, err
	}
	appointments, err := s.repo.GetConfirmedAppointments(ctx, input.ProfessionalID, dayStart, dayEnd)
	if err != nil {
		return nil, err
	}
	filteredAppointments := make([]Appointment, 0, len(appointments))
	for _, appt := range appointments {
		if appt.ID == current.ID {
			continue
		}
		filteredAppointments = append(filteredAppointments, appt)
	}
	blockedPeriods, err := s.repo.GetBlockedPeriods(ctx, input.ProfessionalID, dayStart, dayEnd)
	if err != nil {
		return nil, err
	}

	slots := CalculateSlots(
		bizHours,
		profHours,
		filteredAppointments,
		blockedPeriods,
		nil,
		duration,
		dayStart,
		dayEnd,
		timezone,
		input.ProfessionalID,
	)

	var slotAvailable bool
	for _, slot := range slots {
		if slot.StartsAt.Equal(input.StartsAt.UTC()) && slot.EndsAt.Equal(endsAt.UTC()) {
			slotAvailable = true
			break
		}
	}
	if !slotAvailable {
		return nil, shared.ErrSlotUnavailable
	}

	newAppointment, err := s.repo.RescheduleAppointment(ctx, current, input, endsAt)
	if err != nil {
		return nil, err
	}

	_ = s.InvalidateCache(ctx, establishmentID, current.ProfessionalID, timezone, current.StartsAt)
	_ = s.InvalidateCache(ctx, establishmentID, newAppointment.ProfessionalID, timezone, newAppointment.StartsAt)

	return toPublicAppointmentResult(newAppointment), nil
}

func toPublicAppointmentDetail(appointment *Appointment, minAdvanceCancelHours int) *PublicAppointmentDetail {
	return &PublicAppointmentDetail{
		ID:                    appointment.ID,
		ServiceID:             appointment.ServiceID,
		ProfessionalID:        appointment.ProfessionalID,
		ClientName:            appointment.ClientName,
		ClientEmail:           appointment.ClientEmail,
		ClientPhone:           appointment.ClientPhone,
		ClientBirthDate:       appointment.ClientBirthDate,
		StartsAt:              appointment.StartsAt,
		EndsAt:                appointment.EndsAt,
		Status:                appointment.Status,
		CanCancel:             canCancelAppointment(appointment, minAdvanceCancelHours),
		MinAdvanceCancelHours: minAdvanceCancelHours,
	}
}

func canCancelAppointment(appointment *Appointment, minAdvanceCancelHours int) bool {
	if appointment.Status == "cancelled" || appointment.Status == "completed" || appointment.Status == "no_show" {
		return false
	}
	deadline := appointment.StartsAt.Add(-time.Duration(minAdvanceCancelHours) * time.Hour)
	return time.Now().UTC().Before(deadline)
}

func derefString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

// ── Métodos do gestor (Fase 10) ────────────────────────────────────────────

var validManagerStatuses = map[string]bool{
	"completed": true,
	"no_show":   true,
	"cancelled": true,
}

func (s *Service) ListManagerAppointments(ctx context.Context, establishmentID string, filter AppointmentFilter) ([]ManagerAppointmentRow, int, error) {
	return s.repo.ListManagerAppointments(ctx, establishmentID, filter)
}

func (s *Service) GetManagerAppointment(ctx context.Context, establishmentID, appointmentID string) (*ManagerAppointmentRow, error) {
	if establishmentID == "" || appointmentID == "" {
		return nil, shared.ErrInvalidInput
	}
	return s.repo.GetManagerAppointment(ctx, establishmentID, appointmentID)
}

func (s *Service) UpdateAppointmentStatus(ctx context.Context, input UpdateStatusInput) (*ManagerAppointmentRow, error) {
	if !validManagerStatuses[input.Status] {
		return nil, &shared.DomainError{
			Code:    "INVALID_STATUS",
			Message: "Status inválido. Use: completed, no_show ou cancelled.",
			Status:  400,
		}
	}
	appt, err := s.repo.GetManagerAppointment(ctx, input.EstablishmentID, input.AppointmentID)
	if err != nil {
		return nil, err
	}
	if appt.Status == "cancelled" || appt.Status == "completed" || appt.Status == "no_show" {
		return nil, &shared.DomainError{
			Code:    "STATUS_ALREADY_TERMINAL",
			Message: "O agendamento já está em um estado final e não pode ser alterado.",
			Status:  409,
		}
	}
	updated, err := s.repo.UpdateAppointmentStatus(ctx, input)
	if err != nil {
		return nil, err
	}
	if input.Status == "cancelled" {
		timezone, tzErr := s.repo.GetEstablishmentTimezone(ctx, input.EstablishmentID)
		if tzErr == nil {
			_ = s.InvalidateCache(ctx, input.EstablishmentID, updated.ProfessionalID, timezone, updated.StartsAt)
		}
	}
	return updated, nil
}

func (s *Service) ListManagerBlockedPeriods(ctx context.Context, establishmentID, professionalID, date string) ([]ManagerBlockedPeriod, error) {
	return s.repo.ListManagerBlockedPeriods(ctx, establishmentID, professionalID, date)
}

func (s *Service) CreateBlockedPeriod(ctx context.Context, input CreateBlockedPeriodInput) (*ManagerBlockedPeriod, error) {
	if input.EstablishmentID == "" || input.ProfessionalID == "" {
		return nil, shared.ErrInvalidInput
	}
	if input.EndsAt.Before(input.StartsAt) || input.EndsAt.Equal(input.StartsAt) {
		return nil, &shared.DomainError{
			Code:    "INVALID_PERIOD",
			Message: "ends_at deve ser posterior a starts_at.",
			Status:  400,
		}
	}
	period, err := s.repo.CreateBlockedPeriod(ctx, input)
	if err != nil {
		return nil, err
	}
	timezone, tzErr := s.repo.GetEstablishmentTimezone(ctx, input.EstablishmentID)
	if tzErr == nil {
		_ = s.InvalidateCache(ctx, input.EstablishmentID, input.ProfessionalID, timezone, input.StartsAt)
	}
	return period, nil
}

func (s *Service) DeleteBlockedPeriod(ctx context.Context, establishmentID, blockedPeriodID string) error {
	return s.repo.DeleteBlockedPeriod(ctx, establishmentID, blockedPeriodID)
}
