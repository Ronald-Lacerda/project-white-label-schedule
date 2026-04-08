package calendar

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"schedule/internal/shared"
	"schedule/pkg/crypto"
)

// Repository define as operações de persistência do módulo calendar.
type Repository interface {
	SaveToken(ctx context.Context, token *OAuthToken) error
	GetToken(ctx context.Context, establishmentID string) (*OAuthToken, error)
	DeleteToken(ctx context.Context, establishmentID string) error
	UpdateCalendarID(ctx context.Context, professionalID, calendarID string) error
	ListProfessionals(ctx context.Context, establishmentID string) ([]ProfessionalRef, error)
	SetGoogleCalendarConnected(ctx context.Context, establishmentID string, connected bool) error
}

type repository struct {
	db *sqlx.DB
}

// NewRepository cria uma implementação concreta de Repository.
func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

// SaveToken persiste ou atualiza o token OAuth2 de um estabelecimento.
// Os tokens são criptografados com AES-256-GCM antes de gravar no banco.
func (r *repository) SaveToken(ctx context.Context, token *OAuthToken) error {
	encAccess, err := crypto.Encrypt(token.AccessToken)
	if err != nil {
		return err
	}

	encRefresh, err := crypto.Encrypt(token.RefreshToken)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	_, err = r.db.ExecContext(ctx, `
		INSERT INTO google_oauth_tokens
			(establishment_id, access_token, refresh_token, expiry, scope, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			access_token  = VALUES(access_token),
			refresh_token = VALUES(refresh_token),
			expiry        = VALUES(expiry),
			scope         = VALUES(scope),
			updated_at    = VALUES(updated_at)`,
		token.EstablishmentID,
		encAccess,
		encRefresh,
		token.Expiry.UTC(),
		token.Scope,
		now,
	)
	return err
}

// GetToken busca e descriptografa o token OAuth2 de um estabelecimento.
func (r *repository) GetToken(ctx context.Context, establishmentID string) (*OAuthToken, error) {
	type row struct {
		EstablishmentID string    `db:"establishment_id"`
		AccessToken     string    `db:"access_token"`
		RefreshToken    string    `db:"refresh_token"`
		Expiry          time.Time `db:"expiry"`
		Scope           string    `db:"scope"`
		UpdatedAt       time.Time `db:"updated_at"`
	}

	var enc row
	err := r.db.GetContext(ctx, &enc,
		`SELECT * FROM google_oauth_tokens WHERE establishment_id = ?`,
		establishmentID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.ErrNotFound
		}
		return nil, err
	}

	access, err := crypto.Decrypt(enc.AccessToken)
	if err != nil {
		return nil, err
	}

	refresh, err := crypto.Decrypt(enc.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &OAuthToken{
		EstablishmentID: enc.EstablishmentID,
		AccessToken:     access,
		RefreshToken:    refresh,
		Expiry:          enc.Expiry,
		Scope:           enc.Scope,
		UpdatedAt:       enc.UpdatedAt,
	}, nil
}

// DeleteToken remove o token OAuth2 de um estabelecimento.
func (r *repository) DeleteToken(ctx context.Context, establishmentID string) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM google_oauth_tokens WHERE establishment_id = ?`,
		establishmentID,
	)
	return err
}

// UpdateCalendarID atualiza o google_calendar_id de um profissional.
func (r *repository) UpdateCalendarID(ctx context.Context, professionalID, calendarID string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE professionals SET google_calendar_id = ?, updated_at = ? WHERE id = ?`,
		calendarID, time.Now().UTC(), professionalID,
	)
	return err
}

// ListProfessionals retorna a lista de profissionais ativos de um estabelecimento
// com os campos necessários para gestão de calendários.
func (r *repository) ListProfessionals(ctx context.Context, establishmentID string) ([]ProfessionalRef, error) {
	var result []ProfessionalRef
	err := r.db.SelectContext(ctx, &result,
		`SELECT id, name, google_calendar_id
		 FROM professionals
		 WHERE establishment_id = ? AND active = true
		 ORDER BY display_order, name`,
		establishmentID,
	)
	return result, err
}

// SetGoogleCalendarConnected atualiza a flag google_calendar_connected do estabelecimento.
func (r *repository) SetGoogleCalendarConnected(ctx context.Context, establishmentID string, connected bool) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE establishments SET google_calendar_connected = ?, updated_at = ? WHERE id = ?`,
		connected, time.Now().UTC(), establishmentID,
	)
	return err
}
