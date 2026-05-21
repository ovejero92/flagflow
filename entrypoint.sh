#!/bin/sh
set -e

echo "==> FlagFlow: aplicando migraciones..."
migrate -path=/app/migrations -database "${DATABASE_URL}" up

echo "==> FlagFlow: iniciando servidor en puerto ${PORT:-8080}..."
exec ./feature-flag-service
