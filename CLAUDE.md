# CLAUDE.md — Guia do Projeto

Plataforma SaaS whitelabel de agendamentos para negócios de serviços.
Nicho inicial: salões de beleza e barbearias. Arquitetura pensada para expansão a qualquer nicho baseado em agenda.

---

## Stack

| Camada | Tecnologia |
|---|---|
| Backend | Go (Chi ou Echo) |
| Frontend | Nuxt.js 3 (Vue 3, TypeScript, Composition API) |
| Banco de dados | MySQL 8.x |
| Cache / Filas | Redis 7 |
| Infraestrutura local | Docker Compose |
| Migrations | golang-migrate |
| CSS | Tailwind CSS |

---

## Estrutura de repositório

```
/
├── backend/
│   ├── cmd/
│   │   └── api/
│   │       └── main.go
│   ├── internal/
│   │   ├── tenancy/        # Gestão de estabelecimentos (tenant root)
│   │   ├── auth/           # JWT, refresh tokens, Google OAuth2
│   │   ├── catalog/        # Serviços e profissionais
│   │   ├── scheduling/     # Motor de disponibilidade + agendamentos
│   │   ├── calendar/       # Integração Google Agenda
│   │   ├── whitelabel/     # Configuração visual e slug
│   │   └── shared/         # Middleware, erros, DTOs, paginação
│   ├── pkg/
│   │   ├── gcal/           # Wrapper Google Calendar API
│   │   └── crypto/         # Criptografia de tokens OAuth
│   ├── migrations/
│   ├── go.mod
│   ├── Makefile
│   └── .env.example
│
├── frontend/
│   ├── layouts/
│   │   ├── manager.vue     # Painel do gestor (sidebar + header)
│   │   └── booking.vue     # Página pública (mobile-first, sem chrome)
│   ├── pages/
│   │   ├── login.vue
│   │   ├── dashboard/
│   │   └── p/
│   │       └── [slug]/     # Fluxo público de agendamento
│   ├── composables/
│   ├── components/
│   ├── middleware/
│   │   └── auth.ts
│   └── nuxt.config.ts
│
├── Especificação do negócio.md
├── Arquitetura do Sistema.md
├── MVP Tasks.md
└── CLAUDE.md
```

---

## Comandos de desenvolvimento

```bash
# Subir ambiente local
docker compose up -d

# Backend
cd backend
make run          # go run ./cmd/api
make migrate      # executa migrations pendentes
make build        # compila binário

# Frontend
cd frontend
npm install
npm run dev
```

---

## Convenções de código — Backend (Go)

### Estrutura interna de cada módulo

Cada módulo em `internal/` segue sempre esta estrutura de camadas:

```
handler.go      → decodifica request, chama service, codifica response
service.go      → regras de negócio, orquestra repositórios
repository.go   → queries SQL (interface + implementação concreta)
```

**Handler nunca chama repositório diretamente.** A cadeia é sempre:
`Handler → Service → Repository`

### Injeção de dependência

Usar injeção manual via construtores — sem framework de DI. Exemplo:

```go
type Service struct {
    repo   Repository
    redis  *redis.Client
}

func NewService(repo Repository, redis *redis.Client) *Service {
    return &Service{repo: repo, redis: redis}
}
```

### Contexto e tenant

O `establishment_id` do tenant autenticado **sempre** vem do contexto da request — nunca do body ou query param em rotas autenticadas. Todo repositório recebe `ctx context.Context` e extrai o tenant de lá.

```go
// Extrair tenant do contexto
establishmentID := ctx.Value(tenantKey).(string)
```

### Tratamento de erros

Erros de domínio são tipos customizados definidos em `internal/shared/errors.go`:

```go
var ErrSlotUnavailable = &DomainError{Code: "SLOT_UNAVAILABLE", Status: 409}
var ErrOutsideBusinessHours = &DomainError{Code: "OUTSIDE_BUSINESS_HOURS", Status: 422}
```

Handlers convertem erros de domínio em JSON com o formato padrão. Nunca expor erros internos ao cliente.

### IDs

Usar **ULID** para todos os IDs primários. Importar `github.com/oklog/ulid/v2`.

### Datas e horários

- Sempre armazenar em UTC no banco
- Converter para o timezone do estabelecimento apenas na camada de apresentação
- Usar `time.Time` em Go; nunca strings para representar datas em código interno

### Nomenclatura

- Arquivos: `snake_case.go`
- Structs/interfaces: `PascalCase`
- Funções/métodos exportados: `PascalCase`
- Variáveis e funções internas: `camelCase`
- Constantes: `PascalCase` ou `SCREAMING_SNAKE` para valores de enum

---

## Convenções de código — Frontend (Nuxt/Vue)

### Organização de componentes

```
components/
├── ui/           # Componentes genéricos reutilizáveis (Button, Input, Modal, etc.)
├── manager/      # Componentes do painel do gestor
└── booking/      # Componentes do fluxo de agendamento público
```

### Composables

Toda lógica de estado e chamadas de API fica em composables em `composables/`. Páginas e componentes não fazem `$fetch` diretamente.

```ts
// Correto
const { appointments, fetchAppointments } = useAppointments()

// Errado — não chamar $fetch diretamente em componente
const data = await $fetch('/api/v1/appointments')
```

### Nomenclatura

- Componentes: `PascalCase.vue`
- Composables: `camelCase.ts` com prefixo `use` (ex: `useAppointments.ts`)
- Páginas: `kebab-case.vue` ou `[param].vue`

### Layouts

- `manager.vue` — obrigatório para todas as rotas autenticadas do gestor
- `booking.vue` — obrigatório para todas as rotas públicas `/p/:slug/...`

---

## Multi-tenancy — Regras críticas

1. **Toda query no banco filtra por `establishment_id`** — sem exceção.
2. O `establishment_id` é resolvido por middleware a partir do token JWT (rotas do gestor) ou do slug na URL (rotas públicas).
3. Nunca aceitar `establishment_id` vindo do corpo da request em operações autenticadas.
4. Em testes, sempre criar dados com `establishment_id` explícito para evitar vazamento entre cenários.

---

## Motor de disponibilidade — Regras críticas

O arquivo `internal/scheduling/availability.go` é o coração do produto. Qualquer alteração aqui exige testes.

**Regras que o algoritmo deve respeitar (em ordem de prioridade):**
1. Horário de funcionamento do estabelecimento
2. Jornada individual do profissional no dia
3. Duração do serviço deve caber inteiramente no slot
4. Agendamentos já confirmados (status != 'cancelled')
5. Bloqueios manuais cadastrados pelo gestor
6. Eventos externos do Google Agenda (quando integrado)

**Race condition no booking:** O endpoint `POST /pub/:slug/appointments` usa `SELECT ... FOR UPDATE` dentro de uma transação para garantir exclusividade do slot. Não remover essa proteção.

---

## API — Formato padrão de resposta

```json
// Sucesso com dados
{ "data": { ... } }

// Sucesso com lista
{ "data": [...], "meta": { "page": 1, "per_page": 20, "total": 45 } }

// Erro
{ "error": { "code": "SLOT_UNAVAILABLE", "message": "Mensagem legível", "details": {} } }
```

HTTP status codes:
- `200` — sucesso
- `201` — recurso criado
- `400` — erro de validação do input
- `401` — não autenticado
- `403` — sem permissão
- `409` — conflito (ex: slot indisponível, slug duplicado)
- `422` — erro de regra de negócio
- `500` — erro interno (nunca expor detalhes)

---

## Segurança — Checklist obrigatório

- [ ] Todas as queries SQL usam parâmetros (`?` placeholders) — nunca interpolação de string
- [ ] Tokens OAuth do Google armazenados criptografados em repouso (AES-256-GCM via `pkg/crypto`)
- [ ] CORS configurado explicitamente — nunca `*` em produção
- [ ] Rate limiting ativo nos endpoints públicos de agendamento
- [ ] JWT secret e encryption key nunca no código — sempre via variável de ambiente
- [ ] Soft delete para profissionais e serviços — nunca `DELETE` físico de dados com agendamentos vinculados

---

## Banco de dados — Referência rápida

| Tabela | Propósito |
|---|---|
| `establishments` | Tenant raiz — cada linha é um negócio cliente |
| `whitelabel_configs` | Logo, cores, domínio — 1:1 com establishment |
| `google_oauth_tokens` | Tokens OAuth separados por segurança |
| `users` | Gestores autenticados |
| `business_hours` | Horários de funcionamento por dia da semana |
| `professionals` | Profissionais/funcionários do estabelecimento |
| `professional_hours` | Jornada individual por dia da semana |
| `services` | Catálogo de serviços com duração e preço |
| `professional_services` | N:N entre profissional e serviço |
| `blocked_periods` | Bloqueios manuais de agenda |
| `appointments` | Agendamentos — tabela central do produto |

**Índices críticos em `appointments`:**
- `(establishment_id, starts_at, ends_at)` — queries de disponibilidade por estabelecimento
- `(professional_id, starts_at, ends_at)` — queries de disponibilidade por profissional
- `idempotency_key UNIQUE` — previne agendamentos duplicados

---

## Variáveis de ambiente obrigatórias

```env
# Banco
DATABASE_URL=mysql://user:pass@localhost:3306/schedule

# Redis
REDIS_URL=redis://localhost:6379

# JWT
JWT_SECRET=
JWT_EXPIRY=15m
REFRESH_TOKEN_EXPIRY=168h

# Google OAuth2
GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=
GOOGLE_REDIRECT_URL=

# Criptografia (tokens Google)
ENCRYPTION_KEY=

# App
APP_ENV=development
APP_PORT=8080
BASE_URL=http://localhost:8080
```

---

## Escopo do MVP — O que NÃO implementar agora

Se uma funcionalidade abaixo aparecer como sugestão durante o desenvolvimento, **recusar e registrar para pós-MVP**:

- Notificações (e-mail, WhatsApp, SMS)
- Pagamento ou cobrança antecipada
- Domínio customizado por tenant (wildcard SSL)
- App mobile
- Marketplace público
- Programa de fidelidade
- Gestão financeira ou comissionamento
- Confirmação manual de agendamentos pelo gestor
- Múltiplos gestores por estabelecimento

---

## Documentos de referência

| Arquivo | Conteúdo |
|---|---|
| `Especificação do negócio.md` | Regras de negócio, público-alvo, hipóteses do MVP |
| `Arquitetura do Sistema.md` | Stack, padrões, schema, design de API |
| `MVP Tasks.md` | Lista completa de tarefas com progresso |
