package database

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sudo-JP/Load-Manager/backend/internal/config"
	"context"
)

type Database struct {
	Pool *pgxpool.Pool 
}

func DatabaseConnection() (*Database, error) {
	URL, err := config.DatabaseConfig()
	if err != nil {
		return nil, err
	}

	conn, err := pgxpool.New(context.Background(), URL) 
	if err != nil {
		return nil, err
	}

	db := &Database { Pool: conn }

	return db, nil 
}

func (db *Database) Close() {
	db.Pool.Close()
}
