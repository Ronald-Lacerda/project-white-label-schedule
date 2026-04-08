package tenancy

import "time"

type Establishment struct {
	ID                      string    `db:"id"                        json:"id"`
	Name                    string    `db:"name"                      json:"name"`
	Slug                    string    `db:"slug"                      json:"slug"`
	Timezone                string    `db:"timezone"                  json:"timezone"`
	ContactEmail            *string   `db:"contact_email"             json:"contact_email"`
	ContactPhone            *string   `db:"contact_phone"             json:"contact_phone"`
	MinAdvanceCancelHours   int       `db:"min_advance_cancel_hours"  json:"min_advance_cancel_hours"`
	Active                  bool      `db:"active"                    json:"active"`
	GoogleCalendarConnected bool      `db:"google_calendar_connected" json:"google_calendar_connected"`
	CreatedAt               time.Time `db:"created_at"                json:"created_at"`
	UpdatedAt               time.Time `db:"updated_at"                json:"updated_at"`
}

type BusinessHour struct {
	ID              string `db:"id"               json:"id"`
	EstablishmentID string `db:"establishment_id" json:"-"`
	DayOfWeek       int    `db:"day_of_week"      json:"day_of_week"`
	OpenTime        string `db:"open_time"        json:"open_time"`
	CloseTime       string `db:"close_time"       json:"close_time"`
	IsClosed        bool   `db:"is_closed"        json:"is_closed"`
}
