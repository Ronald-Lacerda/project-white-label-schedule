# Planejamento de Ajustes do Frontend

## Objetivo

Mapear os problemas atuais do frontend e transformar os ajustes em uma fila de trabalho clara, com foco em:

- consistencia visual do painel e do fluxo publico
- correção de interacoes fora do padrao esperado
- reducao de regressao em formularios e modais
- consolidacao do design system ja existente

## Resumo do diagnostico

O frontend ja possui uma base visual boa em `frontend/assets/css/main.css` e componentes como `AppButton`, `AppSurface`, `AppStatusPill` e `AppModal`. O problema principal hoje nao e falta de base, e sim uso inconsistente dessa base entre as telas.

Pelo codigo atual, `Novo profissional` e `Novo servico` ja estao renderizando dentro de `AppModal` em:

- `frontend/pages/dashboard/professionals.vue`
- `frontend/pages/dashboard/services.vue`

Entao a melhoria aqui deve ser tratada como:

- validar se ainda existe regressao visual/em execucao fazendo o formulario aparecer inline
- garantir que criar e editar continuem obrigatoriamente no modal
- remover qualquer caminho alternativo que volte a renderizar formulario direto na pagina

## Problemas mapeados

### 1. Inconsistencia do design system entre telas

Existe um design system declarado, mas ele ainda nao foi aplicado de forma uniforme.

Problemas observados:

- `frontend/pages/login.vue` e `frontend/pages/criar-conta.vue` usam classes visuais isoladas (`bg-gray-*`, `border-gray-*`, botao nativo) em vez de `AppSurface`, `AppButton`, `ds-input` e tokens globais.
- `frontend/layouts/manager.vue` mistura tokens do sistema com cores hardcoded em `slate`, o que quebra a coerencia do whitelabel e da navegacao.
- `frontend/pages/dashboard/hours.vue` ainda usa checkbox bruto e estrutura mais simples que o restante do dashboard.
- `frontend/pages/p/[slug]/index.vue` possui muitos blocos clicaveis montados com `button + style` inline, sem componente reutilizavel de selecao.
- `frontend/pages/dashboard/settings.vue` usa um pseudo-botao para upload de logo em vez de componente padrao de acao/upload.

Impacto:

- experiencia visual irregular
- aumento do custo de manutencao
- maior chance de tela nova fugir do padrao

### 2. Padrao de interacao inconsistente para criar, editar e excluir

As telas de profissionais e servicos ja usam modal para criar/editar, mas outras interacoes ainda dependem de `confirm()` do navegador.

Problemas observados:

- `frontend/pages/dashboard/professionals.vue` usa `confirm()` para remover profissional
- `frontend/pages/dashboard/services.vue` usa `confirm()` para remover servico
- `frontend/pages/dashboard/appointments.vue` usa `confirm()` para remover bloqueio
- `frontend/pages/dashboard/settings/google.vue` usa `confirm()` para desconectar Google Agenda

Impacto:

- experiencia inconsistente com o restante da interface
- visual nativo do navegador quebra a identidade do produto
- baixo controle sobre texto, loading, erro e estados de confirmacao

### 3. Fluxos criticos ainda com UX incompleta

Algumas telas funcionam, mas ainda faltam refinamentos importantes para uso real.

Problemas observados:

- `frontend/pages/dashboard/professionals.vue` e `frontend/pages/dashboard/services.vue` nao deixam explicito no plano de tela se existe algum fallback inline/regressao; isso precisa ser validado manualmente e blindado no codigo.
- `frontend/pages/dashboard/hours.vue` nao oferece uma experiencia de configuracao tao forte quanto o restante do dashboard, especialmente para alternar aberto/fechado e revisar horarios.
- `frontend/pages/dashboard/appointments.vue` mistura interacao automatica por `@change` com botao `Aplicar`, o que gera duplicidade de comportamento.
- `frontend/pages/dashboard/appointments.vue` nao possui modal proprio de confirmacao para acoes destrutivas e nao deixa tao claro o loading por linha/acao.
- `frontend/pages/p/[slug]/index.vue` concentra muitos passos e estados em uma unica pagina grande, com pouca componentizacao e alta chance de regressao visual.

Impacto:

- uso operacional menos fluido
- mais chance de erro do usuario
- dificuldade para evoluir as telas sem quebrar comportamento

### 4. Logica de dados espalhada nas paginas

Ha um padrao desejado no projeto para manter acesso a API em composables, mas ele ainda nao esta aplicado de forma consistente.

Problemas observados:

- `frontend/composables/useAppointments.ts` ainda usa `$fetch` direto para listagem, em vez de centralizar no wrapper de API
- `frontend/pages/dashboard/settings.vue` faz upload de logo com `$fetch` direto dentro da pagina
- `frontend/pages/p/[slug]/index.vue` concentra toda a logica de consulta publica, disponibilidade, cancelamento e reagendamento na propria pagina

Impacto:

- mais acoplamento entre UI e regra de negocio
- manutencao mais lenta
- testes e reaproveitamento mais dificeis

### 5. Risco de bugs de data e fuso horario

Existem trechos usando `toISOString()` para derivar dia ou enviar datas locais.

Pontos de atencao:

- `frontend/pages/dashboard/index.vue`
- `frontend/pages/dashboard/appointments.vue`
- `frontend/pages/p/[slug]/index.vue`

Riscos:

- trocar de dia em fusos negativos/positivos
- listar filtro no dia errado
- criar bloqueio com horario deslocado
- reagendamento/agenda publica exibirem datas inconsistentes

### 6. Falta de componentes reutilizaveis para padroes recorrentes

O projeto ja tem base suficiente para extrair mais componentes e reduzir repeticao.

Faltam componentes como:

- modal de confirmacao padrao
- card de selecao para servico/profissional/horario
- campo de upload estilizado
- toggle/switch padronizado para disponibilidade
- bloco padrao de feedback de erro/sucesso

Impacto:

- repeticao de markup
- divergencia visual entre telas equivalentes
- correcao mais cara quando o padrao muda

### 7. Cobertura de validacao frontend insuficiente

O frontend nao apresenta testes automatizados e o `package.json` nao possui scripts de teste/lint.

Impacto:

- regressao visual e comportamental mais facil
- pouca seguranca ao refatorar modais, filtros e fluxo publico

### 8. Higiene do repositorio frontend

Existem artefatos gerados misturados no frontend, como `.nuxt`, `.output` e `node_modules`, o que torna a revisao mais ruidosa.

Impacto:

- investigacao mais lenta
- maior risco de editar ou revisar artefato gerado por engano
- onboarding tecnico pior

## Tarefas recomendadas

### Prioridade alta

- Validar em execucao o fluxo de `Novo profissional` e `Novo servico` e garantir que criar/editar acontecam somente em modal.
- Remover qualquer resquicio de formulario inline ou fallback visual nas telas de profissionais e servicos.
- Criar um `AppConfirmModal` para substituir todos os `confirm()` nativos do dashboard.
- Padronizar exclusao de profissional, servico, bloqueio e desconexao do Google com modal de confirmacao, loading e mensagem consistente.
- Corrigir o tratamento de datas locais para filtros, agenda publica e criacao de bloqueios, evitando dependencia direta de `toISOString()` para datas de interface.

### Prioridade media

- Migrar `frontend/pages/login.vue` para o design system do projeto.
- Migrar `frontend/pages/criar-conta.vue` para o design system do projeto.
- Revisar `frontend/layouts/manager.vue` para remover hardcodes em `slate` e alinhar a navegacao aos tokens do sistema.
- Evoluir `frontend/pages/dashboard/hours.vue` com controles mais consistentes, melhor hierarquia visual e componente de toggle padronizado.
- Extrair um componente reutilizavel de opcao selecionavel para uso no fluxo publico (`servico`, `profissional`, `data`, `horario`).
- Padronizar feedback de sucesso/erro/loading com componentes ou blocos reutilizaveis.

### Prioridade media/baixa

- Refatorar `frontend/pages/dashboard/settings.vue` para mover upload de logo para composable/servico de API.
- Refatorar `frontend/composables/useAppointments.ts` para usar o mesmo padrao do wrapper `useApi`.
- Quebrar `frontend/pages/p/[slug]/index.vue` em componentes/composables menores.
- Criar um componente de upload estilizado e reutilizavel.
- Revisar o comportamento duplicado do filtro em `appointments.vue` para decidir entre auto-apply ou botao explicito, sem manter os dois ao mesmo tempo.

### Qualidade e seguranca de entrega

- Adicionar pelo menos validacoes automatizadas basicas para os fluxos criticos do frontend.
- Criar um checklist manual de regressao visual para dashboard e fluxo publico.
- Limpar a estrategia de artefatos gerados do frontend para reduzir ruido no repositorio.

## Quebra tecnica sugerida

### Bloco 1. Interacoes e regressao de modais

- revisar `frontend/pages/dashboard/professionals.vue`
- revisar `frontend/pages/dashboard/services.vue`
- criar `frontend/components/ui/AppConfirmModal.vue`
- substituir confirmacoes nativas por modal
- validar criar, editar, cancelar e fechar modal por clique externo/botao

### Bloco 2. Consolidacao visual do dashboard

- revisar `frontend/layouts/manager.vue`
- migrar `frontend/pages/dashboard/hours.vue`
- revisar feedbacks em `appointments.vue`, `professionals.vue`, `services.vue` e `settings.vue`

### Bloco 3. Auth e onboarding

- migrar `frontend/pages/login.vue`
- migrar `frontend/pages/criar-conta.vue`
- alinhar formularios, botoes, erros e links ao design system

### Bloco 4. Fluxo publico

- extrair componentes da pagina `frontend/pages/p/[slug]/index.vue`
- criar componente de opcao selecionavel
- revisar estados de consulta, cancelamento, reagendamento e confirmacao
- revisar responsividade e clareza do fluxo multi-etapas

### Bloco 5. Arquitetura frontend

- mover chamadas de API de pagina para composables dedicados
- padronizar uso de `useApi`
- centralizar helpers de data/hora
- preparar base para testes

## Arquivos mais impactados

- `frontend/pages/dashboard/professionals.vue`
- `frontend/pages/dashboard/services.vue`
- `frontend/pages/dashboard/appointments.vue`
- `frontend/pages/dashboard/hours.vue`
- `frontend/pages/dashboard/settings.vue`
- `frontend/pages/dashboard/settings/google.vue`
- `frontend/pages/login.vue`
- `frontend/pages/criar-conta.vue`
- `frontend/pages/p/[slug]/index.vue`
- `frontend/layouts/manager.vue`
- `frontend/components/ui/AppModal.vue`
- `frontend/components/ui/AppButton.vue`
- `frontend/assets/css/main.css`
- `frontend/composables/useAppointments.ts`
- `frontend/composables/useEstablishment.ts`

## Criterios de aceite para considerar a frente concluida

- todas as acoes de criar/editar profissional e servico acontecem exclusivamente via modal
- nenhuma acao destrutiva depende de `confirm()` nativo do navegador
- login, criar conta e dashboard usam o mesmo idioma visual do design system
- fluxo publico tem componentes reutilizaveis para selecao e menos repeticao de markup
- filtros e bloqueios nao sofrem erro de dia/horario por conversao de fuso
- chamadas HTTP ficam centralizadas em composables/padroes consistentes
- existe checklist manual de regressao e pelo menos uma camada inicial de validacao automatizada

## Ordem sugerida de execucao

1. Blindar modais de profissionais e servicos e substituir `confirm()` nativos.
2. Corrigir helpers de data/fuso e revisar agenda/bloqueios.
3. Migrar login e criar conta para o design system.
4. Padronizar `manager.vue` e `hours.vue`.
5. Componentizar o fluxo publico.
6. Refatorar chamadas de API e adicionar validacoes.

## Status apos implementacao desta rodada

Concluido nesta etapa:

- blindagem de criar/editar profissionais e servicos em modal
- substituicao dos `confirm()` nativos por modal de confirmacao reutilizavel
- correcao do uso de datas locais em dashboard e agenda publica
- migracao visual de login e criar conta para o design system
- padronizacao de `manager.vue` e `hours.vue`
- componentizacao base do fluxo publico com card selecionavel
- centralizacao parcial de chamadas de API em composables
- criacao de checklist manual em `features/checklist-regressao-frontend.md`
- criacao de validacao documentada em `frontend/VALIDATION.md`
- adicao de `.gitignore` para reduzir ruido de artefatos gerados

Pendente para uma rodada futura:

- adicionar `vue-tsc`/typecheck dedicado quando a dependencia estiver disponivel no ambiente
- adicionar testes automatizados de smoke/regressao para fluxos criticos
