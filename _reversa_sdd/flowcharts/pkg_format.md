# Fluxograma — pkg/format

> **Módulo:** `pkg/format`  
> **Gerado em:** 2026-05-29

Este fluxograma ilustra o processamento puro de formatação de strings do hashrate atual com divisões e escalas decimais (H/s, KH/s, MH/s).

```mermaid
flowchart TD
    Start([Chamada FormatHashRate]) --> CompZero{hps == 0?}
    CompZero -->|Sim| ReturnZero[Retornar '0 H/s']
    
    CompZero -->|Não| CompMega{hps >= 1.000.000?}
    CompMega -->|Sim| ReturnMega[Retornar hps / 1.000.000 com sufixo 'MH/s']
    
    CompMega -->|Não| CompKilo{hps >= 1.000?}
    CompKilo -->|Sim| ReturnKilo[Retornar hps / 1.000 com sufixo 'KH/s']
    
    CompKilo -->|Não| ReturnBase[Retornar hps com sufixo 'H/s']
```
