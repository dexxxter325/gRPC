package storage

import (
	"GRPC/internal/domain/models"
	"GRPC/internal/storage/postgres"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	User
	Investment
}

func NewStorage(db *pgxpool.Pool) *Storage {
	return &Storage{
		User:       postgres.NewUserPostgres(db),
		Investment: postgres.NewInvestmentPostgres(db),
	}
}

type User interface {
	SaveUser(ctx context.Context, email string, hashedPassword []byte) (userId int64, err error)
	GetUserByEmail(ctx context.Context, email string) (user models.User, err error)
}

type Investment interface {
	Create(ctx context.Context, amount int64, currency string) (investmentId int64, err error)
	Get(ctx context.Context) (investment models.Investment, err error)
	Delete(ctx context.Context, investmentId int64) error
}
