# Mapeamento do Legado — internal/ui

> **Status:** Mapeado da Especificação Alvo  
> **Módulo:** `internal/ui`  
> **Nível de Confiança:** 🟢 CONFIRMADO

Este módulo implementa a camada de visualização interativa baseada em Bubbletea, contendo o model raiz, keybindings estruturados, telas rotativas puras (Dashboard, Clock, Global Stats) e componentes estilizados com Lipgloss.

---

## 1. Arquivos Mapeados no Legado

Os seguintes arquivos compõem o módulo `internal/ui` com base no blueprint:

| Arquivo Alvo | Arquivo de Origem (Legado) | Linhas / Seção no Legado | Descrição |
|---|---|---|---|
| `internal/ui/app.go` | [nerdtui-spec.md](file:///home/dny8888/workspace/github/nerdminertui/nerdtui-spec.md) | Lines 264-290 (§5.6) | Bubbletea `AppModel` contendo loops de Init, Update e View. |
| `internal/ui/keys.go` | [nerdtui-spec.md](file:///home/dny8888/workspace/github/nerdminertui/nerdtui-spec.md) | Line 96 (§3) | Keybindings nomeados via `key.Map`. |
| `internal/ui/screens/dashboard.go` | [nerdtui-spec.md](file:///home/dny8888/workspace/github/nerdminertui/nerdtui-spec.md) | Lines 394-412 (§8) | Renderização pura da tela de Mining Dashboard. |
| `internal/ui/screens/clock.go` | [nerdtui-spec.md](file:///home/dny8888/workspace/github/nerdminertui/nerdtui-spec.md) | Line 414 (§8) | Renderização pura da tela de Relógio ASCII. |
| `internal/ui/screens/globalstats.go` | [nerdtui-spec.md](file:///home/dny8888/workspace/github/nerdminertui/nerdtui-spec.md) | Line 418 (§8) | Renderização das estatísticas da rede Bitcoin. |
| `internal/ui/components/hashgauge.go` | [nerdtui-spec.md](file:///home/dny8888/workspace/github/nerdminertui/nerdtui-spec.md) | Line 102 (§3) | Barra gráfica do hashrate atual. |
| `internal/ui/components/sparkline.go` | [nerdtui-spec.md](file:///home/dny8888/workspace/github/nerdminertui/nerdtui-spec.md) | Line 103 (§3) | Sparkline gráfico do histórico de hashrate (60s). |
| `internal/ui/components/cpubar.go` | [nerdtui-spec.md](file:///home/dny8888/workspace/github/nerdminertui/nerdtui-spec.md) | Line 104 (§3) | Barra de medição e target de uso de CPU. |
| `internal/ui/components/statusbar.go` | [nerdtui-spec.md](file:///home/dny8888/workspace/github/nerdminertui/nerdtui-spec.md) | Line 105 (§3) | Barra de rodapé com metadados e status de erros. |

---

## 2. Assinaturas e Componentes Mapeados

* **Bubbletea Model Interface (`AppModel`)**:
  * `Init() tea.Cmd` — Inicia ticker de 1s para Uptime e dispara primeiro poll.
  * `Update(tea.Msg) (tea.Model, tea.Cmd)` — Reage a inputs do teclado (`tab`, `+`, `-`, `q`) e mensagens dos workers assíncronos (`HashRateMsg`, `ShareFoundMsg`, `PoolStatsMsg`).
  * `View() string` — Renderiza a tela ativa no terminal dinamicamente respondendo a `tea.WindowSizeMsg`.
* **Visualizações Puras (Screens e Components)**:
  * Todos os métodos possuem a assinatura `Render*(state AppState, w, h int) string` sem efeitos colaterais.
