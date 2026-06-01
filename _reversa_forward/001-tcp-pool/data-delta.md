# Delta no Modelo de Dados

## Modificações na Configuração (internal/config/config.go)

A `struct Config` será expandida para conter os seguintes campos adicionais obrigatórios para a comunicação real TCP:

- `[NOVO] PoolAddress string`: O endereço e porta da pool de mineração. Usará a flag map do viper.
- `[NOVO] WorkerName string`: O identificador da máquina ou rig conectada.

### Comportamento (Fallback)
Ao executar `config.Load()`, caso esses valores estejam em branco (e não estejamos usando `--mock`), o sistema preencherá os defaults combinados:
- `PoolAddress` será inicializado como `"public-pool.io:21496"`.
- `WorkerName` será inicializado como `".nerdtui"`.

Não há mudanças em banco de dados SQLite.
