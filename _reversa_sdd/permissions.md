# Matriz de Permissões — nerdminertui

> **Status:** Mapeamento de Especificações Greenfield (Design)  
> **Nível de Documentação:** COMPLETO  
> **Gerado pelo Detetive em:** 2026-05-29

Este documento formaliza o modelo de controle operacional e de privilégios de execução do **NerdTUI**. Por ser uma aplicação de terminal executável local (`Single-User CLI`), o controle de privilégios se concentra nos modos de execução locais (Flags de Compilação e Permissões do SO) e nas capacidades interativas do Operador do terminal.

---

## 1. Perfis / Papéis Operacionais (Roles)

Identificamos três papéis operacionais que interagem com o ciclo de vida do programa:

1. **Operador do Terminal (User)**: O usuário final que executa e interage com a TUI pelo teclado.
2. **Sistema de CI/CD (Pipeline/Automation)**: Sistema automatizado que compila, testa e valida os Quality Gates.
3. **SO e Processo do Kernel (System/Scheduler)**: O sistema operacional local onde o executável roda.

---

## 2. Matriz de Permissões e Capacidades (CLI Access Matrix)

| Recurso / Ação Operacional | Operador do Terminal | Sistema de CI/CD | SO / Processo | Nível de Confiança |
|---|---|---|---|---|
| **Compilar Binário Estático** | Permitido | **Obrigatório (Gates)** | N/A | 🟢 CONFIRMADO |
| **Executar com Mocking (`--mock`)** | Permitido | Permitido (E2E Smoke) | N/A | 🟢 CONFIRMADO |
| **Acessar Socket Stratum (TCP 21496)** | Permitido | Negado (Quarentena/Local) | Permitido | 🟢 CONFIRMADO |
| **Gravar Histórico SQLite no Disco** | Permitido | Negado (Use `:memory:`) | Permitido | 🟢 CONFIRMADO |
| **Alterar CPUTarget (Teclado `+/-`)** | **Permitido (Interativo)** | Negado (Automático) | N/A | 🟢 CONFIRMADO |
| **Rotacionar Visualização (`tab`)** | **Permitido (Interativo)** | Negado (Automático) | N/A | 🟢 CONFIRMADO |
| **Forçar Encerramento (`q` / `ctrl+c`)**| **Permitido (Interativo)** | Permitido | Permitido (SIGINT/SIGKILL) | 🟢 CONFIRMADO |

---

## 3. Segurança e Restrições Físicas

* **Sandbox CGO**: A restrição `CGO_ENABLED=0` impede a injeção de bibliotecas dinâmicas compartilhadas em C, mitigando vetores de ataque comuns de overflow de buffer no terminal.
* **Persistência de Dados**: O banco de dados SQLite local escreve em `~/.nerdtui/` sob as credenciais e permissões normais do usuário local do Linux. O processo não requer privilégios de administrador (`root` / `sudo`) para rodar e deve ativamente recusar caso seja executado como superusuário para mitigar riscos de segurança local.
