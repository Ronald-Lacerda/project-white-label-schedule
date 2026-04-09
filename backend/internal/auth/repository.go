package auth

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"schedule/internal/shared"
)

type Repository interface {
	FindUserByEmail(ctx context.Context, email string) (*User, error)
	FindUserByID(ctx context.Context, id string) (*User, error)
	CreateAccount(ctx context.Context, input CreateAccountInput) (*User, error)
}

type repository struct {
	db *sqlx.DB
}

type CreateAccountInput struct {
	OwnerName         string
	EstablishmentName string
	Email             string
	PasswordHash      string
	Slug              string
	Timezone          string
	ContactPhone      *string
}

type businessHourSeed struct {
	DayOfWeek int
	OpenTime  string
	CloseTime string
	IsClosed  bool
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindUserByEmail(ctx context.Context, email string) (*User, error) {
	var u User
	err := r.db.GetContext(ctx, &u, `SELECT * FROM users WHERE LOWER(email) = LOWER(?) AND active = true`, strings.TrimSpace(email))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.ErrUnauthorized
		}
		return nil, err
	}
	return &u, nil
}

func (r *repository) FindUserByID(ctx context.Context, id string) (*User, error) {
	var u User
	err := r.db.GetContext(ctx, &u, `SELECT * FROM users WHERE id = ? AND active = true`, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *repository) CreateAccount(ctx context.Context, input CreateAccountInput) (*User, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var emailCount int
	if err := tx.GetContext(ctx, &emailCount, `SELECT COUNT(*) FROM users WHERE LOWER(email) = LOWER(?)`, input.Email); err != nil {
		return nil, err
	}
	if emailCount > 0 {
		return nil, shared.ErrEmailConflict
	}

	var slugCount int
	if err := tx.GetContext(ctx, &slugCount, `SELECT COUNT(*) FROM establishments WHERE slug = ?`, input.Slug); err != nil {
		return nil, err
	}
	if slugCount > 0 {
		return nil, shared.ErrSlugConflict
	}

	now := time.Now().UTC()
	establishmentID := shared.NewID()
	userID := shared.NewID()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO establishments (
			id, name, slug, timezone, contact_email, contact_phone,
			min_advance_cancel_hours, active, google_calendar_connected, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		establishmentID,
		input.EstablishmentName,
		input.Slug,
		input.Timezone,
		input.Email,
		input.ContactPhone,
		0,
		true,
		false,
		now,
		now,
	)
	if err != nil {
		return nil, err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO users (id, establishment_id, name, email, password_hash, role, active, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		userID,
		establishmentID,
		input.OwnerName,
		input.Email,
		input.PasswordHash,
		"owner",
		true,
		now,
	)
	if err != nil {
		return nil, err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO whitelabel_configs (establishment_id, logo_url, primary_color, secondary_color, custom_domain, custom_css)
		VALUES (?, ?, ?, ?, ?, ?)`,
		establishmentID,
		nil,
		"#000000",
		nil,
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}

	for _, hour := range defaultBusinessHourSeeds() {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO business_hours (id, establishment_id, day_of_week, open_time, close_time, is_closed)
			VALUES (?, ?, ?, ?, ?, ?)`,
			shared.NewID(),
			establishmentID,
			hour.DayOfWeek,
			hour.OpenTime,
			hour.CloseTime,
			hour.IsClosed,
		)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &User{
		ID:              userID,
		EstablishmentID: establishmentID,
		Name:            input.OwnerName,
		Email:           input.Email,
		PasswordHash:    input.PasswordHash,
		Role:            "owner",
		Active:          true,
		CreatedAt:       now,
	}, nil
}

func defaultBusinessHourSeeds() []businessHourSeed {
	hours := make([]businessHourSeed, 0, 7)
	for dayOfWeek := 0; dayOfWeek < 7; dayOfWeek++ {
		hours = append(hours, businessHourSeed{
			DayOfWeek: dayOfWeek,
			OpenTime:  "08:00:00",
			CloseTime: "18:00:00",
			IsClosed:  true,
		})
	}
	return hours
}
