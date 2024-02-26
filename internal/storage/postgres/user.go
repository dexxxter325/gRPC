package postgres

import (
	"GRPC/internal/domain/models"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserPostgres struct {
	db *pgxpool.Pool
}

func NewUserPostgres(db *pgxpool.Pool) *UserPostgres {
	return &UserPostgres{db: db}
}

func (p *UserPostgres) SaveUser(ctx context.Context, email string, hashedPassword []byte) (userId int64, err error) {
	panic("implement me !")
}

func (p *UserPostgres) GetUserByEmail(ctx context.Context, email string) (user models.User, err error) {
	panic("implement me!")
}
