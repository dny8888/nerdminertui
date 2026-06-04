# NerdMiner TUI

[![CI](https://github.com/dny8888/nerdminertui/actions/workflows/ci.yml/badge.svg)](https://github.com/dny8888/nerdminertui/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/dny8888/nerdminertui)](https://github.com/dny8888/nerdminertui/releases/latest)

> Um minerador Bitcoin solo otimizado escrito em Go, equipado com uma bela interface de terminal (TUI) interativa e responsiva desenvolvida com Bubbletea.

## Principais Funcionalidades

- ⛏️ **Mining Engine de Alta Performance:** Mineração otimizada sem alocações em tempo real usando `sync.Pool` e Double-SHA256 pré-calculados.
- 💻 **TUI Responsiva & Dinâmica:** Layout interativo utilizando grid dinâmico com navegação e abas.
- 📊 **Dashboards & Estatísticas:** Monitoramento de hashrate e odds em tempo real com gráfico em série temporal persistido em banco SQLite.
- 🌍 **Internacionalizado:** Suporte inicial focado em localização e i18n nativo.
- 🛡️ **Tolerante a Falhas:** Retries automáticos e sistema de toast (alertas) embutido na interface.

## Instalação

### go install (requer Go 1.22+)
```bash
go install github.com/dny8888/nerdminertui/cmd/tui@latest
```

### Download Direto
Você pode encontrar binários para as principais plataformas (Linux, macOS, Windows, ARM) na página de [Releases](https://github.com/dny8888/nerdminertui/releases/latest).

```bash
# Exemplo para Linux AMD64
curl -fsSL https://github.com/dny8888/nerdminertui/releases/latest/download/nerdminertui-linux-amd64.tar.gz | tar xz
sudo mv nerdminertui /usr/local/bin/
```

## Utilização

Simplesmente execute o comando na sua linha de comando:

```bash
nerdminertui
```

### Argumentos de Linha de Comando (Flags)
- `--config <caminho>`: Passa um caminho de configuração YAML alternativo em vez do default em `~/.nerdtui/config.yaml`.
- `--mock`: Inicia a mineração em "Mock Mode" (não é necessário fornecer endereço de Bitcoin). A interface simulará o comportamento do minerador.
- `--no-store`: Inicia a aplicação sem armazenar histórico no SQLite local (utilizando NilStore).

## Atalhos de Teclado (Global)
- **`?`** - Mostra/Oculta painel de Ajuda
- **`q` ou `Ctrl+C`** - Sai do aplicativo
- **`Tab`** - Alterna entre as abas (Dashboard > Global Stats > Settings)
- **`Shift+Tab`** - Alterna na ordem reversa (útil em Settings)
- **`Ctrl+S`** - Salva as configurações de Settings no disco
- **`+ / =`** - Aumenta o uso do limite de CPU em 5%
- **`-`** - Diminui o uso de limite de CPU em 5%
