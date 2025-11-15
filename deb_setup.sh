#!/usr/bin/env bash
set -euo pipefail

echo "=== Updating system packages ==="
sudo apt update && sudo apt upgrade -y

echo "=== Installing base tools (ca-certificates, curl, gnupg, git) ==="
sudo apt install -y ca-certificates curl gnupg lsb-release git

echo "=== Setting up Docker apt repository ==="
sudo install -m 0755 -d /etc/apt/keyrings
if [ ! -f /etc/apt/keyrings/docker.gpg ]; then
  curl -fsSL https://download.docker.com/linux/ubuntu/gpg \
    | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
fi
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] \
  https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" \
  | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

echo "=== Installing Docker Engine + Compose plugin ==="
sudo apt update
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

echo "=== Enabling Docker service ==="
sudo systemctl enable --now docker

# Add current user to docker group (so you can run `docker` without sudo)
if ! getent group docker >/dev/null 2>&1; then
  sudo groupadd docker
fi

echo "=== Adding user '$USER' to 'docker' group ==="
sudo usermod -aG docker "$USER"

echo
echo "=== Docker installation complete ==="
docker --version || true
docker compose version || true

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

echo
echo "IMPORTANT: Log out and log back in so group changes take effect (docker group)."
echo "After relogin you should be able to run 'docker ps' without sudo."
