package util

import "github.com/solists/test_ci/pkg/logger"

func MustInit(err error) {
	if err != nil {
		logger.Fatalf("init failure: %s", err)
	}
}
