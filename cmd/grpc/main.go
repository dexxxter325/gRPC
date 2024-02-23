package main

import (
	"GRPC/gen"
	"context"
	"fmt"
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
		Status:       "good",
	}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatalf("listen failed:%s", err)
	}
	service := &InvestmentServer{}
	server := grpc.NewServer()
	gen.RegisterInvestmentServer(server, service)
	fmt.Println("server started on port 8000!")
	if err = server.Serve(listener); err != nil {
		log.Fatalf("serve failed:%s", err)
	}
}
