# NerdMiner TUI

> Minerador Bitcoin solo em Go com interface TUI (Bubbletea).

## Skills Disponíveis

### release

Conduz o processo de release do projeto seguindo as boas práticas do artigo
"Boas práticas de projetos de código aberto com LLM - O Mínimo" (Fabio Akita).

- **Ativação:** `/release`, `release`, `preparar release`, `publicar versão`
- **Skill:** `.agents/skills/release/SKILL.md`
- **Cobre:** CI/CD com GitHub Actions, releases por tag com changelog, binários multi-plataforma

## Estrutura do Projeto

```
cmd/tui/          → Entrypoint da aplicação
internal/
  config/         → Carregamento de configuração (Viper + env vars)
  model/          → Estado imutável da aplicação (AppState)
  store/          → Persistência SQLite
  ui/             → TUI Bubbletea (app.go, screens/, components/)
  worker/         → Pool client Stratum + mineração (fetcher.go, miner.go)
pkg/
  format/         → Formatação de hashrate e valores
  mining/         → Parser Stratum, montagem de bloco, SHA256d
testutil/         → Helpers para testes
```

## Convenções

- **Linguagem:** Go 1.22+
- **Build:** `make build`, `make run`, `make test`
- **Testes:** `go test ./...` com CGO_ENABLED=0
- **Lint:** golangci-lint (`.golangci.yml`)
- **Config:** `~/.nerdtui/config.yaml` ou variáveis de ambiente `NM_*`
