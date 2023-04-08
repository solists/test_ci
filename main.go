package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
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

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

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
