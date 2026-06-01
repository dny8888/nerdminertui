# ADR 0004: Estado de UI Imutável por Cópia de Valor

> **Status:** APROVADO  
> **Data:** 2026-05-29  
> **Autor:** Detetive (Retroativo)

---

## 1. Contexto

A TUI roda na thread principal do Bubbletea, contudo, a lógica assíncrona de mineração de CPU e requisições HTTP de pool operam em threads de segundo plano (goroutines concorrentes). Se estes componentes acessassem ou gravassem dados no modelo de dados compartilhado concorrentemente, o programa sofreria de inconsistências críticas e pânicos de corrida de dados (`data races`), exigindo travas de mutex pesadas.

## 2. Decisão

Definimos que a estrutura `AppState` que modela a UI é tratada como um valor imutável. As seguintes convenções de arquitetura são aplicadas de forma não-negociável:
1. Nenhuma goroutine ou componente externo tem ponteiros para o `AppState` interno do `AppModel`.
2. A comunicação de novas medições assíncronas é feita estritamente emitindo structs de mensagens tipadas (`tea.Msg`) que o Bubbletea encaminha para a função centralizada `Update()`.
3. Toda modificação no `AppState` gera uma cópia integral do estado alterado. Exemplo funcional: `(s AppState) WithHashRate(hps float64) AppState` retorna a cópia modificada, sem mutar o receiver original.

## 3. Consequências

* **Positivas**:
  * Eliminação completa de corridas de dados no modelo de dados da interface.
  * Facilidade de testar funções de renderização (`screens` e `components`) enviando variações de estados imutáveis em testes unitários rápidos e determinísticos.
* **Negativas / Desafios**:
  * A cópia de estruturas na memória gera alocações frequentes, contudo o impacto é negligenciável dada a dimensão reduzida do `AppState` e o limite de atualização visual restrito a 1s.
