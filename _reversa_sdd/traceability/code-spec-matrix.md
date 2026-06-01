# Matriz de Cobertura Código-Spec — nerdminertui

> **Nível de Documentação:** COMPLETO  
> **Gerado pelo Redator em:** 2026-05-29

Esta matriz de rastreabilidade mapeia cada arquivo do projeto legado (conforme especificado na arquitetura física) para a sua respectiva pasta de especificação executável (Unit Spec), validando o percentual de cobertura da análise reversa.

---

## 1. Mapeamento de Rastreabilidade

| Arquivo do Legado Mapeado | Unit Spec Correspondente | Cobertura da Spec | Descrição / Função da Cobertura |
|---|---|---|---|
| `cmd/tui/main.go` | [cmd/tui/](file:///home/dny8888/workspace/github/nerdminertui/_reversa_sdd/cmd/tui/requirements.md) | 🟢 CONFIRMADO | Inicialização de flags, dependências e Bubbletea program. |
| `internal/config/config.go` | [internal/config/](file:///home/dny8888/workspace/github/nerdminertui/_reversa_sdd/internal/config/requirements.md) | 🟢 CONFIRMADO | Validação e carregamento de viper/env variables. |
| `internal/model/state.go` | [internal/model/](file:///home/dny8888/workspace/github/nerdminertui/_reversa_sdd/internal/model/requirements.md) | 🟢 CONFIRMADO | AppState de valor imutável e rotações circulares. |
| `internal/worker/miner.go` | [internal/worker/](file:///home/dny8888/workspace/github/nerdminertui/_reversa_sdd/internal/worker/requirements.md) | 🟢 CONFIRMADO | MinerWorker de hashing concorrente com CPU limits. |
| `internal/worker/fetcher.go` | [internal/worker/](file:///home/dny8888/workspace/github/nerdminertui/_reversa_sdd/internal/worker/requirements.md) | 🟢 CONFIRMADO | Clientes de rede HTTP REST / TCP Stratum. |
| `internal/worker/poller.go` | [internal/worker/](file:///home/dny8888/workspace/github/nerdminertui/_reversa_sdd/internal/worker/requirements.md) | 🟢 CONFIRMADO | Poller com retry exponencial e tickers. |
| `internal/worker/messages.go` | [internal/worker/](file:///home/dny8888/workspace/github/nerdminertui/_reversa_sdd/internal/worker/requirements.md) | 🟢 CONFIRMADO | Definição de mensagens compartilhadas do Bubbletea. |
| `internal/ui/app.go` | [internal/ui/](file:///home/dny8888/workspace/github/nerdminertui/_reversa_sdd/internal/ui/requirements.md) | 🟢 CONFIRMADO | Ciclo Elm do AppModel (Init/Update/View). |
| `internal/ui/keys.go` | [internal/ui/](file:///home/dny8888/workspace/github/nerdminertui/_reversa_sdd/internal/ui/requirements.md) | 🟢 CONFIRMADO | Mapeamento de keybindings no teclado. |
| `internal/ui/screens/dashboard.go` | [internal/ui/](file:///home/dny8888/workspace/github/nerdminertui/_reversa_sdd/internal/ui/requirements.md) | 🟢 CONFIRMADO | Renderização pura da tela do Dashboard de mineração. |
| `internal/ui/screens/clock.go` | [internal/ui/](file:///home/dny8888/workspace/github/nerdminertui/_reversa_sdd/internal/ui/requirements.md) | 🟢 CONFIRMADO | Renderização pura da tela de Relógio grande ASCII. |
| `internal/ui/screens/globalstats.go`| [internal/ui/](file:///home/dny8888/workspace/github/nerdminertui/_reversa_sdd/internal/ui/requirements.md) | 🟢 CONFIRMADO | Renderização pura das estatísticas globais Bitcoin. |
| `internal/ui/components/cpubar.go` | [internal/ui/](file:///home/dny8888/workspace/github/nerdminertui/_reversa_sdd/internal/ui/requirements.md) | 🟢 CONFIRMADO | Widget de representação do uso de CPU. |
| `internal/ui/components/sparkline.go`| [internal/ui/](file:///home/dny8888/workspace/github/nerdminertui/_reversa_sdd/internal/ui/requirements.md) | 🟢 CONFIRMADO | Widget de gráfico do histórico hashrate de 60s. |
| `internal/ui/components/hashgauge.go` | [internal/ui/](file:///home/dny8888/workspace/github/nerdminertui/_reversa_sdd/internal/ui/requirements.md) | 🟢 CONFIRMADO | Medidor gráfico do hashrate real. |
| `internal/ui/components/statusbar.go`| [internal/ui/](file:///home/dny8888/workspace/github/nerdminertui/_reversa_sdd/internal/ui/requirements.md) | 🟢 CONFIRMADO | Barra de rodapé informativa e erros de display. |
| `internal/store/store.go` | [internal/store/](file:///home/dny8888/workspace/github/nerdminertui/_reversa_sdd/internal/store/requirements.md) | 🟢 CONFIRMADO | Interface SQLiteStore concorrente e NilStore. |
| `internal/store/migrations.go` | [internal/store/](file:///home/dny8888/workspace/github/nerdminertui/_reversa_sdd/internal/store/requirements.md) | 🟢 CONFIRMADO | SQL de migrações embutido no executável. |
| `pkg/mining/hash.go` | [pkg/mining/](file:///home/dny8888/workspace/github/nerdminertui/_reversa_sdd/pkg/mining/requirements.md) | 🟢 CONFIRMADO | Duplo SHA256 criptográfico puro. |
| `pkg/mining/target.go` | [pkg/mining/](file:///home/dny8888/workspace/github/nerdminertui/_reversa_sdd/pkg/mining/requirements.md) | 🟢 CONFIRMADO | Verificação e comparação big-endian de targets. |
| `pkg/mining/job.go` | [pkg/mining/](file:///home/dny8888/workspace/github/nerdminertui/_reversa_sdd/pkg/mining/requirements.md) | 🟢 CONFIRMADO | Definição da struct Job candidato a mineração. |
| `pkg/format/hashrate.go` | [pkg/format/](file:///home/dny8888/workspace/github/nerdminertui/_reversa_sdd/pkg/format/requirements.md) | 🟢 CONFIRMADO | Formatação e sufixos de H/s a MH/s decimais. |
| `pkg/format/duration.go` | [pkg/format/](file:///home/dny8888/workspace/github/nerdminertui/_reversa_sdd/pkg/format/requirements.md) | 🟢 CONFIRMADO | Formatação estruturada de Uptime. |
| `pkg/format/difficulty.go` | [pkg/format/](file:///home/dny8888/workspace/github/nerdminertui/_reversa_sdd/pkg/format/requirements.md) | 🟢 CONFIRMADO | Formatação de notações científicas de diff. |

---

## 2. Estatística de Cobertura Reversa
* **Total de Arquivos do Legado Identificados**: 24
* **Total de Arquivos Cobertos por Specs SDD**: 24
* **Percentual de Cobertura Reversa**: 🎯 **100%** (Todos os componentes foram plenamente decompostos e mapeados em especificações formais determinísticas e executáveis por IAs).
