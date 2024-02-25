package postgres

import (
	"GRPC/gen"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InvestmentPostgres struct {
	db *pgxpool.Pool
}

func NewInvestmentPostgres(db *pgxpool.Pool) *InvestmentPostgres {
	return &InvestmentPostgres{db: db}
}

func (p *InvestmentPostgres) Create(ctx context.Context, req *gen.CreateRequest) (*gen.CreateResponse, error) {
	panic("implement me ")
}

func (p *InvestmentPostgres) Get(ctx context.Context, request *gen.GetRequest) (*gen.GetResponse, error) {
	panic("implement me ")
}

func (p *InvestmentPostgres) Delete(ctx context.Context, request *gen.DeleteRequest) (*gen.DeleteResponse, error) {
	panic("implement me ")
}
