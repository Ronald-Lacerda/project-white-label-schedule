package gcal

import (
	"context"
	"fmt"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

// Event representa um evento a ser criado no Google Agenda.
type Event struct {
	Title     string
	StartTime time.Time
	EndTime   time.Time
	Notes     string
}

// Period representa um intervalo de tempo ocupado.
type Period struct {
	Start time.Time
	End   time.Time
}

// Client e o wrapper da Google Calendar API.
type Client struct {
	svc   *calendar.Service
	token *oauth2.Token
	cfg   *oauth2.Config
}

func requiredEnv(name string) (string, error) {
	value := os.Getenv(name)
	if value == "" {
		return "", fmt.Errorf("gcal: missing required env %s", name)
	}
	return value, nil
}

// ValidateConfig verifica se as variaveis obrigatorias do Google OAuth estao configuradas.
func ValidateConfig() error {
	required := []string{
		"GOOGLE_CLIENT_ID",
		"GOOGLE_CLIENT_SECRET",
		"GOOGLE_REDIRECT_URL",
	}

	for _, key := range required {
		if _, err := requiredEnv(key); err != nil {
			return err
		}
	}

	return nil
}

// oauthConfig le as credenciais do ambiente e retorna a configuracao OAuth2.
func oauthConfig() (*oauth2.Config, error) {
	clientID, err := requiredEnv("GOOGLE_CLIENT_ID")
	if err != nil {
		return nil, err
	}

	clientSecret, err := requiredEnv("GOOGLE_CLIENT_SECRET")
	if err != nil {
		return nil, err
	}

	redirectURL, err := requiredEnv("GOOGLE_REDIRECT_URL")
	if err != nil {
		return nil, err
	}

	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			calendar.CalendarScope,
		},
		Endpoint: google.Endpoint,
	}, nil
}

// NewClient cria um Client autenticado com o token OAuth2 fornecido.
func NewClient(ctx context.Context, token *oauth2.Token) (*Client, error) {
	cfg, err := oauthConfig()
	if err != nil {
		return nil, err
	}

	httpClient := cfg.Client(ctx, token)
	svc, err := calendar.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, fmt.Errorf("gcal: failed to create calendar service: %w", err)
	}

	return &Client{svc: svc, token: token, cfg: cfg}, nil
}

// CreateCalendar cria uma nova agenda secundaria com o nome fornecido e retorna seu ID.
func (c *Client) CreateCalendar(ctx context.Context, name string) (string, error) {
	cal := &calendar.Calendar{
		Summary:  name,
		TimeZone: "UTC",
	}

	created, err := c.svc.Calendars.Insert(cal).Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("gcal: failed to create calendar %q: %w", name, err)
	}

	return created.Id, nil
}

// CreateEvent cria um evento no calendario especificado e retorna o ID do evento criado.
func (c *Client) CreateEvent(ctx context.Context, calendarID string, event Event) (string, error) {
	gEvent := &calendar.Event{
		Summary:     event.Title,
		Description: event.Notes,
		Start: &calendar.EventDateTime{
			DateTime: event.StartTime.UTC().Format(time.RFC3339),
			TimeZone: "UTC",
		},
		End: &calendar.EventDateTime{
			DateTime: event.EndTime.UTC().Format(time.RFC3339),
			TimeZone: "UTC",
		},
	}

	created, err := c.svc.Events.Insert(calendarID, gEvent).Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("gcal: failed to create event: %w", err)
	}

	return created.Id, nil
}

// DeleteEvent remove um evento do calendario. Ignora o erro 404 (ja deletado).
func (c *Client) DeleteEvent(ctx context.Context, calendarID, eventID string) error {
	err := c.svc.Events.Delete(calendarID, eventID).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("gcal: failed to delete event %q: %w", eventID, err)
	}
	return nil
}

// ListBusyPeriods retorna os periodos ocupados em um calendario entre from e to.
func (c *Client) ListBusyPeriods(ctx context.Context, calendarID string, from, to time.Time) ([]Period, error) {
	req := &calendar.FreeBusyRequest{
		TimeMin: from.UTC().Format(time.RFC3339),
		TimeMax: to.UTC().Format(time.RFC3339),
		Items:   []*calendar.FreeBusyRequestItem{{Id: calendarID}},
	}

	resp, err := c.svc.Freebusy.Query(req).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("gcal: failed to query freebusy: %w", err)
	}

	cal, ok := resp.Calendars[calendarID]
	if !ok {
		return nil, nil
	}

	periods := make([]Period, 0, len(cal.Busy))
	for _, busy := range cal.Busy {
		start, err := time.Parse(time.RFC3339, busy.Start)
		if err != nil {
			continue
		}
		end, err := time.Parse(time.RFC3339, busy.End)
		if err != nil {
			continue
		}
		periods = append(periods, Period{Start: start, End: end})
	}

	return periods, nil
}

// RefreshToken renova o access token usando o refresh token fornecido.
func (c *Client) RefreshToken(ctx context.Context, refreshToken string) (*oauth2.Token, error) {
	cfg, err := oauthConfig()
	if err != nil {
		return nil, err
	}

	expiredToken := &oauth2.Token{
		RefreshToken: refreshToken,
		Expiry:       time.Now().Add(-time.Hour),
	}

	tokenSource := cfg.TokenSource(ctx, expiredToken)
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("gcal: failed to refresh token: %w", err)
	}

	return newToken, nil
}

// AuthCodeURL gera a URL de autorizacao OAuth2 com o state fornecido.
func AuthCodeURL(state string) (string, error) {
	cfg, err := oauthConfig()
	if err != nil {
		return "", err
	}

	return cfg.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce), nil
}

// Exchange troca o codigo de autorizacao por um token OAuth2.
func Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	cfg, err := oauthConfig()
	if err != nil {
		return nil, err
	}

	token, err := cfg.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("gcal: failed to exchange auth code: %w", err)
	}
	return token, nil
}
