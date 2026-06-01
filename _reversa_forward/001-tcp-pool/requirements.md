# Requirements: Comunicação TCP Real com a Pool (Stratum)

> Identificador: `001-tcp-pool`
> Data: `2026-05-29`
> Pasta da extração reversa: `_reversa_sdd/`
> Confidência: 🟢 CONFIRMADO, 🟡 INFERIDO, 🔴 LACUNA / DÚVIDA

## 1. Resumo executivo

Esta feature substitui o mock de rede atual por uma implementação real do protocolo Stratum TCP (JSON-RPC 1.0). O objetivo é conectar o NerdTUI a uma pool de mineração real (como a public-pool.io), realizar o handshake com o endereço BTC configurado pelo usuário, e receber trabalhos válidos (`mining.Job`) para que o hashrate gerado contribua legitimamente na busca de blocos na rede Bitcoin.

## 2. Contexto a partir do legado

| Fonte | Trecho relevante | Confidência |
|-------|------------------|-------------|
| `_reversa_sdd/architecture.md#2-topologia-de-integracoes-externas` | A comunicação principal é feita via JSON-RPC 1.0 via TCP Stratum na porta 21496 para handshake e polling de jobs. | 🟢 |
| `_reversa_sdd/domain.md#rn-02-limitacao-de-configuracoes` | O programa exige um endereço BTC real (`MockMining == false`) para se conectar e evitar desperdício. | 🟢 |
| `_reversa_sdd/code-analysis.md#2-4-internal-worker` | As conexões TCP e de rede usam o poller com retry exponencial. | 🟢 |

## 3. Personas e cenários de uso

| Persona | Objetivo | Cenário-chave |
|---------|----------|---------------|
| Minerador Solo | Contribuir com hashrate real na rede | O usuário configura o endereço, o TUI conecta na public-pool.io e mostra hashrate real sem simuladores. |

## 4. Regras de negócio novas ou alteradas

1. **RN-06:** O `fetcher.go` deve realizar a conexão TCP bruta ao servidor da pool, usando o protocolo Stratum. 🟢
   - Tipo: nova
2. **RN-07:** O MinerWorker deve processar e utilizar os valores reais `extranonce1` e `extranonce2` recebidos do job da pool. 🟢
   - Tipo: nova
3. **RN-08:** O envio de shares para a pool deve ocorrer através do mesmo canal TCP ativo (JSON-RPC `mining.submit`). 🟢
   - Tipo: nova

## 5. Requisitos Funcionais

| ID | Requisito | Prioridade | Critério de aceite | Confidência |
|----|-----------|------------|--------------------|-------------|
| RF-01 | Estabelecer conexão TCP (Stratum) | Must | Conexão não cai e mantem-se ativa com envio periódico de keep-alive se necessário. | 🟢 |
| RF-02 | Handshake e Autenticação (mining.subscribe / mining.authorize) | Must | O TUI envia o endereço BTC e recebe autorização e os dados iniciais do extranonce. | 🟢 |
| RF-03 | Recepção de Tarefas (mining.notify) | Must | O TUI escuta e repassa jobs recebidos da pool para a thread mineradora local. | 🟢 |
| RF-04 | Submissão de Shares (mining.submit) | Must | Ao encontrar um share válido localmente, o TUI o envia via TCP para a pool e verifica o `result: true`. | 🟢 |

## 6. Requisitos Não Funcionais

| Tipo | Requisito | Evidência ou justificativa | Confidência |
|------|-----------|----------------------------|-------------|
| Confiabilidade | Reconexão automática e backoff exponencial | Se a pool cair, o `poller` de rede não deve travar a UI; deve tentar de novo sem bloquear os goroutines. | 🟢 |
| Observabilidade | Emissão de mensagens (tea.Msg) claras | Eventos de erro de socket TCP devem virar alertas no StatusBar da TUI. | 🟡 |

## 7. Critérios de Aceitação

```gherkin
Cenário: Conexão bem sucedida ao iniciar o app
  Dado que o app é iniciado sem a flag --mock e com o endereço BTC válido
  Quando o worker de rede entra em ação
  Então ele estabelece a conexão TCP Stratum
  E a TUI exibe "Conectado" na barra de status

Cenário: Queda de conexão com a pool
  Dado que a conexão TCP foi estabelecida
  Quando o soquete de rede cai ou dá timeout
  Então a barra de status exibe "Desconectado"
  E o fetcher entra em loop de retry exponencial até reconectar
```

## 8. Prioridade MoSCoW

| Item | MoSCoW | Justificativa |
|------|--------|---------------|
| RF-01 | Must | Sem TCP não tem mineração real |
| RF-02 | Must | Sem subscrição, o servidor não envia jobs |
| RF-03 | Must | Precisa ouvir trabalhos da pool |
| RF-04 | Must | O share precisa chegar na rede para virar recompensa |

## 9. Esclarecimentos

### Sessão 2026-05-29
- **Q:** Qual endereço Stratum deve ser o padrão para a conexão da pool?
  **R:** Usar "public-pool.io:21496" como padrão (fallback), mas ler primeiro do `config.yaml` caso o usuário queira configurar uma pool diferente.
- **Q:** Como deve ser composto o nome do worker que enviaremos para a pool no momento da autenticação?
  **R:** Adicionar um novo campo `WorkerName` no `config.yaml`, mas usar `.nerdtui` como padrão se ficar vazio.

## 10. Lacunas

Nenhuma dúvida pendente.

## 11. Histórico de alterações

| Data | Alteração | Autor |
|------|-----------|-------|
| 2026-05-29 | Versão inicial gerada por `/reversa-requirements` | reversa |
