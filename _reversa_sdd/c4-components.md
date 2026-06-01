# C4 Diagrama de Componentes (Nível 3) — nerdminertui

> **Módulo:** Arquitetura Global  
> **Nível de Documentação:** COMPLETO  
> **Gerado pelo Arquiteto em:** 2026-05-29

Este diagrama de nível 3 decompõe os componentes e pacotes Go mapeados dentro dos containers de execução do **NerdTUI**.

```mermaid
C4Component
    title Diagrama de Componentes Go do executável NerdTUI (Nível 3)
    
    Container_Boundary(tui_boundary, "Bubbletea UI Container") {
        Component(app, "ui.AppModel", "Go struct", "Controlador central do loop Elm (Init, Update, View).")
        Component(screens, "ui.screens", "Go functions puras", "Renderizadores de visualizações completas (Dashboard, Clock, GlobalStats).")
        Component(components, "ui.components", "Go functions puras", "Renderizadores de sub-widgets (CPUBar, Sparkline, HashGauge, StatusBar).")
        Component(model, "model.AppState", "Go struct (valor)", "Modelo de estado imutável contendo os dados de rendering da UI.")
        Component(keys, "ui.keys", "Go configuration", "Mapeador de keybindings nomeados da aplicação.")
    }
    
    Container_Boundary(worker_boundary, "Background Hashing Engine & Fetcher") {
        Component(miner_worker, "worker.MinerWorker", "Go struct & thread", "Loop de mineração concorrente em batches com CPU sleep scheduler.")
        Component(fetcher, "worker.Fetcher & Poller", "Go client & ticker", "Polla blocos a cada 5s de APIs externas e de conexões Stratum.")
        Component(messages, "worker.messages", "Go types", "Mensagens do Bubbletea (tea.Msg) compartilhadas pelos pacotes.")
    }
    
    Container_Boundary(store_boundary, "Local Metrics Store") {
        Component(db_store, "store.SQLiteStore", "Go struct & driver", "Interface e implementação concreta do banco SQLite em Go puro.")
        Component(db_migrations, "store.Migrations", "Go SQL embed", "Cria tabelas e indexadores em memória RAM ou arquivo local.")
    }
    
    Container_Boundary(pkg_boundary, "pkg / Core Matemática") {
        Component(pkg_mining, "pkg.mining", "Go pure functions", "Cálculos matemáticos puros e algoritmos de duplo SHA-256 (SHA256d).")
        Component(pkg_format, "pkg.format", "Go pure functions", "Formatação de HPS, tempo decorrido e strings numéricas.")
    }
    
    Rel(app, model, "Lê e clona novas instâncias na mutação", "Pass por cópia")
    Rel(app, keys, "Consulta keybindings", "Memória")
    Rel(app, screens, "Delega renderings de tela inteira", "Memory Call")
    Rel(screens, components, "Compõe layouts com sub-widgets", "Memory Call")
    Rel(screens, pkg_format, "Formata valores para string", "Memory Call")
    
    Rel(miner_worker, pkg_mining, "Calcula duplo hash SHA256 e verifica targets", "Memory Call")
    Rel(miner_worker, messages, "Escreve HashRateMsg e ShareFoundMsg no outCh", "Go Channels")
    Rel(fetcher, messages, "Escreve PoolStatsMsg no outCh", "Go Channels")
    
    Rel(messages, app, "Processado de forma serial no Update()", "Bubbletea internal queue")
    Rel(app, db_store, "Grava histórico a cada 1s", "SQL API")
    Rel(db_store, db_migrations, "Inicializa tabelas no startup", "SQL execution")
```
