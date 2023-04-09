package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
)

const (
	LocalEnv = "local"
)

var logger *zap.Logger

func init() {
	var err error
	if os.Getenv("ENV") == LocalEnv {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		panic(err)
	}

	defer func(logger *zap.Logger) {
		err = logger.Sync()
		if err != nil {
			logger.Error(fmt.Sprintf("logger sync: %v", err))
		}
	}(logger)
}

func Infof(format string, args ...interface{}) {
	logger.Info(fmt.Sprintf(format, args...))
}

func Info(msg interface{}) {
	logger.Info(fmt.Sprint(msg))
}

func Warnf(format string, args ...interface{}) {
	logger.Warn(fmt.Sprintf(format, args...))
}

func Warn(msg interface{}) {
	logger.Warn(fmt.Sprint(msg))
}

func Errorf(format string, args ...interface{}) {
	logger.Error(fmt.Sprintf(format, args...))
}

func Error(msg interface{}) {
	logger.Error(fmt.Sprint(msg))
}

func Fatalf(format string, args ...interface{}) {
	logger.Fatal(fmt.Sprintf(format, args...))
}

func Fatal(msg interface{}) {
	logger.Fatal(fmt.Sprint(msg))
}
