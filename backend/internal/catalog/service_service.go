package catalog

import (
	"context"
	"time"

	"schedule/internal/shared"
)

type SvcService struct {
	repo SvcRepository
}

func NewSvcService(repo SvcRepository) *SvcService {
	return &SvcService{repo: repo}
}

func (s *SvcService) List(ctx context.Context, establishmentID string) ([]Svc, error) {
	return s.repo.List(ctx, establishmentID)
}

func (s *SvcService) Get(ctx context.Context, id, establishmentID string) (*Svc, error) {
	return s.repo.FindByID(ctx, id, establishmentID)
}

func (s *SvcService) Create(ctx context.Context, establishmentID string, in SvcInput) (*Svc, error) {
	svc := &Svc{
		ID:              shared.NewID(),
		EstablishmentID: establishmentID,
		Name:            in.Name,
		Description:     in.Description,
		DurationMinutes: in.DurationMinutes,
		PriceCents:      in.PriceCents,
		DisplayOrder:    in.DisplayOrder,
		Active:          true,
		CreatedAt:       time.Now().UTC(),
	}
	if err := s.repo.Create(ctx, svc); err != nil {
		return nil, err
	}
	return svc, nil
}

func (s *SvcService) Update(ctx context.Context, id, establishmentID string, in SvcInput) (*Svc, error) {
	svc, err := s.repo.FindByID(ctx, id, establishmentID)
	if err != nil {
		return nil, err
	}

	svc.Name = in.Name
	svc.Description = in.Description
	svc.DurationMinutes = in.DurationMinutes
	svc.PriceCents = in.PriceCents
	svc.DisplayOrder = in.DisplayOrder

	if err := s.repo.Update(ctx, svc); err != nil {
		return nil, err
	}
	return svc, nil
}

func (s *SvcService) Delete(ctx context.Context, id, establishmentID string) error {
	if _, err := s.repo.FindByID(ctx, id, establishmentID); err != nil {
		return err
	}
	return s.repo.SoftDelete(ctx, id, establishmentID)
}

type SvcInput struct {
	Name            string
	Description     *string
	DurationMinutes int
	PriceCents      *int
	DisplayOrder    int
}
