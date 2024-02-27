package service

import (
	"GRPC/internal/domain/models"
	"GRPC/internal/storage"
	"context"
	"github.com/sirupsen/logrus"
)

type Service struct {
	Auth
	Investment
}

func NewService(storage *storage.Storage, logger *logrus.Logger) *Service {
	return &Service{
		Auth:       NewAuthService(storage.User, logger),
		Investment: NewInvestmentService(storage.Investment, logger),
	}
}

type Auth interface {
	Register(ctx context.Context, email, password string) (userId int64, err error)
	Login(ctx context.Context, email, password string) (token string, err error)
}

type Investment interface {
	Create(ctx context.Context, amount int64, currency string) (investmentId int64, err error)
	Get(ctx context.Context) (investment []models.Investment, err error)
	Delete(ctx context.Context, investmentId int64) error
}
