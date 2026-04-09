# Plano de saneamento de acentuacao e textos visiveis

## Objetivo

Corrigir de forma eficiente todos os textos do sistema exibidos para usuarios finais e gestores que hoje estao:

- sem acentuacao
- com grafia inconsistente
- com sinais de encoding quebrado, como `Â`, `Ã` ou separadores corrompidos

O objetivo desta iniciativa nao e apenas "colocar acento", mas reduzir retrabalho futuro criando uma base minima para concentrar textos visiveis e evitar que novas telas repitam o problema.

## Diagnostico atual

O problema aparece em duas frentes principais:

### Frontend

Os textos visiveis estao majoritariamente inline em paginas, layouts, componentes e composables. Isso faz com que a correcao manual fique espalhada e sujeita a regressao.

Arquivos com maior concentracao de ocorrencias identificadas:

- `frontend/pages/page/[slug]/index.vue`
- `frontend/pages/dashboard/appointments.vue`
- `frontend/pages/dashboard/index.vue`
- `frontend/pages/dashboard/professionals.vue`
- `frontend/pages/dashboard/services.vue`
- `frontend/pages/dashboard/settings/index.vue`
- `frontend/layouts/manager.vue`
- `frontend/components/booking/PublicHeader.vue`
- `frontend/components/ui/AppConfirmModal.vue`
- `frontend/composables/useAppointments.ts`
- `frontend/composables/useEstablishment.ts`
- `frontend/composables/useAuth.ts`

### Backend

As mensagens de erro de dominio e algumas mensagens retornadas pela API tambem estao sem acento. Como essas respostas sao mostradas na interface, elas fazem parte do problema de UX.

Pontos principais:

- `backend/internal/shared/errors.go`
- handlers e fluxos de agendamento/autenticacao que propagam mensagens da API para a UI

### Encoding

Foi identificado sinal de encoding inconsistente em trechos como:

- `Â·` em listas da UI
- comentarios e textos com `Ã` ou caracteres corrompidos

Isso indica que parte do problema nao esta apenas na escolha das palavras, mas tambem na codificacao dos arquivos.

## Resultado esperado

Apos o saneamento:

- todos os textos visiveis ao gestor e ao cliente devem estar com acentuacao correta em portugues
- mensagens de erro retornadas pela API devem estar legiveis e padronizadas
- o frontend deve ter um ponto unico ou no minimo um agrupamento central para os textos mais recorrentes
- arquivos de interface e dominio devem permanecer em UTF-8 sem caracteres corrompidos
- novos fluxos devem seguir um padrao que desestimule strings inline sem revisao

## Escopo

### Inclui

- revisar e corrigir textos de navegacao, botoes, titulos, descricoes, placeholders, estados vazios, loadings e confirmacoes
- revisar e corrigir mensagens de erro exibidas ao usuario
- revisar textos do fluxo publico de agendamento
- revisar textos do painel do gestor
- corrigir casos evidentes de mojibake e encoding quebrado
- criar uma estrutura simples para centralizacao de copy do frontend
- definir checklist de validacao manual por fluxo

### Nao inclui

- internacionalizacao completa com multi-idioma
- revisao de documentacao interna que nao aparece ao usuario final ou gestor
- reescrita ampla de UX copy alem do necessario para legibilidade e padronizacao
- alteracoes de naming interno de tipos, ids, campos de banco ou contratos de API que nao impactem texto exibido

## Estrategia recomendada

### Fase 1 - Base estrutural minima

Criar uma camada simples para evitar que a correcao vire manutencao dispersa.

Recomendacao:

- criar um arquivo como `frontend/constants/copy.ts` ou `frontend/i18n/ptBR.ts`
- mover para esse arquivo os textos compartilhados ou recorrentes
- usar o catalogo primeiro em layout, navegacao, modais e labels reaproveitados

Essa fase nao precisa ser uma solucao completa de i18n. O foco e reduzir novas strings soltas e concentrar o vocabulario do sistema.

### Fase 2 - Correcao dos fluxos de maior impacto

Corrigir primeiro as areas mais visiveis e com maior densidade de texto:

1. fluxo publico de agendamento
2. dashboard inicial do gestor
3. layout/navegacao do painel
4. agenda do gestor
5. configuracoes, profissionais, servicos e horarios

### Fase 3 - Backend e mensagens de dominio

Padronizar as mensagens de erro que sobem da API:

- revisar `backend/internal/shared/errors.go`
- garantir acentuacao correta
- manter consistencia verbal entre frontend e backend
- revisar se composables exibem `error.message` da API sempre que existir

### Fase 4 - Saneamento de encoding

Executar uma varredura final em arquivos alterados para confirmar:

- UTF-8 preservado
- ausencia de `Â`, `Ã`, `â€”`, `â€¢` e similares em textos visiveis
- ausencia de separadores corrompidos em comentarios e strings

## Ordem de execucao

### Bloco 1 - Navegacao e estrutura global

Arquivos alvo:

- `frontend/layouts/manager.vue`
- `frontend/components/ui/AppConfirmModal.vue`
- `frontend/components/booking/PublicHeader.vue`

Objetivo:

- corrigir labels principais do sistema
- corrigir textos de modal e navegacao
- padronizar termos recorrentes como:
  - Início
  - Visão operacional
  - Serviços
  - Catálogo
  - Horários
  - Configurações
  - Integrações
  - Confirmação

### Bloco 2 - Fluxo publico

Arquivos alvo:

- `frontend/pages/page/[slug]/index.vue`
- `frontend/pages/p/[slug]/index.vue` se ainda estiver em uso
- `frontend/composables` ligados ao fluxo publico, se retornarem mensagens ao usuario

Objetivo:

- corrigir todo o wizard de agendamento
- revisar estados de loading, vazio, erro e confirmacao
- revisar placeholders e rotulos de dados pessoais
- corrigir textos de resumo e confirmacao final

### Bloco 3 - Dashboard inicial

Arquivos alvo:

- `frontend/pages/dashboard/index.vue`

Objetivo:

- corrigir cards, checklist operacional e CTA do link publico
- eliminar casos de separador corrompido como `Â·`

### Bloco 4 - Agenda do gestor

Arquivos alvo:

- `frontend/pages/dashboard/appointments.vue`
- `frontend/composables/useAppointments.ts`

Objetivo:

- corrigir filtros, titulos, tabela, cards mobile, modal de detalhe e mensagens de erro
- revisar acoes de bloqueio, cancelamento e paginacao

### Bloco 5 - Cadastros e configuracoes

Arquivos alvo:

- `frontend/pages/dashboard/professionals.vue`
- `frontend/pages/dashboard/services.vue`
- `frontend/pages/dashboard/hours.vue`
- `frontend/pages/dashboard/settings/index.vue`
- `frontend/pages/dashboard/settings/google.vue`
- `frontend/composables/useEstablishment.ts`
- `frontend/composables/useServices.ts`
- `frontend/composables/useProfessionals.ts`
- `frontend/composables/useBusinessHours.ts`
- `frontend/composables/useAuth.ts`

Objetivo:

- corrigir formularios, status, sucesso/erro, campos e orientacoes operacionais

### Bloco 6 - Backend

Arquivos alvo:

- `backend/internal/shared/errors.go`
- `backend/internal/shared/response.go`
- handlers que retornem mensagens de validacao de forma direta

Objetivo:

- padronizar mensagens retornadas pela API
- garantir que mensagens exibidas no frontend estejam com acentuacao correta

## Proposta de estrutura tecnica

### Frontend

Criar um ponto central de copy:

- `frontend/constants/copy.ts`

Estrutura sugerida:

```ts
export const copy = {
  nav: {
    dashboard: 'Início',
    services: 'Serviços',
    settings: 'Configurações',
  },
  common: {
    loading: 'Carregando...',
    save: 'Salvar',
    cancel: 'Cancelar',
    confirm: 'Confirmar',
  },
  booking: {
    notFound: 'Estabelecimento não encontrado.',
  },
}
```

Nao e obrigatorio migrar tudo de uma vez. A recomendacao e usar essa estrutura primeiro nos pontos globais e nos textos repetidos.

### Backend

Manter mensagens de dominio centralizadas em:

- `backend/internal/shared/errors.go`

Padrao recomendado:

- frases curtas
- portugues natural
- acentuacao correta
- tom consistente com o frontend

Exemplos esperados:

- "Horário fora do expediente do estabelecimento."
- "Dados inválidos na requisição."
- "Não autenticado."
- "Sem permissão para esta operação."

## Checklist tecnica

- [ ] Criar arquivo central simples para copy do frontend
- [ ] Corrigir labels de navegacao do layout do gestor
- [ ] Corrigir textos compartilhados de modal e stepper publico
- [ ] Corrigir integralmente o fluxo publico de agendamento
- [ ] Corrigir dashboard inicial do gestor
- [ ] Corrigir tela de agendamentos e seu composable
- [ ] Corrigir telas de profissionais, servicos, horarios e configuracoes
- [ ] Corrigir mensagens visiveis nos composables de autenticacao e estabelecimento
- [ ] Revisar e corrigir `backend/internal/shared/errors.go`
- [ ] Revisar mensagens de erro propagadas pela API para a UI
- [ ] Eliminar casos de `Â`, `Ã` e caracteres corrompidos em arquivos alterados
- [ ] Validar os principais fluxos manualmente em ambiente local

## Validacao manual

### Fluxo publico

- abrir o link publico de agendamento
- validar titulos, descricoes e labels de cada passo
- validar estados de loading, vazio, erro e sucesso
- concluir um agendamento e revisar a tela final

### Gestor

- fazer login
- revisar menu lateral e navegacao mobile
- revisar dashboard inicial
- revisar agenda, bloqueios e modal de detalhe
- revisar profissionais, servicos, horarios e configuracoes

### Erros

- provocar erros de validacao conhecidos
- confirmar que a mensagem exibida esta acentuada e legivel
- confirmar que nao ha regressao visual por aumento de comprimento de texto

## Criterios de aceite

- nenhum texto visivel ao gestor ou ao cliente permanece sem acento quando deveria ter
- nenhum texto visivel nos fluxos revisados apresenta mojibake ou encoding corrompido
- labels e mensagens repetidas passam a ter fonte central ou padrao claro de reutilizacao
- mensagens de dominio do backend estao coerentes com o restante da interface
- o fluxo publico e o painel principal passam por validacao manual sem regressao textual visivel

## Riscos e cuidados

- corrigir apenas visualmente no frontend e esquecer mensagens vindas da API
- editar muitos arquivos sem padronizar encoding e reintroduzir caracteres quebrados
- ajustar texto demais e misturar correcao ortografica com redesign de copy
- deixar o catalogo de copy pela metade, criando dois padroes ao mesmo tempo

## Recomendacao de rollout

Para reduzir risco e facilitar revisao:

1. estruturar o arquivo central de copy
2. atacar navegacao + fluxo publico + dashboard inicial
3. revisar agenda e configuracoes
4. fechar com mensagens do backend e saneamento final de encoding

Essa ordem entrega valor rapido, diminui a percepcao imediata de produto "sem acabamento" e evita retrabalho nas telas mais acessadas.
