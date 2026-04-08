## Rotas — Módulo de Disponibilidade (Fase 8)

```
GET /pub/{slug}/availability → schedulingHandler.GetAvailability (público, usa SlugTenantMiddleware)
```

### Dependências para main.go

```go
schedulingRepo    := scheduling.NewRepository(db)
schedulingSvc     := scheduling.NewService(schedulingRepo, rdb)
schedulingHandler := scheduling.NewHandler(schedulingSvc)
```

### Registro das rotas (Chi)

```go
r.Route("/pub/{slug}", func(r chi.Router) {
    r.Use(shared.SlugTenantMiddleware(tenancySvc.ResolveSlug))
    r.Get("/availability", schedulingHandler.GetAvailability)
    // POST /appointments → Fase 9
})
```

### Parâmetros da rota GET /pub/{slug}/availability

| Parâmetro       | Tipo   | Obrigatório | Descrição                              |
|-----------------|--------|-------------|----------------------------------------|
| service_id      | string | sim         | ID do serviço a agendar               |
| professional_id | string | não         | Filtra por profissional específico     |
| date_from       | string | sim         | Data inicial no formato YYYY-MM-DD    |
| date_to         | string | sim         | Data final no formato YYYY-MM-DD      |
| timezone        | string | não         | Fuso horário (default: America/Sao_Paulo) |

Restrição: `date_to - date_from <= 31 dias`

### Invalidação de cache (Fase 9)

Ao criar ou cancelar agendamentos, chamar:

```go
schedulingSvc.InvalidateCache(ctx, establishmentID, professionalID, timezone, date)
```
