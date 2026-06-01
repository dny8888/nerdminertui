# Mapeamento do Legado — pkg/mining

> **Status:** Mapeado da Especificação Alvo  
> **Módulo:** `pkg/mining`  
> **Nível de Confiança:** 🟢 CONFIRMADO (código matemático puro)

Este módulo implementa a lógica matemática e as primitivas criptográficas da mineração de blocos Bitcoin, empacotando o processamento sem efeitos colaterais.

---

## 1. Arquivos Mapeados no Legado

Os seguintes arquivos compõem o módulo `pkg/mining` com base no blueprint:

| Arquivo Alvo | Arquivo de Origem (Legado) | Linhas / Seção no Legado | Descrição |
|---|---|---|---|
| `pkg/mining/hash.go` | [nerdtui-spec.md](file:///home/dny8888/workspace/github/nerdminertui/nerdtui-spec.md) | Lines 330-337 (§5.8) | Duplo SHA-256 (`SHA256d`) usando pacotes nativos da stdlib. |
| `pkg/mining/target.go` | [nerdtui-spec.md](file:///home/dny8888/workspace/github/nerdminertui/nerdtui-spec.md) | Lines 338-340 (§5.8) | Lógica de validação de target e dificuldade. |
| `pkg/mining/job.go` | [nerdtui-spec.md](file:///home/dny8888/workspace/github/nerdminertui/nerdtui-spec.md) | Line 148 (§4) | Definição da struct `Job` (header, target, extranonce, etc.). |

---

## 2. Assinaturas e Componentes Mapeados

* **Duplo SHA-256 (`SHA256d`)**:
  * **Função**: `SHA256d(data []byte) [32]byte`
  * **Função**: `HashHeader(header []byte, nonce uint32) [32]byte`
* **Target e Dificuldade**:
  * **Função**: `MeetsTarget(hash, target [32]byte) bool` — Retorna verdadeiro se o hash Big-Endian for estritamente menor que o target.
  * **Função**: `DifficultyFromHash(hash [32]byte) float64` — Calcula a dificuldade relativa baseada no hash do bloco original em relação ao target do bloco de gênese do Bitcoin.
