# ADR 0001 — Tratamento de cancelamentos de requisições HTTP

- Status: Aceito
- Data: 2026-10-03

## Contexto

Durante testes de carga com autocannon, foram observadas operações
interrompidas com context.Canceled e contexto HTTP também cancelado.
Os eventos se concentraram no encerramento dos testes, comportamento
compatível com o fechamento das conexões pelo cliente.

A aplicação tratava esses eventos como erros internos, registrando
stack traces e tentando responder com HTTP 500. Ao deixar de escrever
essa resposta, os middlewares classificavam a ausência de status como
200, produzindo uma indicação incorreta de sucesso.

O defer cancel() nas operações de acesso ao banco está correto:
ele libera recursos do contexto derivado ao sair da função e não
cancela o contexto HTTP original.

## Decisão

Centralizar em errors.go o tratamento de operações que retornem
context.Canceled quando o contexto HTTP também estiver cancelado.

Nesses casos, registrar o evento em nível Info e encerrar o fluxo
sem tentar escrever uma resposta JSON. Falhas inesperadas continuam
sendo tratadas como erros internos.

Nos logs de conclusão:

- Incluir request_id para correlacionar os eventos.
- Registrar context_state separadamente do status.
- Preservar qualquer status já escrito pelo handler.
- Usar status 0 quando nenhum status tiver sido escrito e o contexto
  estiver cancelado. Esse valor existe apenas na observabilidade.

Nas métricas, utilizar status="canceled" para cancelamentos sem status
HTTP escrito. Preservar os códigos já escritos nos demais casos.

Manter a observação da duração e o decremento do Gauge de requisições
em andamento também nos caminhos de cancelamento.

## Consequências

Cancelamentos reconhecidos deixam de gerar stack traces de erro e
não são automaticamente classificados como HTTP 500 ou como sucesso.

Status escrito pela API não comprova entrega da resposta ao cliente.
Uma requisição pode registrar status 200 e context_state="canceled".

A série status="canceled" representa somente cancelamentos sem status
escrito, não todos os contextos cancelados.

A validação em carga confirmou a classificação nos logs e métricas
e o retorno de http_requests_in_flight a zero após o processamento.

Expiração de prazo é um cenário distinto e não está abrangida por
esta decisão.
