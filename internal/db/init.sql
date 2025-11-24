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

-- =========================
-- GARAGE STORAGE TABLES
-- =========================

CREATE TABLE IF NOT EXISTS garage_spaces (
    id          BIGSERIAL PRIMARY KEY,
    name        TEXT NOT NULL,
    description TEXT,
    location    TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS garage_items (
    id          BIGSERIAL PRIMARY KEY,
    space_id    BIGINT NOT NULL REFERENCES garage_spaces(id) ON DELETE CASCADE,
    name        TEXT NOT NULL,
    quantity    INT NOT NULL DEFAULT 1,
    notes       TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
