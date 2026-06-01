# Actions: Comunicação TCP Real com a Pool (Stratum)

> Identificador: `001-tcp-pool`
> Data: `2026-05-29`
> Roadmap: `_reversa_forward/001-tcp-pool/roadmap.md`

## Resumo

| Métrica | Valor |
|---------|-------|
| Total de ações | 6 |
| Paralelizáveis (`[//]`) | 1 |
| Maior cadeia de dependência | 3 (T001 -> T002 -> T003 -> T004) |

## Fase 1, Preparação

| ID | Descrição | Dependências | Paralelismo | Arquivo alvo | Confidência | Status |
|----|-----------|--------------|-------------|--------------|-------------|--------|
| T001 | Adicionar `PoolAddress` e `WorkerName` à struct `Config` com default values (`public-pool.io:21496` e `.nerdtui`) | - | `[//]` | `internal/config/config.go` | 🟢 | `[X]` |

## Fase 2, Testes

Nenhum teste isolado planejado antes do core. A mineração em si depende do acoplamento com o worker local e socket vivo.

## Fase 3, Núcleo

| ID | Descrição | Dependências | Paralelismo | Arquivo alvo | Confidência | Status |
|----|-----------|--------------|-------------|--------------|-------------|--------|
| T002 | Criar structs JSON-RPC V1 (Subscribe, Authorize, Submit, Notify) em um arquivo de domínio Stratum. | T001 | - | `internal/worker/stratum.go` | 🟢 | `[X]` |
| T003 | Refatorar fetcher para abrir TCP (`net.Dial`), ler loop `bufio.Scanner` e executar Subscribe/Authorize com a config. | T002 | - | `internal/worker/fetcher.go` | 🟢 | `[X]` |

## Fase 4, Integração

| ID | Descrição | Dependências | Paralelismo | Arquivo alvo | Confidência | Status |
|----|-----------|--------------|-------------|--------------|-------------|--------|
| T004 | Conectar o envio de jobs (`mining.notify`) à goroutine do Miner e canalizar submits de shares de volta pro fetcher TCP. | T003 | - | `internal/worker/miner.go` | 🟢 | `[X]` |
| T005 | Adaptar o poller para capturar quebra de EOF no socket TCP e acionar redial com backoff exponencial. | T003 | - | `internal/worker/poller.go` | 🟢 | `[X]` |

## Fase 5, Polimento

| ID | Descrição | Dependências | Paralelismo | Arquivo alvo | Confidência | Status |
|----|-----------|--------------|-------------|--------------|-------------|--------|
| T006 | Emitir mensagens `tea.Msg` adequadas (ConnectionStatusMsg) para atualizar a UI entre "Conectando", "Conectado" e "Desconectado". | T003 | - | `internal/worker/messages.go` | 🟢 | `[X]` |

## Notas de execução


## Histórico de alterações

| Data | Alteração | Autor |
|------|-----------|-------|
| 2026-05-29 | Versão inicial gerada por `/reversa-to-do` | reversa |
