# Pending Architectural Gaps — NerdTUI

The following gaps were identified during the Reversa Review phase.

---

## Active Gaps List

*No pending gaps remaining. All identified architectural ambiguities have been resolved and documented in `questions.md`.*

---

## Completed and Resolved Gaps

### 1. HTTP REST & Stratum TCP Co-existence
- **Módulo:** `internal/worker`
- **Severidade:** CRÍTICA
- **Status:** ✅ RESOLVIDO
- **Resolução:** HTTP Client e Stratum TCP rodam de forma concorrente no modo real de mineração. O Stratum mantêm a conexão viva de jobs e o HTTP REST busca estatísticas gerais da pool a cada 5s para exibição.

### 2. Expansão de Atalhos de Caminho (Tilde Expansion)
- **Módulo:** `internal/config`
- **Severidade:** MODERADA
- **Status:** ✅ RESOLVIDO
- **Resolução:** A expansão de atalhos (`~/`) será tratada dinamicamente pela função `ExpandPath` em `internal/config/paths.go` durante a execução de `config.Load()`.

### 3. Persistência de Métricas Acumuladas
- **Módulo:** `internal/store`
- **Severidade:** MODERADA
- **Status:** ✅ RESOLVIDO
- **Resolução:** `SharesFound` e `BestDifficulty` são transientes na v1, reiniciando com 0 a cada execução. Somente o histórico do hashrate de 60s será persistido em SQLite.

### 4. Submissão Virtual de Shares no Stratum
- **Módulo:** `internal/worker`
- **Severidade:** CRÍTICA
- **Status:** ✅ RESOLVIDO
- **Resolução:** A submissão física via `mining.submit` JSON-RPC é realizada ao encontrar shares, com timeout de 10s. O status da submissão é exibido em tempo real no Dashboard.
