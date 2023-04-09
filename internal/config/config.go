package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBDSN  string
	Env    string
	APIKey string
}

const (
	LocalEnv               = "local"
	PostgresDriver         = "postgres"
	PostgresMigrationsPath = "migrations"
	DbgPort                = 8084
	MainPort               = 8080
	GrpcPort               = 8082
)

func GetConfig() *Config {
	dbHost := "postgres-service.postgres.svc.cluster.local"
	dbPort := "5432"
	dbName := os.Getenv("POSTGRES_DB")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPass := os.Getenv("POSTGRES_PASSWORD")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbName)
	connStr = "postgres://bxteam:pricemonitoring@127.0.0.1:5432/postgres?search_path=test_ci&sslmode=disable"
	return &Config{
		DBDSN:  connStr,
		Env:    os.Getenv("ENV"),
		APIKey: os.Getenv("OPENAI_API_KEY"),
	}
}
