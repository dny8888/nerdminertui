# Reconstruction Plan — NerdTUI

**Fonte:** original
**Stack:** Go 1.26, Bubbletea + Lipgloss + Bubbles, Viper, modernc.org/sqlite, testify, goleak
**Gerado em:** 2026-05-29
**Status:** 13 tarefas | 13 concluídas | 0 pendentes

---

## Alertas de pré-voo

Nenhum gap crítico 🔴 pendente. Todos os 4 gaps arquiteturais foram resolvidos na fase de revisão. Pode iniciar com segurança.

---

## Tarefas

### Tarefa 01 — Scaffolding do Projeto
**Status:** done ✅
**Lê:** `_reversa_sdd/architecture.md`, `_reversa_sdd/dependencies.md`
**Constrói:** `go.mod`, `Makefile`, `.golangci.yml`, `AGENTS.md`, estrutura de diretórios vazia (`cmd/tui/`, `internal/config/`, `internal/model/`, `internal/worker/`, `internal/ui/`, `internal/ui/screens/`, `internal/ui/components/`, `internal/store/`, `pkg/mining/`, `pkg/format/`, `testutil/`)
**Pronto quando:** `go mod tidy` executa sem erros, `CGO_ENABLED=0 go build ./...` compila (mesmo sem código funcional), e a árvore de diretórios reflete a arquitetura especificada

---

### Tarefa 02 — Schema do Banco de Dados
**Status:** done ✅
**Lê:** `_reversa_sdd/erd-complete.md`, `_reversa_sdd/data-dictionary.md`, `_reversa_sdd/internal/store/design.md`
**Constrói:** `internal/store/migrations.go` (DDL embeddado via `//go:embed`), schema SQL com tabela `hashrate_history`, índice e pragmas WAL
**Pronto quando:** Schema compila e pode ser executado em SQLite `:memory:` criando a tabela com tipos, constraints e índice corretos

---

### Tarefa 03 — Entidades de Domínio (internal/model)
**Status:** done ✅
**Lê:** `_reversa_sdd/domain.md`, `_reversa_sdd/data-dictionary.md`, `_reversa_sdd/internal/model/requirements.md`, `_reversa_sdd/internal/model/design.md`, `_reversa_sdd/internal/model/tasks.md`
**Constrói:** `internal/model/state.go` — `AppState` struct (valor imutável), `ScreenID` type, constantes de domínio, método `WithHashRate`
**Pronto quando:** `AppState` é um value type sem ponteiros internos, `WithHashRate` retorna cópia com histórico rotacionado, testes unitários cobrem imutabilidade e FIFO rotation (100% branch coverage)

---

### Tarefa 04 — Máquinas de Estado
**Status:** done ✅
**Lê:** `_reversa_sdd/state-machines.md`
**Constrói:** Validação de transições no `ScreenID` (integrado ao `internal/model/state.go` ou como método dedicado)
**Pronto quando:** Transições cíclicas `Dashboard→Clock→GlobalStats→Dashboard` implementadas e testadas

---

### Tarefa 05 — pkg/mining (folha, sem dependências internas)
**Status:** done ✅
**Lê:** `_reversa_sdd/pkg/mining/requirements.md`, `_reversa_sdd/pkg/mining/design.md`, `_reversa_sdd/pkg/mining/tasks.md`
**Constrói:** `pkg/mining/hash.go` (`SHA256d`, `HashHeader`), `pkg/mining/target.go` (`MeetsTarget`, `DifficultyFromHash`), `pkg/mining/job.go` (`Job` struct)
**Pronto quando:** Testes com vetores NIST, comparação big-endian correta, boundary hash==target retorna false, ≥90% branch coverage

---

### Tarefa 06 — pkg/format (folha, sem dependências internas)
**Status:** done ✅
**Lê:** `_reversa_sdd/pkg/format/requirements.md`, `_reversa_sdd/pkg/format/design.md`, `_reversa_sdd/pkg/format/tasks.md`
**Constrói:** `pkg/format/hashrate.go`, `pkg/format/duration.go`, `pkg/format/difficulty.go`
**Pronto quando:** FormatHashRate trata 0/K/M/G, FormatUptime exibe d/h/m, FormatBlockHeight usa separador, ≥90% branch coverage

---

### Tarefa 07 — internal/config (depende de: nenhum interno)
**Status:** done ✅
**Lê:** `_reversa_sdd/internal/config/requirements.md`, `_reversa_sdd/internal/config/design.md`, `_reversa_sdd/internal/config/tasks.md`
**Constrói:** `internal/config/config.go` (`Config` struct, `Load`, `Validate`), `internal/config/paths.go` (`ExpandPath`)
**Pronto quando:** Viper carrega defaults e env vars, ExpandPath resolve `~/`, Validate rejeita CPUTarget fora de [0.05,1.0] e BTCAddress vazio sem mock, ≥80% branch coverage

---

### Tarefa 08 — internal/store (depende de: internal/config para paths)
**Status:** done ✅
**Lê:** `_reversa_sdd/internal/store/requirements.md`, `_reversa_sdd/internal/store/design.md`, `_reversa_sdd/internal/store/tasks.md`
**Constrói:** `internal/store/store.go` (`Store` interface, `SQLiteStore`, `NilStore`), `internal/store/migrations.go` (já criado na Tarefa 02, integrar)
**Pronto quando:** AppendHashRate+QueryHashRateHistory round-trip funciona em `:memory:`, NilStore não panics, WAL mode ativo, concurrent append sem data race com `-race`, ≥75% branch coverage

---

### Tarefa 09 — internal/worker/messages (depende de: internal/model, pkg/mining)
**Status:** done ✅
**Lê:** `_reversa_sdd/internal/worker/requirements.md`, `_reversa_sdd/internal/worker/design.md`, `_reversa_sdd/internal/worker/tasks.md`
**Constrói:** `internal/worker/messages.go` — definição de todos os `tea.Msg` types (`HashRateMsg`, `ShareFoundMsg`, `PoolStatsMsg`, `MinerErrorMsg`, `PoolErrorMsg`)
**Pronto quando:** Todos os tipos compilam com campos tipados, sem `interface{}`/`any`

---

### Tarefa 10 — internal/worker (depende de: messages, pkg/mining, internal/store)
**Status:** done ✅
**Lê:** `_reversa_sdd/internal/worker/requirements.md`, `_reversa_sdd/internal/worker/design.md`, `_reversa_sdd/internal/worker/tasks.md`
**Constrói:** `internal/worker/miner.go` (`MinerWorker`), `internal/worker/fetcher.go` (`PoolClient` interface, `HTTPPoolClient`, `StratumPoolClient`, `MockPoolClient` com `SubmitShare`), `internal/worker/poller.go` (`PollCmd`)
**Pronto quando:** MinerWorker throttle math correto a 25/50/75/100%, shares encontradas são submetidas via SubmitShare, goroutine leak zerado via goleak, ≥75% branch coverage

---

### Tarefa 11 — internal/ui (depende de: internal/model, internal/worker, pkg/format)
**Status:** done ✅
**Lê:** `_reversa_sdd/internal/ui/requirements.md`, `_reversa_sdd/internal/ui/design.md`, `_reversa_sdd/internal/ui/tasks.md`
**Constrói:** `internal/ui/app.go` (`AppModel` — Init/Update/View), `internal/ui/keys.go`, `internal/ui/screens/dashboard.go`, `internal/ui/screens/clock.go`, `internal/ui/screens/globalstats.go`, `internal/ui/components/hashgauge.go`, `internal/ui/components/sparkline.go`, `internal/ui/components/cpubar.go`, `internal/ui/components/statusbar.go`
**Pronto quando:** Tab rotaciona telas ciclicamente, +/- ajustam CPUTarget com clamp, q/ctrl+c emite tea.Quit, telas respondem a WindowSizeMsg, screens são funções puras, ≥60% branch coverage no app + ≥75% nas screens

---

### Tarefa 12 — cmd/tui (main) (depende de: todos os pacotes anteriores)
**Status:** done ✅
**Lê:** `_reversa_sdd/cmd/tui/requirements.md`, `_reversa_sdd/cmd/tui/design.md`, `_reversa_sdd/cmd/tui/tasks.md`
**Constrói:** `cmd/tui/main.go` — thin wiring: parse flags → config.Load → store.New → MinerWorker + channels → AppModel → tea.NewProgram
**Pronto quando:** Binary compila com `CGO_ENABLED=0`, inicia com `--mock`, exibe dashboard, responde a `q` para sair

---

### Tarefa 13 — Fluxos de Usuário e Smoke Test E2E
**Status:** done ✅
**Lê:** `_reversa_sdd/user-stories/mining-flow.md`
**Constrói:** `testutil/fixtures.go`, teste E2E `TestMain_StartsAndQuits_NoGoroutineLeak`, validação final com `make ci-pr`
**Pronto quando:** Programa inicia com `--mock`, recebe `q`, termina limpo em <5s, sem goroutine leaks, `make ci-pr` passa integralmente
