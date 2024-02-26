package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthPostgres struct {
	db *pgxpool.Pool
}

func NewAuthPostgres(db *pgxpool.Pool) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (p *AuthPostgres) Register(ctx context.Context, email, password string) (userId int, err error) {
	panic("implement me !")
}

func (p *AuthPostgres) Login(ctx context.Context, email, password string) (token string, err error) {
	panic("implement me!")
}
