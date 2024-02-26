package service

import (
	"GRPC/internal/storage"
	"context"
)

type InvestmentService struct {
	storage storage.Investment
}

func NewInvestmentService(storage storage.Investment) *InvestmentService {
	return &InvestmentService{storage: storage}
}

func (s *InvestmentService) Create(ctx context.Context, amount int64, currency string) (investmentId int, err error) {
	panic("implement me")
}

func (s *InvestmentService) Get(ctx context.Context) (amount int64, currency string, err error) {
	panic("implement me")
}

func (s *InvestmentService) Delete(ctx context.Context, investmentId int) error {
	panic("implement me")
}
