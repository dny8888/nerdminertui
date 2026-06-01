# Inventário do Projeto — nerdminertui

> **Status:** Mapeamento de Especificações Greenfield (Design)  
> **Gerado pelo Scout em:** 2026-05-29

---

## 1. Visão Geral Física do Repositório

O repositório físico está atualmente em estágio de **especificação/greenfield**, contendo as definições do projeto a ser construído:

| Arquivo/Pasta | Tipo | Descrição |
|---|---|---|
| [nerdtui-spec.md](file:///home/dny8888/workspace/github/nerdminertui/nerdtui-spec.md) | Arquivo (MD) | Especificação completa de design, arquitetura, testes e quality gates do NerdTUI. |
| [AGENTS.md](file:///home/dny8888/workspace/github/nerdminertui/AGENTS.md) | Arquivo (MD) | Definição de ativação do framework de Engenharia Reversa. |
| `.reversa/` | Diretório | Arquivos de configuração, estado e templates do framework Reversa. |
| `.agents/` | Diretório | Skills e configurações dos agentes de IA do projeto. |

---

## 2. Arquitetura Alvo do Projeto (conforme `nerdtui-spec.md`)

Abaixo está o inventário de pacotes e componentes planejados para o **NerdTUI**, conforme a arquitetura estipulada em `nerdtui-spec.md`:

```
nerdtui/
├── AGENTS.md                       # Diretrizes de navegação e restrições para agentes
├── Makefile                        # Automação de builds, testes, lint e coverage
├── go.mod                          # Arquivo de módulo Go (github.com/user/nerdtui)
│
├── cmd/tui/
│   └── main.go                     # Ponto de entrada (wiring, parse flags, run app)
│
├── internal/
│   ├── config/
│   │   └── config.go               # Configuração e validação via Viper
│   │
│   ├── model/
│   │   └── state.go                # AppState (imutável) e constantes de domínio
│   │
│   ├── worker/
│   │   ├── fetcher.go              # Polling de estatísticas e Stratum TCP
│   │   ├── poller.go               # Ticker de polling com retry exponencial
│   │   ├── miner.go                # Loop de mineração (CPU throttling, channels)
│   │   └── messages.go             # Todos os tea.Msg (mensagens tipadas)
│   │
│   ├── ui/
│   │   ├── app.go                  # AppModel (Bubbletea Model raiz: Init/Update/View)
│   │   ├── keys.go                 # Keybindings nomeados (key.Map)
│   │   ├── screens/
│   │   │   ├── dashboard.go        # Renderização do Dashboard principal (função pura)
│   │   │   ├── clock.go            # Renderização do Relógio ASCII (função pura)
│   │   │   └── globalstats.go      # Renderização das estatísticas globais (função pura)
│   │   └── components/
│   │       ├── hashgauge.go        # Medidor visual de hashrate (função pura)
│   │       ├── sparkline.go        # Gráfico sparkline do histórico de hashrate
│   │       ├── cpubar.go           # Barra de progresso da CPU (CPUTarget vs CPUActual)
│   │       └── statusbar.go        # Barra de status do rodapé (conexão, uptime, etc.)
│   │
│   └── store/
│       ├── store.go                # Interface e implementação SQLiteStore (WAL mode)
│       └── migrations.go           # SQL de criação de tabelas (embutido)
│
└── pkg/
    ├── mining/
    │   ├── hash.go                 # Primitivas de duplo SHA256 (SHA256d)
    │   ├── target.go               # MeetsTarget e cálculo de dificuldade
    │   └── job.go                  # Estrutura Job (header, target, extranonce, etc.)
    └── format/
        ├── hashrate.go             # Formatação de hashrate (ex: "12.4 KH/s", "1.2 MH/s")
        ├── duration.go             # Formatação de uptime (ex: "2d 03h 14m")
        └── difficulty.go           # Formatação de dificuldade (ex: "1.23e+12")
```

### Detalhes de Organização e Restrições de Importação

* **pkg/**: Contém código de utilidade geral e primitivas puras e matemáticas. **Proibido importar `internal/` ou `cmd/`**.
* **internal/**: Encapsula a lógica da aplicação, TUI, trabalhadores de background e repositório. **Proibido importar `cmd/`**.
* **cmd/**: Contém apenas o wiring principal e inicialização.

---

## 3. Mapeamento de Módulos para os Próximos Agentes

Para as próximas fases de engenharia reversa e geração de especificações SDD, usaremos o mapeamento de módulos extraído do design de pacotes acima:

1. **cmd/tui** — Ponto de entrada e bootstrapping.
2. **internal/config** — Validação e parsing de configurações e variáveis de ambiente.
3. **internal/model** — Definição do estado imutável da TUI.
4. **internal/worker** — Gerenciamento do loop de hashing, Stratum TCP, HTTP polling e sincronização Bubbletea por mensagens (`tea.Msg`).
5. **internal/ui** — Design, componentes visuais puros e loop principal (Model/Update/View).
6. **internal/store** — Armazenamento SQLite local em Go puro.
7. **pkg/mining** — Primitivas matemáticas e criptográficas de mineração Bitcoin.
8. **pkg/format** — Formatação de unidades no terminal.
