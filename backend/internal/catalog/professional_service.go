package catalog

import (
	"context"
	"time"

	"schedule/internal/shared"
)

// CalendarProvisioner é implementado por calendar.Service e injetado opcionalmente.
type CalendarProvisioner interface {
	ProvisionProfessional(ctx context.Context, establishmentID, professionalID, professionalName string) error
	DeprovisionProfessional(ctx context.Context, establishmentID, calendarID string) error
}

type ProfessionalService struct {
	repo    ProfessionalRepository
	calSvc  CalendarProvisioner
}

func NewProfessionalService(repo ProfessionalRepository) *ProfessionalService {
	return &ProfessionalService{repo: repo}
}

func (s *ProfessionalService) WithCalendar(cal CalendarProvisioner) *ProfessionalService {
	s.calSvc = cal
	return s
}

func (s *ProfessionalService) List(ctx context.Context, establishmentID string) ([]Professional, error) {
	return s.repo.List(ctx, establishmentID)
}

func (s *ProfessionalService) Get(ctx context.Context, id, establishmentID string) (*Professional, error) {
	return s.repo.FindByID(ctx, id, establishmentID)
}

func (s *ProfessionalService) Create(ctx context.Context, establishmentID string, in ProfessionalInput) (*Professional, error) {
	now := time.Now().UTC()
	p := &Professional{
		ID:              shared.NewID(),
		EstablishmentID: establishmentID,
		Name:            in.Name,
		AvatarURL:       in.AvatarURL,
		Phone:           in.Phone,
		DisplayOrder:    in.DisplayOrder,
		Active:          true,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}

	if s.calSvc != nil {
		_ = s.calSvc.ProvisionProfessional(ctx, establishmentID, p.ID, p.Name)
	}

	return p, nil
}

func (s *ProfessionalService) Update(ctx context.Context, id, establishmentID string, in ProfessionalInput) (*Professional, error) {
	p, err := s.repo.FindByID(ctx, id, establishmentID)
	if err != nil {
		return nil, err
	}

	p.Name = in.Name
	p.AvatarURL = in.AvatarURL
	p.Phone = in.Phone
	p.DisplayOrder = in.DisplayOrder
	p.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *ProfessionalService) Delete(ctx context.Context, id, establishmentID string) error {
	p, err := s.repo.FindByID(ctx, id, establishmentID)
	if err != nil {
		return err
	}

	if err := s.repo.SoftDelete(ctx, id, establishmentID); err != nil {
		return err
	}

	if s.calSvc != nil && p.GoogleCalendarID != nil && *p.GoogleCalendarID != "" {
		_ = s.calSvc.DeprovisionProfessional(ctx, establishmentID, *p.GoogleCalendarID)
	}

	return nil
}

func (s *ProfessionalService) GetHours(ctx context.Context, id, establishmentID string) ([]ProfessionalHour, error) {
	if _, err := s.repo.FindByID(ctx, id, establishmentID); err != nil {
		return nil, err
	}
	return s.repo.GetHours(ctx, id)
}

func (s *ProfessionalService) UpdateHours(ctx context.Context, id, establishmentID string, inputs []ProfessionalHourInput) ([]ProfessionalHour, error) {
	if _, err := s.repo.FindByID(ctx, id, establishmentID); err != nil {
		return nil, err
	}

	existing, _ := s.repo.GetHours(ctx, id)
	existingByDay := map[int]string{}
	for _, h := range existing {
		existingByDay[h.DayOfWeek] = h.ID
	}

	hours := make([]ProfessionalHour, 0, len(inputs))
	for _, in := range inputs {
		h := ProfessionalHour{
			ProfessionalID: id,
			DayOfWeek:      in.DayOfWeek,
			StartTime:      in.StartTime,
			EndTime:        in.EndTime,
			IsUnavailable:  in.IsUnavailable,
		}
		if existID, ok := existingByDay[in.DayOfWeek]; ok {
			h.ID = existID
		} else {
			h.ID = shared.NewID()
		}
		hours = append(hours, h)
	}

	if err := s.repo.UpsertHours(ctx, hours); err != nil {
		return nil, err
	}
	return hours, nil
}

func (s *ProfessionalService) UpdateServices(ctx context.Context, id, establishmentID string, serviceIDs []string) error {
	if _, err := s.repo.FindByID(ctx, id, establishmentID); err != nil {
		return err
	}
	return s.repo.SetServices(ctx, id, serviceIDs)
}

type ProfessionalInput struct {
	Name         string
	AvatarURL    *string
	Phone        *string
	DisplayOrder int
}

type ProfessionalHourInput struct {
	DayOfWeek     int
	StartTime     string
	EndTime       string
	IsUnavailable bool
}
