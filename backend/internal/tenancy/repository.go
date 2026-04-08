package tenancy

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"schedule/internal/shared"
)

type Repository interface {
	FindByID(ctx context.Context, id string) (*Establishment, error)
	FindBySlug(ctx context.Context, slug string) (*Establishment, error)
	Update(ctx context.Context, e *Establishment) error
	SlugExists(ctx context.Context, slug, excludeID string) (bool, error)
	GetBusinessHours(ctx context.Context, establishmentID string) ([]BusinessHour, error)
	UpsertBusinessHours(ctx context.Context, hours []BusinessHour) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindByID(ctx context.Context, id string) (*Establishment, error) {
	var e Establishment
	err := r.db.GetContext(ctx, &e, `SELECT * FROM establishments WHERE id = ?`, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.ErrNotFound
		}
		return nil, err
	}
	return &e, nil
}

func (r *repository) FindBySlug(ctx context.Context, slug string) (*Establishment, error) {
	var e Establishment
	err := r.db.GetContext(ctx, &e, `SELECT * FROM establishments WHERE slug = ? AND active = true`, slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.ErrNotFound
		}
		return nil, err
	}
	return &e, nil
}

func (r *repository) SlugExists(ctx context.Context, slug, excludeID string) (bool, error) {
	var count int
	err := r.db.GetContext(ctx, &count,
		`SELECT COUNT(*) FROM establishments WHERE slug = ? AND id != ?`,
		slug, excludeID,
	)
	return count > 0, err
}

func (r *repository) Update(ctx context.Context, e *Establishment) error {
	e.UpdatedAt = time.Now().UTC()
	_, err := r.db.ExecContext(ctx, `
		UPDATE establishments SET
			name = ?, slug = ?, timezone = ?,
			contact_email = ?, contact_phone = ?,
			min_advance_cancel_hours = ?, updated_at = ?
		WHERE id = ?`,
		e.Name, e.Slug, e.Timezone,
		e.ContactEmail, e.ContactPhone,
		e.MinAdvanceCancelHours, e.UpdatedAt,
		e.ID,
	)
	return err
}

func (r *repository) GetBusinessHours(ctx context.Context, establishmentID string) ([]BusinessHour, error) {
	var hours []BusinessHour
	err := r.db.SelectContext(ctx, &hours,
		`SELECT * FROM business_hours WHERE establishment_id = ? ORDER BY day_of_week`,
		establishmentID,
	)
	return hours, err
}

func (r *repository) UpsertBusinessHours(ctx context.Context, hours []BusinessHour) error {
	if len(hours) == 0 {
		return nil
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, h := range hours {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO business_hours (id, establishment_id, day_of_week, open_time, close_time, is_closed)
			VALUES (?, ?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE
				open_time = VALUES(open_time),
				close_time = VALUES(close_time),
				is_closed = VALUES(is_closed)`,
			h.ID, h.EstablishmentID, h.DayOfWeek, h.OpenTime, h.CloseTime, h.IsClosed,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
