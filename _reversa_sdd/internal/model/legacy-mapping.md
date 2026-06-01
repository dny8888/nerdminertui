# Mapeamento do Legado — internal/model

> **Status:** Mapeado da Especificação Alvo  
> **Módulo:** `internal/model`  
> **Nível de Confiança:** 🟢 CONFIRMADO

Este módulo é responsável por definir o estado global da aplicação e garantir a imutabilidade do `AppState` durante as transições de ciclo do Bubbletea.

---

## 1. Arquivos Mapeados no Legado

Os seguintes arquivos compõem o módulo `internal/model` com base no blueprint:

| Arquivo Alvo | Arquivo de Origem (Legado) | Linhas / Seção no Legado | Descrição |
|---|---|---|---|
| `internal/model/state.go` | [nerdtui-spec.md](file:///home/dny8888/workspace/github/nerdminertui/nerdtui-spec.md) | Lines 126-156 (§4), Lines 195-210 (§5.3) | Definição da struct `AppState`, constantes globais e funções puras de atualização de estado. |

---

## 2. Assinaturas e Componentes Mapeados

* **Struct de Estado (`AppState`)**:
  * Contém dados estruturados por cópia (sem ponteiros) para garantir imutabilidade.
* **Constantes e Enums**:
  * `ScreenID` (ScreenDashboard=0, ScreenClock=1, ScreenGlobalStats=2)
  * `NumScreens = 3`, `MinCPUTarget = 0.05`, `MaxCPUTarget = 1.00`, `CPUStep = 0.05`, `HashHistoryLen = 60`
* **Método Funcional (`WithHashRate`)**:
  * **Função**: `(s AppState) WithHashRate(hps float64) AppState`
  * **Comportamento**: Retorna uma cópia modificada com `HashRate` atualizado e rola o histórico FIFO `[60]float64` sem sofrer mutações no receptor original.
