# Refatoracao do fluxo publico mobile-first

## Objetivo

Simplificar a jornada publica de agendamento para um fluxo linear, claro e mobile-first:

1. Dados pessoais
2. Servico
3. Profissional
4. Horario
5. Confirmacao

O objetivo principal desta fase e reduzir atrito na experiencia publica sem reabrir o escopo de consulta, cancelamento e reagendamento na home do link publico.

## Escopo implementado

- Remocao visual da secao de consulta, cancelamento e reagendamento da tela publica principal.
- Reorganizacao da pagina `frontend/pages/p/[slug]/index.vue` como wizard progressivo.
- Inclusao dos campos obrigatorios:
  - nome
  - email
  - telefone
  - data de nascimento
- Confirmacao final antes do envio do agendamento.
- Persistencia de `client_email` e `client_birth_date` em `appointments`.
- Adaptacao do backend publico para validar e retornar os novos campos.
- Adaptacao do detalhe do agendamento no painel para exibir email e data de nascimento quando presentes.

## Benchmark e direcao de UX

Referencias de mercado usadas para orientar o fluxo:

- Square Appointments: fluxo explicito de selecao operacional e coleta de dados com poucos desvios.
- SimplyBook.me: campos adicionais tratados como extensao do formulario do cliente.
- Booksy: perguntas ao cliente como parte natural do processo de booking.

Direcao aplicada:

- uma etapa por vez no mobile
- um CTA primario por etapa
- resumo compacto do progresso
- confirmacao final separada do formulario
- remocao de ruido operacional secundario da home publica

## Impacto tecnico

### Frontend

- `frontend/pages/p/[slug]/index.vue`
  - passou a operar como wizard em 5 passos
  - removeu lookup/cancel/reschedule da interface principal
  - reorganizou a logica para priorizar navegacao progressiva

- `frontend/components/booking/PublicHeader.vue`
  - stepper atualizado para 5 etapas

- `frontend/composables/useAppointments.ts`
  - tipos do painel atualizados para suportar email e data de nascimento

- `frontend/pages/dashboard/appointments.vue`
  - detalhe do agendamento agora pode mostrar email e nascimento

### Backend

- `backend/migrations/000012_add_client_fields_to_appointments.*`
  - adiciona `client_email`
  - adiciona `client_birth_date`

- `backend/internal/scheduling/model.go`
  - structs publicos e de gestao atualizados com os novos campos

- `backend/internal/scheduling/handler.go`
  - `POST /pub/{slug}/appointments` agora valida email e data de nascimento

- `backend/internal/scheduling/service.go`
  - regras de criacao e resposta ajustadas para os novos dados

- `backend/internal/scheduling/repository.go`
  - selects, inserts e listagens atualizados para persistir e retornar os novos campos

## Contrato atualizado

### POST `/pub/{slug}/appointments`

Payload esperado:

```json
{
  "service_id": "string",
  "professional_id": "string",
  "starts_at": "RFC3339",
  "client_name": "string",
  "client_email": "string",
  "client_phone": "string",
  "client_birth_date": "YYYY-MM-DD",
  "idempotency_key": "string"
}
```

Validacoes:

- `client_name` obrigatorio
- `client_email` obrigatorio e em formato valido
- `client_phone` obrigatorio
- `client_birth_date` obrigatoria no formato `YYYY-MM-DD`

## Checklist tecnica

- [x] Criar migration para novos campos em `appointments`
- [x] Atualizar models do modulo de agendamento
- [x] Atualizar handler publico para validar email e nascimento
- [x] Atualizar service para aceitar, propagar e retornar os novos dados
- [x] Atualizar queries de insert e select no repository
- [x] Refatorar a tela publica para wizard mobile-first
- [x] Remover consulta/cancelamento/reagendamento da home publica
- [x] Criar tela final de confirmacao antes do envio
- [x] Atualizar o painel para exibir os novos dados no detalhe do agendamento
- [x] Documentar a feature nesta pasta

## Criterios de aceite

- O fluxo publico principal nao exibe consulta, cancelamento nem reagendamento.
- O usuario nao avanca do passo 1 sem preencher os 4 campos obrigatorios.
- O profissional e obrigatorio antes da escolha do horario.
- O backend recusa criacao de agendamento sem email ou data de nascimento.
- O agendamento confirmado retorna email e data de nascimento na resposta publica.
- Registros antigos continuam legiveis mesmo com os novos campos nulos.
- O detalhe do agendamento no painel exibe email e nascimento quando disponiveis.

## Fase futura: campos configuraveis por gestor

Esta fase nao implementa configuracao dinamica de campos. O backlog recomendado para a proxima etapa e:

- criar configuracao de formulario por estabelecimento
- permitir campos obrigatorios/opcionais por nicho
- suportar ordem de exibicao e rotulos customizaveis
- versionar a configuracao para nao quebrar historico de agendamentos
- avaliar se passa a fazer sentido introduzir uma entidade propria de cliente
