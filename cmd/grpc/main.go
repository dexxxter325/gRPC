package main

import (
	"GRPC/gen"
	"GRPC/internal/config"
	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

type InvestmentServer struct {
	gen.UnimplementedInvestmentServer //заглушка для сервера-пустая реализация методов интерфейса сервиса
}

func (s *InvestmentServer) Create(context.Context, *gen.CreateRequest) (*gen.CreateResponse, error) {
	return &gen.CreateResponse{
		InvestmentId: 1,
	}, nil
}

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	cfg, err := config.Init()
	if err != nil {
		logrus.Fatalf("config.init failed:%s", err)
	}
	//db, err := postgres.ConnToPostgres(cfg)
	listener, err := net.Listen("tcp", cfg.GRPC.Port)
	if err != nil {
		logrus.Fatalf("listen failed:%s", err)
	}
	service := &InvestmentServer{}
	server := grpc.NewServer()
	gen.RegisterInvestmentServer(server, service)
	logrus.Info("server started on port 8000!")
	if err = server.Serve(listener); err != nil {
		logrus.Fatalf("serve failed:%s", err)
	}
}
