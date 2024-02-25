package postgres

import (
	"GRPC/gen"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthPostgres struct {
	db *pgxpool.Pool
}

func NewAuthPostgres(db *pgxpool.Pool) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (p *AuthPostgres) Register(ctx context.Context, req *gen.RegisterRequest) (*gen.RegisterResponse, error) {
	panic("implement me !")
}

func (p *AuthPostgres) Login(ctx context.Context, req *gen.LoginRequest) (*gen.LoginResponse, error) {
	panic("implement me!")
}
