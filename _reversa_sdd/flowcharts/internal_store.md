# Fluxograma — internal/store

> **Módulo:** `internal/store`  
> **Gerado em:** 2026-05-29

Este fluxograma ilustra o fluxo de gravação e consulta de histórico concorrente no banco de dados SQLite (WAL Mode) do NerdTUI.

```mermaid
flowchart TD
    Start([Chamada de I/O de Store]) --> Operation{Operação?}
    
    Operation -->|AppendHashRate| AppendExec[Executar INSERT com busy_timeout de 5000ms]
    Operation -->|QueryHashRateHistory| QueryExec[Executar SELECT com LIMIT 60 e recorded_at DESC]
    
    AppendExec --> CheckAppendErr{Erro?}
    CheckAppendErr -->|Sim| ReturnAppendErr[Retornar error]
    CheckAppendErr -->|Não| ReturnAppendOK[Retornar nil]
    
    QueryExec --> CheckQueryErr{Erro?}
    CheckQueryErr -->|Sim| ReturnQueryErr[Retornar error]
    CheckQueryErr -->|Não| ParseRows[Construir slice de float64]
    ParseRows --> ReturnQueryOK[Retornar slice e nil]
```
