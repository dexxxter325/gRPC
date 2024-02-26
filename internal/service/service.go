package service

import (
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
	Register(ctx context.Context, email, password string) (userId int64, err error)
	Login(ctx context.Context, email, password string) (token string, err error)
}

type Investment interface {
	Create(ctx context.Context, amount int64, currency string) (investmentId int64, err error)
	Get(ctx context.Context) (amount int64, currency string, err error)
	Delete(ctx context.Context, investmentId int64) error
}
