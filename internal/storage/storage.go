package storage

import (
	"GRPC/gen"
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
	Register(ctx context.Context, req *gen.RegisterRequest) (*gen.RegisterResponse, error)
	Login(ctx context.Context, req *gen.LoginRequest) (*gen.LoginResponse, error)
}

type Investment interface {
	Create(ctx context.Context, req *gen.CreateRequest) (*gen.CreateResponse, error)
	Get(ctx context.Context, request *gen.GetRequest) (*gen.GetResponse, error)
	Delete(ctx context.Context, request *gen.DeleteRequest) (*gen.DeleteResponse, error)
}
