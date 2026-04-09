# Refatoracao da tela de agendamentos

## Objetivo

Simplificar a tela `/dashboard/appointments` para focar no que o gestor realmente precisa no MVP: visualizar todos os agendamentos com filtros persistentes e cancelar um agendamento quando necessario.

## Contexto atual

A tela atual possui mais elementos do que o fluxo atual justifica:

- Cards de resumo (Na fila, Confirmados, Concluidos, Atencao) calculados dinamicamente sobre a lista filtrada
- Acoes na tabela e no modal de detalhe para tres transicoes de status: Concluir, No-show e Cancelar
- No fluxo atual, todos os agendamentos chegam com status `confirmed` — nao ha status `pending` nem confirmacao pelo gestor
- As acoes "Concluir" e "No-show" nao sao necessarias no MVP e adicionam ruido visual sem valor operacional imediato

A tela atual tambem mistura responsabilidades: serve tanto de visualizacao de agenda quanto de central de operacoes de status, o que aumenta a complexidade da interface sem um beneficio claro para o MVP.

## Resultado esperado

Apos a refatoracao, a tela deve:

- Exibir o filtro persistente (data, profissional, status) com os bloqueios ativos
- Listar todos os agendamentos na tabela com paginacao
- Permitir que o gestor cancele um agendamento confirmado diretamente da tabela ou do modal de detalhe
- Remover os cards de resumo da parte superior
- Remover as acoes "Concluir" e "No-show" da tabela e do modal

## Escopo

### Inclui

- remover os quatro cards de resumo do topo da pagina
- remover os botoes "Concluir" e "No-show" da coluna de acoes da tabela (desktop e mobile)
- remover os botoes "Concluir" e "No-show" do footer do modal de detalhe
- manter o botao "Cancelar" visivel para agendamentos com status `confirmed`
- manter o filtro persistente intacto (data, profissional, status)
- manter a secao de bloqueios ativos dentro do filtro
- manter a tabela de agendamentos com paginacao
- manter o modal de detalhe com as informacoes do cliente e do agendamento
- manter o modal de bloqueio de periodo e o fluxo de remocao de bloqueio
- manter o filtro de status com todas as opcoes existentes (inclusive `completed` e `no_show`) para que o gestor consiga consultar historico

### Nao inclui

- confirmacao manual de agendamentos pelo gestor (pos-MVP)
- status `pending` na criacao do agendamento (pos-MVP)
- remocao definitiva de logica de status do backend — apenas esconder da UI
- alteracao de qualquer endpoint de backend
- mudanca no composable `useAppointments`

## Impacto tecnico

### Frontend

#### Pagina de agendamentos

Arquivo: `frontend/pages/dashboard/appointments.vue`

Remover:
- bloco `<div class="grid gap-4 md:grid-cols-2 xl:grid-cols-4">` com os quatro `AppSurface` de resumo (linhas 22-32)
- `computed summaryCards` no script
- botoes "Concluir" e "No-show" na coluna de acoes da tabela desktop (dentro do `v-if="appointment.status === 'confirmed'"`)
- botoes "Concluir" e "No-show" nos cards mobile
- botoes "Concluir" e "No-show" no `template #footer` do modal de detalhe

Manter:
- `AppSurface tone="brand"` com filtros e bloqueios ativos
- `AppSurface tone="default" padding="none"` com a tabela, cards mobile e paginacao
- modal de detalhe com dados do cliente, servico, profissional e data
- modal de bloqueio de periodo
- `AppConfirmModal` de remocao de bloqueio

Simplificar no template da tabela desktop (coluna Acoes):

```html
<!-- antes -->
<div v-if="appointment.status === 'confirmed'" class="flex justify-end gap-2">
  <AppButton size="sm" variant="secondary" @click="changeStatus(appointment, 'completed')">Concluir</AppButton>
  <AppButton size="sm" variant="ghost" @click="changeStatus(appointment, 'no_show')">No-show</AppButton>
  <AppButton size="sm" variant="danger" @click="changeStatus(appointment, 'cancelled')">Cancelar</AppButton>
</div>

<!-- depois -->
<div v-if="appointment.status === 'confirmed'" class="flex justify-end gap-2">
  <AppButton size="sm" variant="danger" @click="changeStatus(appointment, 'cancelled')">Cancelar</AppButton>
</div>
```

Simplificar no footer do modal de detalhe:

```html
<!-- antes -->
<div v-if="detailAppt?.status === 'confirmed'" class="flex flex-wrap justify-end gap-2">
  <AppButton variant="secondary" @click="changeStatus(detailAppt, 'completed')">Marcar concluido</AppButton>
  <AppButton variant="ghost" @click="changeStatus(detailAppt, 'no_show')">No-show</AppButton>
  <AppButton variant="danger" @click="changeStatus(detailAppt, 'cancelled')">Cancelar</AppButton>
</div>

<!-- depois -->
<div v-if="detailAppt?.status === 'confirmed'" class="flex flex-wrap justify-end gap-2">
  <AppButton variant="danger" @click="changeStatus(detailAppt, 'cancelled')">Cancelar agendamento</AppButton>
</div>
```

#### Script

Remover do `<script setup>`:
- `computed summaryCards` inteiro (nao e mais referenciado no template)

Manter intacto:
- todos os refs e reactives de filtro, detalhe e bloqueio
- funcoes `changeStatus`, `applyFilters`, `clearFilters`, `goPage`
- funcoes de bloqueio: `openBlockModal`, `submitBlock`, `promptRemoveBlockedPeriod`, `removeBlockedPeriod`
- helpers de formatacao de data e status

### Backend

Nenhuma alteracao necessaria. O endpoint `PATCH /api/v1/appointments/{id}/status` ja aceita `cancelled` e permanece funcional.

## Arquivos mais provaveis

- `frontend/pages/dashboard/appointments.vue` — unico arquivo alterado

## Quebra de tarefas

### Bloco 1. Remover os cards de resumo

- remover o bloco `<div class="grid gap-4 md:grid-cols-2 xl:grid-cols-4">` com os quatro `AppSurface` de resumo no template
- remover o `computed summaryCards` no script
- validar que a pagina carrega sem erros de referencia

### Bloco 2. Simplificar acoes na tabela desktop

- na coluna "Acoes", remover os botoes "Concluir" e "No-show"
- manter apenas o botao "Cancelar" para status `confirmed`
- manter o label "Finalizado" para status terminais

### Bloco 3. Simplificar acoes nos cards mobile

- remover os botoes "Concluir" e "No-show" dos cards mobile (`xl:hidden`)
- manter o botao "Cancelar" e o botao "Detalhes"

### Bloco 4. Simplificar acoes no modal de detalhe

- no `template #footer`, remover os botoes "Marcar concluido" e "No-show"
- manter apenas o botao "Cancelar agendamento" para status `confirmed`
- validar que o modal continua fechando apos o cancelamento e que o status atualiza na tabela

### Bloco 5. Validacao final

- verificar que a pagina carrega sem erros no console
- verificar filtros (data, profissional, status) funcionam e atualizam a tabela
- verificar que clicar em "Cancelar" muda o status na tabela e no modal
- verificar paginacao
- verificar bloqueio de periodo: criar e remover
- verificar responsividade (mobile e desktop)

## Criterios de aceite

- a pagina nao exibe mais os cards de resumo
- a coluna de acoes da tabela exibe apenas "Cancelar" para agendamentos confirmados
- o modal de detalhe exibe apenas o botao "Cancelar agendamento" para agendamentos confirmados
- o filtro persistente continua funcionando normalmente
- os bloqueios ativos continuam sendo exibidos dentro da secao de filtro
- o cancelamento atualiza o status tanto na tabela quanto no modal sem recarregar a pagina
- nenhuma alteracao de backend e necessaria

## Riscos e atencoes

- o `computed summaryCards` referencia `appointments.value` e `blockedPeriods.value` — garantir que sua remocao nao deixe variaveis orfas em outros pontos do template antes de deletar
- o filtro de status ainda mantem as opcoes `completed` e `no_show` para consulta de historico — nao remover essas opcoes, mesmo que no fluxo atual esses status nao sejam mais gerados pelo gestor
- a funcao `changeStatus` continua no script pois e usada para cancelamento — nao remover

## Consideracoes pos-MVP

Quando o fluxo de confirmacao manual for implementado:

- adicionar status `pending` na criacao do agendamento publico (ajuste no backend e no fluxo publico)
- adicionar botao "Confirmar" na coluna de acoes para agendamentos com status `pending`
- adicionar opcao `pending` no filtro de status
- possivelmente reintroduzir cards de resumo com metricas mais relevantes (ex: pendentes de confirmacao)

Esses pontos nao fazem parte desta refatoracao e devem ser tratados em feature separada.

## Recomendacao de implementacao

Esta refatoracao e inteiramente no frontend e afeta apenas um arquivo. A estrategia ideal e:

1. remover os cards de resumo e o computed associado
2. simplificar as acoes em sequencia (tabela desktop, cards mobile, modal)
3. validar localmente que nenhum erro de referencia foi introduzido
4. confirmar o fluxo de cancelamento de ponta a ponta
