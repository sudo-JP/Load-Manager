package config 

import (
	"os"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
)

type Config struct {
	POSTGRES_USER 		string 
	POSTGRES_PASSWORD 	string 
	POSTGRES_HOST 		string 
	POSTGRES_PORT 		string 
	POSTGRES_DB 		string 
}
func DatabaseConfig() (string, error) {
    // Load .env file (do this once at startup)
    if err := godotenv.Load(); err != nil {
        return "", fmt.Errorf("error loading .env file: %w", err)
    }

    user := os.Getenv("POSTGRES_USER")
    password := os.Getenv("POSTGRES_PASSWORD")
    host := os.Getenv("POSTGRES_HOST")
    port := os.Getenv("POSTGRES_PORT")
    db := os.Getenv("POSTGRES_DB")

    if user == "" || password == "" || host == "" || port == "" || db == "" {
        return "", errors.New("one or more required POSTGRES_* environment variables are missing")
    }

    url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
        user, password, host, port, db,
    )

    return url, nil
}
