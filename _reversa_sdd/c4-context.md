# C4 Diagrama de Contexto (Nível 1) — nerdminertui

> **Módulo:** Arquitetura Global  
> **Nível de Documentação:** COMPLETO  
> **Gerado pelo Arquiteto em:** 2026-05-29

Este diagrama contextualiza o **NerdTUI** em seu ecossistema, mapeando os usuários e os sistemas externos de integração.

```mermaid
C4Context
    title Diagrama de Contexto de Sistema para o NerdTUI (Nível 1)
    
    Person(user, "Operador do Terminal", "Desenvolvedor ou entusiasta de Bitcoin que executa o NerdTUI para acompanhar a mineração solo.")
    
    System(system, "NerdTUI CLI", "Dashboard de mineração solo de terminal. Mina blocos, gerencia CPU e exibe métricas interativas.")
    
    System_Ext(stratum, "Pool Stratum Server", "Servidor de pool solo de mineração (ex: public-pool.io via porta 21496 TCP). Envia trabalhos criptográficos.")
    
    System_Ext(rest, "Public Pool REST API", "Serviço REST de estatísticas públicas de blocos e hashrate global da rede.")
    
    Rel(user, system, "Executa, monitora estatísticas, altera target de CPU e rotaciona telas via teclado")
    Rel(system, stratum, "Recebe blocos candidatos (Jobs) e submete shares encontrados", "JSON-RPC via TCP")
    Rel(system, rest, "Consome estatísticas gerais de blocos e rede a cada 5s", "HTTP REST / JSON")
```
