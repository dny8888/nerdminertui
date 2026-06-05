# 🔍 Análise Completa — NerdMiner TUI

> **Data:** 2026-06-04 | **Arquivos analisados:** 53+ | **Findings totais:** 64

---

## Resumo Executivo

| Área                      | 🔴 Critical | 🟠 High | 🟡 Medium | 🟢 Low  | Total  |
| ------------------------- | :--------: | :----: | :------: | :----: | :----: |
| **⛏️ Mining Engine**       |     2      |   5    |    7     |   6    | **20** |
| **🖥️ TUI / UX**            |     1      |   9    |    6     |   2    | **18** |
| **🏗️ Infra / CI / Config** |     2      |   6    |    10    |   7    | **25** |
| **Total**                 |   **5**    | **20** |  **23**  | **15** | **63** |

---

## 🏆 Top 10 Quick Wins (Impacto × Esforço)

| #   | Finding                                                                                                                                   | Área   | Esforço      | Impacto                       |
| --- | ----------------------------------------------------------------------------------------------------------------------------------------- | ------ | ------------ | ----------------------------- |
| 1   | Canal de throttle bloqueia TUI ([app.go:164](file:///mnt/bkp/workspace/github/nerdminertui/internal/ui/app.go#L164))                      | TUI    | 1 linha      | 🔴 Freeze da UI                |
| 2   | `Validate()` nunca é chamado ([main.go:35](file:///mnt/bkp/workspace/github/nerdminertui/cmd/tui/main.go#L35))                            | Config | 1 linha      | 🔴 Config inválida em produção |
| 3   | `-race` ausente no release workflow ([release.yml:19](file:///mnt/bkp/workspace/github/nerdminertui/.github/workflows/release.yml#L19))   | CI     | 1 palavra    | 🟠 Race conditions em release  |
| 4   | Erros da pool são silenciados ([app.go:263](file:///mnt/bkp/workspace/github/nerdminertui/internal/ui/app.go#L263))                       | UX     | ~20 linhas   | 🟠 Usuário voa às cegas        |
| 5   | Screen IDs são `0,1,2` em vez de constantes ([app.go:307](file:///mnt/bkp/workspace/github/nerdminertui/internal/ui/app.go#L307))         | Code   | 3 linhas     | 🟠 Quebra silenciosa           |
| 6   | Focus index usa magic numbers ([app.go:126](file:///mnt/bkp/workspace/github/nerdminertui/internal/ui/app.go#L126))                       | Code   | 7 constantes | 🟠 Manutenção impossível       |
| 7   | `MeetsTarget` reverte bytes a cada chamada ([target.go:9](file:///mnt/bkp/workspace/github/nerdminertui/pkg/mining/target.go#L9))         | Perf   | ~10 linhas   | 🔴 Hot path do miner           |
| 8   | `FormatUptime` perde horas quando dias=0 ([duration.go:23](file:///mnt/bkp/workspace/github/nerdminertui/pkg/format/duration.go#L23))     | Bug    | 2 linhas     | 🟡 Exibição incorreta          |
| 9   | CPUBar fabrica valor quando actual=0 ([cpubar.go:25](file:///mnt/bkp/workspace/github/nerdminertui/internal/ui/components/cpubar.go#L25)) | UX     | 3 linhas     | 🟡 Dados enganosos             |
| 10  | `-trimpath` ausente no release ([release.yml:59](file:///mnt/bkp/workspace/github/nerdminertui/.github/workflows/release.yml#L59))        | CI     | 1 flag       | 🟠 Vaza paths do build         |

---

## ⛏️ Mining Engine & Performance

### 🔴 CRITICAL

#### C1. `MeetsTarget` reverte endianness em cada chamada no hot loop
**Arquivo:** [target.go:9-24](file:///mnt/bkp/workspace/github/nerdminertui/pkg/mining/target.go#L9-L24)

Chamada ~10.000 vezes por batch por CPU core. A comparação `hash[31-i]` vs `target[i]` é correta mas subótima no caminho mais quente do minerador. O target é **constante para todo o job** — deveria ser pré-revertido uma única vez.

```diff
 // Proposta: Pré-reverter target uma vez, comparar diretamente
+// TargetLE armazenado no Job já em little-endian
 for i := 0; i < 32; i++ {
-    if hash[31-i] < target[i] { return true }
-    if hash[31-i] > target[i] { return false }
+    if hash[i] < targetLE[i] { return true }
+    if hash[i] > targetLE[i] { return false }
 }
```

**Otimização extra:** Para dificuldade padrão de pools, bytes 0-3 do target são `0x00`. Uma comparação `uint32` dos primeiros 4 bytes rejeitaria 99.99%+ dos hashes em **1 instrução de máquina**.

#### C2. `DifficultyFromHash` aloca `big.Int`/`big.Float` pesados
**Arquivo:** [target.go:43-86](file:///mnt/bkp/workspace/github/nerdminertui/pkg/mining/target.go#L43-L86)

4+ alocações heap por chamada. `TargetFromDifficulty` é chamado 1x por job (aceitável), mas `DifficultyFromHash` pode ser chamado em shares encontrados. Considerar `sync.Pool` ou math inteiro para dificuldades comuns.

---

### 🟠 HIGH

#### H1. `math/rand` sem seed explícita e contention entre goroutines
**Arquivos:** [miner.go:57,109](file:///mnt/bkp/workspace/github/nerdminertui/internal/worker/miner.go#L57), [astronomy.go:28](file:///mnt/bkp/workspace/github/nerdminertui/pkg/trivia/astronomy.go#L28)

`rand.Uint32()` é chamado de múltiplas goroutines concorrentes. Em Go 1.20+ o source global é safe, mas compete pelo lock global. Com Go 1.26 deveria usar `math/rand/v2` e instâncias `*rand.Rand` por worker.

#### H2. Erros de `hex.DecodeString` silenciosamente descartados no parser Stratum
**Arquivo:** [stratum_parser.go:39-42](file:///mnt/bkp/workspace/github/nerdminertui/pkg/mining/stratum_parser.go#L39-L42)

```go
versionBytes, _ := hex.DecodeString(versionHex)   // ← erro ignorado
prevhashBytes, _ := hex.DecodeString(prevhashHex)  // ← erro ignorado
```

Se a pool enviar hex malformado, o header resultante será lixo — hashes desperdiçados e potencial panic em `reverseBytes` com slice nil/curto.

#### H3. `MempoolClient.FetchStats` — `defer resp.Body.Close()` dentro de `if`
**Arquivo:** [fetcher.go:62-89](file:///mnt/bkp/workspace/github/nerdminertui/internal/worker/fetcher.go#L62-L89)

Dois requests HTTP sequenciais usam `defer resp.Body.Close()` em blocos `if`. O `defer` executa quando `FetchStats` retorna, não quando o `if` sai. O primeiro body fica aberto durante todo o segundo request.

#### H4. Goroutine de `client.reconnect` sem cancellation via context
**Arquivo:** [fetcher.go:376-392](file:///mnt/bkp/workspace/github/nerdminertui/internal/worker/fetcher.go#L376-L392)

```go
go func() {
    time.Sleep(time.Duration(waitTime) * time.Second) // ← ignora ctx
    // ... modifica c.Address, c.Port, fecha conn
}()
```

Se o app encerrar durante o sleep, esta goroutine vaza e tenta modificar uma conexão já limpa.

#### H5. `sendAndWait` incrementa `reqID` manualmente em vez de usar `nextID()`
**Arquivo:** [fetcher.go:270-271](file:///mnt/bkp/workspace/github/nerdminertui/internal/worker/fetcher.go#L270-L271)

Inconsistente com `send()` que usa `nextID()`. Se `nextID()` ganhar lógica extra, `sendAndWait` não será beneficiado.

---

### 🟡 MEDIUM

| #   | Finding                                                                  | Arquivo                                                                                                    | Impacto                              |
| --- | ------------------------------------------------------------------------ | ---------------------------------------------------------------------------------------------------------- | ------------------------------------ |
| M1  | Segundo round SHA256 aloca via `sha256.Sum256()` (primeiro é zero-alloc) | [hash.go:61](file:///mnt/bkp/workspace/github/nerdminertui/pkg/mining/hash.go#L61)                         | Perf: ~50% das alocações no hot path |
| M2  | `MarshalBinary` error ignorado — panic se falhar                         | [hash.go:37](file:///mnt/bkp/workspace/github/nerdminertui/pkg/mining/hash.go#L37)                         | Estabilidade                         |
| M3  | Merkle branch aloca `make([]byte, 64)` por branch (2 locais duplicados)  | [stratum_parser.go:32,110](file:///mnt/bkp/workspace/github/nerdminertui/pkg/mining/stratum_parser.go#L32) | Perf + DRY violation                 |
| M4  | HPS frágil: assume ticker = 1s, não divide por tempo real                | [miner.go:170-174](file:///mnt/bkp/workspace/github/nerdminertui/internal/worker/miner.go#L170)            | Correção do cálculo                  |
| M5  | Sem backpressure no `outCh` para `HashRateMsg`                           | [miner.go:180-183](file:///mnt/bkp/workspace/github/nerdminertui/internal/worker/miner.go#L180)            | Miner control loop trava             |
| M6  | `RebuildHeaderWithExtraNonce2` duplica `ParseStratumJob` (~30 linhas)    | [stratum_parser.go:92-125](file:///mnt/bkp/workspace/github/nerdminertui/pkg/mining/stratum_parser.go#L92) | Manutenção                           |
| M7  | `FormatUptime` perde horas quando dias=0 (bug)                           | [duration.go:23](file:///mnt/bkp/workspace/github/nerdminertui/pkg/format/duration.go#L23)                 | UI mostra `135m` em vez de `2h 15m`  |

### 🟢 LOW

| #   | Finding                                                                     | Arquivo                                                                                                 |
| --- | --------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------- |
| L1  | `MinerHashState` não tem método `Reset()` — cria garbage a cada job         | [hash.go:15-44](file:///mnt/bkp/workspace/github/nerdminertui/pkg/mining/hash.go#L15)                   |
| L2  | `SubscribeResult`/`NotifyParams` types declarados mas nunca usados          | [stratum.go:27-31](file:///mnt/bkp/workspace/github/nerdminertui/internal/worker/stratum.go#L27)        |
| L3  | Benchmark usa mesmo nonce em toda iteração (resultados irreais)             | [hash_test.go:33,47](file:///mnt/bkp/workspace/github/nerdminertui/pkg/mining/hash_test.go#L33)         |
| L4  | `SpaceWords` pode colidir entre workers (24 palavras, 8 CPUs = ~74% chance) | [astronomy.go:28](file:///mnt/bkp/workspace/github/nerdminertui/pkg/trivia/astronomy.go#L28)            |
| L5  | `MempoolClient` cria novo `http.Client` por chamada (sem reuso TCP)         | [fetcher.go:60](file:///mnt/bkp/workspace/github/nerdminertui/internal/worker/fetcher.go#L60)           |
| L6  | Teste `TestClientsStub` faz HTTP real para mempool.space                    | [worker_test.go:123](file:///mnt/bkp/workspace/github/nerdminertui/internal/worker/worker_test.go#L123) |

---

## 🖥️ TUI / Interface / UX

### 🔴 CRITICAL

#### T1. Canal de throttle bloqueia Update() — TUI congela
**Arquivo:** [app.go:164-172](file:///mnt/bkp/workspace/github/nerdminertui/internal/ui/app.go#L164)

```go
case "+", "=":
    m.state = m.state.WithCPUTarget(0.05)
    if m.throttleCh != nil {
        m.throttleCh <- m.state.CPUTarget // ← BLOQUEIA se buffer cheio
    }
```

Send bloqueante num `chan float64` (buffer 10) dentro de `Update()`. Se o consumer do miner estiver lento, o Bubbletea runtime inteiro congela.

**Fix:** `select { case m.throttleCh <- val: default: }` (non-blocking) ou usar `tea.Cmd`.

---

### 🟠 HIGH (UX)

#### T2. Erros da pool e do miner silenciosamente descartados
**Arquivo:** [app.go:263-266](file:///mnt/bkp/workspace/github/nerdminertui/internal/ui/app.go#L263)

```go
case worker.MinerErrorMsg:
    // Ignoring for now or handle appropriately
case worker.PoolErrorMsg:
    // Ignoring or handle
```

Se a conexão falhar, shares forem rejeitados, ou o parser der erro, o usuário **não vê absolutamente nada**. Precisa de uma área de notificação/toast.

#### T3. Sem feedback visual ao salvar configuração (Ctrl+S)
**Arquivo:** [app.go:96-114](file:///mnt/bkp/workspace/github/nerdminertui/internal/ui/app.go#L96)

Config é salva em background goroutine, resultado nunca é comunicado de volta. Usuário não sabe se salvou ou não.

#### T4. Navegação de telas é unidirecional — só `tab` avança, sem `shift+tab`
**Arquivo:** [app.go:161-162](file:///mnt/bkp/workspace/github/nerdminertui/internal/ui/app.go#L161), [state.go:73-76](file:///mnt/bkp/workspace/github/nerdminertui/internal/model/state.go#L73)

Usuário que ultrapassar a tela desejada precisa ciclar por todas as restantes.

---

### 🟠 HIGH (Code Quality da TUI)

#### T5. Focus index usa magic numbers espalhados por todo app.go
**Arquivo:** [app.go:126-159](file:///mnt/bkp/workspace/github/nerdminertui/internal/ui/app.go#L126)

```go
if m.settings.FocusIndex > 6 { m.settings.FocusIndex = 0 }
if i < 4 && m.settings.FocusIndex == i { isFocused = true }
else if i == 4 && m.settings.FocusIndex == 6 { isFocused = true }
```

Mapeamento 0-6 sem constantes nomeadas é pesadelo de manutenção.

#### T6. `View()` usa `0, 1, 2` em vez de `model.ScreenDashboard` etc.
**Arquivo:** [app.go:307-314](file:///mnt/bkp/workspace/github/nerdminertui/internal/ui/app.go#L307)

Derrota o propósito do type system.

#### T7. Lógica de uptime duplicada entre dashboard e statusbar
**Arquivos:** [dashboard.go:25-35](file:///mnt/bkp/workspace/github/nerdminertui/internal/ui/screens/dashboard.go#L25), [statusbar.go:51-63](file:///mnt/bkp/workspace/github/nerdminertui/internal/ui/components/statusbar.go#L51)

Mesma lógica copy-paste. Deveria usar `pkg/format/FormatUptime`.

#### T8. Formatação de milhares duplicada 2x no mesmo arquivo
**Arquivo:** [globalstats.go:70-82,93-103](file:///mnt/bkp/workspace/github/nerdminertui/internal/ui/screens/globalstats.go#L70)

---

### 🟡 MEDIUM

| #   | Finding                                                                             | Arquivo                                                                                                                                                                                                                                                                                                                                                                                                      |
| --- | ----------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| T9  | Larguras de cards hardcoded (24, 76, 66 cols) — overflow em terminais < 84 cols     | [globalstats.go:22,61](file:///mnt/bkp/workspace/github/nerdminertui/internal/ui/screens/globalstats.go#L22), [settings.go:107](file:///mnt/bkp/workspace/github/nerdminertui/internal/ui/screens/settings.go#L107)                                                                                                                                                                                          |
| T10 | Strings em português hardcoded sem i18n (`"Conectado"`, `"Loteria"`, `"odds hoje"`) | Múltiplos arquivos                                                                                                                                                                                                                                                                                                                                                                                           |
| T11 | CPUBar fabrica valor `actual = target - 0.01` quando actual=0                       | [cpubar.go:25-27](file:///mnt/bkp/workspace/github/nerdminertui/internal/ui/components/cpubar.go#L25)                                                                                                                                                                                                                                                                                                        |
| T12 | Header não preenche largura do terminal — linha de `8` caracteres fixa              | [header.go:29-35](file:///mnt/bkp/workspace/github/nerdminertui/internal/ui/components/header.go#L29)                                                                                                                                                                                                                                                                                                        |
| T13 | Estilos Lipgloss instanciados dentro de `View()` a cada render                      | [app.go:143](file:///mnt/bkp/workspace/github/nerdminertui/internal/ui/app.go#L143), [globalstats.go:52](file:///mnt/bkp/workspace/github/nerdminertui/internal/ui/screens/globalstats.go#L52), [statusbar.go:38](file:///mnt/bkp/workspace/github/nerdminertui/internal/ui/components/statusbar.go#L38), [header.go:17](file:///mnt/bkp/workspace/github/nerdminertui/internal/ui/components/header.go#L17) |
| T14 | O(n) array shift + full scan em HashRateHistory a cada segundo                      | [state.go:60-68](file:///mnt/bkp/workspace/github/nerdminertui/internal/model/state.go#L60), [app.go:206-218](file:///mnt/bkp/workspace/github/nerdminertui/internal/ui/app.go#L206)                                                                                                                                                                                                                         |

### 🟢 LOW

| #   | Finding                                                        | Arquivo                                                                                            |
| --- | -------------------------------------------------------------- | -------------------------------------------------------------------------------------------------- |
| T15 | Block height mostra `#0` antes dos dados do mempool carregarem | [header.go:14](file:///mnt/bkp/workspace/github/nerdminertui/internal/ui/components/header.go#L14) |
| T16 | Doc comment menciona tela "Clock" que não existe               | [screens/doc.go:1-2](file:///mnt/bkp/workspace/github/nerdminertui/internal/ui/screens/doc.go#L1)  |

---

### 🚫 Features Ausentes (TUI)

| #       | Feature                                                                | Impacto                       |
| ------- | ---------------------------------------------------------------------- | ----------------------------- |
| **MF1** | Área de notificação/toast para erros e eventos                         | 🟠 HIGH — erros são invisíveis |
| **MF2** | Tela de ajuda / referência de keybindings (`?`)                        | 🟠 HIGH — discoverability      |
| **MF3** | Feedback de status de reconexão (retries, backoff)                     | 🟡 MEDIUM                      |
| **MF4** | Persistência de estatísticas (store SQLite inicializado mas não usado) | 🟡 MEDIUM                      |
| **MF5** | Validação de input nos campos de settings                              | 🟡 MEDIUM                      |
| **MF6** | Suporte a mouse (scroll, click)                                        | 🟢 LOW                         |

---

## 🏗️ 

### 🔴 CRITICAL

#### I1. `config.Validate()` existe mas nunca é chamado em produção
**Arquivos:** [config.go:28-36](file:///mnt/bkp/workspace/github/nerdminertui/internal/config/config.go#L28), [main.go:35-38](file:///mnt/bkp/workspace/github/nerdminertui/cmd/tui/main.go#L35)

`Load()` não chama `Validate()`. Grep no codebase mostra **zero chamadas** em código de produção. Um usuário pode setar `cpu_target: 99.0` e o app aceita. A validação manual em `main.go:49` duplica parte da lógica mas ignora bounds do CPUTarget.

#### I2. BTC address aceita qualquer string — sem validação de formato
**Arquivo:** [config.go:32-34](file:///mnt/bkp/workspace/github/nerdminertui/internal/config/config.go#L32)

`Validate()` só verifica `!= ""`. Aceita `"hello"`, `"DROP TABLE"`, ou senhas coladas acidentalmente. Endereço malformado = todo trabalho de mineração desperdiçado.

---

### 🟠 HIGH

| #   | Finding                                                                    | Arquivo                                                                                                                                                                          |
| --- | -------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| I3  | Flag `--config` é aceita mas silenciosamente ignorada                      | [main.go:24,34](file:///mnt/bkp/workspace/github/nerdminertui/cmd/tui/main.go#L24)                                                                                               |
| I4  | Permissões inseguras: config dir `0755`, debug.log `0666`                  | [config.go:105](file:///mnt/bkp/workspace/github/nerdminertui/internal/config/config.go#L105), [main.go:127](file:///mnt/bkp/workspace/github/nerdminertui/cmd/tui/main.go#L127) |
| I5  | Sem versionamento de migration — mudanças no schema quebram DBs existentes | [migrations.go:11-14](file:///mnt/bkp/workspace/github/nerdminertui/internal/store/migrations.go#L11)                                                                            |
| I6  | Release workflow roda testes sem `-race`                                   | [release.yml:19](file:///mnt/bkp/workspace/github/nerdminertui/.github/workflows/release.yml#L19)                                                                                |
| I7  | Release workflow sem `-trimpath` — vaza paths do build                     | [release.yml:59](file:///mnt/bkp/workspace/github/nerdminertui/.github/workflows/release.yml#L59)                                                                                |
| I8  | SQLite store inicializado mas nunca utilizado (`_ = st`)                   | [main.go:54-69](file:///mnt/bkp/workspace/github/nerdminertui/cmd/tui/main.go#L54)                                                                                               |

---

### 🟡 MEDIUM

| #   | Finding                                                                         | Arquivo                                                                                                  |
| --- | ------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------- |
| I9  | Testes usam `os.Setenv` em vez de `t.Setenv` — env vars vazam                   | [config_test.go:82-91](file:///mnt/bkp/workspace/github/nerdminertui/internal/config/config_test.go#L82) |
| I10 | Zero `t.Parallel()` no projeto — CI desnecessariamente lenta                    | Todos `*_test.go`                                                                                        |
| I11 | `ExpandPath("~")` sem `/` retorna `"~"` literal                                 | [paths.go:10](file:///mnt/bkp/workspace/github/nerdminertui/internal/config/paths.go#L10)                |
| I12 | Sem validação de porta no config (aceita -1, 99999)                             | [config.go:16](file:///mnt/bkp/workspace/github/nerdminertui/internal/config/config.go#L16)              |
| I13 | `ConnectionStatus` usa strings mágicas em PT (`"Conectado"`) — deveria ser enum | [state.go:42](file:///mnt/bkp/workspace/github/nerdminertui/internal/model/state.go#L42)                 |
| I14 | Sem `govulncheck` no CI                                                         | [ci.yml](file:///mnt/bkp/workspace/github/nerdminertui/.github/workflows/ci.yml)                         |
| I15 | Testes de store usam `cache=shared` — estado compartilhado entre tests          | [store_test.go:16](file:///mnt/bkp/workspace/github/nerdminertui/internal/store/store_test.go#L16)       |
| I16 | `Save()` ignora erros de `ExpandPath` e `MkdirAll`                              | [config.go:104-105](file:///mnt/bkp/workspace/github/nerdminertui/internal/config/config.go#L104)        |
| I17 | `QueryHashRateHistory` dupla alocação — reverse in-place                        | [store.go:73-76](file:///mnt/bkp/workspace/github/nerdminertui/internal/store/store.go#L73)              |
| I18 | Sem validação de `pool_address` (aceita string vazia)                           | [config.go:15](file:///mnt/bkp/workspace/github/nerdminertui/internal/config/config.go#L15)              |

### 🟢 LOW

| #   | Finding                                                                     | Arquivo                                                                                          |
| --- | --------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------ |
| I19 | Makefile não injeta version via `-X main.version=...` em build local        | [Makefile:10](file:///mnt/bkp/workspace/github/nerdminertui/Makefile#L10)                        |
| I20 | `.golangci.yml` sem lista explícita de linters — depende de defaults        | [.golangci.yml](file:///mnt/bkp/workspace/github/nerdminertui/.golangci.yml)                     |
| I21 | `testutil/` é muito fino — falta `NewTestStore(t)`, `NewTestConfig(t)`      | [testutil/fixtures.go](file:///mnt/bkp/workspace/github/nerdminertui/testutil/fixtures.go)       |
| I22 | README mínimo — falta screenshot, descrição, config reference, badges       | [README.md](file:///mnt/bkp/workspace/github/nerdminertui/README.md)                             |
| I23 | Comment `AppState` diz "100% thread safe" mas `app.go` faz mutações diretas | [state.go:29-31](file:///mnt/bkp/workspace/github/nerdminertui/internal/model/state.go#L29)      |
| I24 | `gocyclo` min-complexity = 25 (padrão industria é 10-15)                    | [.golangci.yml:3](file:///mnt/bkp/workspace/github/nerdminertui/.golangci.yml#L3)                |
| I25 | `SchemaDDL` exportado desnecessariamente (internal package)                 | [migrations.go:9](file:///mnt/bkp/workspace/github/nerdminertui/internal/store/migrations.go#L9) |

---

## 📊 Cobertura de Testes

| Pacote            | Arquivos | Testes              | Avaliação                                                 |
| ----------------- | -------- | ------------------- | --------------------------------------------------------- |
| `pkg/mining`      | 4 source | 3 test + benchmarks | ✅ **Boa**                                                 |
| `internal/worker` | 5 source | 1 test (13 cases)   | ✅ **Boa**                                                 |
| `internal/ui`     | 7 source | 1 test (screens)    | ⚠️ **Parcial** — `app.go` (320 linhas) tem **zero testes** |
| `internal/config` | 3 source | 1 test              | ⚠️ **Parcial** — `Save()` não testado                      |
| `internal/store`  | 3 source | 1 test              | ✅ **Boa** — inclui concurrency test                       |
| `internal/model`  | 1 source | 1 test              | ✅ **Boa** — imutabilidade testada                         |
| `pkg/format`      | 3 source | 1 test              | ⚠️ **Parcial** — bug M7 não coberto                        |
| `pkg/trivia`      | 1 source | 1 test              | ✅ **Adequada**                                            |

### Gaps Críticos de Teste

- `app.go` Update/View cycle — **0% coverage** no arquivo mais crítico
- `RebuildHeaderWithExtraNonce2` — nova função sem nenhum teste
- `MempoolClient.FetchStats` — teste real contra internet (flaky)
- `RenderSettings` — tela mais complexa sem teste de render
- Sem fuzz tests para `ParseStratumJob` com hex malformado

---

## 🏛️ Arquitetura — Pontos Positivos

> [!TIP]
> O projeto tem uma base arquitetural **sólida**:

- ✅ Separação limpa: `pkg/mining` é computação pura, `internal/worker` trata I/O
- ✅ Midstate optimization em `MinerHashState` — bem projetado e correto
- ✅ Interface `PoolClient` permite teste via `MockPoolClient`
- ✅ Detecção de goroutine leak via `goleak` no `TestMain`
- ✅ Abordagem extranonce2 por worker (merkle roots únicos) elimina coordenação de nonces
- ✅ CI + Release automáticos com GitHub Actions

## 🏛️ Arquitetura — Preocupações

- ⚠️ `StratumPoolClient` tem 380+ linhas com responsabilidades misturadas (TCP I/O, JSON parsing, state machine, reconexão). Considerar split em `StratumConnection` + `StratumProtocol`.
- ⚠️ Struct `Job` com 14 campos caminha para god-object. Separar `StratumJobMeta` (params de submissão) de `MiningJob` (header + target para hashing).
- ⚠️ `app.go` concentra muita lógica — 320 linhas com message handling, keyboard dispatch e state manipulation.

---

## 🗺️ Roadmap Sugerido

### Sprint 1: Quick Wins & Estabilidade (1-2 dias)
- [ ] Fix T1: Non-blocking throttle channel
- [ ] Fix I1: Chamar `Validate()` em main.go
- [ ] Fix I6/I7: `-race` e `-trimpath` no release workflow
- [ ] Fix T6: Usar constantes `ScreenID` em View()
- [ ] Fix M7: Bug do FormatUptime
- [ ] Fix T11: CPUBar não fabricar valor

### Sprint 2: UX & Feedback (2-3 dias)
- [ ] Impl MF1: Componente de toast/notificação
- [ ] Fix T2: Surfaçar erros da pool/miner
- [ ] Fix T3: Feedback visual para Ctrl+S
- [ ] Impl MF2: Tela de help (`?`)
- [ ] Fix T4: Shift+Tab para navegação reversa
- [ ] Fix T9: Larguras responsivas baseadas em terminal width

### Sprint 3: Performance do Mining (1-2 dias)
- [ ] Fix C1: Pré-reverter target, early-exit uint32
- [ ] Fix M1: Segundo hasher pré-alocado para double-SHA256
- [ ] Fix H1: `math/rand/v2` com instâncias per-worker
- [ ] Fix L4: Atribuição round-robin de SpaceWords
- [ ] Adicionar benchmark para `MeetsTarget`

### Sprint 4: Code Quality & Testes (2-3 dias)
- [ ] Fix T5: Constantes nomeadas para focus index
- [ ] Fix T7/T8: Extrair formatação duplicada
- [ ] Testes para `app.go` Update/View cycle
- [ ] Testes para `RebuildHeaderWithExtraNonce2`
- [ ] Fuzz tests para `ParseStratumJob`
- [ ] Fix I9/I10: `t.Setenv` + `t.Parallel()`
