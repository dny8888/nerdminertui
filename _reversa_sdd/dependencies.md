# Dependências do Projeto — nerdminertui

> **Gerado pelo Scout em:** 2026-05-29

Este documento lista as tecnologias básicas, frameworks, bibliotecas e dependências de desenvolvimento mapeados a partir da especificação `nerdtui-spec.md` e do arquivo `AGENTS.md`.

---

## 1. Ambiente de Execução Principal

| Requisito | Detalhe / Versão | Origem da Informação | Notas |
|---|---|---|---|
| **Linguagem** | Go 1.23+ | `nerdtui-spec.md` (§3, §12) | Requisito fundamental de runtime |
| **Compilação** | `CGO_ENABLED=0` | `nerdtui-spec.md` (Restrição R1) | Geração de binários 100% estáticos sem GCC |

---

## 2. Bibliotecas de Runtime (Dependencies)

Mapeadas a partir do ecossistema especificado para o **NerdTUI**:

| Dependência | Categoria | Propósito | Notas |
|---|---|---|---|
| `github.com/charmbracelet/bubbletea` | TUI Framework | Arquitetura MUV (Model-Update-View) do terminal | Componente central da interface |
| `github.com/charmbracelet/lipgloss` | Styling UI | Definição de cores, bordas, layouts e estilos de terminal | Utilizado para renderização e temas |
| `github.com/charmbracelet/bubbles` | UI Components | Componentes interativos prontos (ex: spinners, statusbar) | Agiliza o desenvolvimento de elementos padrão |
| `github.com/spf13/viper` | Config Manager | Carregamento de arquivos YAML/JSON e sobreposição de Env Vars | Usado em `internal/config` |
| `modernc.org/sqlite` | Banco de Dados | Driver SQLite implementado em Go puro (sem CGO) | Compatível com a restrição `CGO_ENABLED=0` |

---

## 3. Ferramentas e Dependências de Testes

| Dependência / Ferramenta | Tipo | Propósito | Notas |
|---|---|---|---|
| `github.com/stretchr/testify` | Test Assertions | Asserções fluídas e limpas nos testes unitários | Especificado como `testify/assert` |
| `go.uber.org/goleak` | Leak Detector | Evitar goroutine leaks no loop do `MinerWorker` | Rodado automaticamente em todo `TestMain` |
| `github.com/charmbracelet/x/teatest` | TUI Integration | Testar assincronamente os loops e updates do Bubbletea | Utilizado para testes no `internal/ui/app` |
| `httptest` | Mock HTTP | Testar o `HTTPPoolClient` da pool | Parte da biblioteca padrão do Go |

---

## 4. Ferramentas de Linting & Quality Gates (CI/CD)

As ferramentas necessárias no pipeline de integração contínua (CI):

| Ferramenta | Propósito | Configuração | Severidade |
|---|---|---|---|
| `golangci-lint` (v1.59+) | Análise estática do código | `.golangci.yml` na raiz | Bloqueante |
| `govulncheck` | Verificação de falhas de segurança conhecidas nas dependências | Integrado via Makefile | Bloqueante (zero High/Critical) |
| `gocyclo` | Validação de complexidade ciclomática | Threshold: máximo de 15 por função | Bloqueante |
| `goleak` | Prevenção de goroutine leaks de produção | Rodado em testes unitários | Bloqueante |
