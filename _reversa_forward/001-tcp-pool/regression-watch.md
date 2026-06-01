# Regression Watch: 001-tcp-pool

| ID | Origem (arquivo, seção) | Regra esperada após mudança | Tipo de verificação | Sinal de violação |
|----|-------------------------|-----------------------------|---------------------|-------------------|
| W001 | `domain.md` (RN-02) | A configuração de rede usa `PoolAddress` e `WorkerName` (não mais `PoolURL`). | ausência | Referência a `PoolURL` no pacote `config` ou UI. |
| W002 | `architecture.md` (Worker) | A integração Stratum opera em loop assíncrono mantendo a conexão TCP em vida (`Run`). | presença | `Fetcher` implementando loops periódicos curtos (HTTP style) invés de TCP longo. |
| W003 | `architecture.md` (State) | O MUV Pattern deve continuar intocável, sem que os retornos do Stratum interfiram no estado. | preservação | Uso de mutex ou estado não puro dentro da struct `AppState`. |

## Histórico de re-extrações
<!-- Preenchido automaticamente nas re-extrações do Reversa -->

## Arquivadas
