# Especificação de Negócio — Plataforma Whitelabel de Agendamentos

## 1. Visão do produto

A proposta é criar uma plataforma **whitelabel** de agendamentos de serviços.  
O produto deve permitir que empresas usem o sistema com sua própria marca, oferecendo um link de agendamento para o cliente final e um painel de gestão para o estabelecimento.

O nicho inicial de validação será **salões de beleza e barbearias**, mas o produto não deve nascer limitado a esse segmento.  
A visão de negócio é que a solução possa ser adaptada para qualquer nicho baseado em agenda, como estética, clínica, consultoria, manutenção, bem-estar e outros.

## 2. Motivação do negócio

A oportunidade nasce da percepção de que muitos estabelecimentos ainda possuem meios ruins de agendar serviços.  
Isso normalmente aparece em formatos como:
- Agenda em papel.
- Atendimento manual via WhatsApp.
- Falta de padronização na confirmação.
- Dificuldade para visualizar horários livres.
- Dependência de alguém responder mensagens para concluir o agendamento.

Em salões e barbearias, esse problema é relevante porque o atendimento frequentemente ocorre ao mesmo tempo em que a equipe precisa responder clientes e administrar horários, o que gera atrito operacional e perda de eficiência [web:4][web:2].  
Também existe impacto em atrasos, cancelamentos e dificuldade de encaixe quando a operação não possui um fluxo organizado de agenda [web:2][web:5].

## 3. Objetivo do MVP

O MVP deve provar que existe demanda real por uma solução simples de agendamento com marca própria para negócios de serviços.  
O foco inicial não é construir um sistema extremamente completo, e sim validar se:
- Os estabelecimentos enxergam valor em profissionalizar o agendamento.
- Os clientes finais preferem agendar por link em vez de depender de conversa manual.
- O gestor consegue organizar múltiplos funcionários em um fluxo centralizado.
- O modelo whitelabel é percebido como diferencial comercial.
- A integração com Google Agenda resolve bem o problema operacional inicial.

## 4. Posicionamento inicial

O produto será vendido como uma solução de agendamento com a marca do próprio negócio cliente.  
A proposta de valor inicial é:

- Organizar a agenda de múltiplos profissionais.
- Reduzir agendamentos manuais.
- Dar mais autonomia ao cliente final.
- Melhorar a percepção de profissionalismo do estabelecimento.
- Permitir implantação rápida sem necessidade de um app próprio.

## 5. Público-alvo

### 5.1 Público inicial de validação
- Salões de beleza.
- Barbearias.

### 5.2 Público futuro
- Clínicas e consultórios.
- Estúdios de estética.
- Spas.
- Prestadores de serviços com agenda.
- Empresas que atendem por horário marcado.

### 5.3 Perfil do comprador
O comprador principal é o **gestor do negócio**.  
É ele quem contrata, configura a conta, cadastra equipe, define regras de funcionamento e acompanha a operação.

### 5.4 Perfil do usuário final
O usuário final é o **cliente do estabelecimento**, que acessa um link para:
- Ver serviços.
- Escolher profissional.
- Escolher dia e horário.
- Confirmar o agendamento.

## 6. Princípios de negócio

O produto deve seguir estes princípios:

1. A experiência de agendamento precisa ser simples.
2. O gestor deve manter controle central da operação.
3. O cliente final não deve precisar criar conta no MVP para agendar.
4. O sistema deve funcionar para equipes com múltiplos funcionários.
5. A solução deve ser genérica o suficiente para atender outros nichos no futuro.
6. O whitelabel deve ser percebido como parte central do produto, não como detalhe opcional.
7. O MVP deve priorizar clareza operacional acima de sofisticação.

## 7. Modelo operacional do MVP

### 7.1 Estrutura de contas
No MVP, cada estabelecimento terá uma conta própria dentro da plataforma.  
Essa conta representa um negócio cliente.

Cada conta terá:
- Um gestor principal.
- Uma marca própria exibida ao cliente final.
- Uma lista de profissionais/funcionários.
- Uma lista de serviços.
- Um link público de agendamento.

### 7.2 Gestão centralizada pelo gestor
No MVP, os funcionários **não** conectam contas próprias de agenda.  
Toda a operação será centralizada na conta do gestor.

Isso significa que:
- O gestor conecta uma única conta Google ao sistema.
- Os calendários dos funcionários são administrados a partir dessa estrutura central.
- O estabelecimento mantém controle sobre a agenda da equipe.
- A operação continua organizada mesmo quando há troca de funcionários.

Esse modelo é adequado para o MVP porque simplifica a gestão e favorece o controle do salão sobre os agendamentos [web:22][web:29].

## 8. Integração com Google Agenda

A integração com Google Agenda faz parte do MVP.  
Ela existe para permitir sincronização prática da agenda operacional do estabelecimento.

### 8.1 Regra principal
A conta Google conectada será a conta do gestor do estabelecimento.

### 8.2 Funcionamento esperado
Para cada funcionário, deverá existir uma agenda vinculada dentro da estrutura administrada pelo gestor.  
Na prática de negócio, cada profissional terá sua própria agenda operacional, ainda que tudo seja controlado centralmente [web:22][web:29].

### 8.3 Objetivos da integração
- Evitar conflitos de horário.
- Refletir os agendamentos em uma agenda real do negócio.
- Facilitar acompanhamento pelo gestor.
- Permitir que a operação do salão não dependa exclusivamente do painel da plataforma.

### 8.4 Regra de propriedade
Os dados operacionais de agenda pertencem ao estabelecimento cliente, e não ao funcionário individual.  
Esse princípio é importante porque o comprador da solução é o negócio, não o colaborador.

## 9. Whitelabel

O produto deve ser percebido pelo cliente final como pertencente ao estabelecimento que contratou a plataforma.  
Por isso, o whitelabel deve estar presente desde o MVP.

### 9.1 Elementos esperados
- Nome do estabelecimento.
- Logo.
- Cores principais.
- Link público identificável.
- Comunicação visual coerente com a marca do cliente.

### 9.2 Objetivo de negócio
O whitelabel não é apenas estética.  
Ele serve para:
- Aumentar confiança do cliente final.
- Reforçar profissionalismo.
- Diferenciar o produto de soluções genéricas.
- Permitir revenda e expansão para diferentes nichos.

## 10. Perfis de acesso

### 10.1 Gestor
O gestor terá acesso administrativo ao sistema.  
Esse perfil deve conseguir:
- Configurar a conta.
- Conectar o Google Agenda.
- Cadastrar serviços.
- Cadastrar profissionais.
- Definir horários de funcionamento.
- Disponibilizar o link público.
- Visualizar e administrar agendamentos.

### 10.2 Cliente final
O cliente final acessará um link de agendamento.  
Esse perfil deve conseguir:
- Ver os serviços disponíveis.
- Escolher profissional, quando aplicável.
- Ver horários disponíveis.
- Realizar agendamento.
- Receber confirmação.
- Cancelar ou reagendar, caso essa funcionalidade seja habilitada no fluxo do negócio.

## 11. Regras de negócio do agendamento

### 11.1 Jornada básica do cliente final
A jornada principal do cliente final deve seguir esta ordem:

1. Acessar o link do estabelecimento.
2. Escolher o serviço.
3. Escolher o profissional, ou optar por atendimento com qualquer profissional disponível.
4. Escolher data e horário.
5. Informar dados básicos.
6. Confirmar o agendamento.

### 11.2 Oferta de horários
O sistema deve mostrar apenas horários realmente disponíveis.  
A disponibilidade exibida deve respeitar:
- Horário de funcionamento do estabelecimento.
- Jornada de trabalho do profissional.
- Duração do serviço.
- Conflitos já existentes na agenda.
- Regras de bloqueio configuradas.

### 11.3 Profissional específico ou qualquer um
O estabelecimento poderá decidir se o cliente:
- Escolhe um profissional específico.
- Pode agendar com qualquer profissional disponível.

Essa regra é importante porque alguns negócios vendem a experiência do especialista, enquanto outros vendem apenas o atendimento dentro do horário disponível.

### 11.4 Duração do serviço
Todo serviço deve possuir duração configurável.  
A duração do serviço impacta diretamente a oferta de horários.

### 11.5 Regras futuras previstas
Mesmo que não estejam completas no MVP, a modelagem do negócio deve considerar a possibilidade futura de:
- Tempo de intervalo entre atendimentos.
- Duração diferente por profissional.
- Recursos compartilhados.
- Confirmação manual.
- Políticas de antecedência mínima.
- Limite de cancelamento.

## 12. Estrutura de dados de negócio

### 12.1 Estabelecimento
Cada estabelecimento deve possuir:
- Nome comercial.
- Identidade visual.
- Dados de contato.
- Link público de agendamento.
- Configurações de agenda.

### 12.2 Profissional
Cada profissional deve possuir:
- Nome.
- Status ativo/inativo.
- Serviços que atende.
- Jornada de trabalho.
- Agenda vinculada.
- Ordem de exibição no fluxo de agendamento.

### 12.3 Serviço
Cada serviço deve possuir:
- Nome.
- Descrição curta.
- Duração.
- Preço opcional.
- Profissionais aptos a executar.
- Status ativo/inativo.

### 12.4 Agendamento
Cada agendamento deve registrar:
- Estabelecimento.
- Serviço.
- Profissional.
- Cliente.
- Data e hora.
- Status do agendamento.
- Origem do agendamento.
- Referência de calendário.

## 13. Status do agendamento

O negócio deve considerar pelo menos estes status:

- Pendente.
- Confirmado.
- Cancelado.
- Concluído.
- Não compareceu.

Mesmo que o MVP trabalhe com uma versão simplificada, esses estados ajudam a evoluir o produto para gestão operacional mais completa.

## 14. Multi-funcionários

O suporte a múltiplos funcionários é obrigatório no MVP.  
Essa é uma regra central do negócio e não um extra.

### 14.1 O que isso significa
- Um mesmo estabelecimento pode ter vários profissionais.
- Cada profissional tem sua própria disponibilidade operacional.
- O cliente final pode agendar com base no profissional.
- O gestor precisa ter visão central de toda a equipe.

### 14.2 Implicação de negócio
O produto não deve ser desenhado como agenda de profissional autônomo individual.  
Ele deve nascer com mentalidade de operação de equipe, mesmo que também atenda negócios menores.

## 15. Generalização para outros nichos

Embora o nicho inicial seja salão/barbearia, o produto deve evitar linguagem excessivamente específica na sua essência.  
A plataforma deve permitir adaptação para outros segmentos de agendamento.

### 15.1 O que pode ser específico no início
- Comunicação comercial.
- Exemplos do site.
- Casos de uso.
- Estratégia de validação.

### 15.2 O que não deve ser preso ao nicho
- Modelo de serviços.
- Estrutura de profissionais.
- Regras de disponibilidade.
- Marca do cliente.
- Link de agendamento.
- Fluxo principal de reserva.

## 16. Hipóteses de validação

O negócio parte de algumas hipóteses que precisam ser testadas:

1. Estabelecimentos sentem dor real com o agendamento atual.
2. Essa dor é frequente o suficiente para justificar contratação.
3. Um link de autoagendamento melhora a operação.
4. O gestor valoriza centralização de equipe.
5. O whitelabel aumenta atratividade comercial.
6. O nicho inicial aceita integrar sua rotina ao Google Agenda.
7. O produto pode depois ser expandido para outros mercados sem reconstrução total.

## 17. Critérios de sucesso do MVP

O MVP será considerado promissor se conseguir demonstrar:

- Interesse real de estabelecimentos em testar.
- Uso frequente do link de agendamento.
- Capacidade do gestor de configurar a operação sem fricção excessiva.
- Funcionamento consistente para múltiplos profissionais.
- Redução da dependência de atendimento manual para marcar horários.
- Percepção clara de valor no modelo whitelabel.

## 18. Limites do MVP

Para evitar dispersão, o MVP não precisa nascer com tudo.  
A prioridade é resolver o agendamento centralizado com marca própria.

Itens que podem ficar fora da primeira versão:
- Aplicativo mobile nativo.
- Marketplace aberto ao público.
- Programa de fidelidade.
- Gestão financeira profunda.
- Estoque.
- Comissionamento avançado.
- Automação complexa de marketing.
- Múltiplas integrações além do essencial.

## 19. Riscos de negócio

Os principais riscos percebidos são:

- O mercado já possui concorrentes consolidados, então o produto precisa ter posicionamento claro [web:11][web:8].
- O whitelabel sozinho pode não ser suficiente como diferencial se a operação básica não for muito boa.
- Alguns estabelecimentos podem preferir continuar no WhatsApp por hábito.
- A configuração inicial pode virar barreira se for complexa demais.
- A expansão para qualquer nicho pode enfraquecer a proposta comercial se acontecer cedo demais.

## 20. Estratégia recomendada de validação

A validação deve começar com poucos clientes piloto do nicho inicial.  
O objetivo não é escalar imediatamente, e sim aprender com uso real.

### 20.1 Perguntas a validar com pilotos
- Como eles agendam hoje?
- Onde perdem mais tempo?
- Quem organiza a agenda?
- Existe problema com múltiplos funcionários?
- O cliente final consegue usar o link sem ajuda?
- O gestor entende valor em centralizar tudo?
- O Google Agenda se encaixa bem na rotina deles?

### 20.2 Sinais positivos
- Uso recorrente da plataforma.
- Redução de conversa manual para marcar horário.
- Aprovação da experiência do cliente final.
- Facilidade do gestor em operar o sistema.
- Pedido por continuidade após o teste.

## 21. Definição resumida do produto

Produto SaaS whitelabel de agendamentos para negócios de serviços, com foco inicial em salões de beleza e barbearias, suporte a múltiplos funcionários, gestão centralizada pelo gestor e integração com Google Agenda.

## 22. Decisões já tomadas

- O nicho inicial de entrada será salão de beleza/barbearia.
- O produto deve ser adaptável para qualquer nicho no futuro.
- O sistema será whitelabel.
- Haverá login de gestor.
- Haverá link para o usuário final.
- O MVP deve suportar múltiplos funcionários.
- No MVP, os funcionários serão gerenciados pela conta do gestor.
- A integração base será com Google Agenda.
- O foco do documento é regra de negócio, não desenho técnico.

## 23. Questões em aberto

Estas decisões ainda podem precisar de definição antes do planejamento técnico detalhado:

- O cliente final poderá cancelar sozinho?
- O cliente final poderá reagendar sozinho?
- Haverá confirmação manual ou automática?
- Haverá cobrança antecipada ou não?
- O cliente precisará informar apenas nome e telefone ou mais dados?
- O agendamento será confirmado imediatamente ou dependerá de aprovação do salão?
- O link público terá domínio customizado no MVP ou apenas subpágina padrão?
- O sistema enviará notificações por WhatsApp, e-mail ou ambos?
