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

	"github.com/fullstorydev/grpcui/standalone"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/solists/test_ci/pkg/logger"
	v1 "github.com/solists/test_ci/pkg/pb/myapp/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

var (
	requestCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "myapp_requests_total",
			Help: "Total number of requests received",
		},
		[]string{"method", "host"},
	)
)

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
	reflection.Register(server)

	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}

	serveMux := runtime.NewServeMux()
	mustInit(v1.RegisterCalculatorHandlerFromEndpoint(ctx, serveMux, lis.Addr().String(),
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}))

	go func() {
		logger.Info("started grpc gateway on 8082 port")
		mustInit(server.Serve(lis))
	}()

	dbgMux := mux.NewRouter()
	dbgMux.Use(loggingMiddleware)
	dbgMux.Use(metricMiddleware)
	serveGrpcUI(dbgMux)
	serveSwagger(dbgMux)
	dbgMux.Handle("/metrics", promhttp.Handler())

	go func() {
		logger.Info("started gateway on localhost:8084")
		mustInit(http.ListenAndServe(":8084", dbgMux))
	}()

	logger.Info("started gateway on localhost:8080")
	mustInit(http.ListenAndServe(":8080", serveMux))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Infof("Request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func metricMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCounter.With(prometheus.Labels{"method": r.Method, "host": r.Host}).Inc()
		next.ServeHTTP(w, r)
	})
}

func serveGrpcUI(mux *mux.Router) {
	conn, err := grpc.Dial(":8082",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	uiHandler, err := standalone.HandlerViaReflection(context.Background(), conn, "/grpcui")
	if err != nil {
		panic(err)
	}

	mux.PathPrefix("/grpcui/").Handler(http.StripPrefix("/grpcui", uiHandler))
}

func serveSwagger(mux *mux.Router) {
	mux.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "api/api.swagger.json")
	})
}

func mustInit(err error) {
	if err != nil {
		logger.Fatalf("init failure: %s", err)
	}
}
