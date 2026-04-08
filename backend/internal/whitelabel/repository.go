package whitelabel

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"schedule/internal/shared"
)

type Repository interface {
	FindByEstablishment(ctx context.Context, establishmentID string) (*Config, error)
	Upsert(ctx context.Context, cfg *Config) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindByEstablishment(ctx context.Context, establishmentID string) (*Config, error) {
	var c Config
	err := r.db.GetContext(ctx, &c,
		`SELECT * FROM whitelabel_configs WHERE establishment_id = ?`,
		establishmentID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Retorna config padrão se ainda não existe
			return &Config{
				EstablishmentID: establishmentID,
				PrimaryColor:    "#000000",
			}, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *repository) Upsert(ctx context.Context, cfg *Config) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO whitelabel_configs (establishment_id, logo_url, primary_color, secondary_color, custom_domain, custom_css)
		VALUES (?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			logo_url = VALUES(logo_url),
			primary_color = VALUES(primary_color),
			secondary_color = VALUES(secondary_color),
			custom_domain = VALUES(custom_domain),
			custom_css = VALUES(custom_css)`,
		cfg.EstablishmentID, cfg.LogoURL, cfg.PrimaryColor,
		cfg.SecondaryColor, cfg.CustomDomain, cfg.CustomCSS,
	)
	return err
}

// Verifica isolamento — nunca retorna config de outro tenant
func (r *repository) assertTenant(ctx context.Context, establishmentID string) error {
	contextID := shared.EstablishmentIDFromContext(ctx)
	if contextID != "" && contextID != establishmentID {
		return shared.ErrForbidden
	}
	return nil
}
