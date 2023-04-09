package main

import (
	"context"
	"mymod/internal/config"
	"mymod/internal/controller"
	"mymod/internal/migrate"
	"mymod/internal/repository"
	"mymod/internal/service"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/solists/test_ci/pkg/logger"
	v1 "github.com/solists/test_ci/pkg/pb/myapp/v1"
	httpSwagger "github.com/swaggo/http-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	cfg := config.GetConfig()

	db, err := sqlx.Connect(config.PostgresDriver, cfg.DBDSN)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %s", err)
	}

	mustInit(migrate.Migrate(
		cfg,
		db,
		config.PostgresDriver,
		config.PostgresMigrationsPath,
		false,
	))

	repo := repository.NewRepository(db)
	ctrl := controller.NewController(repo, cfg)
	serviceImpl := service.NewService(ctrl)
	server := grpc.NewServer()
	v1.RegisterCalculatorServer(server, serviceImpl)

	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}

	serveMux := runtime.NewServeMux()
	mustInit(v1.RegisterCalculatorHandlerFromEndpoint(ctx, serveMux, lis.Addr().String(),
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}))

	dbgMux := mux.NewRouter()
	serveSwagger(dbgMux)

	go func() {
		mustInit(http.ListenAndServe(":8084", dbgMux))
		logger.Info("started gateway on localhost:8084")
	}()

	mustInit(server.Serve(lis))
	logger.Info("started grpc gateway on 8082 port")
	mustInit(http.ListenAndServe(":8080", serveMux))
	logger.Info("started gateway on localhost:8080")
}

func serveSwagger(mux *mux.Router) {
	mux.Handle("/swagger", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8084/swagger.json"), // URL to your swagger.json
	))
	mux.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "api/api.swagger.json")
	})
}

func mustInit(err error) {
	if err != nil {
		logger.Fatalf("init failure: %s", err)
	}
}
