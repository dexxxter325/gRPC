package main

import (
	"GRPC/internal/app"
	"GRPC/internal/config"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	cfg, err := config.Init()
	if err != nil {
		logger.Fatalf("init config failed:%s", err)
	}
	app.RunGrpcGateway(logger, cfg)
	app.RunGRPC(logger, cfg)
}
