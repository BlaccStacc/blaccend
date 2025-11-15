#!/usr/bin/env bash
set -euo pipefail
if [ -f docker-compose.yml ]; then
  echo
  echo "=== Starting project with: docker compose up -d --build ==="
  docker compose up -d --build
else
  echo
  echo "No docker-compose.yml found in current directory."
  echo "Run this script from the project root, or start later with:"
  echo "  docker compose up -d --build"
fi