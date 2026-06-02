# Análise de Otimização: NerdMiner C++ vs NerdMinerTUI Go
Data: 02 de Junho de 2026

Ao analisar o código fonte original em C++ do projeto NerdMiner (especialmente `mining.cpp` e `stratum.cpp`), foram identificados 5 pontos críticos de melhoria para o `NerdMinerTUI` escrito em Go. A aplicação destas técnicas pode aumentar exponencialmente a taxa de hashes por segundo (hashrate) e a estabilidade geral do minerador.

## 1. Zero-Allocation (Redução de Garbage Collection)
No pacote `pkg/mining/hash.go`, a função `HashHeader` estava alocando um novo slice dinamicamente a cada cálculo de hash:
```go
payload := make([]byte, len(header)+4)
copy(payload, header)
```
**Impacto:** Em altas taxas de hash (ex: 100 kH/s), isso significa alocar megabytes de memória por segundo que são imediatamente descartados. Isso sobrecarrega enormemente o Garbage Collector do Go, derrubando o hashrate para limpar a memória.
**Solução:** O NerdMiner sobrescreve o nonce diretamente no mesmo buffer existente na memória (`((uint32_t*)(sha_buffer+64+12))[0] = nonce;`). Em Go, passaremos a reaproveitar um único array estático `[80]byte` por worker, modificando apenas os 4 bytes finais a cada tentativa.

## 2. Otimização de SHA-256 (Midstate Baking)
**Impacto:** O Go processava os 80 bytes inteiros do zero via `sha256.Sum256` em cada iteração de nonce. Como o algoritmo SHA-256 ingere blocos de 64 bytes, 80 bytes exigem o processamento do primeiro bloco de 64, do segundo com os 16 restantes (padding), e de um último bloco para o hash duplo (SHA256d).
**Solução:** Sabemos que os primeiros 64 bytes do block header nunca mudam dentro de um mesmo "Job". O C++ pré-calcula (`nerd_mids` e `nerd_sha256_bake`) o estado interno (midstate) do SHA-256 para esses 64 bytes uma única vez. No loop principal, só processa do 65º byte em diante. Replicar esse Midstate Caching em Go (escrevendo um SHA256 interno que aceite estado salvo ou importando código dedicado) reduzirá o processamento de hashing em 50%, praticamente dobrando o hashrate.

## 3. Nonce Inicial Aleatório (Anti-Colisão)
**Impacto:** Se múltiplos usuários de NerdMinerTUI entrarem em uma mesma pool, frequentemente receberão o mesmo Job ID. Iniciar o nonce no `0` como antes faz com que todos computem de forma idêntica e concorram uns com os outros (trabalho duplicado).
**Solução:** O NerdMiner implementa um `RandomGet()` para definir em qual nonce o worker vai iniciar a busca. Em Go, podemos obter a mesma vantagem iniciando a varredura com um nonce pseudoaleatório sempre que um novo Job for recebido.

## 4. Concorrência Plena (Multi-Goroutines)
**Impacto:** O código anterior (`miner.go`) continha apenas um worker na goroutine principal realizando batches de 50.000 hashes.
**Solução:** O Go foi feito para paralelismo brutal. Ao ler o número de threads disponíveis com `runtime.NumCPU()`, podemos instanciar múltiplos workers simultâneos, onde cada um verifica um intervalo de nonces distinto. O hashrate será instantaneamente escalado pelo número de núcleos da máquina.

## 5. Resiliência de Conexão Stratum (Keep-Alive Ativo)
**Impacto:** Problemas de inatividade do socket não eram detectados corretamente se a conexão TCP não caísse no nível do SO (TCP Keep-Alive ineficiente contra freezes de aplicação na pool).
**Solução:** O NerdMiner ativamente envia uma mensagem Stratum `mining.suggest_difficulty` e reinicia o loop caso os dados parem de vir por mais de 5 minutos. Em Go, o Stratum client ganhará rotinas de ping e inatividade (read timeouts) para forçar o fail-fast e garantir novos jobs de modo perpétuo.
