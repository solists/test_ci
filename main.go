package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"log"
	"mymod/internal/config"
	"net/http"
)

var (
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "myapp_requests_total",
			Help: "Total number of requests received",
		},
		[]string{"method", "host"},
	)
)

func init() {
	prometheus.MustRegister(requestCounter)
}

const (
	dbDriver = "postgres"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	cfg := config.GetConfig()

	db, err := sqlx.Connect(dbDriver, cfg.DBURI)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Run the database migrations using goose.
	if err := goose.Up(db.DB, "./migrations"); err != nil {
		log.Fatalf("failed to apply database migrations: %v", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Request received", zap.String("path", r.URL.Path))
		fmt.Fprintf(w, "Hello, world!")
		requestCounter.With(prometheus.Labels{"method": r.Method, "host": r.Host}).Inc()
	})

	if err = http.ListenAndServe(":8080", nil); err != nil {
		logger.Fatal("exited unexpectedly")
		return
	}
}
