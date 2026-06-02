# Benchmarks de Otimização SHA-256 (NerdMinerTUI)
Data: 02 de Junho de 2026

Este documento detalha os resultados empíricos da pesquisa de otimização da função de *hashing* duplo SHA-256 (SHA256d) utilizada no núcleo de mineração do **NerdMinerTUI**.

## O Problema
Na implementação original, o minerador alocava um novo array de *bytes* para concatenar o cabeçalho do bloco (76 bytes) com o *nonce* atual (4 bytes) e então rodava a rotina completa do `crypto/sha256` em todos os 80 bytes. O algoritmo do SHA-256 processa dados em pedaços (blocos) de 64 bytes. Logo, os primeiros 64 bytes geram um estado intermediário (midstate) que é **idêntico** para qualquer tentativa de *nonce* dentro do mesmo Job da *pool*.

O objetivo do benchmark foi comparar três abordagens para reaproveitar este *midstate* sem recorrer à códigos complexos e inseguros em Assembly (como CGO ou injeção forçada).

## Metodologia
Utilizamos o pacote `testing` e a suíte `go test -bench=. -benchmem` avaliando as seguintes abordagens:

1. **FullAlloc (Baseline)**: A cópia trivial (toda iteração cria um buffer de 80 bytes). Devido à otimização avançada de escape analysis do compilador do Go atual (1.22+), esta abordagem já conseguia atingir *zero-allocation* stack-based, mas sofria penalidade de tempo calculando os primeiros 64 bytes do SHA-256.
2. **MidstateUnsafeReuse**: Utilização de pacotes `unsafe` e `reflect` para injetar artificialmente o cache da struct de digestão não exportada (`digest`) da biblioteca padrão `crypto/sha256` do Go a cada iteração.
3. **MidstateMarshalReuseZeroAlloc**: A abordagem finalmente escolhida. Utiliza a interface `encoding.BinaryMarshaler` e `UnmarshalBinary` nativa da biblioteca `crypto/sha256`, reaproveitando uma única instância segura do `hash.Hash` para cada *worker*.

## Resultados da Bateria

Ambiente de teste: `goos: linux`, `goarch: amd64`, `cpu: Intel(R) Core(TM) i5-7200U CPU @ 2.50GHz`

| Abordagem | Iterações | Tempo por Operação (ns/op) | Alocações (B/op) | Número de Alocações |
| :--- | :--- | :--- | :--- | :--- |
| **FullAlloc** (Original) | `1,591,611` | `802.7 ns/op` | `0 B/op` | `0 allocs/op` |
| **MidstateUnsafeReuse** | `1,657,598` | `712.9 ns/op` | `0 B/op` | `0 allocs/op` |
| **MidstateMarshalReuse** | `2,184,531` | `644.8 ns/op` | `0 B/op` | `0 allocs/op` |

### Análise e Conclusões
- O compilador do Go já é formidável em não gerar lixo na memória (GC pause) para buffers de hash que não escapam, o que impediu os "Gargalos de Memória" teorizados inicialmente (0 allocs/op confirmados).
- A clonagem de ponteiros via `unsafe` apresentou *overhead* de Reflexão a cada iteração, consumindo parte do tempo salvo por omitir os primeiros 64 bytes.
- A função nativa `UnmarshalBinary` possui suporte otimizado em baixo nível no Go para injetar estados arbitrários na instância do *hasher*. Ela nos deu o melhor rendimento puro: **644.8 nanosegundos por hash**, o que equivale a um aumento de aproximadamente **~19.6% no Hashrate individual de cada core de CPU**, sem abrir mão da segurança multiplataforma.

Aliado à execução multi-thread (`runtime.NumCPU()`), as otimizações trazem a perfomance do NerdMinerTUI em Go para uma margem bastante competitiva em comparação à sua versão nativa original em C/C++.
