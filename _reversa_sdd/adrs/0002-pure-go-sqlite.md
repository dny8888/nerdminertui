# ADR 0002: Persistência Local via modernc.org/sqlite

> **Status:** APROVADO  
> **Data:** 2026-05-29  
> **Autor:** Detetive (Retroativo)

---

## 1. Contexto

Para prover uma experiência visual fluida na inicialização do NerdTUI, é desejável persistir as estatísticas de hashrate geradas nos últimos 60 segundos. Isso previne que as sparklines gráficas iniciem zeradas ("cold start") após reiniciar a aplicação. O sistema exige a restrição estrita de compilação sem CGO (`CGO_ENABLED=0`) para permitir builds estáticos ultra portáveis e pipelines de CI simplificados.

## 2. Decisão

Adotamos a biblioteca **`modernc.org/sqlite`** como driver do banco de dados SQLite. Diferente de outras implementações comuns como `mattn/go-sqlite3` que dependem de cabeçalhos C compilados pelo GCC local, esta biblioteca é implementada em Go puro, convertida diretamente dos fontes originais em C do SQLite.

Habilitamos o modo WAL (Write-Ahead Logging) para garantir leitura concorrente de alto desempenho das métricas sem travar gravação de novos hashes.

## 3. Consequências

* **Positivas**:
  * Compilação cruzada nativa impecável com `CGO_ENABLED=0` sem necessidade de ferramentas GCC externas no pipeline.
  * O sparkline do Dashboard carrega imediatamente o histórico persistido no startup.
* **Negativas / Desafios**:
  * Performance ligeiramente inferior quando comparada ao driver original compilado com C, embora a diferença seja insignificante no nosso volume de transações métricas.
  * Maior tamanho de binário final devido ao código traduzido de C-para-Go embarcado no executável.
