package scheduling

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"schedule/internal/shared"
)

// BusinessHourRow contém os dados de horário de funcionamento do estabelecimento
// necessários para o algoritmo de disponibilidade.
type BusinessHourRow struct {
	DayOfWeek int    `db:"day_of_week"`
	OpenTime  string `db:"open_time"`
	CloseTime string `db:"close_time"`
	IsClosed  bool   `db:"is_closed"`
}

// ProfessionalHourRow contém os dados de jornada individual do profissional.
type ProfessionalHourRow struct {
	DayOfWeek     int    `db:"day_of_week"`
	StartTime     string `db:"start_time"`
	EndTime       string `db:"end_time"`
	IsUnavailable bool   `db:"is_unavailable"`
}

// ProfessionalRef é uma referência leve ao profissional (sem todos os campos).
type ProfessionalRef struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

// Repository define as operações de acesso a dados necessárias ao motor de disponibilidade.
type Repository interface {
	// GetBusinessHours retorna os horários de funcionamento do estabelecimento.
	GetBusinessHours(ctx context.Context, establishmentID string) ([]BusinessHourRow, error)

	// GetProfessionalHours retorna a jornada individual do profissional por dia.
	GetProfessionalHours(ctx context.Context, professionalID string) ([]ProfessionalHourRow, error)

	// GetConfirmedAppointments retorna agendamentos confirmados (status != 'cancelled')
	// no intervalo [from, to] para o profissional informado.
	GetConfirmedAppointments(ctx context.Context, professionalID string, from, to time.Time) ([]Appointment, error)

	// GetBlockedPeriods retorna bloqueios manuais no intervalo [from, to].
	GetBlockedPeriods(ctx context.Context, professionalID string, from, to time.Time) ([]BlockedPeriod, error)

	// GetActiveProfessionals retorna profissionais ativos do estabelecimento.
	// Se serviceID não for vazio, filtra apenas profissionais que oferecem esse serviço.
	GetActiveProfessionals(ctx context.Context, establishmentID string, serviceID string) ([]ProfessionalRef, error)

	// GetServiceDuration retorna a duração em minutos do serviço.
	// Retorna ErrNotFound se o serviço não pertencer ao estabelecimento.
	GetServiceDuration(ctx context.Context, serviceID, establishmentID string) (int, error)

	GetEstablishmentTimezone(ctx context.Context, establishmentID string) (string, error)
	GetMinAdvanceCancelHours(ctx context.Context, establishmentID string) (int, error)
	ValidateProfessionalForService(ctx context.Context, establishmentID, professionalID, serviceID string) error
	FindAppointmentByIdempotencyKey(ctx context.Context, establishmentID, idempotencyKey string) (*Appointment, error)
	FindAppointmentByIDAndPhone(ctx context.Context, establishmentID, appointmentID, phone string) (*Appointment, error)
	CreateAppointment(ctx context.Context, input CreateAppointmentInput, endsAt time.Time) (*Appointment, error)
	CancelAppointment(ctx context.Context, establishmentID, appointmentID string) (*Appointment, error)
	RescheduleAppointment(ctx context.Context, current *Appointment, input CreateAppointmentInput, endsAt time.Time) (*Appointment, error)

	// Métodos do gestor (Fase 10)
	ListManagerAppointments(ctx context.Context, establishmentID string, filter AppointmentFilter) ([]ManagerAppointmentRow, int, error)
	GetManagerAppointment(ctx context.Context, establishmentID, appointmentID string) (*ManagerAppointmentRow, error)
	UpdateAppointmentStatus(ctx context.Context, input UpdateStatusInput) (*ManagerAppointmentRow, error)
	ListManagerBlockedPeriods(ctx context.Context, establishmentID, professionalID, date string) ([]ManagerBlockedPeriod, error)
	CreateBlockedPeriod(ctx context.Context, input CreateBlockedPeriodInput) (*ManagerBlockedPeriod, error)
	DeleteBlockedPeriod(ctx context.Context, establishmentID, blockedPeriodID string) error
}

type repository struct {
	db *sqlx.DB
}

// NewRepository cria uma implementação concreta de Repository.
func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetBusinessHours(ctx context.Context, establishmentID string) ([]BusinessHourRow, error) {
	var rows []BusinessHourRow
	err := r.db.SelectContext(ctx, &rows,
		`SELECT day_of_week, open_time, close_time, is_closed
		   FROM business_hours
		  WHERE establishment_id = ?
		  ORDER BY day_of_week`,
		establishmentID,
	)
	return rows, err
}

func (r *repository) GetProfessionalHours(ctx context.Context, professionalID string) ([]ProfessionalHourRow, error) {
	var rows []ProfessionalHourRow
	err := r.db.SelectContext(ctx, &rows,
		`SELECT day_of_week, start_time, end_time, is_unavailable
		   FROM professional_hours
		  WHERE professional_id = ?
		  ORDER BY day_of_week`,
		professionalID,
	)
	return rows, err
}

func (r *repository) GetConfirmedAppointments(ctx context.Context, professionalID string, from, to time.Time) ([]Appointment, error) {
	var rows []Appointment
	err := r.db.SelectContext(ctx, &rows,
		`SELECT id, establishment_id, professional_id, service_id,
		        client_name, client_phone, starts_at, ends_at,
		        status, source, google_event_id, notes, idempotency_key,
		        created_at, updated_at
		   FROM appointments
		  WHERE professional_id = ?
		    AND status != 'cancelled'
		    AND starts_at < ?
		    AND ends_at   > ?`,
		professionalID, to, from,
	)
	return rows, err
}

func (r *repository) GetBlockedPeriods(ctx context.Context, professionalID string, from, to time.Time) ([]BlockedPeriod, error) {
	var rows []BlockedPeriod
	err := r.db.SelectContext(ctx, &rows,
		`SELECT id, professional_id, starts_at, ends_at, reason, google_event_id
		   FROM blocked_periods
		  WHERE professional_id = ?
		    AND starts_at < ?
		    AND ends_at   > ?`,
		professionalID, to, from,
	)
	return rows, err
}

func (r *repository) GetActiveProfessionals(ctx context.Context, establishmentID string, serviceID string) ([]ProfessionalRef, error) {
	var refs []ProfessionalRef

	if serviceID == "" {
		err := r.db.SelectContext(ctx, &refs,
			`SELECT id, name
			   FROM professionals
			  WHERE establishment_id = ?
			    AND active = true
			  ORDER BY display_order, name`,
			establishmentID,
		)
		return refs, err
	}

	err := r.db.SelectContext(ctx, &refs,
		`SELECT p.id, p.name
		   FROM professionals p
		   JOIN professional_services ps ON ps.professional_id = p.id
		  WHERE p.establishment_id = ?
		    AND p.active = true
		    AND ps.service_id = ?
		  ORDER BY p.display_order, p.name`,
		establishmentID, serviceID,
	)
	return refs, err
}

func (r *repository) GetServiceDuration(ctx context.Context, serviceID, establishmentID string) (int, error) {
	var duration int
	err := r.db.GetContext(ctx, &duration,
		`SELECT duration_minutes
		   FROM services
		  WHERE id = ?
		    AND establishment_id = ?
		    AND active = true`,
		serviceID, establishmentID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, shared.ErrNotFound
		}
		return 0, err
	}
	return duration, nil
}

func (r *repository) GetEstablishmentTimezone(ctx context.Context, establishmentID string) (string, error) {
	var timezone string
	err := r.db.GetContext(ctx, &timezone,
		`SELECT timezone FROM establishments WHERE id = ? AND active = true`,
		establishmentID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", shared.ErrNotFound
		}
		return "", err
	}
	return timezone, nil
}

func (r *repository) GetMinAdvanceCancelHours(ctx context.Context, establishmentID string) (int, error) {
	var hours int
	err := r.db.GetContext(ctx, &hours,
		`SELECT min_advance_cancel_hours FROM establishments WHERE id = ? AND active = true`,
		establishmentID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, shared.ErrNotFound
		}
		return 0, err
	}
	return hours, nil
}

func (r *repository) ValidateProfessionalForService(ctx context.Context, establishmentID, professionalID, serviceID string) error {
	var count int
	err := r.db.GetContext(ctx, &count, `
		SELECT COUNT(*)
		  FROM professionals p
		  JOIN professional_services ps ON ps.professional_id = p.id
		  JOIN services s ON s.id = ps.service_id
		 WHERE p.id = ?
		   AND p.establishment_id = ?
		   AND p.active = true
		   AND s.id = ?
		   AND s.establishment_id = ?
		   AND s.active = true`,
		professionalID, establishmentID, serviceID, establishmentID,
	)
	if err != nil {
		return err
	}
	if count == 0 {
		return shared.ErrNotFound
	}
	return nil
}

func (r *repository) FindAppointmentByIdempotencyKey(ctx context.Context, establishmentID, idempotencyKey string) (*Appointment, error) {
	var appt Appointment
	err := r.db.GetContext(ctx, &appt, `
		SELECT id, establishment_id, professional_id, service_id,
		       client_name, client_email, client_phone, client_birth_date, starts_at, ends_at,
		       status, source, google_event_id, notes, idempotency_key,
		       created_at, updated_at
		  FROM appointments
		 WHERE establishment_id = ?
		   AND idempotency_key = ?`,
		establishmentID, idempotencyKey,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.ErrNotFound
		}
		return nil, err
	}
	return &appt, nil
}

func (r *repository) FindAppointmentByIDAndPhone(ctx context.Context, establishmentID, appointmentID, phone string) (*Appointment, error) {
	var appt Appointment
	err := r.db.GetContext(ctx, &appt, `
		SELECT id, establishment_id, professional_id, service_id,
		       client_name, client_email, client_phone, client_birth_date, starts_at, ends_at,
		       status, source, google_event_id, notes, idempotency_key,
		       created_at, updated_at
		  FROM appointments
		 WHERE establishment_id = ?
		   AND id = ?
		   AND client_phone = ?`,
		establishmentID, appointmentID, phone,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.ErrNotFound
		}
		return nil, err
	}
	return &appt, nil
}

func (r *repository) CreateAppointment(ctx context.Context, input CreateAppointmentInput, endsAt time.Time) (*Appointment, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var existing Appointment
	err = tx.GetContext(ctx, &existing, `
		SELECT id, establishment_id, professional_id, service_id,
		       client_name, client_email, client_phone, client_birth_date, starts_at, ends_at,
		       status, source, google_event_id, notes, idempotency_key,
		       created_at, updated_at
		  FROM appointments
		 WHERE establishment_id = ?
		   AND idempotency_key = ?`,
		input.EstablishmentID, input.IdempotencyKey,
	)
	if err == nil {
		if commitErr := tx.Commit(); commitErr != nil {
			return nil, commitErr
		}
		return &existing, nil
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	var conflictID string
	err = tx.GetContext(ctx, &conflictID, `
		SELECT id
		  FROM appointments
		 WHERE professional_id = ?
		   AND starts_at < ?
		   AND ends_at > ?
		   AND status != 'cancelled'
		 LIMIT 1
		 FOR UPDATE`,
		input.ProfessionalID, endsAt, input.StartsAt,
	)
	if err == nil {
		return nil, shared.ErrSlotUnavailable
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	now := time.Now().UTC()
	idempotencyKey := input.IdempotencyKey
	appointment := &Appointment{
		ID:              shared.NewID(),
		EstablishmentID: input.EstablishmentID,
		ProfessionalID:  input.ProfessionalID,
		ServiceID:       input.ServiceID,
		ClientName:      input.ClientName,
		ClientEmail:     nullableString(input.ClientEmail),
		ClientPhone:     input.ClientPhone,
		ClientBirthDate: nullableString(input.ClientBirthDate),
		StartsAt:        input.StartsAt.UTC(),
		EndsAt:          endsAt.UTC(),
		Status:          "confirmed",
		Source:          "booking_link",
		IdempotencyKey:  &idempotencyKey,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO appointments (
			id, establishment_id, professional_id, service_id,
			client_name, client_email, client_phone, client_birth_date, starts_at, ends_at,
			status, source, idempotency_key, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		appointment.ID, appointment.EstablishmentID, appointment.ProfessionalID, appointment.ServiceID,
		appointment.ClientName, appointment.ClientEmail, appointment.ClientPhone, appointment.ClientBirthDate, appointment.StartsAt, appointment.EndsAt,
		appointment.Status, appointment.Source, appointment.IdempotencyKey, appointment.CreatedAt, appointment.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return appointment, nil
}

func (r *repository) CancelAppointment(ctx context.Context, establishmentID, appointmentID string) (*Appointment, error) {
	now := time.Now().UTC()
	_, err := r.db.ExecContext(ctx, `
		UPDATE appointments
		   SET status = 'cancelled', updated_at = ?
		 WHERE establishment_id = ?
		   AND id = ?
		   AND status != 'cancelled'`,
		now, establishmentID, appointmentID,
	)
	if err != nil {
		return nil, err
	}

	var appt Appointment
	err = r.db.GetContext(ctx, &appt, `
		SELECT id, establishment_id, professional_id, service_id,
		       client_name, client_email, client_phone, client_birth_date, starts_at, ends_at,
		       status, source, google_event_id, notes, idempotency_key,
		       created_at, updated_at
		  FROM appointments
		 WHERE establishment_id = ?
		   AND id = ?`,
		establishmentID, appointmentID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.ErrNotFound
		}
		return nil, err
	}
	return &appt, nil
}

func (r *repository) RescheduleAppointment(ctx context.Context, current *Appointment, input CreateAppointmentInput, endsAt time.Time) (*Appointment, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var lockedCurrent Appointment
	err = tx.GetContext(ctx, &lockedCurrent, `
		SELECT id, establishment_id, professional_id, service_id,
		       client_name, client_email, client_phone, client_birth_date, starts_at, ends_at,
		       status, source, google_event_id, notes, idempotency_key,
		       created_at, updated_at
		  FROM appointments
		 WHERE establishment_id = ?
		   AND id = ?
		 FOR UPDATE`,
		current.EstablishmentID, current.ID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.ErrNotFound
		}
		return nil, err
	}

	var conflictID string
	err = tx.GetContext(ctx, &conflictID, `
		SELECT id
		  FROM appointments
		 WHERE professional_id = ?
		   AND starts_at < ?
		   AND ends_at > ?
		   AND status != 'cancelled'
		   AND id != ?
		 LIMIT 1
		 FOR UPDATE`,
		input.ProfessionalID, endsAt, input.StartsAt, current.ID,
	)
	if err == nil {
		return nil, shared.ErrSlotUnavailable
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	now := time.Now().UTC()
	_, err = tx.ExecContext(ctx, `
		UPDATE appointments
		   SET status = 'cancelled', updated_at = ?
		 WHERE id = ?
		   AND establishment_id = ?`,
		now, current.ID, current.EstablishmentID,
	)
	if err != nil {
		return nil, err
	}

	idempotencyKey := input.IdempotencyKey
	newAppointment := &Appointment{
		ID:              shared.NewID(),
		EstablishmentID: input.EstablishmentID,
		ProfessionalID:  input.ProfessionalID,
		ServiceID:       input.ServiceID,
		ClientName:      input.ClientName,
		ClientEmail:     nullableString(input.ClientEmail),
		ClientPhone:     input.ClientPhone,
		ClientBirthDate: nullableString(input.ClientBirthDate),
		StartsAt:        input.StartsAt.UTC(),
		EndsAt:          endsAt.UTC(),
		Status:          "confirmed",
		Source:          "booking_link",
		IdempotencyKey:  &idempotencyKey,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO appointments (
			id, establishment_id, professional_id, service_id,
			client_name, client_email, client_phone, client_birth_date, starts_at, ends_at,
			status, source, idempotency_key, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		newAppointment.ID, newAppointment.EstablishmentID, newAppointment.ProfessionalID, newAppointment.ServiceID,
		newAppointment.ClientName, newAppointment.ClientEmail, newAppointment.ClientPhone, newAppointment.ClientBirthDate, newAppointment.StartsAt, newAppointment.EndsAt,
		newAppointment.Status, newAppointment.Source, newAppointment.IdempotencyKey, newAppointment.CreatedAt, newAppointment.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return newAppointment, nil
}

const managerAppointmentSelect = `
	SELECT a.id, a.establishment_id, a.professional_id, p.name AS professional_name,
	       a.service_id, s.name AS service_name, s.duration_minutes,
	       a.client_name, a.client_email, a.client_phone, a.client_birth_date, a.starts_at, a.ends_at,
	       a.status, a.source, a.notes, a.created_at
	  FROM appointments a
	  JOIN professionals p ON p.id = a.professional_id
	  JOIN services s ON s.id = a.service_id`

func (r *repository) ListManagerAppointments(ctx context.Context, establishmentID string, filter AppointmentFilter) ([]ManagerAppointmentRow, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PerPage < 1 {
		filter.PerPage = 20
	}

	where := " WHERE a.establishment_id = ?"
	args := []any{establishmentID}

	if filter.Date != "" {
		where += " AND DATE(CONVERT_TZ(a.starts_at, '+00:00', '+00:00')) = ?"
		args = append(args, filter.Date)
	}
	if filter.ProfessionalID != "" {
		where += " AND a.professional_id = ?"
		args = append(args, filter.ProfessionalID)
	}
	if filter.Status != "" {
		where += " AND a.status = ?"
		args = append(args, filter.Status)
	}

	var total int
	countQuery := "SELECT COUNT(*) FROM appointments a" + where
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.PerPage
	query := managerAppointmentSelect + where + " ORDER BY a.starts_at ASC LIMIT ? OFFSET ?"
	args = append(args, filter.PerPage, offset)

	var rows []ManagerAppointmentRow
	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func nullableString(value string) *string {
	if value == "" {
		return nil
	}
	v := value
	return &v
}

func (r *repository) GetManagerAppointment(ctx context.Context, establishmentID, appointmentID string) (*ManagerAppointmentRow, error) {
	var row ManagerAppointmentRow
	query := managerAppointmentSelect + " WHERE a.establishment_id = ? AND a.id = ?"
	if err := r.db.GetContext(ctx, &row, query, establishmentID, appointmentID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.ErrNotFound
		}
		return nil, err
	}
	return &row, nil
}

func (r *repository) UpdateAppointmentStatus(ctx context.Context, input UpdateStatusInput) (*ManagerAppointmentRow, error) {
	now := time.Now().UTC()
	result, err := r.db.ExecContext(ctx, `
		UPDATE appointments
		   SET status = ?, updated_at = ?
		 WHERE establishment_id = ?
		   AND id = ?
		   AND status NOT IN ('cancelled', 'completed', 'no_show')`,
		input.Status, now, input.EstablishmentID, input.AppointmentID,
	)
	if err != nil {
		return nil, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		// Pode ser not found ou status já terminal
		return r.GetManagerAppointment(ctx, input.EstablishmentID, input.AppointmentID)
	}
	return r.GetManagerAppointment(ctx, input.EstablishmentID, input.AppointmentID)
}

func (r *repository) ListManagerBlockedPeriods(ctx context.Context, establishmentID, professionalID, date string) ([]ManagerBlockedPeriod, error) {
	where := " WHERE p.establishment_id = ?"
	args := []any{establishmentID}

	if professionalID != "" {
		where += " AND bp.professional_id = ?"
		args = append(args, professionalID)
	}
	if date != "" {
		where += " AND DATE(bp.starts_at) = ?"
		args = append(args, date)
	}

	query := `
		SELECT bp.id, bp.professional_id, p.name AS professional_name,
		       bp.starts_at, bp.ends_at, bp.reason
		  FROM blocked_periods bp
		  JOIN professionals p ON p.id = bp.professional_id` + where + `
		 ORDER BY bp.starts_at ASC`

	var rows []ManagerBlockedPeriod
	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *repository) CreateBlockedPeriod(ctx context.Context, input CreateBlockedPeriodInput) (*ManagerBlockedPeriod, error) {
	// Valida que o profissional pertence ao estabelecimento
	var profName string
	err := r.db.GetContext(ctx, &profName,
		`SELECT name FROM professionals WHERE id = ? AND establishment_id = ? AND active = true`,
		input.ProfessionalID, input.EstablishmentID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.ErrNotFound
		}
		return nil, err
	}

	id := shared.NewID()

	var reason *string
	if input.Reason != "" {
		reason = &input.Reason
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO blocked_periods (id, professional_id, starts_at, ends_at, reason)
		VALUES (?, ?, ?, ?, ?)`,
		id, input.ProfessionalID, input.StartsAt.UTC(), input.EndsAt.UTC(), reason,
	)
	if err != nil {
		return nil, err
	}

	return &ManagerBlockedPeriod{
		ID:               id,
		ProfessionalID:   input.ProfessionalID,
		ProfessionalName: profName,
		StartsAt:         input.StartsAt.UTC(),
		EndsAt:           input.EndsAt.UTC(),
		Reason:           reason,
	}, nil
}

func (r *repository) DeleteBlockedPeriod(ctx context.Context, establishmentID, blockedPeriodID string) error {
	result, err := r.db.ExecContext(ctx, `
		DELETE bp FROM blocked_periods bp
		  JOIN professionals p ON p.id = bp.professional_id
		 WHERE bp.id = ?
		   AND p.establishment_id = ?`,
		blockedPeriodID, establishmentID,
	)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return shared.ErrNotFound
	}
	return nil
}
