package util

import (
	"context"
	"fmt"
	"net/http"

	"github.com/fullstorydev/grpcui/standalone"
	"github.com/gorilla/mux"
	"github.com/solists/test_ci/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ServeGrpcUI(mux *mux.Router) {
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

func ServeSwagger(mux *mux.Router) {
	mux.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "api/api.swagger.json")
	})
}

func StartMux(mux *mux.Router, port int) {
	go func() {
		logger.Infof("started on localhost:%v", port)
		MustInit(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
	}()
}
