# Impacto no Legado: 001-tcp-pool
Data: 2026-05-30

## Arquivos Afetados

| Arquivo afetado | Componente | Tipo | Severidade | Justificativa |
|-----------------|------------|------|------------|---------------|
| `internal/config/config.go` | Data | `regra-alterada` | MEDIUM | Adição de dados de Stratum e Worker Name, renomeando PoolURL. |
| `internal/worker/stratum.go` | Worker | `componente-novo` | HIGH | Inserção do formato JSON-RPC Stratum V1. |
| `internal/worker/fetcher.go` | Worker | `regra-alterada` | CRITICAL | Substituição do stub HTTP pelo dialer TCP Stratum e laço bufio. |
| `internal/worker/miner.go` | Worker | `regra-alterada` | HIGH | Canalização dos notifies de jobs (`jobCh`) e submits pela network. |
| `cmd/tui/main.go` | Bootstrapper | `regra-alterada` | HIGH | Instanciação e execução em goroutine do novo StratumClient e injeção do canal de jobs. |
| `internal/model/state.go` | State | `regra-alterada` | MEDIUM | Adição do `ConnectionStatus` e `WorkerName`. |

## Diff Conceitual

A arquitetura Worker, que era baseada em um HTTP client ou mock (pull architecture), sofreu a maior transformação, ancorando-se ativamente num socket TCP através do `net.Dial` com backoff reativo (push architecture). O `MinerWorker` passou a escutar um channel (`jobCh`) bloqueante quando desocupado para consumir os novos blocos provenientes do método `mining.notify` da pool. Os eventos foram adicionados diretamente ao barramento do Bubbletea.

## Preservadas
- **MUV Pattern**: A interface `model.AppState` continua sendo imutável e repassada pelo loop de updates da UI.
- **Isolamento de Threads**: TUI rodando no thread principal (Bubbletea) continua isolado dos Workers de mineração que rodam via Goroutines sem travar a renderização.
- **Graceful Shutdown**: Todos os métodos de IO respeitam o `context.Done()` original injetado na `main`.

## Modificadas
- **RN-02**: A validação e inicialização de configurações agora depende do par `PoolAddress` e `WorkerName`, extinguindo a variável `PoolURL`.
- O pool client agora precisa receber um loop passivo independente (`Run(ctx)`) invés de ser triggado pela `tea.Tick`.
