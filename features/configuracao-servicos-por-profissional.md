# Configuracao de servicos por profissional

## Objetivo

Permitir que o gestor configure, no cadastro de cada profissional, quais servicos ele atende.

Essa configuracao e necessaria para que o fluxo publico liste corretamente os profissionais disponiveis depois que o usuario escolhe um servico.

## Contexto atual

O backend ja possui base parcial pronta para essa feature:

- existe `PUT /api/v1/professionals/{id}/services`
- existe `UpdateServices` no service de profissionais
- existe persistencia na tabela `professional_services`

O gargalo atual esta principalmente no frontend e no contrato de leitura do detalhe do profissional:

- a tela `frontend/pages/dashboard/professionals.vue` hoje so permite editar nome e telefone
- nao existe UI para selecionar servicos atendidos
- o gestor nao enxerga claramente se o profissional esta ou nao vinculado a algum servico

Como consequencia, o fluxo publico pode ficar sem profissionais ao escolher um servico, mesmo com equipe cadastrada.

## Resultado esperado

Ao criar ou editar um profissional, o gestor deve conseguir:

- visualizar os servicos ativos do estabelecimento
- marcar quais servicos aquele profissional atende
- salvar os vinculos
- reabrir o cadastro e ver os servicos previamente selecionados

No fluxo publico, ao selecionar um servico, devem aparecer apenas os profissionais ativos vinculados a esse servico.

## Escopo

### Inclui

- adicionar interface no dashboard para configurar servicos por profissional
- carregar os servicos ativos no modal/pagina de profissional
- exibir os vinculos atuais ao editar
- salvar os `service_ids` vinculados ao profissional
- validar no backend se os servicos pertencem ao mesmo estabelecimento
- melhorar a visibilidade operacional do vinculo no painel

### Nao inclui

- regras avancadas como duracao por profissional
- precificacao por profissional
- agendas diferentes por servico
- ordenacao por servico dentro do fluxo publico
- configuracoes de habilidade, nivel ou especialidade alem do vinculo servico-profissional

## Impacto tecnico

### Frontend

#### Tela de profissionais

Evoluir `frontend/pages/dashboard/professionals.vue` para incluir:

- carregamento da lista de servicos ativos
- estado local de `selectedServiceIds`
- UI de selecao multipla no modal de criar/editar profissional
- pre-carregamento dos servicos ja vinculados no modo edicao
- salvamento do vinculo apos criar/editar o profissional

#### Composables

Revisar:

- `frontend/composables/useProfessionals.ts`
- `frontend/composables/useServices.ts`

Necessidades provaveis:

- metodo para buscar detalhe do profissional com `service_ids`
- metodo para atualizar servicos vinculados
- reaproveitar lista de servicos ativos para montar a selecao

#### UX recomendada

- permitir selecao multipla por cards ou checkboxes estilizados
- mostrar texto auxiliar como:
  - "Selecione os servicos que este profissional atende"
- se nenhum servico for selecionado, mostrar aviso:
  - "Este profissional nao aparecera no agendamento publico ate ter ao menos um servico vinculado"

### Backend

#### Handler

Revisar `backend/internal/catalog/professional_handler.go` para garantir:

- leitura do detalhe do profissional com `service_ids`, se ainda nao existir
- manutencao do endpoint `PUT /api/v1/professionals/{id}/services`

#### Service

Revisar `backend/internal/catalog/professional_service.go` para:

- validar que o profissional pertence ao estabelecimento
- opcionalmente validar consistencia dos `service_ids`

#### Repository

Revisar `backend/internal/catalog/repository.go` para:

- buscar `service_ids` atuais do profissional
- validar que todos os servicos enviados pertencem ao mesmo estabelecimento
- impedir vinculos inconsistentes cross-tenant

## Endpoints

### GET `/api/v1/professionals/{id}`

Objetivo:

- retornar os dados do profissional para edicao

Resposta esperada:

```json
{
  "data": {
    "id": "string",
    "name": "string",
    "phone": "string",
    "active": true,
    "service_ids": ["svc_1", "svc_2"]
  }
}
```

Se o endpoint ainda nao retornar `service_ids`, esta feature deve ampliar o contrato.

### PUT `/api/v1/professionals/{id}/services`

Payload:

```json
{
  "service_ids": ["svc_1", "svc_2"]
}
```

Regras:

- todos os servicos devem pertencer ao mesmo estabelecimento
- idealmente considerar apenas servicos ativos
- lista vazia pode ser aceita, mas deve deixar o profissional invisivel no fluxo publico

## Arquivos mais provaveis

- `frontend/pages/dashboard/professionals.vue`
- `frontend/composables/useProfessionals.ts`
- `frontend/composables/useServices.ts`
- `backend/internal/catalog/professional_handler.go`
- `backend/internal/catalog/professional_service.go`
- `backend/internal/catalog/repository.go`

## Quebra de tarefas

### Bloco 1. Contrato de leitura do profissional

- verificar se `GET /api/v1/professionals/{id}` ja retorna `service_ids`
- se nao retornar, atualizar model, handler, service e repository para incluir esses ids
- validar que a resposta de detalhe funciona para o formulario de edicao

### Bloco 2. Validacao de vinculos no backend

- revisar `SetServices` para garantir consistencia por estabelecimento
- validar que `service_ids` enviados pertencem ao tenant correto
- decidir se servicos inativos podem ou nao ser vinculados
- retornar erro apropriado para ids invalidos

### Bloco 3. UI de selecao no dashboard

- carregar servicos ativos no modal de profissional
- adicionar estado local de selecao multipla
- renderizar lista de servicos com boa usabilidade
- mostrar feedback claro quando nenhum servico estiver selecionado

### Bloco 4. Fluxo de salvar

- no modo criacao:
  - criar profissional
  - salvar `service_ids`
- no modo edicao:
  - atualizar profissional
  - atualizar `service_ids`
- garantir loading, erro e rollback visual coerentes

### Bloco 5. Visibilidade operacional

- exibir na listagem de profissionais um resumo dos servicos vinculados ou contagem
- deixar claro quando o profissional esta sem servicos
- reduzir chance de o gestor achar que o profissional esta disponivel no publico quando nao esta

### Bloco 6. Validacao final

- testar criar profissional com servicos
- testar editar profissional trocando servicos
- testar profissional sem servicos
- testar fluxo publico escolhendo servico com profissionais vinculados
- testar fluxo publico escolhendo servico sem profissionais vinculados

## Criterios de aceite

- o gestor consegue configurar quais servicos um profissional atende
- os servicos vinculados reaparecem corretamente ao editar o profissional
- o backend rejeita vinculos com servicos de outro estabelecimento
- o painel deixa claro quando um profissional nao possui servicos vinculados
- o fluxo publico lista apenas profissionais vinculados ao servico selecionado

## Riscos e atencoes

- se o painel continuar sem indicar servicos vinculados, o gestor pode interpretar a ausencia no fluxo publico como bug
- se `SetServices` nao validar tenant corretamente, pode haver risco de vinculo inconsistente entre estabelecimentos
- se a criacao e a vinculacao forem feitas em duas chamadas separadas sem tratamento de erro adequado, a UX pode ficar inconsistente

## Recomendacao de implementacao

A melhor estrategia para esta fase e aproveitar a base existente e completar apenas o que falta:

- manter a tabela `professional_services`
- manter o endpoint de update de servicos
- ampliar leitura do detalhe do profissional
- concluir a UI no modal de profissionais

Isso fecha a necessidade atual sem introduzir arquitetura paralela nem aumentar demais o escopo.
