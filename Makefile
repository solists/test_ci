PROJECT_NAME=myapp


.PHONY: run-local
run-local:
	make build && ./bin/$(PROJECT_NAME)


.PHONY: mocks
mocks:
	go generate ./...


.PHONY: build
build:
	$(GOENV) CGO_ENABLED=0 go build -v -ldflags "$(LDFLAGS)" -o ./bin/$(PROJECT_NAME) ./main.go

