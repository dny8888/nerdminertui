# Glossário e Regras de Domínio — nerdminertui

> **Status:** Mapeamento de Especificações Greenfield (Design)  
> **Nível de Documentação:** COMPLETO  
> **Gerado pelo Detetive em:** 2026-05-29

Este documento especifica o vocabulário de domínio e as regras de negócio implícitas que guiam a lógica do **NerdTUI**.

---

## 1. Glossário de Domínio

### 1.1 Mineração Solo (Solo Mining)
O processo de tentar encontrar um bloco Bitcoin de forma independente. No NerdTUI, a mineração real de shares é simulada ou conectada a uma pool pública (public-pool.io) que opera no modelo solo, onde o minerador recebe 99% da recompensa caso encontre um bloco.

### 1.2 Hashrate
A velocidade com que o loop de mineração executa a operação criptográfica SHA256d. Medido em hashes por segundo (H/s), kilohashes por segundo (KH/s) ou megahashes por segundo (MH/s).

### 1.3 Share
Um bloco candidato cujo hash resultante é menor do que o target local estabelecido pela pool de mineração. Encontrar shares é uma evidência estatística de que o minerador está operando e contribuindo com poder de processamento.

### 1.4 Target (Alvo)
Um número de 256 bits (32 bytes) que define a dificuldade de mineração. Para que um hash gerado pelo minerador seja considerado válido, o valor numérico do hash (interpretado como um inteiro em formato big-endian) deve ser estritamente menor do que o target.

### 1.5 CPU Throttling (Estrangulamento de CPU)
Lógica que insere micro-sonecas (`time.Sleep`) calculadas dinamicamente no loop de hashing para reduzir o uso efetivo da CPU do terminal do usuário a um target especificado.

---

## 2. Regras de Negócio Implícitas (RN)

### RN-01: Validação do Target Local (pkg/mining)
* **Regra**: Um hash gerado por nonce satisfaz o target se, e somente se, for numericamente menor em comparação byte-a-byte big-endian.
* **Fórmula**: 
  $$\text{MeetsTarget}(\text{hash}, \text{target}) \iff \text{hash} < \text{target}$$
* **Nível de Confiança**: 🟢 CONFIRMADO (especificado no algoritmo fundamental)

### RN-02: Limitação de Configurações (internal/config)
* **Regra**: O programa impede a execução em modo real sem um endereço BTC definido para evitar desperdício de hashrate.
* **Restrição**: Se `MockMining == false` e `BTCAddress == ""`, o validador rejeita a inicialização com erro.
* **Nível de Confiança**: 🟢 CONFIRMADO

### RN-03: Limites do Throttle de CPU (internal/model)
* **Regra**: O target de uso de CPU (`CPUTarget`) deve ser controlado de forma rígida para evitar travamento total da CPU ou inutilidade da mineração.
* **Valores Permitidos**: Limite mínimo de $0.05$ ($5\%$) e limite máximo de $1.00$ ($100\%$), variando em passos fixos de $0.05$ ($5\%$).
* **Nível de Confiança**: 🟢 CONFIRMADO

### RN-04: Cálculo de Hashes por Segundo (internal/worker)
* **Regra**: O hashrate real exibido na UI deve ser calculado a partir dos hashes executados na última janela exata de 1 segundo, impedindo médias acumulativas lentas que distorçam as flutuações rápidas.
* **Nível de Confiança**: 🟢 CONFIRMADO

### RN-05: Imutabilidade do Modelo Bubbletea (internal/model)
* **Regra**: O `AppState` na UI é uma estrutura de valor puro. Não é permitido que nenhuma goroutine assíncrona modifique campos diretamente. Modificações só podem ocorrer via eventos Bubbletea (`tea.Msg`) que resultam em uma nova cópia do modelo de dados.
* **Nível de Confiança**: 🟢 CONFIRMADO
