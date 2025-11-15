#!/usr/bin/env bash
set -euo pipefail

echo "=== Updating system (pacman -Syu) ==="
sudo pacman -Syu --noconfirm

echo "=== Installing Docker, Docker Compose, buildx, git ==="
sudo pacman -S --noconfirm docker docker-compose docker-buildx git

echo "=== Enabling and starting Docker service ==="
sudo systemctl enable --now docker.service

# Add current user to docker group
if ! getent group docker >/dev/null 2>&1; then
  sudo groupadd docker
fi

echo "=== Adding user '$USER' to 'docker' group ==="
sudo gpasswd -a "$USER" docker

echo
echo "=== Docker installation complete ==="
docker --version || true
# On Arch, docker-compose v1 is 'docker-compose', but you can still use it
docker-compose --version || true

if [ -f docker-compose.yml ]; then
  echo
  echo "=== Starting project with: docker compose up -d --build (or docker-compose up -d) ==="
  # Prefer the modern CLI if available
  if docker compose version >/dev/null 2>&1; then
    docker compose up -d --build
  else
    docker-compose up -d --build
  fi
else
  echo
  echo "No docker-compose.yml found in current directory."
  echo "Run this script from the project root, or start later with:"
  echo "  docker compose up -d --build"
  echo "or (if only v1 is available):"
  echo "  docker-compose up -d --build"
fi

echo
echo "IMPORTANT: Log out and log back in (or 'newgrp docker') so the docker group applies."
echo "After that, 'docker ps' should work without sudo."
