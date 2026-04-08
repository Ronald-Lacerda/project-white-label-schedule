# White Label Schedule Design System

## Direção

Este design system foi montado para equilibrar dois contextos:

- `manager`: clareza operacional, leitura rápida, filtros e status visíveis
- `booking`: fluxo linear, baixa fricção e confiança na marca do estabelecimento

A referência de comportamento é:

- Google Agenda para hierarquia temporal e previsibilidade operacional
- Calendly para progressão passo a passo e booking sem atrito
- plataformas como Square, Fresha e Booksy para linguagem de gestão de agenda

## Foundations

Os foundations vivem em `frontend/assets/css/main.css` e `frontend/composables/useBrandTheme.ts`.

### Tokens principais

- Tipografia: `--font-sans`, `--font-display`
- Raio: `--radius-xs` a `--radius-xl`
- Sombra: `--shadow-soft`, `--shadow-card`, `--shadow-float`
- Base neutra: `--color-canvas`, `--color-surface`, `--color-surface-muted`, `--color-border`, `--color-text`
- Semânticos: `--color-success`, `--color-warning`, `--color-danger`, `--color-info`
- Marca: `--color-brand-primary`, `--color-brand-secondary`, `--color-brand-accent`, `--color-brand-surface`, `--color-brand-on-primary`

### Contrato de marca

Persistimos apenas `primary_color` e `secondary_color`.

No client, o restante é derivado:

- `accent`
- `surface`
- `surfaceMuted`
- `onPrimary`
- canais RGB para sombras, fundos e overlays

Isso mantém o whitelabel simples no backend e mais robusto visualmente no frontend.

## Componentes base

### UI

- `AppButton`
- `AppSurface`
- `AppStatusPill`
- `AppModal`

### Booking / Brand

- `PublicHeader`
- `BookingSummary`
- `BrandPreview`

## Classes globais

As classes utilitárias do sistema ficam em `main.css`.

As mais usadas:

- `ds-page`
- `ds-kicker`
- `ds-title`
- `ds-display-title`
- `ds-subtitle`
- `ds-surface`
- `ds-surface-brand`
- `ds-button-*`
- `ds-input`, `ds-select`, `ds-textarea`
- `ds-status-*`
- `ds-stepper`, `ds-step`

## Regras de uso

### Painel do gestor

- base neutra dominante
- cor da marca apenas em destaques e foco
- status sempre com pílulas semânticas
- filtros e ações principais no topo da superfície

### Página pública

- marca visível no topo e nas escolhas ativas
- uma decisão por vez
- CTA principal sempre associado à cor primária
- estados vazios e erro sem interromper a leitura

## Páginas piloto

As referências implementadas neste ciclo são:

- `frontend/pages/dashboard/appointments.vue`
- `frontend/pages/p/[slug]/index.vue`
- `frontend/pages/dashboard/settings.vue`

## Próximos passos

- migrar `services`, `professionals`, `hours` e `dashboard/index` para `AppSurface`, `AppButton` e `AppStatusPill`
- criar componentes adicionais para tabela híbrida, filtros persistentes e cards de entidade
- consolidar um kit visual equivalente no Figma usando os mesmos tokens
