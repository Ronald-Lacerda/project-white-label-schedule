package whitelabel

import "context"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Get(ctx context.Context, establishmentID string) (*Config, error) {
	return s.repo.FindByEstablishment(ctx, establishmentID)
}

func (s *Service) Update(ctx context.Context, establishmentID string, input UpdateInput) (*Config, error) {
	existing, err := s.repo.FindByEstablishment(ctx, establishmentID)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		EstablishmentID: establishmentID,
		LogoURL:         existing.LogoURL,
		PrimaryColor:    existing.PrimaryColor,
		SecondaryColor:  existing.SecondaryColor,
		CustomCSS:       existing.CustomCSS,
		CustomDomain:    existing.CustomDomain,
	}

	if input.LogoURL != nil {
		cfg.LogoURL = input.LogoURL
	}

	if input.PrimaryColor != "" {
		cfg.PrimaryColor = input.PrimaryColor
	}

	if input.SecondaryColor != nil {
		cfg.SecondaryColor = input.SecondaryColor
	}

	if input.CustomCSS != nil {
		cfg.CustomCSS = input.CustomCSS
	}

	if cfg.PrimaryColor == "" {
		cfg.PrimaryColor = "#000000"
	}

	if err := s.repo.Upsert(ctx, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

type UpdateInput struct {
	LogoURL        *string
	PrimaryColor   string
	SecondaryColor *string
	CustomCSS      *string
}

func (s *Service) UpdateLogoURL(ctx context.Context, establishmentID, logoURL string) (*Config, error) {
	existing, err := s.repo.FindByEstablishment(ctx, establishmentID)
	if err != nil {
		return nil, err
	}
	existing.LogoURL = &logoURL
	if err := s.repo.Upsert(ctx, existing); err != nil {
		return nil, err
	}
	return existing, nil
}
