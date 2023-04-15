package main

import (
	"context"
	"fmt"
	"io"
	"mymod/internal/client/openai"
	"mymod/internal/config"
	"mymod/internal/controller"
	"mymod/internal/middleware"
	"mymod/internal/migrate"
	"mymod/internal/repository"
	"mymod/internal/service"
	"mymod/internal/service/tgservice"
	"mymod/internal/util"
	"mymod/pkg/audit"
	"net"
	"net/http"
	"os"

	"github.com/go-telegram/bot"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/solists/test_ci/pkg/logger"
	v1 "github.com/solists/test_ci/pkg/pb/myapp/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	cfg := config.GetConfig()

	downloader := tgservice.NewDownloader("")
	fileReader, err := downloader.Download(ctx, "voice/file_1.oga")
	if err != nil {
		logger.Errorf("err download voice: %v, user: %v", err, 123)
		return
	}
	file, err := os.CreateTemp("", "voice-*.mp3")
	if err != nil {
		logger.Errorf("err create temp: %v, user: %v", err, 123)
		return
	}

	_, err = io.Copy(file, fileReader)
	if err != nil {
		file.Close()
		logger.Errorf("err copy to file voice: %v, user: %v", err, 123)
		return
	}
	file.Close()

	return

	db, err := sqlx.Connect(config.PostgresDriver, cfg.DBDSN)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %s", err)
	}

	util.MustInit(migrate.Migrate(
		db,
		config.PostgresDriver,
		config.PostgresMigrationsPath,
		false,
	))

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.GrpcPort))
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}
	serveMux := runtime.NewServeMux()
	util.MustInit(v1.RegisterTgServiceHandlerFromEndpoint(ctx, serveMux, lis.Addr().String(),
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}))

	auditLogService := audit.NewAuditService(db)
	repo := repository.NewRepository(db)
	openaiClient := openai.NewClient(cfg)
	ctrl := controller.NewController(repo, cfg, auditLogService, openaiClient)
	serviceImpl := service.NewService(ctrl)
	tgService := tgservice.NewService(repo, ctrl, tgservice.NewDownloader(cfg.TGAPIKey))
	server := grpc.NewServer()
	v1.RegisterTgServiceServer(server, serviceImpl)
	reflection.Register(server)

	opts := []bot.Option{
		bot.WithDefaultHandler(tgService.Handler),
	}

	b, _ := bot.New(cfg.TGAPIKey, opts...)
	if cfg.NeedToRebindWebHook {
		if _, err = b.SetWebhook(ctx, &bot.SetWebhookParams{
			URL:            cfg.WebHookHost,
			AllowedUpdates: []string{"message", "inline_query"},
		}); err != nil {
			logger.Errorf("SetWebhook: %v", err)
		}
	}

	go func() {
		util.MustInit(http.ListenAndServe(":2000", b.WebhookHandler()))
	}()

	go func() {
		logger.Infof("started grpc gateway on %d port", config.GrpcPort)
		util.MustInit(server.Serve(lis))
	}()

	dbgMux := mux.NewRouter()
	dbgMux.Use(middleware.LoggingMiddleware)
	dbgMux.Use(middleware.MetricMiddleware)
	util.ServeGrpcUI(dbgMux)
	util.ServeSwagger(dbgMux)
	dbgMux.Handle("/metrics", promhttp.Handler())

	util.StartMux(dbgMux, config.DbgPort)

	go b.StartWebhook(ctx)

	logger.Infof("started gateway on localhost:%d", config.MainPort)
	util.MustInit(http.ListenAndServe(fmt.Sprintf(":%d", config.MainPort), serveMux))
}
