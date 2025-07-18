# 🧠 API Monitoring Services

Projeto de backend em Go para monitoramento de Serviços e APIs, você registra seu serviço nessa API e a mesma fica monitorando o status da sua API.

## Próximos passos

Normalmente as pessoas costumas achar que health check é apenas ver se está retornando 200 (por ora essa API está apenas fazendo o mesmo) mas o que pode acontecer é o serviço estar retornando 200 e mesmo assim pode estar apresentando problemas que ficam invisíveis. A ideia do projeto é tentar deixar mais nítido a saúde real da aplicação.

---

## ✅ Requisitos

Antes de iniciar, certifique-se de que você possui instalado:

- [Go 1.24+](https://golang.org/dl/)
- [Docker + Docker Compose](https://docs.docker.com/get-docker/)
- Um cliente PostgreSQL (ex: [TablePlus](https://tableplus.com/), [pgAdmin](https://www.pgadmin.org/), [psql CLI](https://www.postgresql.org/docs/current/app-psql.html)) — para visualizar/manipular os dados
- `make` (geralmente já incluso em distros Unix/Linux e Mac)
- `migrate` — será instalado automaticamente pelo script `init.sh`

---

## ⚙️ Variáveis de Ambiente

O projeto utiliza o seguinte `.env` (gerado automaticamente na inicialização):

```env
DATABASE_URL="postgres://admin:123@localhost:5432/monitoringdb?sslmode=disable"
POSTGRES_USER=admin
POSTGRES_PASSWORD=123
POSTGRES_DB=monitoringdb
```

---

## 🚀 Como iniciar o projeto localmente

1. Clone o repositório:

   ```bash
   git clone git@github.com:aandresantos/api-monitoring-services.git
   cd api-monitoring-services
   ```

2. Execute o script de inicialização:

   ```bash
   ./init.sh
   ```

   Este script irá:

   - Baixar os módulos Go (`go mod download`)
   - Instalar o CLI do `migrate` (se necessário)
   - Criar o arquivo `.env`
   - Subir o banco de dados PostgreSQL via Docker Compose
   - Executar as migrações com `make migrate-up`

---

## 🗃️ Migrações

As migrações estão localizadas em:  
`./internal/database/migrations`

### Comandos úteis:

```bash
make migrate-new name=[nome_da_migration]

# Subir todas as migrations
make migrate-up DATABASE_URL=[a url do banco]

# Descer uma migration
make migrate-down DATABASE_URL=[a url do banco]
```

---

## 🧪 Testando a conexão com o banco

Após rodar o `init.sh`, o banco estará acessível em:

```
Host: localhost
Porta: 5432
Usuário: admin
Senha: 123
Database: monitoringdb
```

--

## Encerrando o projeto

Você pode utilizar `down.sh` que irá remover todos os containers e todo o ambiente que estiver rodando do projeto.
