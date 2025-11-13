#!/usr/bin/env bash
set -e

# Load environment variables from .env
if [ -f .env ]; then
    export $(grep -v '^#' .env | xargs)
fi

echo "=== Tidying Go modules ==="
go mod tidy

echo "=== Running Go backend ==="
go run ./cmd/server
