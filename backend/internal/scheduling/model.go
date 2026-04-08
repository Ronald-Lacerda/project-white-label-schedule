package scheduling

import "time"

// AppointmentFilter agrupa os filtros opcionais para listagem de agendamentos pelo gestor.
type AppointmentFilter struct {
	Date           string // YYYY-MM-DD, opcional
	ProfessionalID string // opcional
	Status         string // opcional: confirmed, completed, no_show, cancelled
	Page           int
	PerPage        int
}

// ManagerAppointmentRow é o agendamento enriquecido retornado ao gestor.
type ManagerAppointmentRow struct {
	ID               string    `db:"id" json:"id"`
	EstablishmentID  string    `db:"establishment_id" json:"-"`
	ProfessionalID   string    `db:"professional_id" json:"professional_id"`
	ProfessionalName string    `db:"professional_name" json:"professional_name"`
	ServiceID        string    `db:"service_id" json:"service_id"`
	ServiceName      string    `db:"service_name" json:"service_name"`
	DurationMinutes  int       `db:"duration_minutes" json:"duration_minutes"`
	ClientName       string    `db:"client_name" json:"client_name"`
	ClientPhone      string    `db:"client_phone" json:"client_phone"`
	StartsAt         time.Time `db:"starts_at" json:"starts_at"`
	EndsAt           time.Time `db:"ends_at" json:"ends_at"`
	Status           string    `db:"status" json:"status"`
	Source           string    `db:"source" json:"source"`
	Notes            *string   `db:"notes" json:"notes,omitempty"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
}

// UpdateStatusInput é usado pelo gestor para alterar o status de um agendamento.
type UpdateStatusInput struct {
	EstablishmentID string
	AppointmentID   string
	Status          string // completed, no_show, cancelled
}

// CreateBlockedPeriodInput é usado pelo gestor para criar um bloqueio de agenda.
type CreateBlockedPeriodInput struct {
	EstablishmentID string
	ProfessionalID  string
	StartsAt        time.Time
	EndsAt          time.Time
	Reason          string
}

// ManagerBlockedPeriod é o bloqueio enriquecido com nome do profissional.
type ManagerBlockedPeriod struct {
	ID               string    `db:"id" json:"id"`
	ProfessionalID   string    `db:"professional_id" json:"professional_id"`
	ProfessionalName string    `db:"professional_name" json:"professional_name"`
	StartsAt         time.Time `db:"starts_at" json:"starts_at"`
	EndsAt           time.Time `db:"ends_at" json:"ends_at"`
	Reason           *string   `db:"reason" json:"reason,omitempty"`
}

type Appointment struct {
	ID              string    `db:"id"`
	EstablishmentID string    `db:"establishment_id"`
	ProfessionalID  string    `db:"professional_id"`
	ServiceID       string    `db:"service_id"`
	ClientName      string    `db:"client_name"`
	ClientPhone     string    `db:"client_phone"`
	StartsAt        time.Time `db:"starts_at"`
	EndsAt          time.Time `db:"ends_at"`
	Status          string    `db:"status"`
	Source          string    `db:"source"`
	GoogleEventID   *string   `db:"google_event_id"`
	Notes           *string   `db:"notes"`
	IdempotencyKey  *string   `db:"idempotency_key"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}

type CreateAppointmentInput struct {
	EstablishmentID string
	ServiceID       string
	ProfessionalID  string
	StartsAt        time.Time
	ClientName      string
	ClientPhone     string
	IdempotencyKey  string
}

type BlockedPeriod struct {
	ID             string    `db:"id"`
	ProfessionalID string    `db:"professional_id"`
	StartsAt       time.Time `db:"starts_at"`
	EndsAt         time.Time `db:"ends_at"`
	Reason         *string   `db:"reason"`
	GoogleEventID  *string   `db:"google_event_id"`
}

type Slot struct {
	StartsAt       time.Time `json:"starts_at"`
	EndsAt         time.Time `json:"ends_at"`
	ProfessionalID string    `json:"professional_id"`
}

type PublicAppointmentResult struct {
	ID             string    `json:"id"`
	ServiceID      string    `json:"service_id"`
	ProfessionalID string    `json:"professional_id"`
	ClientName     string    `json:"client_name"`
	ClientPhone    string    `json:"client_phone"`
	StartsAt       time.Time `json:"starts_at"`
	EndsAt         time.Time `json:"ends_at"`
	Status         string    `json:"status"`
}

type PublicAppointmentDetail struct {
	ID                    string    `json:"id"`
	ServiceID             string    `json:"service_id"`
	ProfessionalID        string    `json:"professional_id"`
	ClientName            string    `json:"client_name"`
	ClientPhone           string    `json:"client_phone"`
	StartsAt              time.Time `json:"starts_at"`
	EndsAt                time.Time `json:"ends_at"`
	Status                string    `json:"status"`
	CanCancel             bool      `json:"can_cancel"`
	MinAdvanceCancelHours int       `json:"min_advance_cancel_hours"`
}
