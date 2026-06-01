# NerdTUI — Software Design & Test Specification

> Versão 1.0 | Go 1.23 | Bubbletea + Lipgloss + Bubbles

---

## Índice

1. [Visão do Produto](#1-visão-do-produto)
2. [Restrições e Princípios de Design](#2-restrições-e-princípios-de-design)
3. [Arquitetura de Pacotes](#3-arquitetura-de-pacotes)
4. [Modelo de Domínio](#4-modelo-de-domínio)
5. [Especificação de Componentes](#5-especificação-de-componentes)
   - 5.1 [cmd/tui](#51-cmdtui)
   - 5.2 [internal/config](#52-internalconfig)
   - 5.3 [internal/model](#53-internalmodel)
   - 5.4 [internal/worker/fetcher](#54-internalworkerfetcher)
   - 5.5 [internal/worker/miner](#55-internalworkerminer)
   - 5.6 [internal/ui](#56-internalui)
   - 5.7 [internal/store](#57-internalstore)
   - 5.8 [pkg/mining](#58-pkgmining)
   - 5.9 [pkg/format](#59-pkgformat)
6. [Contratos de Mensagens (tea.Msg)](#6-contratos-de-mensagens-teamsg)
7. [Contratos de Interface](#7-contratos-de-interface)
8. [Especificação de Telas](#8-especificação-de-telas)
9. [Estratégia de Throttle](#9-estratégia-de-throttle)
10. [Estratégia de Testes (TDD)](#10-estratégia-de-testes-tdd)
    - 10.1 [Pirâmide e Cobertura-Alvo](#101-pirâmide-e-cobertura-alvo)
    - 10.2 [Catálogo de Testes por Componente](#102-catálogo-de-testes-por-componente)
    - 10.3 [Fixtures e Builders](#103-fixtures-e-builders)
    - 10.4 [Contratos de Mock](#104-contratos-de-mock)
11. [Quality Gates](#11-quality-gates)
    - 11.1 [Visão Geral e Pipeline de CI](#111-visão-geral-e-pipeline-de-ci)
    - 11.2 [Gates por Estágio](#112-gates-por-estágio)
    - 11.3 [Thresholds por Pacote](#113-thresholds-por-pacote)
    - 11.4 [Linters e Regras Estáticas](#114-linters-e-regras-estáticas)
    - 11.5 [Gestão de Flaky Tests](#115-gestão-de-flaky-tests)
    - 11.6 [Automação de Cobertura por Pacote](#116-automação-de-cobertura-por-pacote)
12. [AGENTS.md embutido](#12-agentsmd-embutido)
13. [Makefile Spec](#13-makefile-spec)

---

## 1. Visão do Produto

**NerdTUI** é um dashboard de mineração Bitcoin solo para terminal, inspirado no NerdMiner_v2 (ESP32). Roda em qualquer OS Unix, exibe três telas rotativas (mining stats, clock, global stats), e possui um loop de hashing SHA256d com controle de CPU em tempo real.

**Não-objetivos v1:**
- Submissão real de shares válidas a pools mainnet
- Múltiplos workers paralelos (adicionável em v2)
- GUI ou modo web

---

## 2. Restrições e Princípios de Design

| # | Restrição | Justificativa |
|---|-----------|---------------|
| R1 | `CGO_ENABLED=0` obrigatório | Binary estático, sem `gcc` em CI |
| R2 | Zero `interface{}` / `any` em domínio | Tipos explícitos → agente não precisa inferir |
| R3 | Arquivos ≤ 300 linhas, funções ≤ 30 linhas | Token budget de agente |
| R4 | Estado da UI imutável — copiado por valor | Elimina data races no Model |
| R5 | Workers comunicam exclusivamente via `tea.Cmd`/`tea.Msg` | Sem goroutines escrevendo estado diretamente |
| R6 | Nenhum `panic` em código de produção | Erros como valores, sempre |
| R7 | `context.Context` como primeiro parâmetro em toda função de I/O | Cancelamento propagado |
| R8 | `modernc.org/sqlite` para store (pure Go) | CGO=0 compatível |

---

## 3. Arquitetura de Pacotes

```
nerdtui/
├── AGENTS.md
├── Makefile
├── go.mod                          (module github.com/user/nerdtui)
│
├── cmd/tui/
│   └── main.go                     # thin: parse flags → wire → tea.NewProgram().Run()
│
├── internal/
│   ├── config/
│   │   └── config.go               # Config struct + Load() via viper
│   │
│   ├── model/
│   │   └── state.go                # AppState (valor imutável), constantes de domínio
│   │
│   ├── worker/
│   │   ├── fetcher.go              # FetchPoolStats() → tea.Cmd
│   │   ├── poller.go               # PollCmd() — ticker + retry exponencial
│   │   ├── miner.go                # MinerWorker struct + Run() goroutine
│   │   └── messages.go             # todos os tea.Msg deste pacote
│   │
│   ├── ui/
│   │   ├── app.go                  # AppModel: Model/Update/View raiz
│   │   ├── keys.go                 # key.Map com keybindings nomeados
│   │   ├── screens/
│   │   │   ├── dashboard.go        # RenderDashboard(AppState, w, h) string
│   │   │   ├── clock.go            # RenderClock(AppState, w, h) string
│   │   │   └── globalstats.go      # RenderGlobalStats(AppState, w, h) string
│   │   └── components/
│   │       ├── hashgauge.go        # RenderHashGauge(hps float64, w int) string
│   │       ├── sparkline.go        # RenderSparkline([]float64, w int) string
│   │       ├── cpubar.go           # RenderCPUBar(target float64, w int) string
│   │       └── statusbar.go        # RenderStatusBar(AppState, w int) string
│   │
│   └── store/
│       ├── store.go                # Store interface + SQLiteStore impl
│       └── migrations.go           # SQL de criação de tabelas (embed)
│
└── pkg/
    ├── mining/
    │   ├── hash.go                 # SHA256d(header, nonce) [32]byte
    │   ├── target.go               # MeetsTarget(), DifficultyFromHash()
    │   └── job.go                  # Job struct (header, target, extranonce)
    └── format/
        ├── hashrate.go             # FormatHashRate(hps float64) string
        ├── duration.go             # FormatUptime(d time.Duration) string
        └── difficulty.go           # FormatDifficulty(d float64) string
```

**Regra de importação:** `pkg/` não importa `internal/`. `internal/` não importa `cmd/`. Ciclos são erro de compilação — não há exceção.

---

## 4. Modelo de Domínio

### AppState

Struct de valor (sem ponteiros internos mutáveis). Passada por cópia entre Update() e View().

| Campo | Tipo | Descrição |
|-------|------|-----------|
| `HashRate` | `float64` | Hashes por segundo — medido no último intervalo de 1s |
| `HashRateHistory` | `[60]float64` | Últimos 60 snapshots — sparkline |
| `SharesFound` | `uint64` | Shares que passaram no target local |
| `BestDifficulty` | `float64` | Maior dificuldade já encontrada |
| `BlockHeight` | `uint32` | Altura do bloco atual (vindo da pool) |
| `CPUTarget` | `float64` | Fração 0.05–1.0 desejada de CPU |
| `CPUActual` | `float64` | Fração real medida no último intervalo |
| `PoolConnected` | `bool` | Estado da conexão TCP Stratum |
| `PoolURL` | `string` | Pool configurada |
| `Uptime` | `time.Duration` | Tempo desde o start |
| `StartedAt` | `time.Time` | Timestamp de início |
| `Screen` | `ScreenID` | Tela ativa |
| `Error` | `string` | Último erro não-fatal (exibido na statusbar) |

### Job

```
Job { Header []byte, Target [32]byte, ExtraNonce uint32, Height uint32 }
```

Gerado pelo fetcher (pool real) ou pelo mock job generator (modo loteria local).

---

## 5. Especificação de Componentes

### 5.1 `cmd/tui`

**Responsabilidade:** wiring exclusivo. Sem lógica de negócio.

**Sequência de inicialização:**
1. Parse flags (`--config`, `--no-mine`, `--cpu`)
2. `config.Load()` → `Config`
3. Instanciar `store.New(cfg.StorePath)` → `Store`
4. Instanciar `MinerWorker` com channels
5. Construir `AppModel` injetando deps
6. `tea.NewProgram(model, tea.WithAltscreen()).Run()`

**Invariante:** `main()` nunca retorna erro — usa `log.Fatal` apenas para falha de inicialização irrecuperável.

---

### 5.2 `internal/config`

**Responsabilidade:** carregar e validar configuração. Única fonte de verdade para parâmetros externos.

| Campo | Tipo | Default | Fonte |
|-------|------|---------|-------|
| `PoolURL` | `string` | `"public-pool.io"` | env `NM_POOL_URL` / yaml |
| `PoolPort` | `int` | `21496` | env `NM_POOL_PORT` |
| `BTCAddress` | `string` | `""` | env `NM_BTC_ADDRESS` |
| `PollInterval` | `time.Duration` | `5s` | env `NM_POLL_INTERVAL` |
| `CPUTarget` | `float64` | `0.5` | env `NM_CPU_TARGET` |
| `StorePath` | `string` | `"~/.nerdtui/metrics.db"` | env `NM_STORE_PATH` |
| `Theme` | `string` | `"dark"` | env `NM_THEME` |
| `MockMining` | `bool` | `false` | flag `--mock` |

**Validação:** `Config.Validate() error` — falha se `BTCAddress` vazio e `MockMining=false`; falha se `CPUTarget` fora de [0.05, 1.0].

---

### 5.3 `internal/model`

**Responsabilidade:** definir `AppState` e constantes de domínio. Zero lógica de negócio.

```
ScreenID  type = uint8  { ScreenDashboard=0, ScreenClock=1, ScreenGlobalStats=2 }
NumScreens = 3
MinCPUTarget = 0.05
MaxCPUTarget = 1.00
CPUStep      = 0.05
HashHistoryLen = 60
```

Método único permitido: `AppState.WithHashRate(hps float64) AppState` — retorna cópia com `HashRate` e `HashRateHistory` atualizados. Padrão funcional: não modifica o receiver.

---

### 5.4 `internal/worker/fetcher`

**Responsabilidade:** buscar estado da pool via HTTP (stats REST) ou Stratum TCP (para job real).

**Interface consumida:**

```go
type PoolClient interface {
    FetchStats(ctx context.Context) (PoolStats, error)
    FetchJob(ctx context.Context) (mining.Job, error)
}
```

**Implementações:**
- `HTTPPoolClient` — REST polling em `/api/stats` (public-pool.io pattern)
- `StratumPoolClient` — TCP Stratum JSON-RPC (v1)
- `MockPoolClient` — retorna dados fixos; usado em testes e `--mock`

**Comportamento de erro:** retorna `ErrPoolUnreachable` (sentinela). Nunca panics.

---

### 5.5 `internal/worker/miner`

**Responsabilidade:** loop de hashing SHA256d com controle de CPU. Comunica via channels.

**Estrutura:**

```
MinerWorker {
    throttleCh chan float64   // recebe novo CPUTarget
    outCh      chan tea.Msg   // emite HashRateMsg, ShareFoundMsg, MinerErrorMsg
    job        atomic.Value   // armazena mining.Job atual (atualizado pelo fetcher)
}
```

**Invariantes do loop:**
- Nunca escreve em nenhuma variável fora de seu goroutine sem usar `atomic` ou channel
- Mede `workDuration` a cada `batchSize` iterações; dorme `workDuration × (1-P)/P`
- Emite `HashRateMsg` a cada 1s via ticker interno separado
- Responde a `throttleCh` sem interromper o batch (lê no `select` apenas entre batches)
- Termina limpo ao `ctx.Done()`

**Constantes configuráveis via `Config`:**

| Constante | Default | Significado |
|-----------|---------|-------------|
| `BatchSize` | `50_000` | Hashes por ciclo antes de checar throttle |
| `MetricsInterval` | `1s` | Frequência de emissão de `HashRateMsg` |

---

### 5.6 `internal/ui`

#### `AppModel` (Bubbletea Model)

**State machine de telas:**

```
ScreenDashboard ──tab──► ScreenClock ──tab──► ScreenGlobalStats ──tab──► ScreenDashboard
```

**Update dispatch:**

| Msg recebida | Ação em Update() |
|-------------|-----------------|
| `worker.HashRateMsg` | `state = state.WithHashRate(msg.HPS)` |
| `worker.ShareFoundMsg` | incrementa `SharesFound`, atualiza `BestDifficulty` |
| `worker.PoolStatsMsg` | atualiza `BlockHeight`, `PoolConnected` |
| `worker.MinerErrorMsg` | seta `state.Error` |
| `tea.KeyMsg` `tab` / `→` | avança tela |
| `tea.KeyMsg` `+` / `=` | `CPUTarget += CPUStep`, envia em `throttleCh` |
| `tea.KeyMsg` `-` / `_` | `CPUTarget -= CPUStep`, envia em `throttleCh` |
| `tea.KeyMsg` `q` / `ctrl+c` | `tea.Quit` |
| `tea.WindowSizeMsg` | atualiza `width`, `height` |
| `tickMsg` (1s) | incrementa `Uptime`, retorna próximo tick cmd |

**View:** delega para `screens.Render*(state, width, height)`. Compõe com `components.RenderStatusBar(state, width)` fixo no rodapé.

#### `screens/` — funções puras

Cada `Render*(AppState, w, h int) string` é **função pura** — mesmo input, mesmo output. Zero side effects. Testável com simples string comparison.

#### `components/` — funções puras

Mesmo contrato. Cada componente recebe apenas o que precisa — não o `AppState` inteiro quando só precisa de um campo.

---

### 5.7 `internal/store`

**Responsabilidade:** persistir histórico de hashrate para sparklines sobreviverem a restarts.

**Interface:**

```go
type Store interface {
    AppendHashRate(ctx context.Context, hps float64, at time.Time) error
    QueryHashRateHistory(ctx context.Context, limit int) ([]float64, error)
    Close() error
}
```

**Schema (SQLite):**

```sql
CREATE TABLE hashrate_history (
    id         INTEGER PRIMARY KEY,
    hps        REAL    NOT NULL,
    recorded_at INTEGER NOT NULL  -- unix timestamp
);
CREATE INDEX idx_hashrate_recorded ON hashrate_history(recorded_at DESC);
```

**Implementação `NilStore`** — no-op, para `--no-store`. Mesma interface, métodos retornam `nil`.

---

### 5.8 `pkg/mining`

**Responsabilidade:** primitivas criptográficas e de job. Zero dependências externas além de `crypto/sha256`.

| Função | Contrato |
|--------|----------|
| `SHA256d(data []byte) [32]byte` | Duplo SHA256 |
| `HashHeader(header []byte, nonce uint32) [32]byte` | Monta buffer e chama `SHA256d` |
| `MeetsTarget(hash, target [32]byte) bool` | `hash < target` byte-wise big-endian |
| `DifficultyFromHash(hash [32]byte) float64` | `diff = target_genesis / hash_as_bigint` |

**Nenhuma função neste pacote tem side effects.**

---

### 5.9 `pkg/format`

| Função | Exemplo de output |
|--------|-------------------|
| `FormatHashRate(hps float64) string` | `"12.4 KH/s"`, `"1.2 MH/s"` |
| `FormatUptime(d time.Duration) string` | `"2d 03h 14m"` |
| `FormatDifficulty(d float64) string` | `"p2pool: 1.23e+12"` |
| `FormatBlockHeight(h uint32) string` | `"#892.441"` com separador de milhar |

---

## 6. Contratos de Mensagens (tea.Msg)

Todas as msgs vivem em `internal/worker/messages.go`. Tipos nomeados, nunca `interface{}`.

| Msg | Campos | Emissor |
|-----|--------|---------|
| `HashRateMsg` | `HPS float64`, `CPUActual float64` | `MinerWorker` (1s ticker) |
| `ShareFoundMsg` | `Nonce uint32`, `Hash [32]byte`, `Difficulty float64` | `MinerWorker` |
| `PoolStatsMsg` | `BlockHeight uint32`, `Connected bool`, `NetworkHashRate float64` | `fetcher` |
| `MinerErrorMsg` | `Err error` | `MinerWorker` |
| `PoolErrorMsg` | `Err error` | `fetcher` |
| `tickMsg` | `At time.Time` | `AppModel` (interno) |

---

## 7. Contratos de Interface

Interfaces com **um único propósito**. Seguem io.Reader como modelo — pequenas o suficiente para ser satisfeitas por mocks triviais.

```
PoolClient        FetchStats(ctx) (PoolStats, error)
                  FetchJob(ctx) (mining.Job, error)

Store             AppendHashRate(ctx, hps, at) error
                  QueryHashRateHistory(ctx, limit) ([]float64, error)
                  Close() error

HashRenderer      RenderSparkline(data []float64, width int) string

ThrottleWriter    SetCPUTarget(target float64)
```

Todas as interfaces são definidas **no pacote que as consome**, não no pacote que as implementa (Go idiomático).

---

## 8. Especificação de Telas

### Tela 1 — Dashboard (ScreenDashboard)

```
┌─────────────────────────────────────────┐
│  ⛏ NerdTUI              #892.441        │  ← statusbar topo
├─────────────────────────────────────────┤
│                                         │
│  Hashrate:  12.4 KH/s                   │
│  ▄▂▃▅▆▄▃▂▁▂▃▄▆▅▄▃▂▃▄▅▄▅▆▇▆▅▄▃▂  (60s)  │  ← sparkline
│                                         │
│  CPU: [████████░░]  CPU 80% → 40%       │  ← cpubar + ajuste +/-
│                                         │
│  Shares found:  42                      │
│  Best diff:     2.14e+08                │
│                                         │
├─────────────────────────────────────────┤
│  public-pool.io ● connected  up 2d 03h  │  ← statusbar rodapé
└─────────────────────────────────────────┘
```

### Tela 2 — Clock (ScreenClock)

Hora atual em ASCII art grande centralizada. Abaixo: data, uptime, hashrate compacto.

### Tela 3 — Global Stats (ScreenGlobalStats)

Network hashrate, estimated difficulty, last block time. Dados da pool REST.

**Requisito de layout:** todas as telas respondem a `tea.WindowSizeMsg` — sem hardcode de dimensões.

---

## 9. Estratégia de Throttle

### Modelo matemático

```
sleep = workDuration × (1 - P) / P

onde P = CPUTarget ∈ [0.05, 1.0]
```

| CPUTarget | Proporção trabalho/sleep | CPU efetivo (aprox.) |
|-----------|--------------------------|----------------------|
| 0.25 | 1:3 | ~25% |
| 0.50 | 1:1 | ~50% |
| 0.75 | 3:1 | ~75% |
| 1.00 | sem sleep | ~100% |

### Ajuste pelo usuário

- Teclas `+`/`-` alteram `CPUTarget` em passos de `0.05`
- Novo valor enviado em `throttleCh` entre batches
- UI reflete imediatamente via `AppState.CPUTarget` (otimista); `CPUActual` é confirmado pelo próximo `HashRateMsg`

### Medição de CPUActual

```
CPUActual = workDuration / (workDuration + sleepDuration)
```

Calculado no miner, enviado em cada `HashRateMsg`. UI exibe discrepância quando `|CPUActual - CPUTarget| > 0.05`.

---

## 10. Estratégia de Testes (TDD)

### 10.1 Pirâmide e Cobertura-Alvo

```
           ┌──────────────┐
           │   E2E (tea)  │  ~5% — smoke test do programa completo
         ┌─┴──────────────┴─┐
         │   Integration    │  ~20% — fetcher real vs mock server
       ┌─┴──────────────────┴─┐
       │   Unit Tests         │  ~75% — tudo que é puro: mining, format, model, screens
     ──┴────────────────────────
```

| Pacote | Tier | Cobertura-alvo | Justificativa |
|--------|------|---------------|---------------|
| `pkg/mining` | Gold | ≥ 90% branch | Primitivas criptográficas — falha silenciosa = desastre |
| `pkg/format` | Gold | ≥ 90% branch | Output visual — regressão perceptível ao usuário |
| `internal/model` | Gold | 100% | Imutabilidade do estado é invariante central |
| `internal/config` | Silver | ≥ 80% branch | Validação com paths de erro explícitos |
| `internal/worker/miner` | Silver | ≥ 75% branch | Throttle é comportamental, não só estrutural |
| `internal/ui/screens` | Silver | ≥ 75% branch | Funções puras — testáveis sem bubbletea |
| `internal/store` | Silver | ≥ 75% branch | SQLite in-memory em testes |
| `internal/ui/app` | Bronze | ≥ 60% branch | Update() testável via `tea/teatest` |

---

### 10.2 Catálogo de Testes por Componente

#### `pkg/mining` — unit

| Caso de teste | Comportamento verificado |
|--------------|--------------------------|
| `TestSHA256d_KnownVector` | Output corresponde ao vetor NIST para duplo SHA256 |
| `TestSHA256d_NilInput` | Retorna hash do empty bytes, não panic |
| `TestMeetsTarget_HashBelowTarget_ReturnsTrue` | Comparação big-endian correta |
| `TestMeetsTarget_HashAboveTarget_ReturnsFalse` | — |
| `TestMeetsTarget_HashEqualsTarget_ReturnsFalse` | Boundary: igual não é menor |
| `TestHashHeader_ChangingNonceChangeHash` | Determinismo + sensibilidade ao nonce |
| `TestDifficultyFromHash_MaxHash_ReturnsOne` | Hash = target genesis → diff = 1.0 |
| `TestDifficultyFromHash_HalfHash_ReturnsTwo` | Relação inversamente proporcional |

#### `pkg/format` — unit

| Caso de teste | Comportamento verificado |
|--------------|--------------------------|
| `TestFormatHashRate_Below1K` | `999 H/s` — sem prefixo |
| `TestFormatHashRate_Kilo` | `1.0 KH/s` — boundary exato |
| `TestFormatHashRate_Mega` | `1.5 MH/s` |
| `TestFormatHashRate_Zero` | `"0 H/s"` — sem NaN |
| `TestFormatUptime_UnderMinute` | `"0m 42s"` |
| `TestFormatUptime_OverDay` | `"1d 02h 03m"` |
| `TestFormatBlockHeight_Separators` | `"#892.441"` |

#### `internal/model` — unit

| Caso de teste | Comportamento verificado |
|--------------|--------------------------|
| `TestAppState_WithHashRate_UpdatesField` | Campo `HashRate` atualizado na cópia |
| `TestAppState_WithHashRate_DoesNotMutateOriginal` | Imutabilidade |
| `TestAppState_WithHashRate_RotatesHistory` | `[60]float64` rola FIFO corretamente |
| `TestAppState_WithHashRate_FirstEntry` | Slice não-vazio sem panic |

#### `internal/config` — unit

| Caso de teste | Comportamento verificado |
|--------------|--------------------------|
| `TestConfig_Validate_EmptyAddress_MockFalse_ReturnsError` | — |
| `TestConfig_Validate_EmptyAddress_MockTrue_OK` | Mock mode dispensa address |
| `TestConfig_Validate_CPUTargetBelowMin_ReturnsError` | < 0.05 |
| `TestConfig_Validate_CPUTargetAboveMax_ReturnsError` | > 1.0 |
| `TestConfig_Load_EnvVarOverridesDefault` | `NM_CPU_TARGET=0.3` |
| `TestConfig_Load_MissingFileUsesDefaults` | Sem arquivo de config |

#### `internal/worker/miner` — unit (clock mockado)

| Caso de teste | Comportamento verificado |
|--------------|--------------------------|
| `TestThrottleSleepDuration_50Percent` | `sleep = work × 1.0` |
| `TestThrottleSleepDuration_25Percent` | `sleep = work × 3.0` |
| `TestThrottleSleepDuration_100Percent` | `sleep = 0` |
| `TestThrottleSleepDuration_BoundaryMin` | `CPUTarget = 0.05` não produz sleep negativo |
| `TestMinerWorker_ReceivesThrottleUpdate_MidRun` | Canal atualiza target sem parar loop |
| `TestMinerWorker_CancelContext_StopsGoroutine` | Sem goroutine leak |
| `TestMinerWorker_EmitsHashRateMsgEverySecond` | Ticker interno |
| `TestMinerWorker_ShareFound_EmitsMsg` | Job com target fácil (todos-FF) |

**Padrão para testar o goroutine:** injetar `clock` fake e `batchSize` pequeno (100). Usar `goleak.VerifyNone(t)` no `TestMain`.

#### `internal/ui/screens` — unit (puras, sem tea)

| Caso de teste | Comportamento verificado |
|--------------|--------------------------|
| `TestRenderDashboard_ContainsHashRate` | String contém `FormatHashRate` output |
| `TestRenderDashboard_RespondsToWidth` | Largura 80 vs largura 40 — sem panic |
| `TestRenderDashboard_ZeroState_NoNaN` | `AppState{}` não produz `NaN` na tela |
| `TestRenderClock_ContainsCurrentHour` | Hora está presente no output |
| `TestRenderStatusBar_ConnectedShowsBullet` | `●` presente quando `PoolConnected=true` |
| `TestRenderStatusBar_DisconnectedShowsX` | `✗` quando `false` |
| `TestRenderCPUBar_FullBar` | 100% → barra cheia |
| `TestRenderCPUBar_EmptyBar` | 5% → quase vazia |

#### `internal/store` — integration (SQLite `:memory:`)

| Caso de teste | Comportamento verificado |
|--------------|--------------------------|
| `TestStore_AppendAndQuery_RoundTrip` | Inserção e leitura consistentes |
| `TestStore_QueryLimit_ReturnsMostRecent` | `LIMIT 60` retorna últimos 60, não primeiros |
| `TestStore_ConcurrentAppend_NoDataRace` | `-race` com 3 goroutines escrevendo |
| `TestStore_Close_SecondCloseNoError` | Idempotência |
| `TestNilStore_AppendDoesNotPanic` | `NilStore` implementa interface sem efeito |

#### `internal/worker/fetcher` — integration (httptest.Server)

| Caso de teste | Comportamento verificado |
|--------------|--------------------------|
| `TestHTTPPoolClient_FetchStats_ParsesJSON` | JSON real do public-pool.io parseado corretamente |
| `TestHTTPPoolClient_FetchStats_ServerError_ReturnsErrPoolUnreachable` | 500 → sentinela |
| `TestHTTPPoolClient_FetchStats_Timeout_ContextCanceled` | ctx com 1ms de deadline |
| `TestMockPoolClient_AlwaysReturnsConfiguredJob` | Testa que o mock é determinístico |

#### `internal/ui/app` — integration (teatest)

| Caso de teste | Comportamento verificado |
|--------------|--------------------------|
| `TestAppModel_TabKeyRotatesScreen` | `tab` três vezes volta ao dashboard |
| `TestAppModel_PlusKeyIncreasesCPUTarget` | `+` incrementa em `CPUStep` |
| `TestAppModel_CPUTargetClampedAtMax` | `+` acima de 1.0 permanece em 1.0 |
| `TestAppModel_MinusKeyDecreasesCPUTarget` | — |
| `TestAppModel_CPUTargetClampedAtMin` | `-` abaixo de 0.05 permanece em 0.05 |
| `TestAppModel_HashRateMsg_UpdatesState` | Msg processada, estado reflete novo HPS |
| `TestAppModel_QuitKey_StopsProgram` | `q` emite `tea.Quit` |

#### E2E — smoke test

| Caso de teste | Comportamento verificado |
|--------------|--------------------------|
| `TestMain_StartsAndQuits_NoGoroutineLeak` | Programa inicia com `--mock`, recebe `q`, termina limpo em < 5s |

---

### 10.3 Fixtures e Builders

**`testutil/fixtures.go`** — funções de construção reutilizáveis:

```
DefaultAppState() AppState              — state zerado com campos válidos
AppStateWithHashRate(hps float64) AppState
AppStateConnected() AppState
AppStateDisconnected() AppState
EasyJob() mining.Job                   — Target = [0xFF * 32] (qualquer hash passa)
HardJob() mining.Job                   — Target = [0x00 * 32] (nenhum hash passa)
DefaultConfig() config.Config
```

Nenhuma fixture chama `t.Fatal` — retornam valores, não manipulam teste diretamente.

---

### 10.4 Contratos de Mock

| Mock | Interface satisfeita | Pacote |
|------|---------------------|--------|
| `MockPoolClient` | `PoolClient` | `internal/worker` |
| `SpyThrottleWriter` | `ThrottleWriter` | `internal/ui` |
| `InMemoryStore` | `Store` | `internal/store` |
| `FakeClock` | `clock.Clock` (própria) | `testutil` |

Todos gerados manualmente (sem mockgen) — interfaces pequenas não justificam geração.

---

---

## 11. Quality Gates

### 11.1 Visão Geral e Pipeline de CI

O pipeline tem quatro estágios em sequência. Um estágio só inicia se o anterior passou integralmente — sem exceções manuais.

```
[pre-commit local]  →  [commit / PR]  →  [post-build]  →  [merge / main]
       < 30s               < 5min           < 10min           < 15min
```

**Ferramentas requeridas no ambiente CI:**
- `go 1.23+`
- `golangci-lint v1.59+`
- `govulncheck` (golang.org/x/vuln)
- `go.uber.org/goleak` (via testes, não CLI)

Nenhuma ferramenta externa paga é exigida — tudo open-source, instalável via `go install` ou binário no PATH.

---

### 11.2 Gates por Estágio

#### Estágio 1 — Pre-commit (local, < 30s)

Executado via `make pre-commit`. Desenvolvedor não commita se falhar.

| Gate | Comando | Tipo |
|------|---------|------|
| Formatação | `gofmt -l . \| grep .` → deve retornar vazio | **Bloqueador** |
| Imports organizados | `goimports -l . \| grep .` → deve retornar vazio | **Bloqueador** |
| Vet rápido | `go vet ./...` | **Bloqueador** |
| Build compila | `CGO_ENABLED=0 go build ./...` | **Bloqueador** |

#### Estágio 2 — Commit / PR (CI, < 5min)

Executado em todo push para PR. Falha bloqueia merge.

| Gate | Comando / Critério | Tipo |
|------|--------------------|------|
| Testes unitários 100% pass | `go test -count=1 ./...` | **Bloqueador** |
| Race detector zerado | `go test -race -count=1 ./...` | **Bloqueador** |
| Goroutine leaks zerados | `goleak.VerifyNone` em todo `TestMain` | **Bloqueador** |
| Lint sem warnings | `golangci-lint run ./...` | **Bloqueador** |
| Vulnerabilidades em deps | `govulncheck ./...` — zero HIGH/CRITICAL | **Bloqueador** |
| Build estático sucesso | `CGO_ENABLED=0 go build -o /dev/null ./cmd/tui` | **Bloqueador** |
| Cobertura geral | branch coverage total ≥ 75% | **Bloqueador** |
| Cobertura `pkg/mining` | branch coverage ≥ 90% | **Bloqueador** |
| Cobertura `internal/model` | branch coverage = 100% | **Bloqueador** |
| Cyclomatic complexity | nenhuma função com CC > 15 (`gocyclo`) | **Bloqueador** |
| Arquivos > 300 linhas | nenhum (`wc -l` via script) | **Aviso** |

#### Estágio 3 — Post-build (CI, < 10min)

Executado apenas em PRs para `main` ou tags.

| Gate | Critério | Tipo |
|------|----------|------|
| Testes de integração com `-race` | `go test -race -tags=integration ./...` | **Bloqueador** |
| Cobertura após integração | não regride vs estágio 2 | **Bloqueador** |
| Benchmark de hashing não regride | `go test -bench=BenchmarkSHA256d` — P50 dentro de ±15% da baseline salva em `testdata/bench_baseline.txt` | **Aviso** |
| Smoke test E2E | `TestMain_StartsAndQuits_NoGoroutineLeak` passa em < 5s | **Bloqueador** |

#### Estágio 4 — Merge / main

Gate final antes de fechar o PR.

| Gate | Critério | Tipo |
|------|----------|------|
| Todos os estágios anteriores passaram | sem bypass manual | **Bloqueador** |
| AGENTS.md atualizado | se pacotes foram adicionados ou renomeados | **Aviso** (checklist de PR) |
| Sem TODO/FIXME em código novo | `git diff main --unified=0 \| grep '+.*TODO\|+.*FIXME'` | **Aviso** |

---

### 11.3 Thresholds por Pacote

Classificados em três tiers conforme criticidade de falha silenciosa.

| Pacote | Tier | Branch Coverage Mínimo | Justificativa |
|--------|------|------------------------|---------------|
| `pkg/mining` | **Gold** | ≥ 90% | Bug criptográfico é invisível ao usuário — programa roda, nunca encontra share |
| `pkg/format` | **Gold** | ≥ 90% | Regressão visual direta — usuário vê número errado |
| `internal/model` | **Gold** | 100% | Imutabilidade é a invariante central do sistema |
| `internal/config` | **Silver** | ≥ 80% | Validação com múltiplos paths de erro explícitos |
| `internal/worker/miner` | **Silver** | ≥ 75% | Throttle é comportamental — cobertura estrutural não garante correção da fórmula |
| `internal/ui/screens` | **Silver** | ≥ 75% | Funções puras, mas layout pode silenciosamente truncar dados |
| `internal/store` | **Silver** | ≥ 75% | Paths de erro SQLite devem ser exercitados |
| `internal/ui/app` | **Bronze** | ≥ 60% | `Update()` testado via teatest; cobertura complementada pelos testes de integração |
| `internal/worker/fetcher` | **Bronze** | ≥ 60% | Cobertura complementada pelos testes de integração com `httptest` |
| `cmd/tui` | **Bronze** | ≥ 40% | Wiring puro — cobertura real vem do smoke test E2E |

**Política de regressão:** nenhum PR pode reduzir a cobertura de qualquer pacote abaixo do threshold do seu tier. Regressão é bloqueador mesmo que a cobertura total permaneça acima de 75%.

---

### 11.4 Linters e Regras Estáticas

Configuração via `.golangci.yml` na raiz. Todos os linters abaixo são **bloqueadores** — zero findings tolerados.

| Linter | O que detecta | Bloqueador/Aviso |
|--------|--------------|-----------------|
| `errcheck` | `error` retornado e ignorado sem `_` explícito | **Bloqueador** |
| `staticcheck` | bugs, APIs depreciadas, código morto | **Bloqueador** |
| `govet` | erros detectáveis pelo compilador Go | **Bloqueador** |
| `gocyclo` | CC > 15 por função | **Bloqueador** |
| `revive` | estilo Go idiomático (substitui `golint`) | **Bloqueador** |
| `gosimple` | simplificações idiomáticas disponíveis | **Bloqueador** |
| `goconst` | string literal repetida ≥ 3× sem constante | **Bloqueador** |
| `godot` | comentários de exported symbols sem ponto final | **Aviso** |
| `noctx` | chamadas HTTP sem `context.Context` | **Bloqueador** |
| `bodyclose` | `http.Response.Body` não fechado | **Bloqueador** |
| `exhaustive` | `switch` em tipo enum sem `default` e sem todos os casos | **Bloqueador** (apenas em `ScreenID`) |
| `forbidigo` | uso de `fmt.Print*` fora de `main.go` | **Bloqueador** |
| `gocognit` | cognitive complexity > 20 por função | **Aviso** |

**Exclusões permitidas:**

```yaml
# .golangci.yml — exclusões documentadas
issues:
  exclude-rules:
    - path: "_test.go"
      linters: [errcheck, goconst]   # testes podem ignorar erros em helpers
    - path: "cmd/tui/main.go"
      linters: [forbidigo]           # main.go é o único lugar permitido para fmt.Print
```

---

### 11.5 Gestão de Flaky Tests

**Definição de flaky:** teste que falha em ≥ 2 execuções consecutivas sem mudança de código.

**Política em três níveis:**

| Severidade | Critério | SLA | Ação |
|------------|---------|-----|------|
| P1 | Falha em > 20% dos runs de CI | 48h | Quarentena imediata + issue obrigatória |
| P2 | Falha em 5–20% dos runs | 1 semana | Tag `//nolint:flaky` + issue registrada |
| P3 | Falha em < 5% dos runs | Próximo sprint | Monitoramento |

**Causas conhecidas neste projeto e mitigações:**

| Causa | Onde ocorre | Mitigação |
|-------|------------|-----------|
| `time.Sleep` real em loop de throttle | `TestMinerWorker_*` | Injetar `FakeClock` — não usar `time.Sleep` real em testes |
| Port binding em `httptest.Server` | `TestHTTPPoolClient_*` | Usar `httptest.NewServer` (porta 0, sempre disponível) |
| Ordem de leitura em SQLite WAL | `TestStore_ConcurrentAppend_*` | `_busy_timeout=5000` no DSN de teste |
| Timing de goroutine em teatest | `TestAppModel_*` | Usar `teatest.WaitFor` com timeout explícito de 2s |

**Rerun policy:** máximo 1 rerun automático no CI. Segunda falha consecutiva = alerta obrigatório, nunca silencioso.

---

### 11.6 Automação de Cobertura por Pacote

Script `scripts/check-coverage.sh` — executado no estágio 2. Falha o CI se qualquer pacote estiver abaixo do threshold do seu tier.

```
Pacotes Gold   → falha se branch coverage < 90% (mining, format) ou < 100% (model)
Pacotes Silver → falha se branch coverage < 75%
Pacotes Bronze → falha se branch coverage < 40%
Relatório salvo em coverage-report.txt — artefato de CI preservado por 30 dias
```

O script lê os thresholds de `scripts/coverage-thresholds.json` — não hardcoded no shell. Isso permite que o agente atualize os thresholds sem editar lógica de shell.

---

## 12. AGENTS.md embutido

> Este bloco deve existir como `AGENTS.md` na raiz do repositório.

```markdown
# AGENTS.md — NerdTUI

## Stack
- Go 1.23, CGO_ENABLED=0
- TUI: github.com/charmbracelet/bubbletea + lipgloss + bubbles
- Config: github.com/spf13/viper
- SQLite: modernc.org/sqlite (pure Go)
- Test: testify/assert + go.uber.org/goleak

## Comandos Essenciais
- `make test`       → todos os testes com -race
- `make cover`      → cobertura por pacote + HTML
- `make check-coverage` → valida thresholds por tier
- `make lint`       → golangci-lint
- `make vuln`       → govulncheck (zero HIGH/CRITICAL)
- `make build`      → binary estático em bin/nerdtui
- `make run-mock`   → roda com --mock (sem pool real)
- `make pre-commit` → gates locais (< 30s)
- `make ci-pr`      → gates completos do estágio PR

## Onde está cada coisa
- Lógica criptográfica    → pkg/mining/
- Formatação de display   → pkg/format/
- Estado da aplicação     → internal/model/state.go
- Loop de mineração       → internal/worker/miner.go
- Telas (puras)           → internal/ui/screens/
- Mensagens entre workers → internal/worker/messages.go
- Wiring                  → cmd/tui/main.go

## Convenções
- AppState: sempre copiado por valor, nunca ponteiro
- tea.Msg: definir em messages.go do pacote emissor, nunca inline
- Funções de render: puras — input → string, sem side effects
- context.Context: primeiro parâmetro em qualquer função de I/O
- Erros: retornar como último valor; sentinelas prefixados com Err

## Proibido
- CGO (qualquer import que exija gcc)
- interface{} / any no domínio
- Goroutines escrevendo em AppState diretamente
- panic() fora de TestMain
- Arquivos > 300 linhas
- Funções > 30 linhas
- Merge sem `make ci-pr` passar localmente
- Reduzir cobertura abaixo do threshold do tier do pacote (ver §11.3)

## Caveats
- pkg/mining/hash.go usa crypto/sha256 da stdlib — não trocar por implementação custom sem benchmark comprovado
- internal/store usa WAL mode — não mudar pragma sem rever testes de concorrência
- MinerWorker.Run() DEVE ser chamado em goroutine separada — bloqueia até ctx cancelado
```

---

## 13. Makefile Spec

```makefile
.PHONY: build test cover lint run run-mock clean pre-commit ci check-coverage

BIN      = bin/nerdtui
VERSION  = $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS  = -ldflags="-s -w -X main.version=$(VERSION)"

build:
	CGO_ENABLED=0 go build $(LDFLAGS) -o $(BIN) ./cmd/tui

test:
	go test -race -count=1 ./...

cover:
	go test -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -func=coverage.out | tee coverage-report.txt
	go tool cover -html=coverage.out -o coverage.html

check-coverage:
	@bash scripts/check-coverage.sh  # lê thresholds de scripts/coverage-thresholds.json

lint:
	golangci-lint run ./...

vuln:
	govulncheck ./...

bench:
	go test -bench=BenchmarkSHA256d -benchmem -count=3 ./pkg/mining/... \
	  | tee /tmp/bench_current.txt
	@if [ -f testdata/bench_baseline.txt ]; then \
	    benchstat testdata/bench_baseline.txt /tmp/bench_current.txt; \
	fi

pre-commit:
	@gofmt -l . | grep . && echo "❌ gofmt: formatar antes de commitar" && exit 1 || true
	@go vet ./...
	@CGO_ENABLED=0 go build ./...
	@echo "✅ pre-commit OK"

run:
	go run ./cmd/tui

run-mock:
	go run ./cmd/tui --mock --cpu 0.3

clean:
	rm -rf bin/ coverage.out coverage.html coverage-report.txt

# Estágio 2 (PR) — ordem importa: lint → test+race → cover → vuln → build
ci-pr: lint test check-coverage vuln build

# Estágio 3 (post-build) — integração + bench
ci-post: ci-pr bench
	go test -race -count=1 -tags=integration ./...

# Alias usado pela action principal
ci: ci-pr
```

---

*Documento gerado para uso com Claude Code / agente. Consulte `AGENTS.md` para navegação rápida.*
