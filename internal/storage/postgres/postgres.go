package postgres

import (
	"GRPC/internal/config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnToPostgres(cfg *config.Config) (*pgxpool.Pool, error) {
	data := fmt.Sprintf("host=%s port=%v user=%s dbname=%s password=%s sslmode=%s", cfg.DB.Postgres.Host, cfg.DB.Postgres.Port, cfg.DB.Postgres.User,
		cfg.DB.Postgres.DbName, cfg.DB.Postgres.Password, cfg.DB.Postgres.Sslmode)
	conn, err := pgxpool.New(context.Background(), data)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database:%s", err)
	}
	return conn, nil
}
