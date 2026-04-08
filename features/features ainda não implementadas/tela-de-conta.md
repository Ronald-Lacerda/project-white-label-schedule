# Feature: Tela de Conta

## Objetivo

Criar uma nova tela autenticada no painel do gestor para exibir as informacoes do usuario logado e permitir a alteracao de senha mediante confirmacao da senha atual.

## Contexto atual

O projeto ja possui:

- autenticacao via JWT e refresh token
- estado do usuario no frontend por meio de `useAuth`
- layout autenticado do painel em `frontend/layouts/manager.vue`
- middleware de autenticacao em `frontend/middleware/auth.ts`

Atualmente nao existe:

- rota dedicada para visualizar os dados da conta do usuario logado
- rota dedicada para alterar senha
- pagina especifica de conta no dashboard

## Escopo da feature

### Deve incluir

- nova tela autenticada em `/dashboard/account`
- exibicao de nome, e-mail, perfil e estabelecimento vinculado do usuario logado
- formulario para alteracao de senha
- obrigatoriedade de informar a senha atual antes de definir a nova senha
- confirmacao da nova senha no frontend
- feedback visual de sucesso, erro e estado de carregamento

### Nao faz parte do escopo inicial

- edicao de nome ou e-mail do usuario
- recuperacao de senha por e-mail
- invalidacao global de todas as sessoes apos troca de senha
- historico de alteracoes de conta

## Experiencia esperada

### Navegacao

- adicionar item `Conta` no menu lateral do layout do gestor
- manter `Configuracoes` para dados do estabelecimento e whitelabel

### Layout da pagina

Sugestao de estrutura:

1. Card `Informacoes da conta`
2. Card `Alterar senha`

### Informacoes exibidas

- nome
- e-mail
- perfil
- identificador do estabelecimento ou nome do estabelecimento, se desejado

### Formulario de alteracao de senha

Campos:

- senha atual
- nova senha
- confirmar nova senha

Comportamentos:

- desabilitar botao durante envio
- validar nova senha com minimo de 8 caracteres
- validar confirmacao da senha no cliente
- limpar campos de senha apos sucesso

## Endpoints

### GET `/api/v1/auth/me`

Retorna os dados seguros do usuario autenticado.

Resposta esperada:

```json
{
  "data": {
    "id": "user_id",
    "name": "Nome do usuario",
    "email": "usuario@empresa.com",
    "role": "owner",
    "establishment_id": "establishment_id"
  }
}
```

### PUT `/api/v1/auth/password`

Altera a senha do usuario autenticado.

Payload esperado:

```json
{
  "current_password": "senha-atual",
  "new_password": "nova-senha"
}
```

Resposta esperada:

- `204 No Content` ou `200 OK` com mensagem simples de sucesso

## Regras de negocio

- apenas usuario autenticado pode acessar os endpoints
- a senha atual deve ser validada antes da troca
- a nova senha deve ter no minimo 8 caracteres
- o sistema nao deve expor `password_hash` em nenhuma resposta
- opcionalmente, a nova senha pode ser impedida de ser igual a senha atual

## Estrutura tecnica por arquivo

### Backend

#### `backend/internal/auth/model.go`

- criar DTO seguro para resposta do usuario autenticado
- criar estrutura para entrada de alteracao de senha, se fizer sentido

#### `backend/internal/auth/repository.go`

- adicionar metodo `UpdatePassword(ctx, userID, passwordHash) error`
- reutilizar `FindUserByID` para buscar o usuario autenticado

#### `backend/internal/auth/service.go`

- implementar `GetMe(ctx, userID)`
- implementar `ChangePassword(ctx, userID, currentPassword, newPassword)`
- validar senha atual com bcrypt
- validar a nova senha
- gerar novo hash e persistir

#### `backend/internal/auth/handler.go`

- criar handler `Me`
- criar handler `ChangePassword`
- ler `user_id` do contexto via `shared.UserIDFromContext`

#### `backend/internal/shared/errors.go`

- adicionar erro especifico para senha atual invalida
- opcionalmente adicionar erro especifico para senha fraca

#### `backend/cmd/api/main.go`

- registrar:
  - `GET /api/v1/auth/me`
  - `PUT /api/v1/auth/password`

### Frontend

#### `frontend/composables/useAuth.ts`

- adicionar `fetchMe()`
- adicionar `changePassword(currentPassword, newPassword)`
- atualizar `state.user` quando necessario

#### `frontend/pages/dashboard/account.vue`

- criar a nova pagina autenticada
- renderizar os dados da conta
- implementar formulario de troca de senha
- tratar estados de loading, erro e sucesso

#### `frontend/layouts/manager.vue`

- adicionar link `Conta` apontando para `/dashboard/account`

## Checklist tecnica

- criar DTO seguro para resposta da conta
- adicionar metodo de update de senha no repositorio
- implementar `GetMe` no service
- implementar `ChangePassword` no service
- criar handlers `Me` e `ChangePassword`
- registrar novas rotas autenticadas
- estender `useAuth.ts` com os novos metodos
- criar `frontend/pages/dashboard/account.vue`
- adicionar item `Conta` na sidebar
- validar fluxo completo em ambiente local

## Criterios de aceite

- usuario autenticado consegue acessar `/dashboard/account`
- a tela exibe corretamente nome, e-mail e perfil
- a troca de senha exige a senha atual
- a troca falha com senha atual incorreta
- a troca falha com nova senha invalida
- a troca funciona com dados validos
- apos sucesso, os campos de senha sao limpos
- nenhuma resposta da API expoe `password_hash`

## Riscos e pontos de atencao

- o tipo `User` atual no backend contem `PasswordHash`, portanto e importante usar uma resposta sanitizada
- a feature mexe com autenticacao e exige cuidado extra com validacao e mensagens de erro
- se a troca de senha nao invalidar sessoes anteriores, isso deve ser aceito conscientemente como decisao de MVP
