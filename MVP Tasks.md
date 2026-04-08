# MVP - Lista de Tarefas

> Acompanhamento de progresso do desenvolvimento. Marque `[x]` ao concluir cada item.
> Ultima atualizacao: Abril 2026 — Fase 10 concluida.

---

## Fase 1 - Scaffolding e Infraestrutura

### 1.1 Infraestrutura local
- [x] Criar `docker-compose.yml` com MySQL 8 e Redis 7
- [x] Criar `.env.example` com todas as variaveis documentadas
- [x] Validar que `docker compose up` sobe o ambiente sem erros

### 1.2 Backend - Estrutura inicial (Go)
- [x] Inicializar `go.mod` com nome do modulo
- [x] Adicionar dependencias base: router (Chi ou Echo), MySQL driver (`sqlx` ou `pgx`-like), Redis (`go-redis`), `golang-migrate`, `godotenv`
- [x] Criar estrutura de pastas: `cmd/api/`, `internal/`, `pkg/`, `migrations/`
- [x] Criar `cmd/api/main.go` com servidor HTTP subindo na porta configuravel via env
- [x] Criar `Makefile` com comandos: `make run`, `make migrate`, `make build`

### 1.3 Banco de dados - Migrations
- [x] Configurar `golang-migrate` para rodar via `make migrate`
- [x] Migration 001: tabela `establishments`
- [x] Migration 002: tabela `whitelabel_configs`
- [x] Migration 003: tabela `google_oauth_tokens`
- [x] Migration 004: tabela `users`
- [x] Migration 005: tabela `business_hours`
- [x] Migration 006: tabela `professionals`
- [x] Migration 007: tabela `professional_hours`
- [x] Migration 008: tabela `services`
- [x] Migration 009: tabela `professional_services`
- [x] Migration 010: tabela `blocked_periods`
- [x] Migration 011: tabela `appointments`

### 1.4 Frontend - Estrutura inicial (Nuxt.js)
- [x] Criar projeto Nuxt 3 com TypeScript e Vue 3 Composition API
- [x] Configurar Tailwind CSS
- [x] Criar layout `layouts/manager.vue` (sidebar + header autenticado)
- [x] Criar layout `layouts/booking.vue` (minimalista, mobile-first, sem chrome do sistema)
- [x] Configurar modulo de API client (`$fetch` ou `ofetch`) com base URL do backend via env
- [x] Criar store de autenticacao (`composables/useAuth.ts`)
- [x] Configurar rota guard para area do gestor (`middleware/auth.ts`)
- [x] Estrutura de rotas:
- [x] `/login` - pagina de login do gestor
- [x] `/dashboard` - pagina inicial do painel
- [x] `/p/:slug` - raiz da pagina publica de agendamento

---

## Fase 2 - Autenticacao do Gestor

### 2.1 Backend
- [x] Implementar `POST /api/v1/auth/login` - valida credenciais, retorna JWT + refresh token
- [x] Implementar `POST /api/v1/auth/logout` - invalida refresh token no Redis
- [x] Implementar `POST /api/v1/auth/refresh` - troca refresh token por novo par
- [x] Middleware JWT que valida token e injeta `user_id` + `establishment_id` no contexto
- [x] Seed de usuario gestor inicial para desenvolvimento

### 2.2 Frontend
- [x] Pagina `/login` com formulario (e-mail + senha)
- [x] Logica de login: POST, salvar tokens, redirecionar para `/dashboard`
- [x] Logica de refresh automatico de token antes do vencimento
- [x] Logica de logout + limpeza de estado

---

## Fase 3 - Gestao do Estabelecimento e Whitelabel

### 3.1 Backend
- [x] Implementar `GET /api/v1/establishment` - retorna dados do estabelecimento do gestor logado
- [x] Implementar `PUT /api/v1/establishment` - atualiza nome, contato, timezone
- [x] Implementar `GET /api/v1/whitelabel` - retorna configuracao visual
- [x] Implementar `PUT /api/v1/whitelabel` - atualiza logo e cores
- [x] Middleware de tenant para API publica: resolve `establishment_id` a partir do slug na URL
- [x] Validacao de slug unico ao criar/atualizar
- [ ] Upload de logo (armazenamento local ou S3-compatible)

### 3.2 Frontend - Painel
- [x] Pagina `/dashboard/settings` - formulario de dados do estabelecimento
- [ ] Secao de whitelabel: preview de cores em tempo real, upload de logo
- [ ] Exibir link publico gerado (`/p/:slug`) com botao de copiar

---

## Fase 4 - Horarios de Funcionamento

### 4.1 Backend
- [x] Implementar `GET /api/v1/establishment/business-hours`
- [x] Implementar `PUT /api/v1/establishment/business-hours` - salva horario dos 7 dias (suporte a "fechado")

### 4.2 Frontend - Painel
- [x] Pagina `/dashboard/hours` - grade semanal com campos de abertura/fechamento por dia
- [x] Toggle "fechado" por dia da semana

---

## Fase 5 - Profissionais

### 5.1 Backend
- [x] Implementar `GET /api/v1/professionals`
- [x] Implementar `POST /api/v1/professionals`
- [x] Implementar `GET /api/v1/professionals/:id`
- [x] Implementar `PUT /api/v1/professionals/:id`
- [x] Implementar `DELETE /api/v1/professionals/:id` - soft delete (`active = false`)
- [x] Implementar `PUT /api/v1/professionals/:id/hours` - jornada individual
- [x] Implementar `PUT /api/v1/professionals/:id/services` - associar servicos ao profissional

### 5.2 Frontend - Painel
- [x] Pagina `/dashboard/professionals` - listagem com status ativo/inativo
- [x] Modal/pagina de criacao e edicao de profissional
- [ ] Secao de jornada de trabalho por dia da semana
- [ ] Secao de selecao de servicos atendidos pelo profissional

---

## Fase 6 - Servicos

### 6.1 Backend
- [x] Implementar `GET /api/v1/services`
- [x] Implementar `POST /api/v1/services`
- [x] Implementar `PUT /api/v1/services/:id`
- [x] Implementar `DELETE /api/v1/services/:id` - soft delete

### 6.2 Frontend - Painel
- [x] Pagina `/dashboard/services` - listagem com duracao e preco opcional
- [x] Modal/pagina de criacao e edicao de servico
- [x] Campo de duracao com opcoes pre-definidas (15, 30, 45, 60, 90, 120 min) + custom

---

## Fase 7 - Integracao Google Agenda

### 7.1 Backend
- [x] Implementar `GET /api/v1/google/auth-url` - gera URL de consentimento OAuth2
- [x] Implementar `GET /api/v1/google/callback` - troca code por tokens, salva criptografado
- [x] Implementar `DELETE /api/v1/google/disconnect` - revoga acesso e remove tokens
- [x] Criar `pkg/gcal/` - wrapper da Google Calendar API (criar evento, deletar evento, listar eventos)
- [x] Ao conectar Google: criar automaticamente uma agenda por profissional ativo
- [ ] Ao criar profissional apos conexao: criar agenda no Google automaticamente
- [ ] Job assincrono (Redis queue): criar evento ao confirmar agendamento
- [ ] Job assincrono: deletar/atualizar evento ao cancelar agendamento
- [ ] Renovacao automatica de access token (job periodico a cada 50 minutos)
- [ ] Implementar endpoint de webhook do Google para invalidar cache ao detectar mudanca externa
- [ ] Armazenar `google_event_id` no agendamento apos sync bem-sucedido

### 7.2 Frontend - Painel
- [x] Pagina/secao `/dashboard/settings/google` - botao "Conectar Google Agenda"
- [x] Fluxo de redirect OAuth2 e retorno com feedback de sucesso/erro
- [x] Exibir status da conexao + botao de desconectar
- [x] Exibir lista de calendarios criados por profissional

---

## Fase 8 - Motor de Disponibilidade

### 8.1 Backend
- [x] Implementar `internal/scheduling/availability.go` - algoritmo puro e testavel
- [x] Logica de divisao da janela de trabalho em slots
- [x] Subtracao de agendamentos confirmados
- [x] Subtracao de bloqueios manuais
- [ ] Subtracao de eventos do Google Agenda (quando conectado)
- [x] Suporte a "qualquer profissional disponivel" com execucao paralela via goroutines
- [x] Cache de disponibilidade no Redis (TTL 30s por profissional/dia)
- [ ] Invalidacao de cache ao criar/cancelar agendamento ou bloqueio
- [x] Implementar `GET /pub/:slug/availability` com filtros de servico, profissional e intervalo de datas
- [x] Testes unitarios do algoritmo de disponibilidade (cenarios: slot livre, slot ocupado, fora do horario, duracao nao cabe)

---

## Fase 9 - Agendamento Publico

### 9.1 Backend
- [x] Implementar `GET /pub/:slug` - dados publicos do estabelecimento (whitelabel + info)
- [x] Implementar `GET /pub/:slug/services` - lista servicos ativos
- [x] Implementar `GET /pub/:slug/professionals` - lista profissionais ativos (com filtro por servico)
- [x] Implementar `POST /pub/:slug/appointments` - criacao de agendamento
- [x] Validar disponibilidade no momento da criacao (`SELECT FOR UPDATE`)
- [x] Idempotency key para prevenir duplicatas
- [x] `status = 'confirmed'` direto no INSERT
- [ ] Enfileirar job de sync Google Agenda apos commit
- [x] Implementar `GET /pub/:slug/appointments/:id?phone=:phone` - consulta com auth leve
- [x] Implementar `PATCH /pub/:slug/appointments/:id/cancel` - cancelamento pelo cliente
- [x] Validar regra de antecedencia minima (`min_advance_cancel_hours`)
- [x] Implementar `PATCH /pub/:slug/appointments/:id/reschedule` - reagendamento pelo cliente
- [ ] Rate limiting por IP no endpoint de criacao (60 req/min)

### 9.2 Frontend - Pagina Publica
- [x] Pagina `/p/:slug` - exibir nome, logo e cores do estabelecimento
- [x] Step 1: Selecao de servico (cards com nome, duracao e preco opcional)
- [x] Step 2: Selecao de profissional (ou "qualquer disponivel")
- [x] Step 3: Selecao de data (calendario) e horario (slots disponiveis via API)
- [x] Step 4: Formulario do cliente (nome + telefone)
- [x] Step 5: Tela de confirmacao com resumo do agendamento
- [x] Tela de erro para slot indisponivel (com opcao de escolher outro horario)
- [x] Pagina de consulta de agendamento por ID + telefone
- [x] Fluxo de cancelamento com validacao de antecedencia
- [x] Fluxo de reagendamento (mostra novos slots disponiveis)
- [ ] Responsividade mobile-first em todos os steps

---

## Fase 10 - Gestao de Agendamentos (Painel)

### 10.1 Backend
- [x] Implementar `GET /api/v1/appointments` - com filtros: data, profissional, status
- [x] Implementar `GET /api/v1/appointments/:id`
- [x] Implementar `PATCH /api/v1/appointments/:id/status` - gestor altera status manualmente
- [x] Implementar `GET /api/v1/blocked-periods` - lista bloqueios
- [x] Implementar `POST /api/v1/blocked-periods` - gestor bloqueia periodo
- [x] Implementar `DELETE /api/v1/blocked-periods/:id` - remove bloqueio

### 10.2 Frontend - Painel
- [x] Pagina `/dashboard` - visao de hoje: agendamentos do dia por profissional
- [x] Pagina `/dashboard/appointments` - agenda com filtros navegavel
- [x] Filtros por profissional e status
- [x] Detalhe do agendamento: dados do cliente, servico, status, acoes (concluir, cancelar, no-show)
- [x] Modal de criacao de bloqueio de periodo (por profissional, data/hora, motivo)

---

## Fase 11 - Qualidade e Deploy

### 11.1 Qualidade
- [x] Testes unitarios do motor de disponibilidade (cobertura dos cenarios criticos)
- [x] Testes de integracao dos endpoints de criacao e cancelamento de agendamento
- [ ] Validacao de que nao e possivel criar agendamento em slot ocupado (race condition)
- [ ] Revisao de seguranca: queries parametrizadas, tokens criptografados, CORS configurado
- [ ] Revisao de isolamento de tenant: garantir que nenhum endpoint vaza dados cross-tenant

### 11.2 Deploy MVP
- [ ] Provisionar VPS (Hetzner, DigitalOcean ou similar - minimo 2 vCPU / 4 GB RAM)
- [ ] Configurar Nginx como reverse proxy com SSL (Let's Encrypt)
- [ ] Criar `Dockerfile` para o backend Go (multi-stage build)
- [ ] Configurar Docker Compose de producao (app + nginx + mysql + redis)
- [ ] Configurar backup automatico do MySQL
- [ ] Configurar GitHub Actions: build + deploy via SSH no push para `main`
- [ ] Smoke test pos-deploy: criar estabelecimento, cadastrar profissional, realizar agendamento

---

## Resumo de progresso

| Fase | Tarefas | Concluidas |
|---|---|---|
| 1 - Scaffolding | 22 | 21 |
| 2 - Autenticacao | 8 | 8 |
| 3 - Estabelecimento e Whitelabel | 10 | 7 |
| 4 - Horarios de Funcionamento | 4 | 4 |
| 5 - Profissionais | 11 | 9 |
| 6 - Servicos | 6 | 6 |
| 7 - Google Agenda | 15 | 9 |
| 8 - Motor de Disponibilidade | 10 | 8 |
| 9 - Agendamento Publico | 23 | 21 |
| 10 - Gestao de Agendamentos | 11 | 11 |
| 11 - Qualidade e Deploy | 12 | 2 |
| **Total** | **132** | **106** |
