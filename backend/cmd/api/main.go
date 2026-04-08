package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"

	"schedule/internal/auth"
	"schedule/internal/calendar"
	"schedule/internal/catalog"
	"schedule/internal/scheduling"
	"schedule/internal/shared"
	"schedule/internal/tenancy"
	"schedule/internal/whitelabel"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, reading from environment")
	}

	db := mustConnectMySQL()
	defer db.Close()

	rdb := mustConnectRedis()
	defer rdb.Close()

	// ── Repositories ──────────────────────────────────────────────────────────
	authRepo := auth.NewRepository(db)
	tenancyRepo := tenancy.NewRepository(db)
	wlRepo := whitelabel.NewRepository(db)
	profRepo := catalog.NewProfessionalRepository(db)
	svcRepo := catalog.NewSvcRepository(db)
	calendarRepo := calendar.NewRepository(db)
	schedulingRepo := scheduling.NewRepository(db)

	// ── Services ──────────────────────────────────────────────────────────────
	jwtSecret := mustEnv("JWT_SECRET")
	jwtExpiry := parseDuration(os.Getenv("JWT_EXPIRY"), 15*time.Minute)
	rtExpiry := parseDuration(os.Getenv("REFRESH_TOKEN_EXPIRY"), 7*24*time.Hour)

	authSvc := auth.NewService(authRepo, rdb, jwtSecret, jwtExpiry, rtExpiry)
	tenancySvc := tenancy.NewService(tenancyRepo)
	wlSvc := whitelabel.NewService(wlRepo)
	profSvc := catalog.NewProfessionalService(profRepo)
	svcSvc := catalog.NewSvcService(svcRepo)
	calendarSvc := calendar.NewService(calendarRepo)
	schedulingSvc := scheduling.NewService(schedulingRepo, rdb)
	profSvc.WithCalendar(calendarSvc)

	// ── Handlers ──────────────────────────────────────────────────────────────
	uploadsDir := "uploads/logos"
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	authHandler := auth.NewHandler(authSvc)
	authMiddleware := auth.NewMiddleware(authSvc)
	tenancyHandler := tenancy.NewHandler(tenancySvc)
	wlHandler := whitelabel.NewHandler(wlSvc, uploadsDir, baseURL)
	profHandler := catalog.NewProfessionalHandler(profSvc)
	svcHandler := catalog.NewSvcHandler(svcSvc)
	calendarHandler := calendar.NewHandler(calendarSvc)
	schedulingHandler := scheduling.NewHandler(schedulingSvc)

	// ── Router ────────────────────────────────────────────────────────────────
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(corsMiddleware)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Arquivos estáticos (logos enviados pelo gestor)
	r.Handle("/uploads/*", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	// ── API do Gestor ─────────────────────────────────────────────────────────
	r.Route("/api/v1", func(r chi.Router) {
		// Auth (público)
		r.Post("/auth/login", authHandler.Login)
		r.Post("/auth/register", authHandler.Register)
		r.Post("/auth/refresh", authHandler.Refresh)
		r.Get("/google/callback", calendarHandler.Callback)

		// Rotas autenticadas
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.Authenticate)

			r.Post("/auth/logout", authHandler.Logout)

			// Estabelecimento
			r.Get("/establishment", tenancyHandler.Get)
			r.Put("/establishment", tenancyHandler.Update)
			r.Get("/establishment/business-hours", tenancyHandler.GetBusinessHours)
			r.Put("/establishment/business-hours", tenancyHandler.UpdateBusinessHours)

			// Whitelabel
			r.Get("/whitelabel", wlHandler.Get)
			r.Put("/whitelabel", wlHandler.Update)
			r.Post("/whitelabel/logo", wlHandler.UploadLogo)

			// Profissionais
			r.Get("/professionals", profHandler.List)
			r.Post("/professionals", profHandler.Create)
			r.Get("/professionals/{id}", profHandler.Get)
			r.Put("/professionals/{id}", profHandler.Update)
			r.Delete("/professionals/{id}", profHandler.Delete)
			r.Put("/professionals/{id}/hours", profHandler.UpdateHours)
			r.Put("/professionals/{id}/services", profHandler.UpdateServices)

			// Serviços
			r.Get("/services", svcHandler.List)
			r.Post("/services", svcHandler.Create)
			r.Put("/services/{id}", svcHandler.Update)
			r.Delete("/services/{id}", svcHandler.Delete)

			// Google Agenda
			r.Get("/google/auth-url", calendarHandler.GetAuthURL)
			r.Delete("/google/disconnect", calendarHandler.Disconnect)
			r.Get("/google/status", calendarHandler.GetStatus)

			// Gestão de agendamentos (Fase 10)
			r.Get("/appointments", schedulingHandler.ListManagerAppointments)
			r.Get("/appointments/{id}", schedulingHandler.GetManagerAppointment)
			r.Patch("/appointments/{id}/status", schedulingHandler.UpdateAppointmentStatus)

			// Bloqueios de agenda (Fase 10)
			r.Get("/blocked-periods", schedulingHandler.ListManagerBlockedPeriods)
			r.Post("/blocked-periods", schedulingHandler.CreateBlockedPeriod)
			r.Delete("/blocked-periods/{id}", schedulingHandler.DeleteBlockedPeriod)
		})
	})

	// ── API Pública ───────────────────────────────────────────────────────────
	slugResolver := func(ctx context.Context, slug string) (string, error) {
		e, err := tenancySvc.GetBySlug(ctx, slug)
		if err != nil {
			return "", err
		}
		return e.ID, nil
	}

	r.Route("/pub/{slug}", func(r chi.Router) {
		r.Use(shared.SlugTenantMiddleware(slugResolver))

		// Endpoints públicos — implementados na Fase 9
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			estID := shared.EstablishmentIDFromContext(r.Context())
			e, err := tenancySvc.GetByID(r.Context(), estID)
			if err != nil {
				shared.JSONError(w, err)
				return
			}
			wl, _ := wlSvc.Get(r.Context(), estID)
			shared.JSON(w, http.StatusOK, map[string]any{
				"establishment": e,
				"whitelabel":    wl,
			})
		})

		r.Get("/services", func(w http.ResponseWriter, r *http.Request) {
			estID := shared.EstablishmentIDFromContext(r.Context())
			list, err := svcSvc.List(r.Context(), estID)
			if err != nil {
				shared.JSONError(w, err)
				return
			}
			shared.JSON(w, http.StatusOK, list)
		})

		r.Get("/professionals", func(w http.ResponseWriter, r *http.Request) {
			estID := shared.EstablishmentIDFromContext(r.Context())
			serviceID := r.URL.Query().Get("service_id")

			if serviceID != "" {
				list, err := schedulingRepo.GetActiveProfessionals(r.Context(), estID, serviceID)
				if err != nil {
					shared.JSONError(w, err)
					return
				}
				shared.JSON(w, http.StatusOK, list)
				return
			}

			list, err := profSvc.List(r.Context(), estID)
			if err != nil {
				shared.JSONError(w, err)
				return
			}
			shared.JSON(w, http.StatusOK, list)
		})

		r.Get("/availability", schedulingHandler.GetAvailability)
		r.Post("/appointments", schedulingHandler.CreatePublicAppointment)
		r.Get("/appointments/{id}", schedulingHandler.GetPublicAppointment)
		r.Patch("/appointments/{id}/cancel", schedulingHandler.CancelPublicAppointment)
		r.Patch("/appointments/{id}/reschedule", schedulingHandler.ReschedulePublicAppointment)
	})

	// ── Servidor ──────────────────────────────────────────────────────────────
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("server listening on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("forced shutdown: %v", err)
	}
	log.Println("server stopped")
}

func corsMiddleware(next http.Handler) http.Handler {
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = strings.Join([]string{
			"http://localhost:3000",
			"http://localhost:3001",
			"http://127.0.0.1:3000",
			"http://127.0.0.1:3001",
		}, ",")
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		for _, allowed := range strings.Split(allowedOrigins, ",") {
			if strings.TrimSpace(allowed) == origin {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func mustConnectMySQL() *sqlx.DB {
	dsn := mustEnv("DATABASE_URL")
	dsn = strings.TrimPrefix(dsn, "mysql://")

	db, err := sqlx.Connect("mysql", dsn+"?parseTime=true&loc=UTC&charset=utf8mb4&collation=utf8mb4_unicode_ci")
	if err != nil {
		log.Fatalf("failed to connect to mysql: %v", err)
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	log.Println("connected to mysql")
	return db
}

func mustConnectRedis() *redis.Client {
	addr := mustEnv("REDIS_URL")
	opt, err := redis.ParseURL(addr)
	if err != nil {
		log.Fatalf("failed to parse REDIS_URL: %v", err)
	}
	rdb := redis.NewClient(opt)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
	log.Println("connected to redis")
	return rdb
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("environment variable %s is required", key)
	}
	return v
}

func parseDuration(s string, fallback time.Duration) time.Duration {
	if s == "" {
		return fallback
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return fallback
	}
	return d
}
