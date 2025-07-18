#!/bin/bash

set -e

echo "🧩 Inicializando ambiente do projeto..."

echo "📦 Baixando dependências com go mod download..."
go mod download

if ! command -v migrate &> /dev/null
then
    echo "🔧 Instalando migrate CLI..."
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
    sudo mv migrate /usr/local/bin/
    echo "✅ migrate instalado com sucesso."
else
    echo "✅ migrate já está instalado."
fi

echo "🔐 Gerando arquivo .env com variáveis de ambiente..."

cat <<EOF > .env
DATABASE_URL="postgres://admin:123@localhost:5432/monitoringdb?sslmode=disable"
POSTGRES_USER=admin
POSTGRES_PASSWORD=123
POSTGRES_DB=monitoringdb
EOF

echo "✅ Arquivo .env criado."

echo "🐘 Subindo o PostgreSQL com docker-compose..."
docker compose -f ./docker/docker-compose.dev.yml up -d postgres

echo "🗃️ Executando as migrations com make migrate-up..."
make migrate-up DATABASE_URL="postgres://admin:123@localhost:5432/monitoringdb?sslmode=disable"

echo "🎉 Ambiente inicializado com sucesso!"
