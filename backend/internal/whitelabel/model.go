package whitelabel

type Config struct {
	EstablishmentID string  `db:"establishment_id" json:"establishment_id"`
	LogoURL         *string `db:"logo_url"         json:"logo_url"`
	PrimaryColor    string  `db:"primary_color"    json:"primary_color"`
	SecondaryColor  *string `db:"secondary_color"  json:"secondary_color"`
	CustomDomain    *string `db:"custom_domain"    json:"custom_domain"`
	CustomCSS       *string `db:"custom_css"       json:"custom_css"`
}
