#!/usr/bin/env bash
set -e

echo "=== Updating system ==="
sudo pacman -Syu --noconfirm

echo "=== Installing PostgreSQL ==="
sudo pacman -S --noconfirm postgresql

echo "=== Initialize PostgreSQL database cluster (if not done yet) ==="
sudo -iu postgres initdb --locale $LANG -D /var/lib/postgres/data || true

echo "=== Starting PostgreSQL ==="
sudo systemctl enable --now postgresql

echo "=== Creating database, user, and table ==="
sudo -i -u postgres psql <<EOF
-- Create database
CREATE DATABASE firstdb;

-- Create user
CREATE USER admin WITH PASSWORD 'admin';

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE firstdb TO admin;

-- Ensure the public schema is owned by admin
\c firstdb
ALTER SCHEMA public OWNER TO admin;

-- Create a simple users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username TEXT NOT NULL,
    email TEXT NOT NULL
);

-- Insert example users
INSERT INTO users (username, email) VALUES
('alice', 'alice@example.com'),
('bob', 'bob@example.com'),
('carol', 'carol@example.com');
\q
EOF

echo "=== Installing Go ==="
sudo pacman -S --noconfirm go

echo "=== Installing git ==="
sudo pacman -S --noconfirm git

echo "=== Setup complete ==="
echo "Database: firstdb"
echo "User: admin / Password: admin"
echo "Example users: alice, bob, carol"
echo "Go version: $(go version)"
