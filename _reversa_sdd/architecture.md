# Arquitetura Geral do Sistema — nerdminertui

> **Status:** Mapeamento de Especificações Greenfield (Design)  
> **Nível de Documentação:** COMPLETO  
> **Gerado pelo Arquiteto em:** 2026-05-29

Este documento sintetiza a arquitetura de sistemas, topologia de containers, integrações externas e modelo de dados do **NerdTUI**.

---

## 1. Visão Geral da Arquitetura

O **NerdTUI** é um aplicativo de terminal autônomo de mineração Bitcoin solo escrito em **Go 1.23**. Ele é projetado como um executável binário estático e portável de dependência zero (`CGO_ENABLED=0`), otimizado para rodar em qualquer sistema Unix sem requerer bibliotecas C compartilhadas ou compiladores externos locais.

A arquitetura do sistema adota uma estrutura em camadas funcionais, isolando rigorosamente a lógica matemática, processamento assíncrono de hardware/rede e a camada visual Elm/MUV (`Model-Update-View` do framework Bubbletea).

---

## 2. Topologia de Integrações Externas

O aplicativo funciona de forma híbrida: minera localmente e consome dados de redes públicas de mineração Bitcoin solo:

```
┌─────────────────┐        TCP Stratum (JSON-RPC)       ┌───────────────────────────┐
│                 │ ◄─────────────────────────────────► │ Stratum Server            │
│                 │                                     │ (ex: public-pool.io:21496)│
│                 │                                     └───────────────────────────┘
│   NerdTUI CLI   │
│                 │            HTTP REST                ┌───────────────────────────┐
│                 │ ◄─────────────────────────────────► │ Public Pool REST API      │
│                 │                                     │ (Estatísticas Globais)    │
└─────────────────┘                                     └───────────────────────────┘
```

### Protocolos de Integração:
1. **JSON-RPC 1.0 via TCP Stratum (Porta 21496)**: Handshake e polling de novos `mining.Job` enviados pela pool Bitcoin.
2. **HTTP/1.1 REST (API REST pública)**: Consulta a cada 5s de estatísticas globais da rede (dificuldade estimada, altura do bloco e taxa total de hash da rede).
3. **Mock Mode (Simulador)**: Permite mineração e geração de estatísticas simuladas localmente (flag `--mock`), dispensando acessos à rede externa.

---

## 3. Estrutura de Camadas e Limitações de Acesso

```
┌─────────────────────────────────────────────────────────┐
│                       cmd/tui                           │  <-- Wiring e Parsing
└───────────┬─────────────────────────┬───────────────────┘
            │ Imports                 │ Imports
┌───────────▼─────────────┐ ┌─────────▼─────────────┐
│       internal/ui       │ │    internal/worker    │  <-- Processamento Assíncrono
└───────────┬─────────────┘ └─────────┬─────────────┘
            │                         │ Imports
            ├─────────────────────────┤
            │ Imports                 │
┌───────────▼─────────────┐ ┌─────────▼─────────────┐
│     internal/model      │ │    internal/store     │  <-- Estado e Banco de Dados
└─────────────────────────┘ └───────────────────────┘
            │                         │
            └───────────┬─────────────┘
                        │ Imports
            ┌───────────▼─────────────┐
            │          pkg/           │  <-- Código Puro / Utilitários (Sem I/O)
            └─────────────────────────┘
```

### Regras Estritas de Dependência:
* O pacote `pkg/` é totalmente agnóstico ao resto do sistema. **Não pode importar nada de `internal/` ou `cmd/`**.
* O diretório `internal/` encapsula lógica interna e banco. **Não pode importar nada de `cmd/`**.
* `cmd/tui` é a raiz de wiring e pode importar qualquer módulo interno.

---

## 4. Dívidas Técnicas e Riscos Mapeados

1. **Scheduling de Micro-Sleeps**: O throttle de CPU por sonecas de threads no Go (`time.Sleep`) depende diretamente do escalonador do SO. Pode haver flutuações e discrepâncias pontuais em máquinas virtuais ou WSL.
2. **I/O Monotônico na SQLite**: Gravações frequentes a cada 1s de hashrate podem causar fadiga de disco ao longo de meses em mídias Flash mais baratas se o local do arquivo de métricas `~/.nerdtui/metrics.db` não for montado em partições temporárias em memória RAM (`tmpfs`). O uso do pragma **WAL** e busy timeout de 5000ms mitiga o risco de contenção de escrita concorrente, mas monitoramento é aconselhado.
