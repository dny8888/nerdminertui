# Mapeamento do Legado — internal/store

> **Status:** Mapeado da Especificação Alvo  
> **Módulo:** `internal/store`  
> **Nível de Confiança:** 🟢 CONFIRMADO

Este módulo implementa a persistência local em SQLite no modo sem CGO (`modernc.org/sqlite`), registrando estatísticas de hashrate e logs de shares locais para garantir sobrevivência do histórico de sparklines a reinicializações.

---

## 1. Arquivos Mapeados no Legado

Os seguintes arquivos compõem o módulo `internal/store` com base no blueprint:

| Arquivo Alvo | Arquivo de Origem (Legado) | Linhas / Seção no Legado | Descrição |
|---|---|---|---|
| `internal/store/store.go` | [nerdtui-spec.md](file:///home/dny8888/workspace/github/nerdminertui/nerdtui-spec.md) | Lines 301-314 (§5.7) | Interface `Store` e sua implementação concreta `SQLiteStore` com pragmas de WAL. |
| `internal/store/migrations.go` | [nerdtui-spec.md](file:///home/dny8888/workspace/github/nerdminertui/nerdtui-spec.md) | Lines 315-324 (§5.7) | SQL DDL embutido para criação de tabelas e indexadores. |

---

## 2. Assinaturas e Componentes Mapeados

* **Interface `Store`**:
  * `AppendHashRate(ctx context.Context, hps float64, at time.Time) error`
  * `QueryHashRateHistory(ctx context.Context, limit int) ([]float64, error)`
  * `Close() error`
* **Esquema de Tabelas (`hashrate_history`)**:
  * ID (INTEGER PK), hps (REAL), recorded_at (INTEGER - Unix timestamp)
  * Indexador: `idx_hashrate_recorded` ordenado por timestamp decrescente para otimizar a busca mais recente.
* **Implementação No-Op (`NilStore`)**:
  * Utilizado para desativar a gravação no disco com a flag `--no-store`.
