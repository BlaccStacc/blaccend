package db

import (
    "database/sql"
    _ "github.com/jackc/pgx/v5/stdlib" // _ = import only for side-effects = import doar ca sa trigger init() din acel pachet, nu pt functii
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