#!/usr/bin/env bash
set -e

echo "=== Updating system ==="
sudo apt update && sudo apt upgrade -y

echo "=== Installing PostgreSQL ==="
sudo apt install -y postgresql postgresql-contrib

echo "=== Starting PostgreSQL ==="
sudo systemctl enable --now postgresql

echo "=== Creating database and user ==="
sudo -i -u postgres psql <<EOF
-- Create database
CREATE DATABASE firstdb;

-- Create user
CREATE USER admin WITH PASSWORD 'admin';

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE firstdb TO admin;

-- Ensure the public schema is owned by admin
ALTER SCHEMA public OWNER TO admin;

-- Create a simple users table
\c firstdb
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
sudo apt install -y golang

echo "=== Installing git ==="
sudo apt install -y git

echo "=== Setup complete ==="
echo "Database: firstdb"
echo "User: admin / Password: admin"
echo "Example users: alice, bob, carol"
echo "Go version: $(go version)"
