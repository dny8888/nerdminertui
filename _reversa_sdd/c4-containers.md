# C4 Diagrama de Containers (Nível 2) — nerdminertui

> **Módulo:** Arquitetura Global  
> **Nível de Documentação:** COMPLETO  
> **Gerado pelo Arquiteto em:** 2026-05-29

Este diagrama detalha os limites lógicos internos do executável do **NerdTUI**, mapeando seus containers funcionais (threads/goroutines e banco).

```mermaid
C4Container
    title Diagrama de Containers do executável NerdTUI (Nível 2)
    
    Person(user, "Operador do Terminal", "Desenvolvedor/Minerador")
    
    System_Boundary(c1, "Executável Estático NerdTUI CLI") {
        Container(tui, "Bubbletea UI Container", "Go, Bubbletea, Lipgloss", "Gerencia o loop MUV de desenho de tela, teclas e dispatch de mensagens na thread principal.")
        
        Container(miner, "Background Hashing Engine", "Go (Goroutine)", "Executa loop de hashing SHA256d de alta performance e calcula o micro-sleep de CPU throttling.")
        
        Container(fetcher, "Statistics & Job Fetcher", "Go (Goroutines)", "Busca estatísticas REST a cada 5s e mantém conexão Stratum TCP aberta para extrair jobs de mineração.")
        
        ContainerDb(sqlite, "Local Metrics Store", "modernc.org/sqlite (Go puro)", "Banco SQLite embarcado operando em modo WAL para gravar histórico de HPS.")
    }
    
    System_Ext(stratum, "Pool Stratum Server", "TCP Stratum JSON-RPC")
    System_Ext(rest, "Public Pool REST API", "HTTP REST API")
    
    Rel(user, tui, "Digita teclas, lê telas do terminal", "Terminal VTY")
    Rel(tui, miner, "Altera CPUTarget", "channels (throttleCh)")
    Rel(miner, tui, "Envia taxas e shares válidos", "tea.Msg (HashRateMsg, ShareFoundMsg)")
    Rel(fetcher, tui, "Envia bloco atual e status de rede", "tea.Msg (PoolStatsMsg)")
    Rel(tui, sqlite, "Grava hashrate (1s) e query para sparklines", "Go SQL driver API")
    
    Rel(fetcher, rest, "GET /api/stats", "HTTP REST / JSON")
    Rel(fetcher, stratum, "JSON-RPC 1.0 (Stratum)", "TCP Socket")
```
