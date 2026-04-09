package catalog

import "time"

type Professional struct {
	ID               string    `db:"id"                 json:"id"`
	EstablishmentID  string    `db:"establishment_id"   json:"establishment_id"`
	Name             string    `db:"name"               json:"name"`
	AvatarURL        *string   `db:"avatar_url"         json:"avatar_url"`
	Phone            *string   `db:"phone"              json:"phone"`
	GoogleCalendarID *string   `db:"google_calendar_id" json:"google_calendar_id"`
	ServiceIDs       []string  `db:"-"                  json:"service_ids"`
	DisplayOrder     int       `db:"display_order"      json:"display_order"`
	Active           bool      `db:"active"             json:"active"`
	CreatedAt        time.Time `db:"created_at"         json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"         json:"updated_at"`
}

type ProfessionalHour struct {
	ID             string `db:"id"              json:"id"`
	ProfessionalID string `db:"professional_id" json:"-"`
	DayOfWeek      int    `db:"day_of_week"     json:"day_of_week"`
	StartTime      string `db:"start_time"      json:"start_time"`
	EndTime        string `db:"end_time"        json:"end_time"`
	IsUnavailable  bool   `db:"is_unavailable"  json:"is_unavailable"`
}

type Svc struct {
	ID              string    `db:"id"               json:"id"`
	EstablishmentID string    `db:"establishment_id" json:"establishment_id"`
	Name            string    `db:"name"             json:"name"`
	Description     *string   `db:"description"      json:"description"`
	DurationMinutes int       `db:"duration_minutes" json:"duration_minutes"`
	PriceCents      *int      `db:"price_cents"      json:"price_cents"`
	Active          bool      `db:"active"           json:"active"`
	DisplayOrder    int       `db:"display_order"    json:"display_order"`
	CreatedAt       time.Time `db:"created_at"       json:"created_at"`
}
