# Fluxograma — cmd/tui

> **Módulo:** `cmd/tui`  
> **Gerado em:** 2026-05-29

Este fluxograma ilustra o processo de inicialização, parsing, wiring e início da aplicação terminal Bubbletea do NerdTUI.

```mermaid
flowchart TD
    Start([Início]) --> ParseFlags[Parse CLI Flags: --config, --no-mine, --cpu]
    ParseFlags --> LoadConfig[Carregar Configurações: config.Load]
    LoadConfig --> ValidateConfig{c.Validate == nil?}
    
    ValidateConfig -->|Não| LogFatal[Falha Fatal: log.Fatal e Encerrar]
    ValidateConfig -->|Sim| InitStore[Instanciar SQLite: store.New]
    
    InitStore --> InitWorker[Instanciar MinerWorker]
    InitWorker --> CreateApp[Instanciar Bubbletea AppModel com dependências]
    CreateApp --> RunApp[Executar tea.NewProgram com WithAltscreen]
    RunApp --> End([Fim do Programa])
```
