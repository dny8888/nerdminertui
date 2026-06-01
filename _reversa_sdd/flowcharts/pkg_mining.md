# Fluxograma — pkg/mining

> **Módulo:** `pkg/mining`  
> **Gerado em:** 2026-05-29

Este fluxograma ilustra o processamento puro e determinístico de duplo SHA-256 e a verificação matemática de target de mineração do NerdTUI.

```mermaid
flowchart TD
    Start([Chamada HashHeader]) --> StructHeader[Estruturar bloco: Header + Nonce]
    StructHeader --> SHA1[Executar primeiro SHA256 no buffer]
    SHA1 --> SHA2[Executar segundo SHA256 no primeiro hash]
    SHA2 --> ReturnHash[Retornar [32]byte de SHA256d]
    
    StartVerify([Chamada MeetsTarget]) --> CompareBytes[Comparar big-endian hash < target byte-a-byte]
    CompareBytes --> ReturnBool[Retornar bool]
```
