package config 

import (
	"os"
	"errors"
	"fmt"
)

type Config struct {
	POSTGRES_USER 		string 
	POSTGRES_PASSWORD 	string 
	POSTGRES_HOST 		string 
	POSTGRES_PORT 		string 
	POSTGRES_DB 		string 
}

func DatabaseConfig() (string, error) {
	POSTGRES_USER := os.Getenv("POSTGRES_USER")
	POSTGRES_PASSWORD := os.Getenv("POSTGRES_PASSWORD")
	POSTGRES_HOST := os.Getenv("POSTGRES_HOST")
	POSTGRES_PORT := os.Getenv("POSTGRES_PORT")
	POSTGRES_DB := os.Getenv("POSTGRES_DB")
	

	if POSTGRES_USER == "" ||
		POSTGRES_PASSWORD == "" ||
		POSTGRES_HOST == "" ||
		POSTGRES_PORT == "" ||
		POSTGRES_DB == "" {
		return "", errors.New("one or more required POSTGRES_* environment variables are missing")
	}

	URL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		POSTGRES_USER,
		POSTGRES_PASSWORD, 
		POSTGRES_HOST,
		POSTGRES_PORT, 
		POSTGRES_DB)

	return URL, nil
}
