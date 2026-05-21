# Scripts para Windows (sin make). Uso: .\scripts.ps1 <comando>
param(
    [Parameter(Position = 0)]
    [ValidateSet("run", "build", "test", "migrate-up", "migrate-down", "docker-up", "docker-down", "help")]
    [string]$Command = "help"
)

$ErrorActionPreference = "Stop"
$ProjectRoot = $PSScriptRoot
$DatabaseUrl = "postgres://flaguser:flagpass@host.docker.internal:5432/feature_flags?sslmode=disable"

function Show-Help {
    Write-Host @"
Feature Flag Service - comandos disponibles:

  .\scripts.ps1 run           - Inicia el servidor (go run)
  .\scripts.ps1 build         - Compila bin/feature-flag-service.exe
  .\scripts.ps1 test          - Ejecuta tests
  .\scripts.ps1 migrate-up    - Aplica migraciones (contenedor migrate)
  .\scripts.ps1 migrate-down  - Revierte 1 migracion
  .\scripts.ps1 docker-up     - docker-compose up -d --build
  .\scripts.ps1 docker-down   - docker-compose down

Nota: 'go run' ya aplica migraciones al arrancar; migrate-up es opcional.
"@
}

Push-Location $ProjectRoot
try {
    switch ($Command) {
        "run" {
            go run ./cmd/server
        }
        "build" {
            New-Item -ItemType Directory -Force -Path bin | Out-Null
            go build -o bin/feature-flag-service.exe ./cmd/server
            Write-Host "Binario: bin/feature-flag-service.exe"
        }
        "test" {
            go test ./...
        }
        "migrate-up" {
            docker run --rm `
                -v "${ProjectRoot}/migrations:/migrations" `
                migrate/migrate `
                -path=/migrations `
                -database $DatabaseUrl up
        }
        "migrate-down" {
            docker run --rm `
                -v "${ProjectRoot}/migrations:/migrations" `
                migrate/migrate `
                -path=/migrations `
                -database $DatabaseUrl down 1
        }
        "docker-up" {
            docker-compose up -d --build
            Write-Host "API: http://localhost:8080"
        }
        "docker-down" {
            docker-compose down
        }
        default {
            Show-Help
        }
    }
}
finally {
    Pop-Location
}
