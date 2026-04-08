# Frontend Validation

## Comandos disponiveis

Execute em `frontend`:

```powershell
npm run build
npm run validate
```

## Observacoes

- `validate` hoje reaproveita o build de producao como validacao executavel no ambiente atual.
- O projeto ainda nao possui `vue-tsc` instalado localmente para um `typecheck` dedicado.
- Quando a dependencia estiver disponivel, o proximo passo recomendado e adicionar:
  - `npm run typecheck`
  - validacao automatizada de smoke/regressao para os fluxos criticos

## Checklist manual

Use junto com:

- `features/checklist-regressao-frontend.md`
