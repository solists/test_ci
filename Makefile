PROJECT_NAME=myapp
LOCAL_BIN=$(CURDIR)/bin
BUF_BIN=$(LOCAL_BIN)/buf
PROTOC_GEN_SWAGGER_BIN=$(LOCAL_BIN)/protoc-gen-swagger
PROTOC_GEN_GRPC_GATEWAY=$(LOCAL_BIN)/protoc-gen-grpc-gateway
PROTOC_GEN_GRPC=$(LOCAL_BIN)/protoc-gen-go-grpc
PROTOC_GEN_GO=$(LOCAL_BIN)/protoc-gen-go
MODTOOLS_BIN=$(LOCAL_BIN)/modtools
GOBINDATA_BIN=$(LOCAL_BIN)/go-bindata

.PHONY: run
run:
	make build && ./bin/$(PROJECT_NAME)


.PHONY: mocks
mocks:
	go generate ./...

.PHONY: test
test:
	$(GOENV) go test -v ./...


.PHONY: build
build:
	$(GOENV) CGO_ENABLED=0 go build -v -ldflags "$(LDFLAGS)" -o ./bin/$(PROJECT_NAME) ./main.go

.PHONY:
generate:
	make install-deps
	PATH=$(LOCAL_BIN):$(PATH) $(BUF_BIN) generate -v --path=api/myapp


.PHONY: install-deps
install-deps:
	GOBIN=$(LOCAL_BIN) go get \
            github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
            github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
            google.golang.org/protobuf/cmd/protoc-gen-go \
            github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger \
            github.com/bufbuild/buf/cmd/buf \
            github.com/kannman/modtools \
            github.com/googleapis/googleapis \
            github.com/go-bindata/go-bindata/... \
            google.golang.org/grpc/cmd/protoc-gen-go-grpc \
            && \
    GOBIN=$(LOCAL_BIN) go install \
            github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
            github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
            google.golang.org/protobuf/cmd/protoc-gen-go \
            github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger \
            github.com/bufbuild/buf/cmd/buf \
            github.com/kannman/modtools \
            github.com/go-bindata/go-bindata/... \
            google.golang.org/grpc/cmd/protoc-gen-go-grpc \
            && \
    $(GOENV) go mod download \
    		&& \
    rm -rf vendor.pb && $(MODTOOLS_BIN) vendor '**/*.proto' && mv vendor vendor.pb


