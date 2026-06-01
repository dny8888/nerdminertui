# ERD Completo (Modelo de Dados) — nerdminertui

> **Módulo:** Persistência de Dados  
> **Nível de Documentação:** COMPLETO  
> **Gerado pelo Arquiteto em:** 2026-05-29

Este diagrama de Relacionamento de Entidades (ERD) documenta as tabelas estruturadas no banco SQLite e as relações lógicas com modelos voláteis de memória.

```mermaid
erDiagram
    %% Tabelas SQLite Físicas
    hashrate_history {
        INTEGER id PK "Auto-incremental"
        REAL hps "Hashes por segundo gerados"
        INTEGER recorded_at "Unix timestamp da medição (Int64)"
    }
    
    %% Estruturas Lógicas em Memória (Voláteis)
    AppState {
        float64 HashRate "Medição instantânea"
        float64_array HashRateHistory "FIFO dos últimos 60 snapshots"
        uint64 SharesFound "Contagem acumulativa"
        float64 BestDifficulty "Maior dificuldade encontrada"
        uint32 BlockHeight "Altura de bloco atual"
        float64 CPUTarget "Fração de target de CPU [0.05, 1.0]"
        float64 CPUActual "Medição real do throttle"
        bool PoolConnected "Estado de conexão Stratum"
        string PoolURL "Endpoint configurado"
        time_Duration Uptime "Uptime desde o startup"
        time_Time StartedAt "Timestamp de início"
        ScreenID Screen "Identificador da tela ativa"
        string Error "Último erro de I/O na statusbar"
    }

    Job {
        byte_array Header "Bloco candidato Bitcoin"
        byte_32_array Target "Target hash de dificuldade"
        uint32 ExtraNonce "Contador nonce extra para worker"
        uint32 Height "Altura do bloco correspondente"
    }

    %% Relações Lógicas / Mapeamentos de Fluxos
    AppState ||--o{ hashrate_history : "Recupera histórico (limite 60) no cold start e insere a cada 1s"
    Job ||--o{ AppState : "Atualiza altura de bloco e target em memória"
```
