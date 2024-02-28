package main

import (
	"GRPC/gen"
	"GRPC/internal/config"
	"GRPC/internal/service"
	"GRPC/internal/storage"
	"GRPC/internal/storage/postgres"
	"GRPC/internal/transport/grpc/handler"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

func main() {
	logger := logrus.New()
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	cfg, err := config.Init()
	if err != nil {
		logrus.Fatalf("config.init failed:%s", err)
	}

	db, err := postgres.ConnToPostgres(cfg)
	storages := storage.NewStorage(db)
	services := service.NewService(storages, logger, cfg)
	handlers := handler.NewHandler(services)

	server := grpc.NewServer()

	investmentRegistrar := handler.NewInvestmentServer(handlers)
	gen.RegisterInvestmentServer(server, investmentRegistrar)

	authRegistrar := handler.NewAuthServer(handlers)
	gen.RegisterAuthServer(server, authRegistrar)

	listener, err := net.Listen("tcp", cfg.GRPC.Port)
	if err != nil {
		logrus.Fatalf("listen failed:%s", err)
	}
	logrus.Info("server started on port 8000!")

	if err = server.Serve(listener); err != nil {
		logrus.Fatalf("serve failed:%s", err)
	}
}
