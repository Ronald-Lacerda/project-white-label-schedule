# Checklist de Regressao do Frontend

## Dashboard

- Login abre com layout novo e envia credenciais com sucesso.
- Criar conta abre com layout novo, gera slug automaticamente e redireciona para configuracoes.
- Layout do manager renderiza menu lateral no desktop e navegacao horizontal no mobile.
- Pagina inicial do dashboard carrega sem erro e mostra data correta do dia.
- Botao `Gerenciar profissionais` navega corretamente.
- Botao `Gerenciar servicos` navega corretamente.

## Profissionais

- `Novo profissional` abre modal, nao formulario inline.
- `Editar` abre modal preenchido com os dados atuais.
- `Cancelar` fecha o modal sem salvar.
- Salvar cria ou atualiza o profissional e fecha o modal.
- `Remover` abre modal de confirmacao padrao.
- Confirmar exclusao remove o profissional.

## Servicos

- `Novo servico` abre modal, nao formulario inline.
- `Editar` abre modal preenchido com os dados atuais.
- `Cancelar` fecha o modal sem salvar.
- Salvar cria ou atualiza o servico e fecha o modal.
- `Remover` abre modal de confirmacao padrao.
- Confirmar exclusao remove o servico.

## Agenda

- Filtros de data, profissional e status aplicam sem quebrar a tela.
- `Limpar filtros` restaura o filtro padrao do dia atual.
- `Bloquear periodo` abre modal com inicio e fim preenchidos localmente.
- Salvar bloqueio cria bloqueio sem deslocamento de horario aparente.
- `Remover` bloqueio abre modal de confirmacao padrao.
- Modal de detalhe do agendamento abre e fecha corretamente.
- Acoes de concluir, no-show e cancelar continuam funcionando.

## Horarios

- Pagina de horarios renderiza toggle de aberto/fechado para todos os dias.
- Ao marcar dia como fechado, os campos de hora somem.
- Ao marcar dia como aberto, os campos voltam.
- `Salvar horarios` persiste as alteracoes e mostra feedback de sucesso.

## Configuracoes

- Upload de logo continua funcionando.
- Preview da marca atualiza com a nova logo.
- Salvar dados do estabelecimento funciona.
- Salvar aparencia funciona.
- Link para Google Agenda navega corretamente.

## Google Agenda

- Tela de integracao carrega sem erro.
- `Conectar Google Agenda` continua iniciando o fluxo.
- `Desconectar conta` abre modal de confirmacao padrao.
- Confirmar desconexao atualiza o status da tela.

## Fluxo Publico

- Pagina publica carrega com header e marca do estabelecimento.
- Selecao de servico funciona.
- Selecao de profissional funciona.
- Selecao de data funciona.
- Selecao de horario funciona.
- Formulario final confirma agendamento com sucesso.
- Consulta por codigo e telefone continua funcionando.
- Cancelamento continua funcionando.
- Reagendamento continua funcionando.

## Validacoes tecnicas

- Rodar `npm run build` em `frontend` conclui sem erro.
- Conferir manualmente que nao restou `confirm()` nativo nas telas.
- Conferir manualmente que criar profissional e criar servico nao renderizam formulario inline.
