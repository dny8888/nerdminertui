# Fluxograma — internal/worker

> **Módulo:** `internal/worker`  
> **Gerado em:** 2026-05-29

Este fluxograma ilustra o loop assíncrono do `MinerWorker` com CPU throttling e batches de mineração Bitcoin de 50.000 iterações.

```mermaid
flowchart TD
    Start([Início da Goroutine do MinerWorker]) --> MainLoop{Contexto cancelado?}
    
    MainLoop -->|Sim| CleanUp[Encerrar goroutine de forma limpa]
    MainLoop -->|Não| SelectThrottle{Lercanal throttleCh?}
    
    SelectThrottle -->|Sim| UpdateCPUTarget[Atualizar target de CPU local]
    SelectThrottle -->|Não| ReadJob[Carregar Job atual atomicamente]
    
    UpdateCPUTarget --> ReadJob
    ReadJob --> GetTime[Registrar timestamp inicial]
    
    GetTime --> LoopBatch[Executar batch de 50.000 hashes SHA256d]
    LoopBatch --> CompareTarget{Algum hash satisfaz target do Job?}
    
    CompareTarget -->|Sim| EmitShare[Enviar ShareFoundMsg no outCh]
    CompareTarget -->|Não| MeasureTime[Medir workDuration e CPUActual]
    
    EmitShare --> MeasureTime
    MeasureTime --> CalculateSleep[Calcular sleep = workDuration * 1-P / P]
    CalculateSleep --> SleepState[Dormir sleepDuration]
    SleepState --> EmitHPS[Enviar HashRateMsg no outCh via ticker]
    EmitHPS --> MainLoop
```
