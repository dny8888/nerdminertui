# Mapeamento do Legado — cmd/tui

> **Status:** Mapeado da Especificação Alvo  
> **Módulo:** `cmd/tui`  
> **Nível de Confiança:** 🟢 CONFIRMADO (especificação de wiring)

Este módulo é o ponto de entrada da aplicação, responsável pelo bootstrapping e wiring das dependências sem conter lógica de negócio.

---

## 1. Arquivos Mapeados no Legado

Os seguintes arquivos compõem o módulo `cmd/tui` com base no blueprint:

| Arquivo Alvo | Arquivo de Origem (Legado) | Linhas / Seção no Legado | Descrição |
|---|---|---|---|
| `cmd/tui/main.go` | [nerdtui-spec.md](file:///home/dny8888/workspace/github/nerdminertui/nerdtui-spec.md) | Lines 160-174 (§5.1) | Sequência de inicialização, parse de flags e execução do Bubbletea Program. |
| `Makefile` | [nerdtui-spec.md](file:///home/dny8888/workspace/github/nerdminertui/nerdtui-spec.md) | Lines 862-923 (§13) | Comandos de build, run e mocks do terminal. |

---

## 2. Assinaturas e Componentes Mapeados

* **Entidade de Inicialização (`main`)**:
  * **Arquivo**: `cmd/tui/main.go`
  * **Fluxo**:
    1. Parse flags (`--config`, `--no-mine`, `--cpu`).
    2. `config.Load()` para obter a configuração.
    3. `store.New(cfg.StorePath)` para obter o SQLite database handler.
    4. Inicia channel de controle do `MinerWorker` com buffers apropriados.
    5. Inicia e injeta dependências no Bubbletea `AppModel`.
    6. Executa `tea.NewProgram(model, tea.WithAltscreen()).Run()`.
