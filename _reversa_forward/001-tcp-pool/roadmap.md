# Roadmap: Comunicação TCP Real com a Pool (Stratum)

> Identificador: `001-tcp-pool`
> Data: `2026-05-29`
> Requirements: `_reversa_forward/001-tcp-pool/requirements.md`
> Confidência: 🟢 CONFIRMADO, 🟡 INFERIDO, 🔴 LACUNA

## 1. Resumo da abordagem

Iremos alterar a camada do `internal/worker/fetcher.go` (atualmente um stub) para que utilize a biblioteca padrão `net` do Go, abrindo uma conexão TCP contínua com a pool configurada. Iremos realizar o handshake do protocolo Stratum V1 enviando as mensagens JSON-RPC de `mining.subscribe` e `mining.authorize`. O `config.yaml` será atualizado para suportar o endereço da pool e o nome do worker, respeitando o fallback que combinamos ("public-pool.io:21496" e ".nerdtui"). Também adaptaremos o loop do `miner.go` para que faça envios reais via TCP ao encontrar um share válido.

## 2. Princípios aplicados

Nenhum princípio documentado em `principles.md` (arquivo inexistente). O projeto base já orienta a não usar dependências externas não necessárias (CGO=0), mantendo a implementação TCP com o pacote `net` nativo.

## 3. Decisões técnicas

| ID | Decisão | Justificativa | Alternativas descartadas | Confidência |
|----|---------|----------------|--------------------------|-------------|
| D-01 | Usar TCP puro `net.Dial` e `bufio.Scanner` | O protocolo Stratum V1 usa NDJSON (Newline Delimited JSON). O scanner nativo lida perfeitamente com a leitura de linhas `\n` sem onerar a CPU ou adicionar dependências. | `gorilla/websocket` (pool não usa WS puro) | 🟢 |
| D-02 | Poller reaproveitado | O poller de retry existente (`internal/worker/poller.go`) pode coordenar a reconexão TCP em vez de só focar em HTTP. | Escrever novo watcher TCP do zero | 🟢 |

## 4. Premissas

Nenhuma premissa inferida (as dúvidas sobre fallbacks foram resolvidas diretamente na sessão de `reversa-clarify`).

## 5. Delta arquitetural

| Componente | Arquivo de origem no legado | Tipo de mudança | Resumo |
|------------|------------------------------|-----------------|--------|
| Config | `_reversa_sdd/architecture.md#2-2-internal-config` | regra-alterada | Inclusão de `PoolAddress` e `WorkerName` com default values. |
| Fetcher | `_reversa_sdd/architecture.md#2-4-internal-worker` | regra-alterada | Substituição do mock por cliente TCP Stratum V1 real, lendo e roteando mensagens. |
| MinerWorker | `_reversa_sdd/architecture.md#2-4-internal-worker` | regra-alterada | Passará a chamar o método `submit` TCP do fetcher e usará o `extranonce` real emitido pela pool. |

## 6. Delta no modelo de dados

- Resumo das mudanças: Configuração do app ganha novos campos para a mineração solo real, permitindo a conexão em qualquer pool.
- Detalhe completo em: `_reversa_forward/001-tcp-pool/data-delta.md`

## 7. Delta de contratos externos

| Contrato | Tipo | Arquivo de detalhe |
|----------|------|--------------------|
| Stratum V1 | TCP/JSON-RPC | `_reversa_forward/001-tcp-pool/interfaces/stratum.md` |

## 8. Plano de migração

1. Adicionar campos na struct `Config` e garantir fallback na inicialização (`config.go`).
2. Criar pacote genérico `stratum` ou embutir no `fetcher.go` as chamadas `net.Dial` e parse de JSON-RPC.
3. Testar localmente a conexão, vendo se a pool aceita o endereço fornecido.
4. Conectar a notificação de novos trabalhos (jobs) da conexão TCP ao canal que o minerador consome.

## 9. Riscos e mitigações

| Risco | Impacto | Probabilidade | Mitigação |
|-------|---------|---------------|-----------|
| Formato JSON-RPC inesperado vindo da pool | alto | baixo | Uso de structs abertas (map[string]interface{}) ou json.RawMessage para partes flexíveis. |
| Queda de conexão quebra o app inteiro (Panic) | alto | médio | Não usar `log.Fatal` fora da main. Propagar erros via canais, para a TUI notificar o erro suavemente. |

## 10. Critério de pronto

- [ ] Todas as ações do `actions.md` marcadas `[X]`
- [ ] `cross-check.md` (se executado) sem CRITICAL nem HIGH
- [ ] `regression-watch.md` gerado

## 11. Histórico de alterações

| Data | Alteração | Autor |
|------|-----------|-------|
| 2026-05-29 | Versão inicial gerada por `/reversa-plan` | reversa |
