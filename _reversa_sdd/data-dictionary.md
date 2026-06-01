# Dicionário de Dados — nerdminertui

> **Status:** Mapeamento de Especificações Greenfield (Design)  
> **Nível de Documentação:** COMPLETO  
> **Gerado pelo Archaeologist em:** 2026-05-29

Este documento detalha todas as estruturas de dados, esquemas de banco de dados, parâmetros de configuração e mensagens do domínio mapeados na especificação do **NerdTUI**.

---

## 1. Modelo de Domínio (Entidades de Memória)

### 1.1 `AppState` (`internal/model/state.go`)
Representa o estado completo e imutável da interface do terminal.

| Campo | Tipo | Obrigatório | Valor Padrão | Descrição |
|---|---|---|---|---|
| `HashRate` | `float64` | Sim | `0.0` | Hashes por segundo calculados no último intervalo de 1s. |
| `HashRateHistory` | `[60]float64` | Sim | `[60]float64{0}` | Histórico circular de hashrate (últimos 60 snapshots). |
| `SharesFound` | `uint64` | Sim | `0` | Número de shares que passaram no target local de mineração. |
| `BestDifficulty` | `float64` | Sim | `0.0` | Maior dificuldade de share encontrada nesta execução. |
| `BlockHeight` | `uint32` | Sim | `0` | Altura do bloco mais recente recebido da pool. |
| `CPUTarget` | `float64` | Sim | `0.5` | Fração pretendida de uso de CPU (intervalo `[0.05, 1.0]`). |
| `CPUActual` | `float64` | Sim | `0.0` | Fração medida em tempo real de uso da CPU. |
| `PoolConnected` | `bool` | Sim | `false` | Status de conexão com a Pool de mineração. |
| `PoolURL` | `string` | Sim | `""` | Endereço (URL) da Pool configurada. |
| `Uptime` | `time.Duration` | Sim | `0` | Tempo decorrido desde o início da execução da TUI. |
| `StartedAt` | `time.Time` | Sim | (Tempo de inicialização) | Timestamp indicando o início da mineração. |
| `Screen` | `ScreenID` | Sim | `ScreenDashboard (0)` | Identificador da tela ativa no terminal. |
| `Error` | `string` | Não | `""` | Mensagem do último erro não-fatal a ser mostrada na statusbar. |

### 1.2 `Job` (`pkg/mining/job.go`)
Representa um trabalho de mineração Bitcoin despachado pela pool ou mockado localmente.

| Campo | Tipo | Obrigatório | Descrição |
|---|---|---|---|
| `Header` | `[]byte` | Sim | Cabeçalho bruto do bloco Bitcoin a ser minerado. |
| `Target` | `[32]byte` | Sim | Target hash Bitcoin (big-endian). Nonce deve gerar hash menor que este. |
| `ExtraNonce` | `uint32` | Sim | Contador de nonce extra para variação de hash em workers concorrentes. |
| `Height` | `uint32` | Sim | Altura associada ao bloco deste trabalho de mineração. |

---

## 2. Configurações (`internal/config/config.go`)

Representa as chaves lidas pelo Viper dos arquivos e variáveis de ambiente:

| Propriedade Config | Tipo | Variável Env | Default | Descrição / Validação |
|---|---|---|---|---|
| `PoolURL` | `string` | `NM_POOL_URL` | `"public-pool.io"` | Endereço do endpoint Stratum/REST. |
| `PoolPort` | `int` | `NM_POOL_PORT` | `21496` | Porta de conexão TCP da pool. |
| `BTCAddress` | `string` | `NM_BTC_ADDRESS` | `""` | Endereço Bitcoin para recebimento de recompensas. Opcional em mock, obrigatório em prod. |
| `PollInterval` | `time.Duration` | `NM_POLL_INTERVAL` | `5s` | Frequência de atualização de estatísticas de pool. |
| `CPUTarget` | `float64` | `NM_CPU_TARGET` | `0.5` | Target inicial de CPU limit. Deve estar entre `[0.05, 1.0]`. |
| `StorePath` | `string` | `NM_STORE_PATH` | `"~/.nerdtui/metrics.db"` | Caminho físico de escrita do banco de dados SQLite. |
| `Theme` | `string` | `NM_THEME` | `"dark"` | Tema de layout de cores da TUI (`dark`, `light`). |
| `MockMining` | `bool` | Flag `--mock` | `false` | Quando ativado, ignora conexão Stratum e mina blocos simulados. |

---

## 3. Banco de Dados SQLite (`internal/store/migrations.go`)

### 3.1 Tabela `hashrate_history`
Persiste dados de hashrate ao longo das execuções do programa para manter as sparklines preenchidas.

```sql
CREATE TABLE hashrate_history (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    hps         REAL    NOT NULL,
    recorded_at INTEGER NOT NULL  -- Unix timestamp de gravação (int64)
);
```

#### Indexadores associados:
```sql
CREATE INDEX idx_hashrate_recorded ON hashrate_history(recorded_at DESC);
```

---

## 4. Contratos de Mensagens Bubbletea (`internal/worker/messages.go`)

Todas as transições de estado assíncronas ocorrem através das mensagens abaixo:

| Estrutura de Mensagem | Atributos internos | Emissor (Worker) | Propósito |
|---|---|---|---|
| `HashRateMsg` | `HPS float64`, `CPUActual float64` | `MinerWorker` | Atualiza o hashrate atual e de CPU real medidos na UI. |
| `ShareFoundMsg` | `Nonce uint32`, `Hash [32]byte`, `Difficulty float64` | `MinerWorker` | Disparado quando um hash passa no target local. |
| `PoolStatsMsg` | `BlockHeight uint32`, `Connected bool`, `NetworkHashRate float64` | `PoolClient` (Fetcher) | Transmite as estatísticas atualizadas da pool de mineração. |
| `MinerErrorMsg` | `Err error` | `MinerWorker` | Propaga um erro ocorrido dentro da thread do minerador para o Model. |
| `PoolErrorMsg` | `Err error` | `PoolClient` (Fetcher) | Propaga um erro ocorrido em comunicações HTTP/TCP da pool. |
| `tickMsg` | `At time.Time` | `AppModel` (interno) | Disparado a cada 1s para atualizar o uptime global do NerdTUI. |
