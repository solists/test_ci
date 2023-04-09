module mymod

go 1.20

require (
	github.com/fullstorydev/grpcui v1.3.1
	github.com/golang/mock v1.4.4
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.15.2
	github.com/jmoiron/sqlx v1.3.5
	github.com/json-iterator/go v1.1.12
	github.com/lib/pq v1.10.7
	github.com/opentracing/opentracing-go v1.2.0
	github.com/ory/dockertest/v3 v3.9.1
	github.com/pressly/goose/v3 v3.10.0
	github.com/prometheus/client_golang v1.14.0
	github.com/sashabaranov/go-openai v1.7.0
	github.com/solists/test_ci/pkg/logger v0.0.0-00010101000000-000000000000
	github.com/solists/test_ci/pkg/pb/myapp v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.8.2
	google.golang.org/grpc v1.54.0
)

require (
	github.com/Nvveen/Gotty v0.0.0-20120604004816-cd527374f1e5 // indirect
	github.com/benbjohnson/clock v1.3.0 // indirect
	github.com/cenkalti/backoff/v4 v4.2.0 // indirect
	github.com/containerd/continuity v0.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/imdario/mergo v0.3.13 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/opencontainers/runc v1.1.4 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v1.2.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)

require (
	github.com/Azure/go-ansiterm v0.0.0-20230124172434-306776ec8161 // indirect
	github.com/Microsoft/go-winio v0.6.0 // indirect
	github.com/docker/cli v23.0.2+incompatible // indirect
	github.com/docker/docker v23.0.2+incompatible // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/moby/term v0.0.0-20221205130635-1aeaba878587 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0-rc2 // indirect
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.9.0 // indirect
	golang.org/x/mod v0.9.0 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/tools v0.7.0 // indirect
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bufbuild/protocompile v0.5.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/fullstorydev/grpcurl v1.8.6 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/jhump/protoreflect v1.15.1 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.37.0 // indirect
	github.com/prometheus/procfs v0.9.0 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.24.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	google.golang.org/genproto v0.0.0-20230403163135-c38d8f061ccd // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/solists/test_ci/pkg/batchprocessor => ./pkg/batchprocessor

replace github.com/solists/test_ci/pkg/logger => ./pkg/logger

replace github.com/solists/test_ci/pkg/pb/myapp => ./pkg/pb/myapp
