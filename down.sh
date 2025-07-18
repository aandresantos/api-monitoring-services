#!/bin/bash

set -e

echo "🧹 Finalizando e limpando o ambiente..."

echo "🛑 Derrubando os containers e removendo volumes..."
docker compose -f docker/docker-compose.dev.yml down -v 

# if [ -f .env ]; then
#   echo "🗑️ Removendo arquivo .env..."
#   rm .env
# fi

if [ -f api-monitoring-services ]; then
  echo "🗑️ Removendo binário compilado (api-monitoring-services)..."
  rm api-monitoring-services
fi

echo "✅ Ambiente limpo com sucesso!"
