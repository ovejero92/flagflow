@echo off
echo Levantando FlagFlow con Docker...
docker-compose up -d --build

echo Esperando que el servicio este listo...
timeout /t 8 /nobreak >nul

start http://localhost:8080/api/v1/health
echo FlagFlow listo en http://localhost:8080
