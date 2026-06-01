# Mapeamento do Legado — pkg/format

> **Status:** Mapeado da Especificação Alvo  
> **Módulo:** `pkg/format`  
> **Nível de Confiança:** 🟢 CONFIRMADO (funções de formatação puras)

Este módulo engloba funções puras de formatação de string de dados de mineração e tempo de execução (uptime) do terminal para o display gráfico.

---

## 1. Arquivos Mapeados no Legado

Os seguintes arquivos compõem o módulo `pkg/format` com base no blueprint:

| Arquivo Alvo | Arquivo de Origem (Legado) | Linhas / Seção no Legado | Descrição |
|---|---|---|---|
| `pkg/format/hashrate.go` | [nerdtui-spec.md](file:///home/dny8888/workspace/github/nerdminertui/nerdtui-spec.md) | Line 349 (§5.9) | Converte HPS em string legível (ex: "12.4 KH/s"). |
| `pkg/format/duration.go` | [nerdtui-spec.md](file:///home/dny8888/workspace/github/nerdminertui/nerdtui-spec.md) | Line 350 (§5.9) | Converte tempo de atividade em string (ex: "2d 03h 14m"). |
| `pkg/format/difficulty.go` | [nerdtui-spec.md](file:///home/dny8888/workspace/github/nerdminertui/nerdtui-spec.md) | Line 351 (§5.9) | Converte floats de dificuldade científica para notação legível. |

---

## 2. Assinaturas e Componentes Mapeados

* **Formatação de Hashrate (`FormatHashRate`)**:
  * **Função**: `FormatHashRate(hps float64) string`
  * **Exemplo**: `999.0` -> `"999 H/s"`, `1000.0` -> `"1.0 KH/s"`, `12400.0` -> `"12.4 KH/s"`.
* **Formatação de Uptime (`FormatUptime`)**:
  * **Função**: `FormatUptime(d time.Duration) string`
  * **Exemplo**: `time.Duration(42 * time.Second)` -> `"0m 42s"`, `time.Duration(49 * time.Hour)` -> `"2d 01h 00m"`.
* **Formatação de Dificuldade (`FormatDifficulty`)**:
  * **Função**: `FormatDifficulty(d float64) string`
* **Formatação de Bloco (`FormatBlockHeight`)**:
  * **Função**: `FormatBlockHeight(h uint32) string`
  * **Exemplo**: `892441` -> `"#892.441"`.
