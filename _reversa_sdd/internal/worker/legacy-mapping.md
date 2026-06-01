# Mapeamento do Legado — internal/worker

> **Status:** Mapeado da Especificação Alvo  
> **Módulo:** `internal/worker`  
> **Nível de Confiança:** 🟢 CONFIRMADO

Este módulo orquestra os trabalhadores assíncronos: a goroutine do minerador (`MinerWorker`) que executa o loop de hashing SHA256d com CPU throttling, e o `fetcher` que busca estatísticas de pool via HTTP REST ou Stratum TCP.

---

## 1. Arquivos Mapeados no Legado

Os seguintes arquivos compõem o módulo `internal/worker` com base no blueprint:

| Arquivo Alvo | Arquivo de Origem (Legado) | Linhas / Seção no Legado | Descrição |
|---|---|---|---|
| `internal/worker/fetcher.go` | [nerdtui-spec.md](file:///home/dny8888/workspace/github/nerdminertui/nerdtui-spec.md) | Lines 212-232 (§5.4) | Clientes HTTP e Stratum para buscar trabalhos e estatísticas. |
| `internal/worker/poller.go` | [nerdtui-spec.md](file:///home/dny8888/workspace/github/nerdminertui/nerdtui-spec.md) | Line 90 (§3) | Ticker e retry exponencial para reconexão da pool. |
| `internal/worker/miner.go` | [nerdtui-spec.md](file:///home/dny8888/workspace/github/nerdminertui/nerdtui-spec.md) | Lines 234-262 (§5.5) | Thread de hashing controlada com CPU throttling e channels. |
| `internal/worker/messages.go` | [nerdtui-spec.md](file:///home/dny8888/workspace/github/nerdminertui/nerdtui-spec.md) | Lines 356-369 (§6) | Mensagens de sincronização do Bubbletea (`tea.Msg`). |

---

## 2. Assinaturas e Componentes Mapeados

* **Interface `PoolClient`**:
  * `FetchStats(ctx) (PoolStats, error)`
  * `FetchJob(ctx) (mining.Job, error)`
* **Minerador e CPU Throttling (`MinerWorker`)**:
  * **Estrutura**: `MinerWorker { throttleCh chan float64, outCh chan tea.Msg, job atomic.Value }`
  * **Algoritmo de Throttling**: `sleep = workDuration * (1 - P) / P` onde P é o `CPUTarget`.
  * **Ação**: Executa batch de 50.000 hashes, mede o tempo, calcula a proporção e dorme, emitindo `HashRateMsg` a cada 1s e `ShareFoundMsg` quando um nonce satisfaz o target do Job.
