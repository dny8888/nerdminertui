# Mapeamento do Legado — internal/config

> **Status:** Mapeado da Especificação Alvo  
> **Módulo:** `internal/config`  
> **Nível de Confiança:** 🟢 CONFIRMADO

Este módulo gerencia o carregamento de configurações do sistema, suporte a arquivos YAML e sobreposição de variáveis de ambiente usando `viper`.

---

## 1. Arquivos Mapeados no Legado

Os seguintes arquivos compõem o módulo `internal/config` com base no blueprint:

| Arquivo Alvo | Arquivo de Origem (Legado) | Linhas / Seção no Legado | Descrição |
|---|---|---|---|
| `internal/config/config.go` | [nerdtui-spec.md](file:///home/dny8888/workspace/github/nerdminertui/nerdtui-spec.md) | Lines 176-193 (§5.2) | Definição da struct `Config`, carregamento via viper e regras de validação. |

---

## 2. Assinaturas e Componentes Mapeados

* **Método de Carregamento e Validação**:
  * **Função**: `Load() (*Config, error)`
  * **Função**: `(c *Config) Validate() error`
  * **Regras de Validação**:
    * Se `MockMining == false` e `BTCAddress == ""`, falha com erro de endereço BTC obrigatório.
    * Se `CPUTarget` estiver fora do intervalo `[0.05, 1.0]`, falha com erro de limite de CPU.
