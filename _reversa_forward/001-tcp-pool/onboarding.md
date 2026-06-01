# Onboarding: Testando o Stratum V1 TCP

1. Compile o NerdTUI:
   `make build` ou `go build -o bin/nerdtui ./cmd/tui`

2. Execute o binário sem o modo mock e forneça um endereço BTC válido para o config:
   No `config.yaml` ou via interface (ou flag, caso adicione), coloque um endereço real da rede principal (iniciado com 1, 3 ou bc1).

3. Inicie o app:
   `./bin/nerdtui`

4. Observe a barra de status inferior.
   Ela deve passar rapidamente de "Conectando..." para "Conectado". Caso haja erro no `PoolAddress`, a conexão falhará e a barra exibirá "Desconectado", tentando novamente de forma invisível.

5. O dashboard deve mostrar hashes reais sendo calculados usando `MeetsTarget` em cima do block target da public-pool.io. A mineração pode demorar para achar shares dependendo de quão rápida for a sua máquina localmente vs dificuldade (que para public-pool costuma ser alta).
