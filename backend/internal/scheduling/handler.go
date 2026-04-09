package scheduling

import (
	"fmt"
	"net/http"
	"net/mail"
	"sort"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"schedule/internal/shared"
)

// Handler expõe os endpoints HTTP do módulo de agendamentos.
type Handler struct {
	svc *Service
}

// NewHandler cria um Handler com injeção de dependência.
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// GetAvailability trata GET /pub/{slug}/availability
//
// Query params:
//
//	service_id=:id        (obrigatório)
//	professional_id=:id   (opcional — vazio = todos os profissionais)
//	date_from=2026-04-08  (obrigatório, formato YYYY-MM-DD)
//	date_to=2026-04-14    (obrigatório, formato YYYY-MM-DD)
//
// Response:
//
//	{
//	  "data": {
//	    "2026-04-08": [
//	      { "starts_at": "...", "ends_at": "...", "professional_id": "..." }
//	    ],
//	    "2026-04-09": []
//	  }
//	}
func (h *Handler) GetAvailability(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	establishmentID := shared.EstablishmentIDFromContext(ctx)

	q := r.URL.Query()

	serviceID := q.Get("service_id")
	if serviceID == "" {
		shared.JSONError(w, &shared.DomainError{
			Code:    "MISSING_SERVICE_ID",
			Message: "O parâmetro service_id é obrigatório.",
			Status:  http.StatusBadRequest,
		})
		return
	}

	dateFromStr := q.Get("date_from")
	dateToStr := q.Get("date_to")

	if dateFromStr == "" || dateToStr == "" {
		shared.JSONError(w, &shared.DomainError{
			Code:    "MISSING_DATE_RANGE",
			Message: "Os parâmetros date_from e date_to são obrigatórios.",
			Status:  http.StatusBadRequest,
		})
		return
	}

	const layout = "2006-01-02"

	dateFrom, err := time.Parse(layout, dateFromStr)
	if err != nil {
		shared.JSONError(w, &shared.DomainError{
			Code:    "INVALID_DATE_FORMAT",
			Message: "Formato inválido para date_from. Use YYYY-MM-DD.",
			Status:  http.StatusBadRequest,
		})
		return
	}

	dateTo, err := time.Parse(layout, dateToStr)
	if err != nil {
		shared.JSONError(w, &shared.DomainError{
			Code:    "INVALID_DATE_FORMAT",
			Message: "Formato inválido para date_to. Use YYYY-MM-DD.",
			Status:  http.StatusBadRequest,
		})
		return
	}

	// Validar que date_to >= date_from.
	if dateTo.Before(dateFrom) {
		shared.JSONError(w, &shared.DomainError{
			Code:    "INVALID_DATE_RANGE",
			Message: "date_to não pode ser anterior a date_from.",
			Status:  http.StatusBadRequest,
		})
		return
	}

	// Validar máximo de 31 dias.
	if dateTo.Sub(dateFrom) > 31*24*time.Hour {
		shared.JSONError(w, &shared.DomainError{
			Code:    "DATE_RANGE_TOO_LARGE",
			Message: "O intervalo máximo permitido é de 31 dias.",
			Status:  http.StatusBadRequest,
		})
		return
	}

	professionalID := q.Get("professional_id")

	// Timezone padrão UTC; idealmente buscado do establishment (simplificado aqui).
	// O service pode receber "UTC" e os dados já estarão corretos para exibição.
	// Em produção, o timezone vem do establishment carregado pelo middleware.
	timezone := q.Get("timezone")
	if timezone == "" {
		timezone = "America/Sao_Paulo" // default para o nicho inicial
	}

	opts := AvailabilityOptions{
		EstablishmentID: establishmentID,
		ServiceID:       serviceID,
		ProfessionalID:  professionalID,
		DateFrom:        dateFrom.UTC(),
		DateTo:          dateTo.Add(24 * time.Hour).UTC(), // inclui o dia de dateTo
		Timezone:        timezone,
	}

	slotsByProf, err := h.svc.GetAvailability(ctx, opts)
	if err != nil {
		shared.JSONError(w, err)
		return
	}

	// Agrupa slots por data (YYYY-MM-DD) no timezone padrão.
	// A resposta é um mapa data → []Slot, incluindo dias sem slots (array vazio).
	loc, _ := time.LoadLocation(timezone)
	if loc == nil {
		loc = time.UTC
	}

	// Inicializa todos os dias do intervalo com arrays vazios.
	byDate := make(map[string][]Slot)
	current := dateFrom
	for !current.After(dateTo) {
		byDate[current.Format(layout)] = []Slot{}
		current = current.Add(24 * time.Hour)
	}

	// Preenche com os slots calculados.
	for _, slots := range slotsByProf {
		for _, slot := range slots {
			dateKey := slot.StartsAt.In(loc).Format(layout)
			byDate[dateKey] = append(byDate[dateKey], slot)
		}
	}

	// Ordena slots dentro de cada dia por starts_at.
	for dateKey := range byDate {
		sort.Slice(byDate[dateKey], func(i, j int) bool {
			si := byDate[dateKey][i]
			sj := byDate[dateKey][j]
			return si.StartsAt.Before(sj.StartsAt)
		})
	}

	shared.JSON(w, http.StatusOK, byDate)
}

func (h *Handler) CreatePublicAppointment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	establishmentID := shared.EstablishmentIDFromContext(ctx)

	var req struct {
		ServiceID       string `json:"service_id"`
		ProfessionalID  string `json:"professional_id"`
		StartsAt        string `json:"starts_at"`
		ClientName      string `json:"client_name"`
		ClientEmail     string `json:"client_email"`
		ClientPhone     string `json:"client_phone"`
		ClientBirthDate string `json:"client_birth_date"`
		IdempotencyKey  string `json:"idempotency_key"`
	}
	if err := shared.Decode(r, &req); err != nil {
		shared.JSONError(w, err)
		return
	}

	startsAt, err := time.Parse(time.RFC3339, req.StartsAt)
	if err != nil {
		shared.JSONError(w, &shared.DomainError{
			Code:    "INVALID_STARTS_AT",
			Message: "Formato invalido para starts_at. Use RFC3339.",
			Status:  http.StatusBadRequest,
		})
		return
	}

	if strings.TrimSpace(req.ClientName) == "" ||
		strings.TrimSpace(req.ClientEmail) == "" ||
		strings.TrimSpace(req.ClientPhone) == "" ||
		strings.TrimSpace(req.ClientBirthDate) == "" {
		shared.JSONError(w, shared.ErrInvalidInput)
		return
	}

	if _, err := mail.ParseAddress(strings.TrimSpace(req.ClientEmail)); err != nil {
		shared.JSONError(w, &shared.DomainError{
			Code:    "INVALID_CLIENT_EMAIL",
			Message: "Formato invalido para client_email.",
			Status:  http.StatusBadRequest,
		})
		return
	}

	if _, err := time.Parse("2006-01-02", strings.TrimSpace(req.ClientBirthDate)); err != nil {
		shared.JSONError(w, &shared.DomainError{
			Code:    "INVALID_CLIENT_BIRTH_DATE",
			Message: "Formato invalido para client_birth_date. Use YYYY-MM-DD.",
			Status:  http.StatusBadRequest,
		})
		return
	}

	result, err := h.svc.CreatePublicAppointment(ctx, CreateAppointmentInput{
		EstablishmentID: establishmentID,
		ServiceID:       req.ServiceID,
		ProfessionalID:  req.ProfessionalID,
		StartsAt:        startsAt.UTC(),
		ClientName:      req.ClientName,
		ClientEmail:     req.ClientEmail,
		ClientPhone:     req.ClientPhone,
		ClientBirthDate: req.ClientBirthDate,
		IdempotencyKey:  req.IdempotencyKey,
	})
	if err != nil {
		shared.JSONError(w, err)
		return
	}

	shared.JSON(w, http.StatusCreated, result)
}

func (h *Handler) GetPublicAppointment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	establishmentID := shared.EstablishmentIDFromContext(ctx)
	appointmentID := chi.URLParam(r, "id")
	phone := r.URL.Query().Get("phone")

	result, err := h.svc.GetPublicAppointment(ctx, establishmentID, appointmentID, phone)
	if err != nil {
		shared.JSONError(w, err)
		return
	}

	shared.JSON(w, http.StatusOK, result)
}

func (h *Handler) CancelPublicAppointment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	establishmentID := shared.EstablishmentIDFromContext(ctx)
	appointmentID := chi.URLParam(r, "id")

	var req struct {
		Phone string `json:"phone"`
	}
	if err := shared.Decode(r, &req); err != nil {
		shared.JSONError(w, err)
		return
	}

	result, err := h.svc.CancelPublicAppointment(ctx, establishmentID, appointmentID, req.Phone)
	if err != nil {
		shared.JSONError(w, err)
		return
	}

	shared.JSON(w, http.StatusOK, result)
}

func (h *Handler) ReschedulePublicAppointment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	establishmentID := shared.EstablishmentIDFromContext(ctx)
	appointmentID := chi.URLParam(r, "id")

	var req struct {
		Phone    string `json:"phone"`
		StartsAt string `json:"starts_at"`
	}
	if err := shared.Decode(r, &req); err != nil {
		shared.JSONError(w, err)
		return
	}

	startsAt, err := time.Parse(time.RFC3339, req.StartsAt)
	if err != nil {
		shared.JSONError(w, &shared.DomainError{
			Code:    "INVALID_STARTS_AT",
			Message: "Formato invalido para starts_at. Use RFC3339.",
			Status:  http.StatusBadRequest,
		})
		return
	}

	result, err := h.svc.ReschedulePublicAppointment(ctx, establishmentID, appointmentID, req.Phone, startsAt.UTC())
	if err != nil {
		shared.JSONError(w, err)
		return
	}

	shared.JSON(w, http.StatusOK, result)
}

// ── Handlers do gestor (Fase 10) ──────────────────────────────────────────

// ListManagerAppointments trata GET /api/v1/appointments
// Query params: date (YYYY-MM-DD), professional_id, status, page, per_page
func (h *Handler) ListManagerAppointments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	establishmentID := shared.EstablishmentIDFromContext(ctx)
	q := r.URL.Query()

	page := parseIntParam(q.Get("page"), 1)
	perPage := parseIntParam(q.Get("per_page"), 20)
	if perPage > 100 {
		perPage = 100
	}

	filter := AppointmentFilter{
		Date:           q.Get("date"),
		ProfessionalID: q.Get("professional_id"),
		Status:         q.Get("status"),
		Page:           page,
		PerPage:        perPage,
	}

	rows, total, err := h.svc.ListManagerAppointments(ctx, establishmentID, filter)
	if err != nil {
		shared.JSONError(w, err)
		return
	}

	shared.JSONList(w, http.StatusOK, rows, shared.Meta{
		Page:    page,
		PerPage: perPage,
		Total:   total,
	})
}

// GetManagerAppointment trata GET /api/v1/appointments/{id}
func (h *Handler) GetManagerAppointment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	establishmentID := shared.EstablishmentIDFromContext(ctx)
	appointmentID := chi.URLParam(r, "id")

	appt, err := h.svc.GetManagerAppointment(ctx, establishmentID, appointmentID)
	if err != nil {
		shared.JSONError(w, err)
		return
	}

	shared.JSON(w, http.StatusOK, appt)
}

// UpdateAppointmentStatus trata PATCH /api/v1/appointments/{id}/status
func (h *Handler) UpdateAppointmentStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	establishmentID := shared.EstablishmentIDFromContext(ctx)
	appointmentID := chi.URLParam(r, "id")

	var req struct {
		Status string `json:"status"`
	}
	if err := shared.Decode(r, &req); err != nil {
		shared.JSONError(w, err)
		return
	}

	updated, err := h.svc.UpdateAppointmentStatus(ctx, UpdateStatusInput{
		EstablishmentID: establishmentID,
		AppointmentID:   appointmentID,
		Status:          req.Status,
	})
	if err != nil {
		shared.JSONError(w, err)
		return
	}

	shared.JSON(w, http.StatusOK, updated)
}

// ListManagerBlockedPeriods trata GET /api/v1/blocked-periods
// Query params: professional_id, date (YYYY-MM-DD)
func (h *Handler) ListManagerBlockedPeriods(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	establishmentID := shared.EstablishmentIDFromContext(ctx)
	q := r.URL.Query()

	periods, err := h.svc.ListManagerBlockedPeriods(ctx, establishmentID, q.Get("professional_id"), q.Get("date"))
	if err != nil {
		shared.JSONError(w, err)
		return
	}

	shared.JSON(w, http.StatusOK, periods)
}

// CreateBlockedPeriod trata POST /api/v1/blocked-periods
func (h *Handler) CreateBlockedPeriod(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	establishmentID := shared.EstablishmentIDFromContext(ctx)

	var req struct {
		ProfessionalID string `json:"professional_id"`
		StartsAt       string `json:"starts_at"`
		EndsAt         string `json:"ends_at"`
		Reason         string `json:"reason"`
	}
	if err := shared.Decode(r, &req); err != nil {
		shared.JSONError(w, err)
		return
	}

	startsAt, err := time.Parse(time.RFC3339, req.StartsAt)
	if err != nil {
		shared.JSONError(w, &shared.DomainError{Code: "INVALID_STARTS_AT", Message: "Formato invalido para starts_at. Use RFC3339.", Status: http.StatusBadRequest})
		return
	}
	endsAt, err := time.Parse(time.RFC3339, req.EndsAt)
	if err != nil {
		shared.JSONError(w, &shared.DomainError{Code: "INVALID_ENDS_AT", Message: "Formato invalido para ends_at. Use RFC3339.", Status: http.StatusBadRequest})
		return
	}

	period, err := h.svc.CreateBlockedPeriod(ctx, CreateBlockedPeriodInput{
		EstablishmentID: establishmentID,
		ProfessionalID:  req.ProfessionalID,
		StartsAt:        startsAt.UTC(),
		EndsAt:          endsAt.UTC(),
		Reason:          req.Reason,
	})
	if err != nil {
		shared.JSONError(w, err)
		return
	}

	shared.JSON(w, http.StatusCreated, period)
}

// DeleteBlockedPeriod trata DELETE /api/v1/blocked-periods/{id}
func (h *Handler) DeleteBlockedPeriod(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	establishmentID := shared.EstablishmentIDFromContext(ctx)
	blockedPeriodID := chi.URLParam(r, "id")

	if err := h.svc.DeleteBlockedPeriod(ctx, establishmentID, blockedPeriodID); err != nil {
		shared.JSONError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseIntParam(s string, fallback int) int {
	if s == "" {
		return fallback
	}
	var n int
	if _, err := fmt.Sscanf(s, "%d", &n); err != nil || n < 1 {
		return fallback
	}
	return n
}
