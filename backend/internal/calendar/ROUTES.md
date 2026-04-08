# Rotas — Módulo Calendar (Google Agenda)

Registrar no `cmd/api/main.go` após as rotas existentes.

## Imports a adicionar

```go
import (
    "schedule/internal/calendar"
)
```

## Inicialização (após os outros repos/services/handlers)

```go
// Repositório e serviço do calendar
calendarRepo := calendar.NewRepository(db)
calendarSvc  := calendar.NewService(calendarRepo)
googleHandler := calendar.NewHandler(calendarSvc)
```

## Registro das rotas

Dentro do bloco `r.Route("/api/v1", ...)`:

```go
// Google Agenda — rota pública (callback OAuth — sem JWT)
r.Get("/google/callback", googleHandler.Callback)

// Google Agenda — rotas autenticadas (dentro do grupo com authMiddleware.Authenticate)
r.Get("/google/auth-url",    googleHandler.GetAuthURL)
r.Delete("/google/disconnect", googleHandler.Disconnect)
r.Get("/google/status",      googleHandler.GetStatus)
```

## Exemplo completo de onde inserir

```go
r.Route("/api/v1", func(r chi.Router) {
    // Auth (público)
    r.Post("/auth/login", authHandler.Login)
    r.Post("/auth/refresh", authHandler.Refresh)

    // Google OAuth callback — público (Google redireciona aqui sem JWT)
    r.Get("/google/callback", googleHandler.Callback)

    // Rotas autenticadas
    r.Group(func(r chi.Router) {
        r.Use(authMiddleware.Authenticate)

        // ... rotas existentes ...

        // Google Agenda
        r.Get("/google/auth-url",      googleHandler.GetAuthURL)
        r.Delete("/google/disconnect", googleHandler.Disconnect)
        r.Get("/google/status",        googleHandler.GetStatus)
    })
})
```

## Variável de ambiente adicional necessária

```env
# URL base do frontend (para redirect pós-callback)
FRONTEND_URL=http://localhost:3000
```
