# ADR 0001: Adoção do Bubbletea + Lipgloss para Terminal UI

> **Status:** APROVADO  
> **Data:** 2026-05-29  
> **Autor:** Detetive (Retroativo)

---

## 1. Contexto

Para o desenvolvimento do **NerdTUI**, que é um dashboard de mineração solo do Bitcoin projetado para rodar diretamente em ambientes UNIX e terminais locais, precisamos de um framework que ofereça:
* Renderização limpa, responsiva e performática.
* Arquitetura organizada para gerenciamento de eventos concorrentes de rede e hashing.
* Suporte robusto a keybindings.
* Design de componentes visuais estilizados.

## 2. Decisão

Adotamos a biblioteca de código aberto **Bubbletea** (`github.com/charmbracelet/bubbletea`) para a arquitetura de interface terminal e o **Lipgloss** (`github.com/charmbracelet/lipgloss`) para a modelagem visual das telas.

O Bubbletea adota o paradigma **Model-Update-View (MUV)** inspirado na linguagem Elm, o que garante excelente organização e previsibilidade do estado visual do programa.

## 3. Consequências

* **Positivas**:
  * Loop de eventos altamente determinístico e estruturado.
  * Facilidade de realizar testes de integração visual assíncronos usando a biblioteca de mock de terminal `teatest`.
  * Separação estrita de responsabilidades: rendering de tela é função pura; interações geram mensagens no canal.
* **Negativas / Desafios**:
  * Exige que os desenvolvedores dominem o fluxo Elm (MUV) e o gerenciamento de goroutines em segundo plano via `tea.Cmd`.
