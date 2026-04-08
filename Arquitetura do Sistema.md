# Arquitetura do Sistema — Plataforma Whitelabel de Agendamentos

> Documento de referência técnica. Versão 1.1 — Abril 2026.

---

## 1. Decisões de Stack

### 1.1 Linguagem backend: Go ✅ Decidido

Go é a escolha definitiva para o backend.

O coração desta plataforma é o **motor de disponibilidade** — calcular slots livres considerando jornada do profissional, duração do serviço, conflitos no Google Agenda e regras de bloqueio. Com goroutines, é possível verificar a disponibilidade de múltiplos profissionais **em paralelo** com latência mínima. Isso impacta diretamente a experiência do cliente final na tela de seleção de horário.

Binário único compilado simplifica deploy, reduz custo de infra no início e elimina dependências de runtime.

**Framework Go recomendado:** [Chi](https://github.com/go-chi/chi) ou **Echo**
- Chi: minimalista, idiomático, bom suporte a middleware
- Echo: mais baterias incluídas, familiaridade com estrutura de rotas mais explícita

### 1.2 Banco de dados principal: MySQL 8.x

Mantém o conhecimento existente. MySQL 8 tem:
- Window functions (útil para relatórios de agenda)
- JSON columns (para configurações flexíveis de whitelabel)
- CTEs e melhores índices compostos
- Suporte nativo a timezone via `CONVERT_TZ`

### 1.3 Cache e filas: Redis

- Cache de slots de disponibilidade (TTL curto, ~30s)
- Sessões de manager
- Fila de jobs (sincronização Google Agenda)
- Rate limiting por tenant no endpoint público de agendamento
- Idempotency keys para criação de agendamentos

### 1.4 Frontend: Vue.js ✅ Decidido

| Superfície | Stack | Motivo |
|---|---|---|
| Painel do gestor | **Nuxt.js** (Vue 3 + SSR opcional) | SPA com auth, dashboard, rotas organizadas por módulo |
| Página pública de agendamento | **Nuxt.js** (mesmo projeto, rota `/p/:slug`) | Reaproveitamento de componentes, SSR para SEO |

Usar Nuxt.js permite um **único projeto frontend** para as duas superfícies, separando por layout:
- `layouts/manager.vue` — painel autenticado
- `layouts/booking.vue` — página pública minimalista, mobile-first

---

## 2. Padrão Arquitetural: Monólito Modular

### Por que não microserviços agora

Microserviços introduzem latência de rede entre módulos, complexidade operacional (service discovery, tracing distribuído) e dificuldade de debug — tudo isso prejudica a velocidade de validação do MVP.

**Monólito Modular** é o padrão certo para este estágio:
- Um único processo deployável
- Módulos internos com fronteiras bem definidas (packages separados)
- Sem chamadas de rede entre módulos (chamadas diretas em memória)
- Fácil de extrair um módulo como serviço separado no futuro se houver necessidade de escala

### Estrutura de módulos internos

```
/
├── cmd/
│   └── api/            # Ponto de entrada — inicializa servidor HTTP
│
├── internal/
│   ├── tenancy/        # Gestão de estabelecimentos (tenant root)
│   ├── auth/           # Autenticação do gestor (JWT + Google OAuth2)
│   ├── catalog/        # Serviços e profissionais
│   ├── scheduling/     # Motor de disponibilidade + criação de agendamentos
│   ├── calendar/       # Integração Google Agenda
│   ├── notifications/  # E-mail e WhatsApp (abstração, implementações plugáveis)
│   ├── whitelabel/     # Configuração visual e domínio
│   └── shared/         # DTOs, erros, utilidades, middleware
│
├── pkg/
│   ├── gcal/           # Wrapper da Google Calendar API
│   ├── mailer/         # Abstração de envio de e-mail
│   └── whatsapp/       # Abstração de envio WhatsApp (Evolution API / Twilio)
│
├── migrations/         # Migrações de banco de dados (golang-migrate)
└── docker-compose.yml
```

---

## 3. Diagrama de Componentes

```
┌─────────────────────────────────────────────────────────────┐
│                        CLIENTE FINAL                        │
│            (acessa link público do estabelecimento)          │
└──────────────────────────┬──────────────────────────────────┘
                           │ HTTPS
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                  NGINX (Reverse Proxy + TLS)                 │
└────────────┬──────────────────────────────┬─────────────────┘
             │                              │
             ▼                              ▼
┌────────────────────────┐    ┌─────────────────────────────┐
│  API Pública           │    │  API do Gestor (autenticada) │
│  /pub/:slug/...        │    │  /api/v1/...                 │
│  (sem auth)            │    │  (JWT)                       │
└────────────┬───────────┘    └──────────────┬──────────────┘
             │                               │
             └───────────────┬───────────────┘
                             │
                    ┌────────▼────────┐
                    │  Go HTTP Server │
                    │  (Chi / Echo)   │
                    └────────┬────────┘
                             │
          ┌──────────────────┼──────────────────┐
          │                  │                  │
    ┌─────▼──────┐   ┌───────▼──────┐   ┌──────▼──────┐
    │  Scheduling │   │   Calendar   │   │   Catalog   │
    │   Module   │   │   Module     │   │   Module    │
    │            │   │  (Google)    │   │             │
    └─────┬──────┘   └───────┬──────┘   └─────────────┘
          │                  │
    ┌─────▼──────────────────▼─────┐
    │           MySQL 8             │
    │   (dados transacionais)       │
    └──────────────────────────────┘

    ┌──────────────────────────────┐
    │            Redis              │
    │  cache · sessões · filas      │
    └──────────────────────────────┘

    ┌──────────────────────────────┐
    │      Google Calendar API      │
    │  (OAuth2 · Events · Webhooks) │
    └──────────────────────────────┘
```

---

## 4. Multi-tenancy

### Estratégia: Shared Database, Shared Schema com tenant_id

Todas as tabelas possuem `establishment_id` como chave de isolamento. Um middleware de tenant resolve o estabelecimento a partir do slug na URL e injeta o ID no contexto da request.

**Fluxo de resolução do tenant:**

```
Request → GET /pub/barbearia-silva/services
           │
           ▼
     TenantMiddleware
           │
     SELECT id FROM establishments WHERE slug = 'barbearia-silva'
           │
     ctx = context.WithValue(ctx, tenantKey, establishmentID)
           │
     Handler recebe ctx com tenant resolvido
```

**Proteção contra vazamento cross-tenant:**

Todo repositório recebe o `establishment_id` do contexto e aplica como filtro obrigatório. Nunca aceita establishment_id vindo do request body para operações de leitura de dados sensíveis.

### Tabela de referência rápida: onde establishment_id aparece

| Tabela | Isolamento |
|---|---|
| professionals | establishment_id (FK) |
| services | establishment_id (FK) |
| appointments | establishment_id (FK) |
| business_hours | via establishment_id |
| blocked_periods | via professional_id → establishment_id |
| whitelabel_configs | establishment_id (PK/FK) |

---

## 5. Modelo de Dados

### 5.1 Diagrama ER simplificado

```
establishments
    │
    ├──< users (managers)
    ├──< business_hours
    ├──< professionals
    │        ├──< professional_hours
    │        ├──< professional_services >──┐
    │        └──< blocked_periods          │
    ├──< services ─────────────────────────┘
    ├──< appointments
    └──  whitelabel_configs
```

### 5.2 DDL das tabelas principais

```sql
-- Estabelecimentos (raiz do tenant)
CREATE TABLE establishments (
    id              CHAR(26) PRIMARY KEY,           -- ULID
    name            VARCHAR(120) NOT NULL,
    slug            VARCHAR(60) NOT NULL UNIQUE,    -- ex: barbearia-silva
    timezone        VARCHAR(50) NOT NULL DEFAULT 'America/Sao_Paulo',
    contact_email   VARCHAR(120),
    contact_phone   VARCHAR(20),
    active          BOOLEAN NOT NULL DEFAULT TRUE,
    google_calendar_connected BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      DATETIME(6) NOT NULL,
    updated_at      DATETIME(6) NOT NULL,
    INDEX idx_slug (slug)
);

-- Configurações de whitelabel (1:1 com establishment)
CREATE TABLE whitelabel_configs (
    establishment_id CHAR(26) PRIMARY KEY,
    logo_url         VARCHAR(500),
    primary_color    VARCHAR(7) NOT NULL DEFAULT '#000000',  -- hex
    secondary_color  VARCHAR(7),
    custom_domain    VARCHAR(120),
    custom_css       TEXT,
    FOREIGN KEY (establishment_id) REFERENCES establishments(id)
);

-- OAuth tokens do Google (separado por segurança)
CREATE TABLE google_oauth_tokens (
    establishment_id CHAR(26) PRIMARY KEY,
    access_token     TEXT NOT NULL,
    refresh_token    TEXT NOT NULL,
    expiry           DATETIME(6) NOT NULL,
    scope            TEXT NOT NULL,
    updated_at       DATETIME(6) NOT NULL,
    FOREIGN KEY (establishment_id) REFERENCES establishments(id)
);

-- Usuários gestores
CREATE TABLE users (
    id               CHAR(26) PRIMARY KEY,
    establishment_id CHAR(26) NOT NULL,
    name             VARCHAR(120) NOT NULL,
    email            VARCHAR(120) NOT NULL,
    password_hash    VARCHAR(255) NOT NULL,
    role             ENUM('owner','manager') NOT NULL DEFAULT 'owner',
    active           BOOLEAN NOT NULL DEFAULT TRUE,
    created_at       DATETIME(6) NOT NULL,
    UNIQUE INDEX idx_email (email),
    INDEX idx_establishment (establishment_id),
    FOREIGN KEY (establishment_id) REFERENCES establishments(id)
);

-- Horários de funcionamento do estabelecimento
CREATE TABLE business_hours (
    id               CHAR(26) PRIMARY KEY,
    establishment_id CHAR(26) NOT NULL,
    day_of_week      TINYINT NOT NULL,              -- 0=Dom, 1=Seg ... 6=Sab
    open_time        TIME NOT NULL,
    close_time       TIME NOT NULL,
    is_closed        BOOLEAN NOT NULL DEFAULT FALSE,
    UNIQUE INDEX idx_establishment_day (establishment_id, day_of_week),
    FOREIGN KEY (establishment_id) REFERENCES establishments(id)
);

-- Profissionais
CREATE TABLE professionals (
    id               CHAR(26) PRIMARY KEY,
    establishment_id CHAR(26) NOT NULL,
    name             VARCHAR(120) NOT NULL,
    avatar_url       VARCHAR(500),
    phone            VARCHAR(20),
    google_calendar_id VARCHAR(250),               -- ID da agenda no Google
    display_order    SMALLINT NOT NULL DEFAULT 0,
    active           BOOLEAN NOT NULL DEFAULT TRUE,
    created_at       DATETIME(6) NOT NULL,
    updated_at       DATETIME(6) NOT NULL,
    INDEX idx_establishment (establishment_id),
    FOREIGN KEY (establishment_id) REFERENCES establishments(id)
);

-- Jornada individual do profissional
CREATE TABLE professional_hours (
    id              CHAR(26) PRIMARY KEY,
    professional_id CHAR(26) NOT NULL,
    day_of_week     TINYINT NOT NULL,
    start_time      TIME NOT NULL,
    end_time        TIME NOT NULL,
    is_unavailable  BOOLEAN NOT NULL DEFAULT FALSE,
    UNIQUE INDEX idx_professional_day (professional_id, day_of_week),
    FOREIGN KEY (professional_id) REFERENCES professionals(id)
);

-- Serviços do catálogo
CREATE TABLE services (
    id               CHAR(26) PRIMARY KEY,
    establishment_id CHAR(26) NOT NULL,
    name             VARCHAR(120) NOT NULL,
    description      VARCHAR(500),
    duration_minutes SMALLINT NOT NULL,             -- ex: 30, 45, 60
    price_cents      INT,                           -- NULL = não exibir preço
    active           BOOLEAN NOT NULL DEFAULT TRUE,
    display_order    SMALLINT NOT NULL DEFAULT 0,
    created_at       DATETIME(6) NOT NULL,
    INDEX idx_establishment (establishment_id),
    FOREIGN KEY (establishment_id) REFERENCES establishments(id)
);

-- Relação profissional ↔ serviço
CREATE TABLE professional_services (
    professional_id CHAR(26) NOT NULL,
    service_id      CHAR(26) NOT NULL,
    PRIMARY KEY (professional_id, service_id),
    FOREIGN KEY (professional_id) REFERENCES professionals(id),
    FOREIGN KEY (service_id) REFERENCES services(id)
);

-- Bloqueios manuais de agenda
CREATE TABLE blocked_periods (
    id              CHAR(26) PRIMARY KEY,
    professional_id CHAR(26) NOT NULL,
    starts_at       DATETIME(6) NOT NULL,
    ends_at         DATETIME(6) NOT NULL,
    reason          VARCHAR(250),
    google_event_id VARCHAR(250),
    INDEX idx_professional_range (professional_id, starts_at, ends_at),
    FOREIGN KEY (professional_id) REFERENCES professionals(id)
);

-- Agendamentos
CREATE TABLE appointments (
    id               CHAR(26) PRIMARY KEY,
    establishment_id CHAR(26) NOT NULL,
    professional_id  CHAR(26) NOT NULL,
    service_id       CHAR(26) NOT NULL,
    client_name      VARCHAR(120) NOT NULL,
    client_phone     VARCHAR(20) NOT NULL,   -- único campo de contato no MVP
    starts_at        DATETIME(6) NOT NULL,
    ends_at          DATETIME(6) NOT NULL,
    status           ENUM('pending','confirmed','cancelled','completed','no_show')
                     NOT NULL DEFAULT 'confirmed',
    source           ENUM('booking_link','manager','api') NOT NULL DEFAULT 'booking_link',
    google_event_id  VARCHAR(250),
    notes            VARCHAR(500),
    idempotency_key  CHAR(36) UNIQUE,              -- previne duplo agendamento
    created_at       DATETIME(6) NOT NULL,
    updated_at       DATETIME(6) NOT NULL,
    INDEX idx_establishment_time (establishment_id, starts_at, ends_at),
    INDEX idx_professional_time (professional_id, starts_at, ends_at),
    INDEX idx_status (status),
    FOREIGN KEY (establishment_id) REFERENCES establishments(id),
    FOREIGN KEY (professional_id) REFERENCES professionals(id),
    FOREIGN KEY (service_id) REFERENCES services(id)
);
```

### 5.3 Notas de modelagem

- **ULID** como ID primário: ordenável por tempo, sem colisões em ambiente distribuído, UUID-compatível. Alternativa: UUID v7.
- **price_cents**: preço em centavos evita erros de ponto flutuante.
- **google_calendar_id** no profissional: cada profissional mapeia para uma agenda específica dentro da conta Google do gestor.
- **idempotency_key** no agendamento: o frontend envia uma chave única por tentativa de booking. Se o usuário clicar duas vezes, o segundo INSERT falha silenciosamente no UNIQUE constraint.

---

## 6. Design de API

### 6.1 Superfícies de API

O sistema expõe **duas superfícies distintas**:

| Superfície | Prefixo | Auth | Consumidor |
|---|---|---|---|
| API do Gestor | `/api/v1/` | JWT Bearer | Painel admin |
| API Pública de Agendamento | `/pub/:slug/` | Nenhuma | Cliente final |

### 6.2 API do Gestor — endpoints principais

```
# Auth
POST   /api/v1/auth/login
POST   /api/v1/auth/logout
POST   /api/v1/auth/refresh

# Setup inicial (onboarding)
GET    /api/v1/establishment
PUT    /api/v1/establishment
PUT    /api/v1/establishment/business-hours

# Whitelabel
GET    /api/v1/whitelabel
PUT    /api/v1/whitelabel

# Google Agenda
GET    /api/v1/google/auth-url
GET    /api/v1/google/callback              # OAuth2 redirect
DELETE /api/v1/google/disconnect

# Profissionais
GET    /api/v1/professionals
POST   /api/v1/professionals
GET    /api/v1/professionals/:id
PUT    /api/v1/professionals/:id
DELETE /api/v1/professionals/:id
PUT    /api/v1/professionals/:id/hours
PUT    /api/v1/professionals/:id/services

# Serviços
GET    /api/v1/services
POST   /api/v1/services
PUT    /api/v1/services/:id
DELETE /api/v1/services/:id

# Agendamentos
GET    /api/v1/appointments                 # com filtros: date, professional, status
GET    /api/v1/appointments/:id
PATCH  /api/v1/appointments/:id/status
POST   /api/v1/appointments/:id/block       # bloquear período

# Bloqueios
GET    /api/v1/blocked-periods
POST   /api/v1/blocked-periods
DELETE /api/v1/blocked-periods/:id
```

### 6.3 API Pública — endpoints do fluxo de agendamento

```
# Dados do estabelecimento (whitelabel + info)
GET    /pub/:slug

# Catálogo
GET    /pub/:slug/services
GET    /pub/:slug/professionals
GET    /pub/:slug/professionals?service_id=:id    # filtra por serviço

# Motor de disponibilidade — endpoint crítico
GET    /pub/:slug/availability
       ?service_id=:id
       &professional_id=:id     # opcional
       &date_from=2026-04-07
       &date_to=2026-04-14
# Resposta: mapa de dias → slots disponíveis [{starts_at, ends_at, professional_id}]

# Criação do agendamento
POST   /pub/:slug/appointments
       Body: { service_id, professional_id?, starts_at, client_name, client_phone,
               idempotency_key }

# Consulta/cancelamento pelo cliente (autenticação leve por telefone)
GET    /pub/:slug/appointments/:id?phone=:phone
PATCH  /pub/:slug/appointments/:id/cancel        # respeita min_advance_cancel_hours
PATCH  /pub/:slug/appointments/:id/reschedule    # cancela o atual e cria novo
```

### 6.4 Formato padrão de resposta

```json
// Sucesso
{
  "data": { ... },
  "meta": { "page": 1, "per_page": 20, "total": 45 }
}

// Erro
{
  "error": {
    "code": "SLOT_UNAVAILABLE",
    "message": "O horário selecionado não está mais disponível.",
    "details": {}
  }
}
```

### 6.5 Códigos de erro de negócio relevantes

| Código | Situação |
|---|---|
| `SLOT_UNAVAILABLE` | Horário já foi reservado por outro cliente |
| `PROFESSIONAL_INACTIVE` | Profissional desativado após seleção |
| `SERVICE_DURATION_CONFLICT` | Serviço não cabe no slot selecionado |
| `OUTSIDE_BUSINESS_HOURS` | Horário fora do expediente |
| `GOOGLE_SYNC_FAILED` | Agendamento criado mas sync com Google falhou (retentativa em background) |

---

## 7. Motor de Disponibilidade

Este é o componente mais crítico do produto. Qualquer erro aqui resulta em duplos agendamentos ou slots fantasma — experiência péssima para gestor e cliente.

### 7.1 Algoritmo

```
função CalcularSlots(profissional, serviço, data_inicio, data_fim):

1. Para cada dia no intervalo [data_inicio, data_fim]:
   a. Verificar se o dia está dentro do horário de funcionamento do estabelecimento
   b. Verificar se o profissional trabalha nesse dia (professional_hours)
   c. Obter janela de trabalho do profissional no dia [start_time, end_time]

2. Para cada janela de trabalho:
   a. Buscar agendamentos confirmados do profissional no intervalo
   b. Buscar bloqueios manuais do profissional no intervalo
   c. Buscar eventos do Google Agenda (cache Redis, TTL 30s)

3. Calcular slots livres:
   - Dividir a janela de trabalho em fatias de [duração_do_serviço] minutos
   - Remover fatias que colidem com qualquer ocupação (agendamentos + bloqueios + Google)
   - Retornar apenas fatias completamente livres

4. Para "qualquer profissional":
   - Executar algoritmo acima em goroutines paralelas para cada profissional
   - Agregar resultados, deduplicando por starts_at
```

### 7.2 Prevenção de race condition no booking

O endpoint `POST /pub/:slug/appointments` deve ser protegido contra agendamentos simultâneos no mesmo slot:

```
1. Iniciar transação MySQL
2. SELECT ... FOR UPDATE em appointments onde:
   professional_id = X AND starts_at < :ends_at AND ends_at > :starts_at AND status != 'cancelled'
3. Se retornar linhas → rollback → retornar SLOT_UNAVAILABLE
4. INSERT no appointments
5. Commit
6. Enfileirar job de sincronização Google Agenda (assíncrono)
```

O `idempotency_key` único garante que tentativas duplicadas do mesmo cliente retornem o agendamento existente sem criar um novo.

### 7.3 Cache de disponibilidade

```
Redis Key:   availability:{establishment_id}:{professional_id}:{date}
TTL:         30 segundos
Invalidação: ao criar/cancelar agendamento, ao criar/remover bloqueio, ao receber webhook do Google
```

---

## 8. Integração Google Agenda

### 8.1 Fluxo OAuth2

```
Gestor → GET /api/v1/google/auth-url
       → Redirect para Google Consent Screen
       → Google → GET /api/v1/google/callback?code=xxx
       → Backend: troca code por access_token + refresh_token
       → Salva tokens em google_oauth_tokens (criptografados)
       → Cria calendários para cada profissional via API
```

**Scopes necessários:**
- `https://www.googleapis.com/auth/calendar`

### 8.2 Estrutura de calendários

```
Conta Google do gestor
├── Calendário: "Profissional: João Silva"  (google_calendar_id salvo no profissional)
├── Calendário: "Profissional: Maria Santos"
└── Calendário: "Profissional: Pedro Lima"
```

### 8.3 Sincronização

**Ao criar agendamento:**
- Job assíncrono no Redis queue cria evento no calendário do profissional
- Em caso de falha: 3 retentativas com backoff exponencial
- Agendamento existe no banco independente do sucesso do Google

**Ao cancelar agendamento:**
- Job assíncrono deleta o evento pelo `google_event_id`

**Webhook do Google (Push Notifications):**
- Google notifica o sistema quando há mudança externa no calendário
- Sistema invalida cache de disponibilidade do profissional afetado
- Evita conflitos com eventos criados diretamente no Google Agenda

**Renovação de tokens:**
- Job periódico (a cada 50 minutos) renova access_tokens com o refresh_token
- Refresh token nunca expira (a não ser que o gestor revogue acesso)

### 8.4 Estratégia de resiliência

O sistema **não depende** do Google para confirmar um agendamento. O agendamento é persistido primeiro no MySQL. A sincronização com Google é sempre assíncrona. Se o Google estiver indisponível, o agendamento existe e uma retentativa ocorrerá. O cliente não experimenta falha.

---

## 9. Sistema de Notificações

**Fora do escopo do MVP.** Nenhuma notificação automática será implementada na primeira versão.

A interface abaixo fica definida para implementação futura, mas nenhum worker ou provider será construído agora:

```go
type NotificationProvider interface {
    SendBookingConfirmation(ctx context.Context, appointment Appointment) error
    SendBookingReminder(ctx context.Context, appointment Appointment) error
    SendBookingCancellation(ctx context.Context, appointment Appointment) error
}
```

Implementações previstas para pós-MVP:
- `EmailProvider` — SMTP via Resend ou SendGrid
- `WhatsAppProvider` — Evolution API (self-hosted) ou Twilio WhatsApp Business API

### Notificações planejadas (pós-MVP)

| Evento | Cliente | Gestor |
|---|---|---|
| Agendamento confirmado | e-mail + WhatsApp | e-mail |
| Agendamento cancelado | e-mail | e-mail |
| Lembrete (D-1) | WhatsApp | — |

---

## 10. Whitelabel

### 10.1 Resolução da marca

Toda request à API pública inclui o slug do estabelecimento. O middleware resolve e injeta:
- Nome do estabelecimento
- Cores (primary/secondary)
- Logo URL
- Domínio customizado (se configurado)

### 10.2 Domínio

MVP: subdirectory simples — `app.seudominio.com.br/p/barbearia-silva`

Nenhuma configuração extra de DNS ou wildcard SSL é necessária para o MVP. O prefixo `/p/` (de "público") distingue as rotas de booking das rotas do painel do gestor.

Pós-MVP: `agenda.barbearia-silva.com.br` via CNAME + wildcard SSL (Let's Encrypt wildcard ou Cloudflare).

---

## 11. Autenticação e Segurança

### 11.1 Auth do gestor

- **JWT** com expiração curta (15 minutos) + **Refresh Token** (7 dias, rotacionado)
- Refresh tokens armazenados no Redis com hash SHA-256
- Rate limiting no endpoint de login: 5 tentativas / 15 minutos por IP

### 11.2 API pública (cliente final)

- Sem autenticação, mas com **rate limiting por tenant**: 60 requests/minuto por IP
- Proteção contra enumeração de horários: resposta sempre em tempo constante (sem timing attacks)
- Campos de entrada sanitizados e validados (client_name, phone)

### 11.3 Proteções gerais

- Todas as queries parametrizadas (sem SQL raw com interpolação)
- Tokens OAuth do Google armazenados criptografados em repouso (AES-256-GCM)
- Headers de segurança via middleware: `X-Content-Type-Options`, `X-Frame-Options`, `Referrer-Policy`
- CORS configurado explicitamente (não `*`)

---

## 12. Infraestrutura e Deploy

### 12.1 MVP — Infraestrutura mínima (< R$ 200/mês)

```
┌─────────────────────────────────────┐
│           VPS (ex: Hetzner CX21)    │
│                                     │
│  Docker Compose                     │
│  ├── app (Go binary)                │
│  ├── nginx (reverse proxy + SSL)    │
│  ├── mysql:8                        │
│  └── redis:7                        │
└─────────────────────────────────────┘
```

**Spec mínima para MVP:** 2 vCPU, 4 GB RAM (comporta ~50 tenants ativos tranquilamente)

**Provedores recomendados:** Hetzner, DigitalOcean, Vultr

### 12.2 Variáveis de ambiente

```env
# Banco
DATABASE_URL=mysql://user:pass@localhost:3306/schedule

# Redis
REDIS_URL=redis://localhost:6379

# JWT
JWT_SECRET=<chave-256-bits>
JWT_EXPIRY=15m
REFRESH_TOKEN_EXPIRY=168h

# Google OAuth2
GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=
GOOGLE_REDIRECT_URL=https://app.exemplo.com/api/v1/google/callback

# Encryption (tokens do Google)
ENCRYPTION_KEY=<chave-256-bits>

# Notificações
SMTP_HOST=
SMTP_FROM=
WHATSAPP_API_URL=
WHATSAPP_TOKEN=

# App
APP_ENV=production
APP_PORT=8080
BASE_URL=https://app.exemplo.com
```

### 12.3 CI/CD recomendado

```
GitHub Actions:
  push main → build binary → SSH deploy → zero-downtime restart
```

---

## 13. Estrutura de Projeto Go (referência)

```
internal/scheduling/
├── handler.go          # HTTP handlers (decode request, call service, encode response)
├── service.go          # Regras de negócio — CalcularSlots, CriarAgendamento, etc.
├── repository.go       # Queries MySQL (interface + implementação)
├── availability.go     # Motor de disponibilidade (algoritmo puro, testável)
└── availability_test.go

internal/calendar/
├── handler.go
├── service.go          # Orquestra OAuth2 + criação de calendários
├── repository.go       # Persiste tokens
└── gcal_client.go      # Wrapper da Google Calendar API
```

**Convenção:** handler não chama repositório diretamente. A cadeia é sempre:
`Handler → Service → Repository`

---

## 14. Decisões das Questões em Aberto ✅

Todas as questões em aberto da especificação foram resolvidas:

| Questão | Decisão | Impacto técnico |
|---|---|---|
| Cliente cancela/reagenda sozinho? | **Sim** | Endpoint público `PATCH /pub/:slug/appointments/:id/cancel` + campo `min_advance_cancel_hours` no estabelecimento |
| Confirmação automática? | **Sim** | `status = 'confirmed'` direto no INSERT, sem fluxo de aprovação |
| Cobrança antecipada? | **Não** | Sem integração de pagamento no MVP |
| Dados do cliente? | **Nome + telefone apenas** | Schema atual suficiente — sem campo de e-mail obrigatório |
| Domínio customizado no MVP? | **Não** | Subdirectory: `/p/:slug` — sem wildcard SSL |
| Notificações? | **Nenhuma no MVP** | Sem workers de e-mail ou WhatsApp agora — interface definida para pós-MVP |

---

## 15. Roadmap Técnico

### MVP (validação)
- [ ] Auth do gestor (JWT)
- [ ] CRUD profissionais e serviços
- [ ] Horários de funcionamento
- [ ] Motor de disponibilidade
- [ ] Fluxo público de agendamento — nome + telefone, sem login
- [ ] Cancelamento/reagendamento público com regra de antecedência mínima
- [ ] Google OAuth2 + criação de calendários por profissional
- [ ] Sync assíncrono de eventos Google Agenda
- [ ] Whitelabel básico (nome, logo, cores primária/secundária)
- [ ] Rota pública por slug (`/p/:slug`)
- [ ] Painel do gestor (Nuxt.js)
- [ ] Página de agendamento pública (Nuxt.js, mobile-first)

### Pós-MVP (escala e receita)
- [ ] Notificações por e-mail e WhatsApp
- [ ] Lembrete automático (D-1)
- [ ] Dashboard de analytics para o gestor
- [ ] Confirmação manual de agendamentos
- [ ] Domínio customizado + wildcard SSL
- [ ] Cobrança antecipada (Stripe / PagSeguro)
- [ ] Múltiplos gestores por estabelecimento
- [ ] App mobile (React Native ou Flutter)

---

## Glossário

| Termo | Definição |
|---|---|
| Tenant | Um estabelecimento cliente da plataforma |
| Slug | Identificador URL-friendly do estabelecimento (ex: `barbearia-silva`) |
| Slot | Janela de tempo disponível para agendamento |
| Motor de disponibilidade | Algoritmo que calcula slots livres respeitando todas as regras |
| ULID | Universally Unique Lexicographically Sortable Identifier |
| Idempotency Key | Chave única por tentativa de booking para evitar duplicatas |
