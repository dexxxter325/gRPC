package handler

import (
	"GRPC/gen"
	"GRPC/internal/service"
	"context"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

/* в сгегерированных файлах у нас есть методы ,которые мы должны реализовать.Это-заглушка,которая позволяет запустить приложение без реализации всех методов */
type InvestmentServer struct {
	gen.UnimplementedInvestmentServer
	Handler *Handler
}

func NewInvestmentServer(handler *Handler) *InvestmentServer {
	return &InvestmentServer{
		UnimplementedInvestmentServer: gen.UnimplementedInvestmentServer{},
		Handler:                       handler,
	}
}

func (s *InvestmentServer) Register(ctx context.Context, req *gen.RegisterRequest) (*gen.RegisterResponse, error) {
	panic("implement me")

}

func (s *InvestmentServer) Login(ctx context.Context, req *gen.LoginRequest) (*gen.LoginResponse, error) {
	panic("implement me")
}

func (s *InvestmentServer) Create(ctx context.Context, req *gen.CreateRequest) (*gen.CreateResponse, error) {
	panic("implement me")
}
func (s *InvestmentServer) Get(ctx context.Context, request *gen.GetRequest) (*gen.GetResponse, error) {
	panic("implement me")
}
func (s *InvestmentServer) Delete(ctx context.Context, request *gen.DeleteRequest) (*gen.DeleteResponse, error) {
	panic("implement me")
}
