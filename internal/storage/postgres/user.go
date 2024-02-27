package postgres

import (
	"GRPC/internal/domain/models"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserPostgres struct {
	db *pgxpool.Pool
}

func NewUserPostgres(db *pgxpool.Pool) *UserPostgres {
	return &UserPostgres{db: db}
}

func (p *UserPostgres) SaveUser(ctx context.Context, email string, hashedPassword []byte) (userId int64, err error) {
	query := `insert into users(email,password) values ($1,$2) returning id`
	row := p.db.QueryRow(ctx, query, email, hashedPassword)
	if err = row.Scan(&userId); err != nil {
		return 0, fmt.Errorf("scan failed in SaveUser:%s", err)
	}
	return userId, nil

}

func (p *UserPostgres) GetUserByEmail(ctx context.Context, email string) (user models.User, err error) {
	query := `select * from users where email=$1`
	row := p.db.QueryRow(ctx, query, email)
	err = row.Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, fmt.Errorf("user not found with email: %s", email)
		}
		return models.User{}, fmt.Errorf("scan failed in getuser:%s", err)
	}
	return user, nil
}
