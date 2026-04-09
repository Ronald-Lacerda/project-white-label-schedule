package catalog

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"schedule/internal/shared"
)

// ─── Professional Repository ─────────────────────────────────────────────────

type ProfessionalRepository interface {
	List(ctx context.Context, establishmentID string) ([]Professional, error)
	FindByID(ctx context.Context, id, establishmentID string) (*Professional, error)
	Create(ctx context.Context, p *Professional) error
	Update(ctx context.Context, p *Professional) error
	SoftDelete(ctx context.Context, id, establishmentID string) error
	GetHours(ctx context.Context, professionalID string) ([]ProfessionalHour, error)
	UpsertHours(ctx context.Context, hours []ProfessionalHour) error
	SetServices(ctx context.Context, professionalID, establishmentID string, serviceIDs []string) error
	ListServiceIDs(ctx context.Context, professionalID string) ([]string, error)
}

type professionalRepo struct{ db *sqlx.DB }

func NewProfessionalRepository(db *sqlx.DB) ProfessionalRepository {
	return &professionalRepo{db: db}
}

func (r *professionalRepo) List(ctx context.Context, establishmentID string) ([]Professional, error) {
	var result []Professional
	err := r.db.SelectContext(ctx, &result,
		`SELECT * FROM professionals WHERE establishment_id = ? AND active = true ORDER BY display_order, name`,
		establishmentID,
	)
	if err != nil {
		return nil, err
	}
	if err := r.attachServiceIDs(ctx, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *professionalRepo) FindByID(ctx context.Context, id, establishmentID string) (*Professional, error) {
	var p Professional
	err := r.db.GetContext(ctx, &p,
		`SELECT * FROM professionals WHERE id = ? AND establishment_id = ?`,
		id, establishmentID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.ErrNotFound
		}
		return nil, err
	}
	serviceIDs, err := r.ListServiceIDs(ctx, p.ID)
	if err != nil {
		return nil, err
	}
	p.ServiceIDs = serviceIDs
	return &p, nil
}

func (r *professionalRepo) Create(ctx context.Context, p *Professional) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO professionals (id, establishment_id, name, avatar_url, phone, display_order, active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		p.ID, p.EstablishmentID, p.Name, p.AvatarURL, p.Phone,
		p.DisplayOrder, p.Active, p.CreatedAt, p.UpdatedAt,
	)
	return err
}

func (r *professionalRepo) Update(ctx context.Context, p *Professional) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE professionals SET name = ?, avatar_url = ?, phone = ?, display_order = ?, updated_at = ?
		WHERE id = ? AND establishment_id = ?`,
		p.Name, p.AvatarURL, p.Phone, p.DisplayOrder, p.UpdatedAt,
		p.ID, p.EstablishmentID,
	)
	return err
}

func (r *professionalRepo) SoftDelete(ctx context.Context, id, establishmentID string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE professionals SET active = false, updated_at = ?
		WHERE id = ? AND establishment_id = ?`,
		time.Now().UTC(), id, establishmentID,
	)
	return err
}

func (r *professionalRepo) GetHours(ctx context.Context, professionalID string) ([]ProfessionalHour, error) {
	var hours []ProfessionalHour
	err := r.db.SelectContext(ctx, &hours,
		`SELECT * FROM professional_hours WHERE professional_id = ? ORDER BY day_of_week`,
		professionalID,
	)
	return hours, err
}

func (r *professionalRepo) UpsertHours(ctx context.Context, hours []ProfessionalHour) error {
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
			INSERT INTO professional_hours (id, professional_id, day_of_week, start_time, end_time, is_unavailable)
			VALUES (?, ?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE
				start_time = VALUES(start_time),
				end_time = VALUES(end_time),
				is_unavailable = VALUES(is_unavailable)`,
			h.ID, h.ProfessionalID, h.DayOfWeek, h.StartTime, h.EndTime, h.IsUnavailable,
		)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *professionalRepo) SetServices(ctx context.Context, professionalID, establishmentID string, serviceIDs []string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if len(serviceIDs) > 0 {
		placeholders := make([]string, 0, len(serviceIDs))
		args := make([]any, 0, len(serviceIDs)+1)
		args = append(args, establishmentID)
		for _, svcID := range serviceIDs {
			placeholders = append(placeholders, "?")
			args = append(args, svcID)
		}

		query := `
			SELECT COUNT(DISTINCT id)
			FROM services
			WHERE establishment_id = ? AND active = true AND id IN (` + strings.Join(placeholders, ",") + `)`

		var matched int
		if err := tx.GetContext(ctx, &matched, query, args...); err != nil {
			return err
		}
		if matched != len(uniqueStrings(serviceIDs)) {
			return shared.ErrInvalidInput
		}
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM professional_services WHERE professional_id = ?`, professionalID); err != nil {
		return err
	}

	for _, svcID := range uniqueStrings(serviceIDs) {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO professional_services (professional_id, service_id) VALUES (?, ?)`,
			professionalID, svcID,
		); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *professionalRepo) ListServiceIDs(ctx context.Context, professionalID string) ([]string, error) {
	var ids []string
	err := r.db.SelectContext(ctx, &ids,
		`SELECT service_id FROM professional_services WHERE professional_id = ?`,
		professionalID,
	)
	return ids, err
}

func (r *professionalRepo) attachServiceIDs(ctx context.Context, professionals []Professional) error {
	if len(professionals) == 0 {
		return nil
	}

	idsByProfessional := make(map[string][]string, len(professionals))
	for _, professional := range professionals {
		serviceIDs, err := r.ListServiceIDs(ctx, professional.ID)
		if err != nil {
			return err
		}
		idsByProfessional[professional.ID] = serviceIDs
	}

	for i := range professionals {
		professionals[i].ServiceIDs = idsByProfessional[professionals[i].ID]
	}

	return nil
}

func uniqueStrings(values []string) []string {
	if len(values) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}

// ─── Service Repository ───────────────────────────────────────────────────────

type SvcRepository interface {
	List(ctx context.Context, establishmentID string) ([]Svc, error)
	FindByID(ctx context.Context, id, establishmentID string) (*Svc, error)
	Create(ctx context.Context, s *Svc) error
	Update(ctx context.Context, s *Svc) error
	SoftDelete(ctx context.Context, id, establishmentID string) error
}

type svcRepo struct{ db *sqlx.DB }

func NewSvcRepository(db *sqlx.DB) SvcRepository {
	return &svcRepo{db: db}
}

func (r *svcRepo) List(ctx context.Context, establishmentID string) ([]Svc, error) {
	var result []Svc
	err := r.db.SelectContext(ctx, &result,
		`SELECT * FROM services WHERE establishment_id = ? ORDER BY display_order, name`,
		establishmentID,
	)
	return result, err
}

func (r *svcRepo) FindByID(ctx context.Context, id, establishmentID string) (*Svc, error) {
	var s Svc
	err := r.db.GetContext(ctx, &s,
		`SELECT * FROM services WHERE id = ? AND establishment_id = ?`,
		id, establishmentID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.ErrNotFound
		}
		return nil, err
	}
	return &s, nil
}

func (r *svcRepo) Create(ctx context.Context, s *Svc) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO services (id, establishment_id, name, description, duration_minutes, price_cents, active, display_order, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		s.ID, s.EstablishmentID, s.Name, s.Description,
		s.DurationMinutes, s.PriceCents, s.Active, s.DisplayOrder, s.CreatedAt,
	)
	return err
}

func (r *svcRepo) Update(ctx context.Context, s *Svc) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE services SET name = ?, description = ?, duration_minutes = ?, price_cents = ?, display_order = ?
		WHERE id = ? AND establishment_id = ?`,
		s.Name, s.Description, s.DurationMinutes, s.PriceCents, s.DisplayOrder,
		s.ID, s.EstablishmentID,
	)
	return err
}

func (r *svcRepo) SoftDelete(ctx context.Context, id, establishmentID string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE services SET active = false WHERE id = ? AND establishment_id = ?`,
		id, establishmentID,
	)
	return err
}
