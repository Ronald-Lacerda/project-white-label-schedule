# Project White Label Schedule

SaaS de agendamento white-label com backend em Go e frontend em Nuxt 3.

O projeto cobre:
- area do gestor com autenticacao
- configuracao do estabelecimento e whitelabel
- cadastro de profissionais e servicos
- pagina publica de agendamento
- motor de disponibilidade
- integracao inicial com Google Agenda

## Stack

- Backend: Go, Chi, SQLX, MySQL, Redis
- Frontend: Nuxt 3, Vue 3, TypeScript, Tailwind CSS
- Infra local: Docker Compose
- CI: GitHub Actions

## Estrutura

```text
.
|-- backend/
|-- frontend/
|-- features/
|-- .github/workflows/tests.yml
|-- docker-compose.yml
|-- MVP Tasks.md
```

## Pre-requisitos

- Docker e Docker Compose
- Go 1.25+
- Node.js 20+
- npm

## Configuracao local

1. Crie os arquivos de ambiente a partir dos exemplos:

```bash
cp .env.example .env
cp backend/.env.example backend/.env
cp frontend/.env.example frontend/.env
```

2. Preencha pelo menos os campos obrigatorios no `backend/.env`:

- `JWT_SECRET`
- `ENCRYPTION_KEY`

3. Suba MySQL e Redis:

```bash
docker compose up -d
```

## Rodando o backend

```bash
cd backend
make migrate
make seed
make run
```

Backend padrao: `http://localhost:8080`

## Rodando o frontend

```bash
cd frontend
npm ci
npm run dev
```

Frontend padrao: `http://localhost:3000`

## Comandos uteis

Backend:

```bash
cd backend
make build
make test
```

Frontend:

```bash
cd frontend
npm run build
```

## Testes

O backend ja possui:
- testes unitarios do motor de disponibilidade
- testes de integracao para criacao e cancelamento de agendamento

Para rodar:

```bash
cd backend
go test ./...
```

Observacao: os testes de integracao usam MySQL e Redis locais em `127.0.0.1:3306` e `127.0.0.1:6379`.

## GitHub Actions

O workflow em `.github/workflows/tests.yml` executa:

- `go test ./...` no backend com MySQL e Redis
- `npm ci` e `npm run build` no frontend

## Documentacao do projeto

- [MVP Tasks.md](./MVP%20Tasks.md)
- [Arquitetura do Sistema.md](./Arquitetura%20do%20Sistema.md)
- [Especificação do negócio.md](./Especifica%C3%A7%C3%A3o%20do%20neg%C3%B3cio.md)

## Status atual

O projeto ja cobre o fluxo principal de MVP e segue evoluindo pelas tarefas listadas no checklist de MVP.
