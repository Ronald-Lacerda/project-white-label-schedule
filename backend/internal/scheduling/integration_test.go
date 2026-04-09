package scheduling

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"

	"schedule/internal/shared"

	_ "github.com/go-sql-driver/mysql"
)

func TestPublicAppointmentIntegration_CreateAndCancel(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ensureTCP(t, "127.0.0.1:3306")
	ensureTCP(t, "127.0.0.1:6379")

	fixture := newSchedulingIntegrationFixture(t)

	created := fixture.createAppointment(t, map[string]any{
		"service_id":        fixture.serviceID,
		"professional_id":   fixture.professionalID,
		"starts_at":         fixture.startsAt.Format(time.RFC3339),
		"client_name":       "Maria Oliveira",
		"client_email":      "maria@example.com",
		"client_phone":      "11999990000",
		"client_birth_date": "1993-05-14",
		"idempotency_key":   "8c77b028-6cba-4425-ab26-9b00d71fd2f2",
	})

	if created["status"] != "confirmed" {
		t.Fatalf("expected confirmed status, got %v", created["status"])
	}
	if created["client_email"] != "maria@example.com" {
		t.Fatalf("expected client_email in response, got %v", created["client_email"])
	}
	if created["client_birth_date"] != "1993-05-14" {
		t.Fatalf("expected client_birth_date in response, got %v", created["client_birth_date"])
	}

	appointmentID, _ := created["id"].(string)
	if appointmentID == "" {
		t.Fatal("expected created appointment id")
	}

	fixture.assertAppointmentStatus(t, appointmentID, "confirmed")

	cancelled := fixture.cancelAppointment(t, appointmentID, map[string]any{
		"phone": "11999990000",
	})

	if cancelled["status"] != "cancelled" {
		t.Fatalf("expected cancelled status, got %v", cancelled["status"])
	}

	canCancel, ok := cancelled["can_cancel"].(bool)
	if !ok {
		t.Fatalf("expected boolean can_cancel, got %T", cancelled["can_cancel"])
	}
	if canCancel {
		t.Fatal("expected cancelled appointment to be non-cancellable")
	}

	fixture.assertAppointmentStatus(t, appointmentID, "cancelled")
}

type schedulingIntegrationFixture struct {
	db              *sqlx.DB
	redis           *redis.Client
	router          http.Handler
	databaseName    string
	establishmentID string
	serviceID       string
	professionalID  string
	slug            string
	startsAt        time.Time
}

func newSchedulingIntegrationFixture(t *testing.T) *schedulingIntegrationFixture {
	t.Helper()

	rootDB, databaseName := openIsolatedMySQLDatabase(t)
	db := connectMySQLDatabase(t, mysqlDSN("root", "rootpassword", databaseName, false))
	applyMySQLMigrations(t, db)

	redisClient := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   15,
	})
	if err := redisClient.FlushDB(context.Background()).Err(); err != nil {
		t.Fatalf("flush redis: %v", err)
	}

	fixture := &schedulingIntegrationFixture{
		db:              db,
		redis:           redisClient,
		databaseName:    databaseName,
		establishmentID: shared.NewID(),
		serviceID:       shared.NewID(),
		professionalID:  shared.NewID(),
		slug:            "studio-integracao",
		startsAt:        time.Date(2030, 4, 8, 10, 0, 0, 0, time.UTC),
	}

	fixture.seed(t)
	fixture.router = fixture.buildRouter()

	t.Cleanup(func() {
		_ = fixture.redis.FlushDB(context.Background()).Err()
		_ = fixture.redis.Close()
		_ = fixture.db.Close()
		_, _ = rootDB.Exec("DROP DATABASE `" + databaseName + "`")
		_ = rootDB.Close()
	})

	return fixture
}

func (f *schedulingIntegrationFixture) buildRouter() http.Handler {
	repo := NewRepository(f.db)
	svc := NewService(repo, f.redis)
	handler := NewHandler(svc)

	resolver := func(ctx context.Context, slug string) (string, error) {
		if slug != f.slug {
			return "", shared.ErrNotFound
		}
		return f.establishmentID, nil
	}

	r := chi.NewRouter()
	r.Route("/pub/{slug}", func(r chi.Router) {
		r.Use(shared.SlugTenantMiddleware(resolver))
		r.Post("/appointments", handler.CreatePublicAppointment)
		r.Patch("/appointments/{id}/cancel", handler.CancelPublicAppointment)
	})

	return r
}

func (f *schedulingIntegrationFixture) seed(t *testing.T) {
	t.Helper()

	now := time.Date(2026, 4, 8, 12, 0, 0, 0, time.UTC)

	_, err := f.db.Exec(`
		INSERT INTO establishments (
			id, name, slug, timezone, contact_email, contact_phone,
			min_advance_cancel_hours, active, google_calendar_connected, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		f.establishmentID, "Studio Integracao", f.slug, "UTC", "contato@studio.test", "11999990000",
		2, true, false, now, now,
	)
	if err != nil {
		t.Fatalf("insert establishment: %v", err)
	}

	for day := 0; day < 7; day++ {
		if _, err := f.db.Exec(`
			INSERT INTO business_hours (id, establishment_id, day_of_week, open_time, close_time, is_closed)
			VALUES (?, ?, ?, ?, ?, ?)`,
			shared.NewID(), f.establishmentID, day, "09:00:00", "18:00:00", false,
		); err != nil {
			t.Fatalf("insert business_hours day %d: %v", day, err)
		}
	}

	_, err = f.db.Exec(`
		INSERT INTO professionals (
			id, establishment_id, name, avatar_url, phone, google_calendar_id,
			display_order, active, created_at, updated_at
		) VALUES (?, ?, ?, NULL, ?, NULL, ?, ?, ?, ?)`,
		f.professionalID, f.establishmentID, "Ana Souza", "11988887777", 1, true, now, now,
	)
	if err != nil {
		t.Fatalf("insert professional: %v", err)
	}

	for day := 0; day < 7; day++ {
		if _, err := f.db.Exec(`
			INSERT INTO professional_hours (id, professional_id, day_of_week, start_time, end_time, is_unavailable)
			VALUES (?, ?, ?, ?, ?, ?)`,
			shared.NewID(), f.professionalID, day, "09:00:00", "18:00:00", false,
		); err != nil {
			t.Fatalf("insert professional_hours day %d: %v", day, err)
		}
	}

	_, err = f.db.Exec(`
		INSERT INTO services (
			id, establishment_id, name, description, duration_minutes, price_cents,
			active, display_order, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		f.serviceID, f.establishmentID, "Corte", "Corte completo", 60, 5000, true, 1, now,
	)
	if err != nil {
		t.Fatalf("insert service: %v", err)
	}

	_, err = f.db.Exec(`
		INSERT INTO professional_services (professional_id, service_id)
		VALUES (?, ?)`,
		f.professionalID, f.serviceID,
	)
	if err != nil {
		t.Fatalf("insert professional_service: %v", err)
	}
}

func (f *schedulingIntegrationFixture) createAppointment(t *testing.T, body map[string]any) map[string]any {
	t.Helper()
	return f.doJSONRequest(t, http.MethodPost, "/pub/"+f.slug+"/appointments", body, http.StatusCreated)
}

func (f *schedulingIntegrationFixture) cancelAppointment(t *testing.T, appointmentID string, body map[string]any) map[string]any {
	t.Helper()
	return f.doJSONRequest(t, http.MethodPatch, "/pub/"+f.slug+"/appointments/"+appointmentID+"/cancel", body, http.StatusOK)
}

func (f *schedulingIntegrationFixture) doJSONRequest(t *testing.T, method, path string, body map[string]any, wantStatus int) map[string]any {
	t.Helper()

	payload, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	req := httptest.NewRequest(method, path, bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	f.router.ServeHTTP(rec, req)

	if rec.Code != wantStatus {
		t.Fatalf("%s %s returned %d: %s", method, path, rec.Code, rec.Body.String())
	}

	var response struct {
		Data map[string]any `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Data == nil {
		t.Fatalf("expected response data, got body %s", rec.Body.String())
	}
	return response.Data
}

func (f *schedulingIntegrationFixture) assertAppointmentStatus(t *testing.T, appointmentID, want string) {
	t.Helper()

	var got string
	err := f.db.Get(&got, `SELECT status FROM appointments WHERE id = ?`, appointmentID)
	if err != nil {
		t.Fatalf("select appointment status: %v", err)
	}
	if got != want {
		t.Fatalf("expected appointment status %q, got %q", want, got)
	}
}

func openIsolatedMySQLDatabase(t *testing.T) (*sql.DB, string) {
	t.Helper()

	root, err := sql.Open("mysql", mysqlDSN("root", "rootpassword", "", true))
	if err != nil {
		t.Fatalf("open mysql root connection: %v", err)
	}
	if err := root.Ping(); err != nil {
		t.Fatalf("ping mysql root connection: %v", err)
	}

	databaseName := fmt.Sprintf("schedule_integration_%d", time.Now().UTC().UnixNano())
	if _, err := root.Exec("CREATE DATABASE `" + databaseName + "` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci"); err != nil {
		_ = root.Close()
		t.Fatalf("create test database: %v", err)
	}

	return root, databaseName
}

func connectMySQLDatabase(t *testing.T, dsn string) *sqlx.DB {
	t.Helper()

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		t.Fatalf("connect mysql database: %v", err)
	}
	return db
}

func applyMySQLMigrations(t *testing.T, db *sqlx.DB) {
	t.Helper()

	files, err := filepath.Glob(filepath.Join("..", "..", "migrations", "*.up.sql"))
	if err != nil {
		t.Fatalf("list migrations: %v", err)
	}
	sort.Strings(files)

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			t.Fatalf("read migration %s: %v", file, err)
		}
		statement := strings.TrimSpace(string(content))
		if statement == "" {
			continue
		}
		if _, err := db.Exec(statement); err != nil {
			t.Fatalf("apply migration %s: %v", file, err)
		}
	}
}

func mysqlDSN(user, password, database string, multiStatements bool) string {
	base := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/", user, password)
	if database != "" {
		base += database
	}

	params := []string{
		"parseTime=true",
		"loc=UTC",
		"charset=utf8mb4",
		"collation=utf8mb4_unicode_ci",
	}
	if multiStatements {
		params = append(params, "multiStatements=true")
	}

	return base + "?" + strings.Join(params, "&")
}

func ensureTCP(t *testing.T, address string) {
	t.Helper()

	conn, err := net.DialTimeout("tcp", address, 2*time.Second)
	if err != nil {
		t.Skipf("integration dependency unavailable at %s: %v", address, err)
	}
	_ = conn.Close()
}
