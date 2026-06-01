# Contrato Externo: Stratum V1

**Tipo:** TCP Socket (Newline Delimited JSON-RPC 1.0)

## Transporte
O protocolo requer a abertura de uma conexão TCP estrita. O NerdTUI deve utilizar `bufio.NewScanner` para varrer linha a linha até cada byte `\n`. Não há keep-alive de socket padrão; a manutenção da sessão TCP viva em firewalls não costuma necessitar de Nagle em pool de mineração agressiva, mas um idle timeout razoável no net.Dial (ex: 3 minutos sem notify) forçaria reconnect em caso de zumbi connection.

## Tipos de Requisição do NerdTUI

### 1. Subscribe
**Sentido**: NerdTUI -> Pool
```json
{"id": 1, "method": "mining.subscribe", "params": ["NerdTUI/1.0"]}
```
**Resposta (Idempotente)**:
```json
{"id": 1, "result": [ [ ["mining.set_difficulty", "b4b6693b72a50c7116db31..."], ["mining.notify", "ae6812eb4cd7735a302a8a9dd95cf71f"] ], "08000002", 4 ], "error": null}
```
`result[1]` é o extranonce1. `result[2]` é o extranonce2_size.

### 2. Authorize
**Sentido**: NerdTUI -> Pool
```json
{"id": 2, "method": "mining.authorize", "params": ["bc1qxxx.nerdtui", "x"]}
```
**Resposta**:
```json
{"id": 2, "result": true, "error": null}
```

### 3. Submit
**Sentido**: NerdTUI -> Pool (Quando um share válido é encontrado)
```json
{"id": 4, "method": "mining.submit", "params": ["bc1qxxx.nerdtui", "job_id", "extranonce2", "ntime", "nonce"]}
```
**Resposta**:
```json
{"id": 4, "result": true, "error": null}
```

## Tipos de Eventos da Pool

A pool envia dados sem id ou requisicao, que devem ser consumidos sempre que chegam.

### 1. Set Difficulty
```json
{"id": null, "method": "mining.set_difficulty", "params": [1024]}
```

### 2. Notify
```json
{"id": null, "method": "mining.notify", "params": ["bf", "4d16b6f85af6e2198f44ae2a6de67f78487ae5611b77c6c0440b921e00000000", ...]}
```
O array de parâmetros dita todos os blocos do merkle tree e versões necessárias.

## Timeout e Falhas
Quedas no socket fecharão o `io.Reader`. Quando `Scanner.Scan()` retornar falso, deve emitir um erro estruturado de forma não fatal (ex: `DisconnectMsg`) para a TUI e ativar um `time.Sleep` com exponencial de 1 a 30 segundos no worker até reiniciar o loop principal do fetcher e criar novo dial.
