package main

import (
	"GRPC/gen"
	"GRPC/internal/config"
	"context"
	"google.golang.org/grpc"
	"log"
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
	cfg, err := config.Init()
	if err != nil {
		log.Fatalf("config.init failed:%s", err)
	}
	//db, err := postgres.ConnToPostgres(cfg)
	listener, err := net.Listen("tcp", cfg.GRPC.Port)
	if err != nil {
		log.Fatalf("listen failed:%s", err)
	}
	service := &InvestmentServer{}
	server := grpc.NewServer()
	gen.RegisterInvestmentServer(server, service)
	log.Println("server started on port 8000!")
	if err = server.Serve(listener); err != nil {
		log.Fatalf("serve failed:%s", err)
	}
}
