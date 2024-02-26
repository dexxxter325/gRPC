package storage

import (
	"GRPC/internal/storage/postgres"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	Auth
	Investment
}

func NewStorage(db *pgxpool.Pool) *Storage {
	return &Storage{
		Auth:       postgres.NewAuthPostgres(db),
		Investment: postgres.NewInvestmentPostgres(db),
	}
}

type Auth interface {
	Register(ctx context.Context, email, password string) (userId int, err error)
	Login(ctx context.Context, email, password string) (token string, err error)
}

type Investment interface {
	Create(ctx context.Context, amount int64, currency string) (investmentId int, err error)
	Get(ctx context.Context) (amount int64, currency string, err error)
	Delete(ctx context.Context, investmentId int) error
}
