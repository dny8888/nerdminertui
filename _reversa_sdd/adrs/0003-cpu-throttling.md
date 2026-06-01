# ADR 0003: Controle de CPU Reativo com Throttle Matemático

> **Status:** APROVADO  
> **Data:** 2026-05-29  
> **Autor:** Detetive (Retroativo)

---

## 1. Contexto

A execução do loop de mineração SHA256d em Go de forma ininterrupta consome 100% de capacidade do núcleo do processador no qual a goroutine é agendada. Isso causa aquecimento excessivo e consumo de energia em computadores portáteis e servidores locais dos usuários finais do NerdTUI. O sistema requer controle de CPU em tempo real diretamente pelo terminal.

## 2. Decisão

Implementamos um controle dinâmico de throttling de CPU por software na goroutine do minerador (`MinerWorker.Run`). 

O minerador realiza hashes em blocos (batches) discretos de `50.000` iterações. Após cada batch, mede-se o tempo decorrido $\text{workDuration}$. O minerador então força uma suspensão voluntária da thread (`time.Sleep`) calculada pela equação:

$$\text{sleep} = \text{workDuration} \times \frac{1 - P}{P}$$

Onde $P = \text{CPUTarget} \in [0.05, 1.0]$. 

Para confirmar a precisão comportamental do controle de recursos, o minerador calcula a métrica de uso real e a devolve no canal:

$$\text{CPUActual} = \frac{\text{workDuration}}{\text{workDuration} + \text{sleepDuration}}$$

## 3. Consequências

* **Positivas**:
  * O usuário consegue resfriar sua CPU imediatamente diminuindo o target via teclado (`+` / `-`).
  * Sem necessidade de privilégios de root para controlar threads nativas do SO.
* **Negativas / Desafios**:
  * A precisão do sleep depende da granularidade do escalonador de tempo do Kernel do Sistema Operacional local (podendo haver divergência de até $\pm 5\%$ do target em alguns sistemas Windows/WSL).
