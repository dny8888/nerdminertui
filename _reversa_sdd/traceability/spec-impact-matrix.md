# Matriz de Impacto e Rastreabilidade — nerdminertui

> **Status:** Mapeamento de Especificações Greenfield (Design)  
> **Nível de Documentação:** COMPLETO  
> **Gerado pelo Arquiteto em:** 2026-05-29

Esta matriz mapeia o impacto de alterações entre os módulos do **NerdTUI**, garantindo rastreabilidade completa durante futuras evoluções ou manutenções de código.

---

## 1. Tabela de Rastreabilidade Cruzada (Módulos vs Componentes)

Esta tabela mapeia quais pacotes Go sofrem impacto direto quando uma especificação ou funcionalidade de outro módulo é alterada.

| Módulo de Origem (Alterado) | cmd/tui | internal/config | internal/model | internal/worker | internal/ui | internal/store | pkg/mining | pkg/format |
|---|---|---|---|---|---|---|---|---|
| **cmd/tui** | — | Baixo | Baixo | Alto | Alto | Médio | Baixo | Baixo |
| **internal/config** | Alto | — | Baixo | Alto | Médio | Alto | Baixo | Baixo |
| **internal/model** | Alto | Baixo | — | Alto | Alto | Médio | Baixo | Baixo |
| **internal/worker** | Alto | Baixo | Alto | — | Alto | Alto | Alto | Baixo |
| **internal/ui** | Alto | Baixo | Alto | Alto | — | Alto | Baixo | Alto |
| **internal/store** | Alto | Baixo | Médio | Baixo | Alto | — | Baixo | Baixo |
| **pkg/mining** | Baixo | Baixo | Baixo | Alto | Baixo | Baixo | — | Baixo |
| **pkg/format** | Baixo | Baixo | Baixo | Baixo | Alto | Baixo | Baixo | — |

---

## 2. Detalhamento de Impacto Crítico

### 2.1 Alterações no Pacote `internal/model` (AppState)
* **Gravidade**: 🔴 CRÍTICA
* **Justificativa**: O `AppState` é a única fonte de verdade da UI e é passado por valor em todo o loop Elm do Bubbletea. Modificações em campos do `AppState` impactam diretamente:
  * Todos os manipuladores em `internal/ui/app.go` (`Update`).
  * As assinaturas e chamadas em `internal/ui/screens/*.go` e `internal/ui/components/*.go`.
  * As mensagens de comunicação assíncrona geradas em `internal/worker/messages.go`.

### 2.2 Alterações no Pacote `pkg/mining`
* **Gravidade**: 🟡 MÉDIA
* **Justificativa**: Embora seja um pacote matemático puro com acoplamento zero a camadas externas, mudanças em suas estruturas (`Job`) ou assinaturas (`MeetsTarget`) impactam diretamente o loop assíncrono em `internal/worker/miner.go` que calcula os batches.

### 2.3 Alterações no Pacote `internal/config`
* **Gravidade**: 🔴 CRÍTICA
* **Justificativa**: Por controlar a validação estrita no startup e prover caminhos para o Mock Mode (`--mock`), qualquer alteração ou nova chave introduzida afeta diretamente a lógica de bootstrapping em `cmd/tui/main.go` e a verificação de concorrência dos testes.
