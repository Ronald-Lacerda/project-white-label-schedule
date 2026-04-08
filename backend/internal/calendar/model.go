package calendar

import "time"

// OAuthToken representa as credenciais OAuth2 de um estabelecimento.
// Os campos AccessToken e RefreshToken são mantidos descriptografados em memória;
// a criptografia é aplicada somente ao persistir no banco.
type OAuthToken struct {
	EstablishmentID string
	AccessToken     string // descriptografado em memória
	RefreshToken    string // descriptografado em memória
	Expiry          time.Time
	Scope           string
	UpdatedAt       time.Time
}

// ProfessionalRef é uma projeção leve do profissional usada nas respostas do calendário.
type ProfessionalRef struct {
	ID               string  `db:"id"                 json:"id"`
	Name             string  `db:"name"               json:"name"`
	GoogleCalendarID *string `db:"google_calendar_id" json:"google_calendar_id"`
}

// StatusResult é o retorno do método GetStatus do service.
type StatusResult struct {
	Connected     bool              `json:"connected"`
	Professionals []ProfessionalRef `json:"professionals"`
}
