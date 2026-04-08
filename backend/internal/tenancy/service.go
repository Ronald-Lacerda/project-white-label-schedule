package tenancy

import (
	"context"

	"schedule/internal/shared"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetByID(ctx context.Context, id string) (*Establishment, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *Service) GetBySlug(ctx context.Context, slug string) (*Establishment, error) {
	return s.repo.FindBySlug(ctx, slug)
}

func (s *Service) Update(ctx context.Context, id string, input UpdateInput) (*Establishment, error) {
	e, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.Slug != e.Slug {
		exists, err := s.repo.SlugExists(ctx, input.Slug, id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, shared.ErrSlugConflict
		}
	}

	e.Name = input.Name
	e.Slug = input.Slug
	e.Timezone = input.Timezone
	e.ContactEmail = input.ContactEmail
	e.ContactPhone = input.ContactPhone
	e.MinAdvanceCancelHours = input.MinAdvanceCancelHours

	if err := s.repo.Update(ctx, e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *Service) GetBusinessHours(ctx context.Context, establishmentID string) ([]BusinessHour, error) {
	return s.repo.GetBusinessHours(ctx, establishmentID)
}

func (s *Service) UpdateBusinessHours(ctx context.Context, establishmentID string, inputs []BusinessHourInput) ([]BusinessHour, error) {
	hours := make([]BusinessHour, 0, len(inputs))
	for _, in := range inputs {
		h := BusinessHour{
			EstablishmentID: establishmentID,
			DayOfWeek:       in.DayOfWeek,
			OpenTime:        in.OpenTime,
			CloseTime:       in.CloseTime,
			IsClosed:        in.IsClosed,
		}
		// Busca ID existente ou gera novo
		existing, _ := s.repo.GetBusinessHours(ctx, establishmentID)
		for _, ex := range existing {
			if ex.DayOfWeek == in.DayOfWeek {
				h.ID = ex.ID
				break
			}
		}
		if h.ID == "" {
			h.ID = shared.NewID()
		}
		hours = append(hours, h)
	}

	if err := s.repo.UpsertBusinessHours(ctx, hours); err != nil {
		return nil, err
	}
	return hours, nil
}

type UpdateInput struct {
	Name                  string
	Slug                  string
	Timezone              string
	ContactEmail          *string
	ContactPhone          *string
	MinAdvanceCancelHours int
}

type BusinessHourInput struct {
	DayOfWeek int
	OpenTime  string
	CloseTime string
	IsClosed  bool
}
