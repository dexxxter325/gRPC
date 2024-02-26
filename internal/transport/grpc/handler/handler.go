package handler

import (
	"GRPC/gen"
	"GRPC/internal/service"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	handler *Handler
}

func NewInvestmentServer(handler *Handler) *InvestmentServer {
	return &InvestmentServer{
		UnimplementedInvestmentServer: gen.UnimplementedInvestmentServer{},
		handler:                       handler,
	}
}

func (s *InvestmentServer) Register(ctx context.Context, req *gen.RegisterRequest) (*gen.RegisterResponse, error) {
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password in required")
	}
	userId, err := s.handler.service.Register(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Register failed in handler:%s", err)
	}
	return &gen.RegisterResponse{UserId: userId}, nil

}

func (s *InvestmentServer) Login(ctx context.Context, req *gen.LoginRequest) (*gen.LoginResponse, error) {
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password in required")
	}
	token, err := s.handler.service.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "login failed in handler:%s", err) //internal-сломалось что-то на стороне сервера /,а не клиента
	}
	return &gen.LoginResponse{Token: token}, nil
}

func (s *InvestmentServer) Create(ctx context.Context, req *gen.CreateRequest) (*gen.CreateResponse, error) {
	if req.GetCurrency() == "" {
		return nil, status.Error(codes.InvalidArgument, "currency in required")
	}
	if req.GetAmount() == 0 {
		return nil, status.Error(codes.InvalidArgument, "amount is required")
	}
	investmentId, err := s.handler.service.Create(ctx, req.GetAmount(), req.GetCurrency())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create failed in handler:%s", err)
	}
	return &gen.CreateResponse{InvestmentId: investmentId}, nil
}
func (s *InvestmentServer) Get(ctx context.Context, req *gen.GetRequest) (*gen.GetResponse, error) {
	amount, currency, err := s.handler.service.Get(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get failed in handler:%s", err)
	}
	return &gen.GetResponse{
		Amount:   amount,
		Currency: currency,
	}, nil
}
func (s *InvestmentServer) Delete(ctx context.Context, req *gen.DeleteRequest) (*gen.DeleteResponse, error) {
	if req.GetInvestmentId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "investment id is required")
	}
	return &gen.DeleteResponse{}, nil
}
