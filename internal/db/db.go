package db

import (
    "database/sql"
    _ "github.com/jackc/pgx/v5/stdlib"
    "github.com/BlaccStacc/blaccend/internal/config"
)

// establishes db connection
func Connect(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.DBURL)
	if err != nil {
		return nil, err
	}

	//cica ping db
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}