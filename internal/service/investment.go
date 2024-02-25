package service

import (
	"GRPC/gen"
	"GRPC/internal/storage"
	"context"
)

type InvestmentService struct {
	storage storage.Investment
}

func NewInvestmentService(storage storage.Investment) *InvestmentService {
	return &InvestmentService{storage: storage}
}

func (s *InvestmentService) Create(ctx context.Context, req *gen.CreateRequest) (*gen.CreateResponse, error) {
	panic("implement me")
}

func (s *InvestmentService) Get(ctx context.Context, request *gen.GetRequest) (*gen.GetResponse, error) {
	panic("implement me")
}

func (s *InvestmentService) Delete(ctx context.Context, request *gen.DeleteRequest) (*gen.DeleteResponse, error) {
	panic("implement me")
}
