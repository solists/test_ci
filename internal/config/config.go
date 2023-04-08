package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBURI string
}

func GetConfig() *Config {
	dbHost := "postgres-service.postgres.svc.cluster.local"
	dbPort := "5432"
	dbName := os.Getenv("POSTGRES_DB")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPass := os.Getenv("POSTGRES_PASSWORD")

	// Construct the database connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName)
	return &Config{
		DBURI: connStr,
	}
}
