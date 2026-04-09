# Ajustes de Usabilidade no Painel e Fluxo Publico

Data: 2026-04-08

## Resumo

Este pacote agrupa 7 ajustes pequenos, mas transversais, no frontend do projeto:

1. remover o botao `Ajustar horarios` da tela inicial
2. substituir a rota publica `/p/:slug` por uma URL mais amigavel
3. remover o header repetido do layout e das paginas do painel
4. mover `Google Agenda conectado` para o topo do checklist inicial
5. aplicar mascara/formatacao nos campos de telefone
6. aplicar formatacao monetaria nos campos de preco
7. padronizar exibicao de horarios no formato brasileiro abreviado, como `15h` e `16h30`

## Contexto atual

Os pontos afetados estao concentrados principalmente no frontend Nuxt:

- o header global do painel esta em `frontend/layouts/manager.vue`
- a tela inicial do painel esta em `frontend/pages/dashboard/index.vue`
- a rota publica hoje esta implementada em `frontend/pages/p/[slug]/index.vue`
- o slug publico aparece tambem em `frontend/pages/dashboard/settings/index.vue`, `frontend/pages/criar-conta.vue` e `frontend/components/brand/BrandPreview.vue`
- os inputs de telefone hoje estao sem formatacao em `frontend/pages/criar-conta.vue`, `frontend/pages/dashboard/settings/index.vue`, `frontend/pages/dashboard/professionals.vue` e no formulario publico em `frontend/pages/p/[slug]/index.vue`
- o input de preco esta sem formatacao em `frontend/pages/dashboard/services.vue`
- a exibicao de horario usa `toLocaleTimeString('pt-BR', { hour: '2-digit', minute: '2-digit' })` em mais de uma tela, o que hoje gera valores como `15:00`

## Escopo

Inclui:

- ajustes de rotas e links do frontend
- limpeza visual do layout do gestor
- reordenacao do checklist da home
- formatacao de campos de formulario no frontend
- criacao ou centralizacao de helpers de horario brasileiro

Nao inclui:

- mudancas de regras de negocio no backend
- mudancas no payload persistido alem de normalizacao de telefone e preco antes do envio
- alteracao de slug ou migracao de dados existentes no banco
- redirecionamento legado `/p/:slug` no backend, a menos que se queira compatibilidade temporaria

## Decisoes recomendadas

### 1. Nova rota publica

Recomendacao principal:

- trocar `/p/:slug` por `/page/:slug`

Motivo:

- exige menos ajuste semantico do que introduzir um nome mais longo por estabelecimento
- evita a percepcao de rota abreviada ou temporaria
- demanda troca simples em pagina, links e preview

Observacao:

- se houver links antigos ja compartilhados, vale decidir se a compatibilidade com `/p/:slug` sera mantida por uma fase de transicao

## Impacto tecnico

### Frontend

- `frontend/layouts/manager.vue`
  - remover o header global com titulo, descricao e selo `Painel do gestor`
- `frontend/pages/dashboard/index.vue`
  - remover o botao `Ajustar horarios`
  - trocar links `/p/${slug}` para `/page/${slug}`
  - mover o item `Google Agenda conectado` para a primeira posicao do checklist
- `frontend/pages/dashboard/settings/index.vue`
  - trocar prefixo visual `/p/` por `/page/`
  - aplicar mascara de telefone no campo de contato
- `frontend/pages/criar-conta.vue`
  - trocar prefixo visual `/p/` por `/page/`
  - aplicar mascara de telefone no onboarding
- `frontend/components/brand/BrandPreview.vue`
  - atualizar preview da URL publica
- `frontend/pages/p/[slug]/index.vue`
  - renomear para a nova rota publica
  - aplicar mascara de telefone no formulario publico
  - adaptar exibicao de horarios para formato abreviado brasileiro
- `frontend/pages/dashboard/professionals.vue`
  - aplicar mascara no telefone do modal
  - exibir telefone listado com formatacao
- `frontend/pages/dashboard/services.vue`
  - trocar o campo numerico cru por campo com formatacao monetaria amigavel
- `frontend/pages/dashboard/hours.vue`
  - mudar o resumo `08:00 as 18:00` para algo como `8h as 18h`
- `frontend/pages/dashboard/appointments.vue`
  - padronizar horarios exibidos como `15h`, `16h30`
- `frontend/pages/dashboard/index.vue`
  - padronizar horarios da agenda do dia

### Backend

Pode ser necessario validar dois pontos:

- se a pagina publica depende apenas do roteamento do frontend, a mudanca pode ser toda no Nuxt
- se existem endpoints, middlewares ou links documentados assumindo `/p/:slug`, sera preciso alinhar documentacao e possivelmente compatibilidade

## Arquivos provaveis de mudanca

- `C:\Users\Ronal\Project White Label Schedule\frontend\layouts\manager.vue`
- `C:\Users\Ronal\Project White Label Schedule\frontend\pages\dashboard\index.vue`
- `C:\Users\Ronal\Project White Label Schedule\frontend\pages\dashboard\settings\index.vue`
- `C:\Users\Ronal\Project White Label Schedule\frontend\pages\dashboard\professionals.vue`
- `C:\Users\Ronal\Project White Label Schedule\frontend\pages\dashboard\services.vue`
- `C:\Users\Ronal\Project White Label Schedule\frontend\pages\dashboard\hours.vue`
- `C:\Users\Ronal\Project White Label Schedule\frontend\pages\dashboard\appointments.vue`
- `C:\Users\Ronal\Project White Label Schedule\frontend\pages\criar-conta.vue`
- `C:\Users\Ronal\Project White Label Schedule\frontend\components\brand\BrandPreview.vue`
- `C:\Users\Ronal\Project White Label Schedule\frontend\pages\p\[slug]\index.vue`

Arquivos auxiliares recomendados:

- `C:\Users\Ronal\Project White Label Schedule\frontend\composables\usePhoneMask.ts` ou helper equivalente
- `C:\Users\Ronal\Project White Label Schedule\frontend\composables\useCurrencyInput.ts` ou helper equivalente
- `C:\Users\Ronal\Project White Label Schedule\frontend\utils\formatters.ts` ou helper equivalente para hora brasileira

## Quebra em tarefas

### Bloco 1. Rota publica amigavel

- criar a nova pagina em `/page/[slug]`
- atualizar todos os links e previews que hoje apontam para `/p/`
- revisar docs e textos da interface que citam a rota antiga
- decidir se `/p/:slug` continua existindo como alias temporario ou se sera removido de vez

### Bloco 2. Limpeza do header do painel

- remover o header global do layout `manager`
- remover os blocos de titulo/descricao duplicados do topo das paginas do dashboard
- garantir que botoes de acao primaria permanecam visiveis no topo de cada tela quando fizer sentido

### Bloco 3. Ajustes da home do gestor

- remover o botao `Ajustar horarios`
- manter apenas a acao de abrir a pagina publica
- mover `Google Agenda conectado` para o topo do checklist
- revisar o resumo de prontidao para manter a mesma ordem visual

### Bloco 4. Formatacao de telefone

- criar helper unico para mascara brasileira
- aplicar a mascara nos inputs de:
  - criar conta
  - configuracoes do estabelecimento
  - modal de profissional
  - formulario publico de agendamento
- normalizar antes de enviar para a API, preferencialmente so com digitos ou no formato esperado pelo backend
- formatar tambem a exibicao de telefones em listas e detalhes quando houver dado salvo

### Bloco 5. Formatacao monetaria

- criar campo de preco com UX brasileira, exibindo `0,00` durante digitacao
- manter conversao segura para `price_cents` ao salvar
- revisar edicao de servico para carregar valor salvo sem perda de centavos

### Bloco 6. Padronizacao de horario

- extrair helper central para horarios como:
  - `15h`
  - `15h30`
  - `8h`
- aplicar o helper em dashboard, agenda, home, horas de funcionamento e fluxo publico
- revisar textos que hoje usam `as`, `ate` ou `:` para manter consistencia

## Riscos e tradeoffs

- a troca de rota publica pode quebrar links antigos se `/p/:slug` deixar de existir imediatamente
- telefone mascarado precisa preservar colagem, edicao no meio do texto e envio consistente para a API
- input monetario mal implementado pode gerar arredondamento incorreto ou perda de foco ruim durante digitacao
- remover headers de cada pagina pode deixar telas sem contexto se os botoes e secoes nao forem reorganizados com cuidado

## Checklist tecnico

- [ ] localizar e trocar todos os usos de `/p/` no frontend
- [ ] decidir estrategia de compatibilidade para a rota antiga
- [ ] remover header global de `frontend/layouts/manager.vue`
- [ ] remover headers locais duplicados das telas do dashboard
- [ ] remover botao `Ajustar horarios` da home
- [ ] mover `Google Agenda conectado` para o topo do checklist
- [ ] criar helper de telefone e aplicar nas telas afetadas
- [ ] criar helper de moeda e aplicar na criacao/edicao de servicos
- [ ] criar helper de horario brasileiro abreviado
- [ ] aplicar novo formato de horario em todas as telas listadas
- [ ] revisar placeholders, labels e previews de URL
- [ ] validar navegacao manual no painel e no fluxo publico

## Criterios de aceite

- a home do gestor nao exibe mais o botao `Ajustar horarios`
- a pagina publica passa a ser acessada por `/page/:slug` em todos os links visiveis do frontend
- o header com titulo, descricao e `Painel do gestor` deixa de aparecer em todas as telas do painel
- o checklist da home mostra `Google Agenda conectado` como primeiro item
- os campos de telefone aceitam digitacao formatada brasileira
- o campo de preco de servico exibe valor em formato brasileiro e continua salvando corretamente em centavos
- os horarios do sistema deixam de aparecer como `15:00` e passam a aparecer como `15h` ou `15h30`

## Ordem sugerida de execucao

1. rota publica e atualizacao dos links
2. remocao do header global e dos headers locais
3. ajustes da home
4. mascara de telefone
5. formatacao monetaria
6. helper central de horario e aplicacao nas telas

## Validacao manual

- abrir dashboard e navegar por `Inicio`, `Agendamentos`, `Profissionais`, `Servicos`, `Horarios` e `Configuracoes`
- confirmar ausencia do header repetido em todas as telas
- abrir a pagina publica a partir da home e das configuracoes
- criar ou editar profissional com telefone mascarado
- editar telefone do estabelecimento
- preencher telefone no onboarding e no fluxo publico
- criar e editar servico com preco `10,50`, `120`, `199,90`
- verificar horarios na home, em agendamentos, em horarios e na pagina publica
