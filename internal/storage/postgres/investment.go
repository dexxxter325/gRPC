package postgres

import (
	"GRPC/internal/domain/models"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InvestmentPostgres struct {
	db *pgxpool.Pool
}

func NewInvestmentPostgres(db *pgxpool.Pool) *InvestmentPostgres {
	return &InvestmentPostgres{db: db}
}

func (p *InvestmentPostgres) Create(ctx context.Context, amount int64, currency string) (investmentId int64, err error) {
	panic("implement me ")
}

func (p *InvestmentPostgres) Get(ctx context.Context) (investment models.Investment, err error) {
	panic("implement me ")
}

func (p *InvestmentPostgres) Delete(ctx context.Context, investmentId int64) error {
	panic("implement me ")
}
