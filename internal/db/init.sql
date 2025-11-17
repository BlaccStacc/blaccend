CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username TEXT NOT NULL,
    email TEXT NOT NULL,
    password_hash TEXT NOT NULL DEFAULT '',
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    verify_token TEXT,
    verify_expires_at TIMESTAMPTZ,
    totp_secret TEXT,
    twofa_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    date_registered TIMESTAMPTZ NOT NULL,
    last_seen TIMESTAMPTZ DEFAULT NULL
);

-- Optional sample users (these won't have usable passwords, just demo data)
INSERT INTO users (username, email)
VALUES
('alice', 'alice@example.com'),
('bob', 'bob@example.com'),
('carol', 'carol@example.com')
ON CONFLICT DO NOTHING;
