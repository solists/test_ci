package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBDSN string
}

func GetConfig() *Config {
	dbHost := "postgres-service.postgres.svc.cluster.local"
	dbPort := "5432"
	dbName := os.Getenv("POSTGRES_DB")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPass := os.Getenv("POSTGRES_PASSWORD")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbName)
	return &Config{
		DBDSN: connStr,
	}
}
