#!/usr/bin/env sh
set -e

echo "Levantando FlagFlow con Docker..."
docker-compose up -d --build

echo "Esperando que el servicio esté listo..."
sleep 8

URL="http://localhost:8080/api/v1/health"

if command -v xdg-open >/dev/null 2>&1; then
  xdg-open "$URL"
elif command -v open >/dev/null 2>&1; then
  open "$URL"
else
  echo "Abrí en el navegador: $URL"
fi

echo "FlagFlow listo en http://localhost:8080"
