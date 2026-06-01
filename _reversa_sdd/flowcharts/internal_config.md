# Fluxograma — internal/config

> **Módulo:** `internal/config`  
> **Gerado em:** 2026-05-29

Este fluxograma ilustra o fluxo de carregamento das configurações do NerdTUI e suas respectivas regras de validação estritas.

```mermaid
flowchart TD
    Start([Chamada Load]) --> LoadViper[Carregar padrões e overrides de envs via Viper]
    LoadViper --> RunValidate[Executar c.Validate]
    
    RunValidate --> CheckAddress{MockMining == false && BTCAddress == ""}
    CheckAddress -->|Sim| ErrAddress[Retornar ErrBTCAddressRequired]
    
    CheckAddress -->|Não| CheckCPU{CPUTarget < 0.05 || CPUTarget > 1.0}
    CheckCPU -->|Sim| ErrCPU[Retornar ErrCPUTargetOutOfBounds]
    
    CheckCPU -->|Não| ReturnConfig[Retornar Config e nil]
    
    ErrAddress --> ReturnErr[Retornar error]
    ErrCPU --> ReturnErr
```
