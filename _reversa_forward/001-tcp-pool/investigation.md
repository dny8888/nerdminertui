# Investigation: Protocolo Stratum V1

## Pesquisa de Fundo

Para se comunicar com uma pool de mineração Bitcoin (como a public-pool.io) na porta 21496, o NerdTUI precisa implementar o **Stratum Mining Protocol**. É um protocolo focado em JSON-RPC 1.0 por cima de TCP bruto, operando linha a linha (cada linha terminada em `\n` é uma mensagem completa JSON).

### Fluxo de Comunicação Esperado

1. **Início e Handshake**
   O cliente conecta no socket TCP e imediatamente envia:
   ```json
   {"id": 1, "method": "mining.subscribe", "params": ["NerdTUI/1.0"]}
   ```
   A pool responde com as configurações da sessão: `extranonce1` e o tamanho do `extranonce2` (`extranonce2_size`).
   
2. **Autorização**
   Com os dados iniciais aceitos, o cliente se autoriza (é assim que a pool sabe a quem enviar a recompensa):
   ```json
   {"id": 2, "method": "mining.authorize", "params": ["<EnderecoBTC>.<WorkerName>", "x"]}
   ```
   A senha normalmente é `x` para mineração solo ou pools que não exigem senha real. A pool responde com `result: true`.

3. **Recepção de Tarefas (Notify)**
   A pool envia (sem o cliente pedir) métodos de notificação:
   ```json
   {"id": null, "method": "mining.notify", "params": ["job_id", "prevhash", "coinb1", "coinb2", ...]}
   ```
   Sempre que um notify chega (especialmente com a flag `clean_jobs=true`), o minerador descarta seu trabalho atual e começa a minerar no novo block header gerado.

4. **Submissão de Tarefas (Submit)**
   Quando o minerador encontra um hash menor que o target, ele submete:
   ```json
   {"id": 4, "method": "mining.submit", "params": ["<EnderecoBTC>.<WorkerName>", "job_id", "extranonce2", "ntime", "nonce"]}
   ```
   A pool avalia o share enviado e devolve se foi aceito (`result: true`).

## Alternativas Avaliadas

* **Usar bibliotecas prontas do Stratum**: Geralmente são velhas, abandonadas ou excessivamente complexas (cobre Stratum V2, Ethereum, etc). Para o nosso escopo simples de envio do block header e hashrate solo, o Go padrão lidará muito bem via TCP socket simples.
* **Manter o miner em HTTP polling**: O HTTP não tem a eficiência necessária para receber novos trabalhos imediatamente (push), atrasando shares e gerando trabalho inválido (stale shares). O Stratum TCP é obrigatório.
