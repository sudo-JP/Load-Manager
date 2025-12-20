package database

import (
    "context"
    "fmt"
    "os"
    "path/filepath"

    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/sudo-JP/Load-Manager/backend/internal/config"
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

    db := &Database{Pool: conn}

    // Run migrations
    if err := db.RunMigrations(); err != nil {
        return nil, fmt.Errorf("migration failed: %w", err)
    }

    return db, nil
}

func (db *Database) Close() {
    db.Pool.Close()
}

func (db *Database) RunMigrations() error {
    // Get the migrations directory path
    migrationsDir := "internal/migrations"

    // Read all SQL files
    files, err := filepath.Glob(filepath.Join(migrationsDir, "*.sql"))
    if err != nil {
        return fmt.Errorf("failed to read migrations directory: %w", err)
    }

    ctx := context.Background()

    // Execute each migration file in order
    for _, file := range files {
        content, err := os.ReadFile(file)
        if err != nil {
            return fmt.Errorf("failed to read migration file %s: %w", file, err)
        }

        _, err = db.Pool.Exec(ctx, string(content))
        if err != nil {
            return fmt.Errorf("failed to execute migration %s: %w", file, err)
        }

    }

    return nil
}
