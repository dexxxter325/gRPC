package service

import (
	"GRPC/gen"
	"GRPC/internal/storage"
	"context"
)

type Service struct {
	Auth
	Investment
}

func NewService(storage *storage.Storage) *Service {
	return &Service{
		Auth:       NewAuthService(storage.Auth),
		Investment: NewInvestmentService(storage.Investment),
	}
}

type Auth interface {
	Register(ctx context.Context, req *gen.RegisterRequest) (*gen.RegisterResponse, error)
	Login(ctx context.Context, req *gen.LoginRequest) (*gen.LoginResponse, error)
}

type Investment interface {
	Create(ctx context.Context, req *gen.CreateRequest) (*gen.CreateResponse, error)
	Get(ctx context.Context, request *gen.GetRequest) (*gen.GetResponse, error)
	Delete(ctx context.Context, request *gen.DeleteRequest) (*gen.DeleteResponse, error)
}
