# Gerador de CRUD Dinâmico em Go

Aplicação web full-stack que gera automaticamente uma interface CRUD completa a partir de um schema JSON.

## 🛠️ Tecnologias

- **Backend:** Go 1.22+ (stdlib `net/http`)
- **Banco de Dados:** MySQL 5.7+
- **Frontend:** HTML, TailwindCSS (CDN), Vanilla JavaScript
- **Driver:** `go-sql-driver/mysql`

## 🚀 Quick Start

### Pré-requisitos
- Go 1.20+
- MySQL 5.7+

### Setup

1. Clone e instale dependências:
```bash
go mod tidy
```

2. Crie o banco de dados:
```bash
CREATE DATABASE crud_app;
```

3. Compile:
```bash
go build -o crud-app .
```

### Executar

**Linux/macOS:**
```bash
DB_HOST=localhost DB_PORT=3306 DB_NAME=crud_app DB_USER=root DB_PSW=root JSON_SCHEMA=schema.json PORT=8080 ./crud-app
```

**Windows:**
```bash
go build -o crud-app.exe .
crud-app.exe --db-host localhost --db-port 3306 --db-user root --db-psw root --db-name crud_app --port 8080 --json-schema schema.json
```

Acesse em: `http://localhost:8080`

## ⚙️ Configuração

A aplicação aceita variáveis de ambiente ou flags:

| Variável | Flag | Padrão | Obrigatório |
|----------|------|--------|-------------|
| `DB_HOST` | `--db-host` | localhost | Não |
| `DB_PORT` | `--db-port` | 3306 | Não |
| `DB_USER` | `--db-user` | - | Sim |
| `DB_PSW` | `--db-psw` | - | Sim |
| `DB_NAME` | `--db-name` | - | Sim |
| `JSON_SCHEMA` | `--json-schema` | - | Sim |
| `PORT` | `--port` | 8080 | Não |

## 📋 Schema JSON

Crie um arquivo `schema.json` definindo a estrutura da tabela:

```json
{
  "table_name": "produtos",
  "fields": [
    { "name": "id", "type": "int", "primary_key": true },
    { "name": "nome", "type": "string", "required": true },
    { "name": "preco", "type": "float", "required": true },
    { "name": "cpf", "type": "string", "mask": "999.999.999-99", "validation": { "type": "cpf" } },
    { "name": "email", "type": "string", "validation": { "type": "email" } }
  ]
}
```

### Tipos de Campo

| Campo | Tipo | Descrição |
|-------|------|-----------|
| `name` | string | Nome da coluna |
| `type` | string | string, int, float, date, text |
| `primary_key` | bool | Define chave primária |
| `required` | bool | Campo obrigatório |
| `mask` | string | Máscara IMask.js (ex: `999.999.999-99`) |
| `validation` | object | Tipo de validação (cpf, cnpj, email, cep, phone) |

### Máscaras

| Símbolo | Significado |
|---------|-------------|
| `9` | Dígito (0-9) |
| `#` | Letra (A-Z, a-z) |
| `*` | Qualquer caractere |

Exemplo: CPF `999.999.999-99`, Telefone `(99) 99999-9999`
