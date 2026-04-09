package calendar

import (
	"context"
	"fmt"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"schedule/internal/shared"
	"schedule/internal/scheduling"
	"schedule/pkg/gcal"
)

// Service encapsula a logica de negocio da integracao Google Agenda.
type Service struct {
	repo Repository
}

// NewService cria um novo Service com as dependencias injetadas.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func mapConfigError(err error) error {
	if err == nil {
		return nil
	}

	if strings.Contains(err.Error(), "missing required env") {
		return shared.ErrIntegrationNotConfigured
	}

	return err
}

// GenerateAuthURL gera a URL de autorizacao OAuth2.
// O state e o establishmentID para ser recuperado no callback.
func (s *Service) GenerateAuthURL(ctx context.Context, establishmentID string) (string, error) {
	if establishmentID == "" {
		return "", shared.ErrUnauthorized
	}

	url, err := gcal.AuthCodeURL(establishmentID)
	if err != nil {
		return "", mapConfigError(err)
	}

	return url, nil
}

// HandleCallback processa o retorno do OAuth2:
// 1. Troca o code por tokens
// 2. Salva os tokens criptografados no banco
// 3. Cria agendas individuais para todos os profissionais do estabelecimento
// 4. Marca o estabelecimento como google_calendar_connected = true
func (s *Service) HandleCallback(ctx context.Context, code, state string) error {
	if code == "" || state == "" {
		return shared.ErrInvalidInput
	}

	establishmentID := state

	token, err := gcal.Exchange(ctx, code)
	if err != nil {
		if mapped := mapConfigError(err); mapped != err {
			return mapped
		}
		return fmt.Errorf("calendar: oauth exchange failed: %w", err)
	}

	scope := ""
	if tokenScope, ok := token.Extra("scope").(string); ok {
		scope = tokenScope
	}

	oauthToken := &OAuthToken{
		EstablishmentID: establishmentID,
		AccessToken:     token.AccessToken,
		RefreshToken:    token.RefreshToken,
		Expiry:          token.Expiry,
		Scope:           scope,
		UpdatedAt:       time.Now().UTC(),
	}

	if err := s.repo.SaveToken(ctx, oauthToken); err != nil {
		return fmt.Errorf("calendar: failed to save token: %w", err)
	}

	if err := s.createCalendarsForProfessionals(ctx, establishmentID, token); err != nil {
		_ = err
	}

	if err := s.repo.SetGoogleCalendarConnected(ctx, establishmentID, true); err != nil {
		return fmt.Errorf("calendar: failed to mark establishment as connected: %w", err)
	}

	return nil
}

// createCalendarsForProfessionals cria um Google Calendar para cada profissional
// ativo que ainda nao possui um.
func (s *Service) createCalendarsForProfessionals(ctx context.Context, establishmentID string, token *oauth2.Token) error {
	client, err := gcal.NewClient(ctx, token)
	if err != nil {
		return mapConfigError(err)
	}

	professionals, err := s.repo.ListProfessionals(ctx, establishmentID)
	if err != nil {
		return err
	}

	for _, p := range professionals {
		if p.GoogleCalendarID != nil && *p.GoogleCalendarID != "" {
			continue
		}

		calID, err := client.CreateCalendar(ctx, p.Name)
		if err != nil {
			continue
		}

		_ = s.repo.UpdateCalendarID(ctx, p.ID, calID)
	}

	return nil
}

// Disconnect revoga a integracao com o Google Agenda.
func (s *Service) Disconnect(ctx context.Context, establishmentID string) error {
	if establishmentID == "" {
		return shared.ErrUnauthorized
	}

	err := s.repo.DeleteToken(ctx, establishmentID)
	if err != nil && err != shared.ErrNotFound {
		return fmt.Errorf("calendar: failed to delete token: %w", err)
	}

	if err := s.repo.SetGoogleCalendarConnected(ctx, establishmentID, false); err != nil {
		return fmt.Errorf("calendar: failed to update establishment status: %w", err)
	}

	return nil
}

// GetStatus retorna se o estabelecimento tem a integracao ativa
// e a lista de profissionais com seus calendarios.
func (s *Service) GetStatus(ctx context.Context, establishmentID string) (*StatusResult, error) {
	if establishmentID == "" {
		return nil, shared.ErrUnauthorized
	}

	professionals, err := s.repo.ListProfessionals(ctx, establishmentID)
	if err != nil {
		return nil, err
	}

	_, tokenErr := s.repo.GetToken(ctx, establishmentID)
	connected := tokenErr == nil

	return &StatusResult{
		Connected:     connected,
		Professionals: professionals,
	}, nil
}

// ProvisionProfessional cria uma agenda Google para um profissional recém-criado,
// se o estabelecimento já tiver a integração ativa.
func (s *Service) ProvisionProfessional(ctx context.Context, establishmentID, professionalID, professionalName string) error {
	token, err := s.repo.GetToken(ctx, establishmentID)
	if err != nil {
		// Estabelecimento sem integração — silencioso.
		return nil
	}

	client, err := gcal.NewClient(ctx, &oauth2.Token{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
	})
	if err != nil {
		return nil
	}

	calID, err := client.CreateCalendar(ctx, professionalName)
	if err != nil {
		return nil
	}

	return s.repo.UpdateCalendarID(ctx, professionalID, calID)
}

// DeprovisionProfessional exclui a agenda Google de um profissional ao removê-lo.
func (s *Service) DeprovisionProfessional(ctx context.Context, establishmentID, calendarID string) error {
	token, err := s.repo.GetToken(ctx, establishmentID)
	if err != nil {
		return nil
	}

	client, err := gcal.NewClient(ctx, &oauth2.Token{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
	})
	if err != nil {
		return nil
	}

	_ = client.DeleteCalendar(ctx, calendarID)
	return nil
}

// RefreshTokenIfNeeded renova o access token se ele expirar nos proximos 5 minutos.
func (s *Service) RefreshTokenIfNeeded(ctx context.Context, establishmentID string) (*OAuthToken, error) {
	token, err := s.repo.GetToken(ctx, establishmentID)
	if err != nil {
		return nil, err
	}

	if time.Until(token.Expiry) > 5*time.Minute {
		return token, nil
	}

	_, err = gcal.NewClient(ctx, &oauth2.Token{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
	})
	if err != nil {
		return nil, mapConfigError(err)
	}

	client, err := gcal.NewClient(ctx, &oauth2.Token{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
	})
	if err != nil {
		return nil, mapConfigError(err)
	}

	refreshed, err := client.RefreshToken(ctx, token.RefreshToken)
	if err != nil {
		if mapped := mapConfigError(err); mapped != err {
			return nil, mapped
		}
		return nil, fmt.Errorf("calendar: token refresh failed: %w", err)
	}

	scope := token.Scope
	if s2, ok := refreshed.Extra("scope").(string); ok && s2 != "" {
		scope = s2
	}

	updated := &OAuthToken{
		EstablishmentID: establishmentID,
		AccessToken:     refreshed.AccessToken,
		RefreshToken:    refreshed.RefreshToken,
		Expiry:          refreshed.Expiry,
		Scope:           scope,
		UpdatedAt:       time.Now().UTC(),
	}

	if updated.RefreshToken == "" {
		updated.RefreshToken = token.RefreshToken
	}

	if err := s.repo.SaveToken(ctx, updated); err != nil {
		return nil, err
	}

	return updated, nil
}

// ListBusyPeriods retorna os periodos ocupados do calendario Google de um profissional.
// Se o estabelecimento nao estiver conectado ou o profissional ainda nao possuir calendario,
// retorna uma lista vazia para nao bloquear o fluxo publico.
func (s *Service) ListBusyPeriods(ctx context.Context, establishmentID, professionalID string, from, to time.Time) ([]scheduling.Period, error) {
	token, err := s.RefreshTokenIfNeeded(ctx, establishmentID)
	if err != nil {
		if err == shared.ErrNotFound {
			return nil, nil
		}
		if mapped := mapConfigError(err); mapped == shared.ErrIntegrationNotConfigured {
			return nil, nil
		} else if mapped != err {
			return nil, mapped
		}
		return nil, err
	}

	calendarID, err := s.repo.GetProfessionalCalendarID(ctx, establishmentID, professionalID)
	if err != nil {
		if err == shared.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}

	client, err := gcal.NewClient(ctx, &oauth2.Token{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
	})
	if err != nil {
		if mapped := mapConfigError(err); mapped == shared.ErrIntegrationNotConfigured {
			return nil, nil
		} else if mapped != err {
			return nil, mapped
		}
		return nil, err
	}

	periods, err := client.ListBusyPeriods(ctx, calendarID, from, to)
	if err != nil {
		return nil, err
	}

	busy := make([]scheduling.Period, 0, len(periods))
	for _, period := range periods {
		busy = append(busy, scheduling.Period{
			StartsAt: period.Start.UTC(),
			EndsAt:   period.End.UTC(),
		})
	}

	return busy, nil
}

func (s *Service) CreateAppointmentEvent(ctx context.Context, appointment *scheduling.Appointment) (string, error) {
	if appointment == nil {
		return "", shared.ErrInvalidInput
	}

	client, calendarID, err := s.appointmentCalendarClient(ctx, appointment.EstablishmentID, appointment.ProfessionalID)
	if err != nil {
		return "", err
	}
	if client == nil || calendarID == "" {
		return "", nil
	}

	title := fmt.Sprintf("Agendamento - %s", appointment.ClientName)
	notes := fmt.Sprintf("Appointment ID: %s\nTelefone: %s", appointment.ID, appointment.ClientPhone)
	serviceName, err := s.repo.GetServiceName(ctx, appointment.ServiceID, appointment.EstablishmentID)
	if err == nil && serviceName != "" {
		notes += fmt.Sprintf("\nServico: %s", serviceName)
	}
	if appointment.ClientEmail != nil && *appointment.ClientEmail != "" {
		notes += fmt.Sprintf("\nEmail: %s", *appointment.ClientEmail)
	}

	return client.CreateEvent(ctx, calendarID, gcal.Event{
		Title:     title,
		StartTime: appointment.StartsAt,
		EndTime:   appointment.EndsAt,
		Notes:     notes,
	})
}

func (s *Service) DeleteAppointmentEvent(ctx context.Context, appointment *scheduling.Appointment) error {
	if appointment == nil || appointment.GoogleEventID == nil || *appointment.GoogleEventID == "" {
		return nil
	}

	client, calendarID, err := s.appointmentCalendarClient(ctx, appointment.EstablishmentID, appointment.ProfessionalID)
	if err != nil {
		return err
	}
	if client == nil || calendarID == "" {
		return nil
	}

	return client.DeleteEvent(ctx, calendarID, *appointment.GoogleEventID)
}

func (s *Service) appointmentCalendarClient(ctx context.Context, establishmentID, professionalID string) (*gcal.Client, string, error) {
	token, err := s.RefreshTokenIfNeeded(ctx, establishmentID)
	if err != nil {
		if err == shared.ErrNotFound {
			return nil, "", nil
		}
		if mapped := mapConfigError(err); mapped == shared.ErrIntegrationNotConfigured {
			return nil, "", nil
		} else if mapped != err {
			return nil, "", mapped
		}
		return nil, "", err
	}

	calendarID, err := s.repo.GetProfessionalCalendarID(ctx, establishmentID, professionalID)
	if err != nil {
		if err == shared.ErrNotFound {
			return nil, "", nil
		}
		return nil, "", err
	}

	client, err := gcal.NewClient(ctx, &oauth2.Token{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
	})
	if err != nil {
		if mapped := mapConfigError(err); mapped == shared.ErrIntegrationNotConfigured {
			return nil, "", nil
		} else if mapped != err {
			return nil, "", mapped
		}
		return nil, "", err
	}

	return client, calendarID, nil
}
