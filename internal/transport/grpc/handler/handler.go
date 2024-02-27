package handler

import (
	"GRPC/gen"
	"GRPC/internal/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

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

type AuthServer struct {
	gen.UnimplementedAuthServer
	handler *Handler
}

func NewAuthServer(handler *Handler) *AuthServer {
	return &AuthServer{
		UnimplementedAuthServer: gen.UnimplementedAuthServer{},
		handler:                 handler,
	}
}
