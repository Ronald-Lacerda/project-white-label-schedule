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

func (s *SvcService) ListActive(ctx context.Context, establishmentID string) ([]Svc, error) {
	return s.repo.ListActive(ctx, establishmentID)
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

func (s *SvcService) Patch(ctx context.Context, id, establishmentID string, in SvcPatchInput) (*Svc, error) {
	svc, err := s.repo.FindByID(ctx, id, establishmentID)
	if err != nil {
		return nil, err
	}

	if in.Name != nil {
		svc.Name = *in.Name
	}
	if in.DescriptionProvided {
		svc.Description = in.Description
	}
	if in.DurationMinutes != nil {
		svc.DurationMinutes = *in.DurationMinutes
	}
	if in.PriceCentsProvided {
		svc.PriceCents = in.PriceCents
	}
	if in.DisplayOrder != nil {
		svc.DisplayOrder = *in.DisplayOrder
	}
	if in.Active != nil {
		svc.Active = *in.Active
	}

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

type SvcPatchInput struct {
	Name                *string
	Description         *string
	DescriptionProvided bool
	DurationMinutes     *int
	PriceCents          *int
	PriceCentsProvided  bool
	DisplayOrder        *int
	Active              *bool
}
