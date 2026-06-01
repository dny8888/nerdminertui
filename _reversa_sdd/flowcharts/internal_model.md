# Fluxograma — internal/model

> **Módulo:** `internal/model`  
> **Gerado em:** 2026-05-29

Este fluxograma ilustra o fluxo de atualização funcional e imutável do receptor `AppState` ao receber novas medições de hashrate.

```mermaid
flowchart TD
    Start([Chamada WithHashRate]) --> CopyState[Copiar receiver AppState por valor para newState]
    CopyState --> SetHashRate[Atribuir newState.HashRate = hps]
    SetHashRate --> FIFOShift[Deslocar newState.HashRateHistory FIFO com nova entrada]
    FIFOShift --> ReturnNewState[Retornar newState]
```
