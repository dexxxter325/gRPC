package handler

import (
	"GRPC/gen"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *InvestmentServer) Create(ctx context.Context, req *gen.CreateRequest) (*gen.CreateResponse, error) {
	if req.GetCurrency() == "" {
		return nil, status.Error(codes.InvalidArgument, "currency in required")
	}
	if req.GetAmount() == 0 {
		return nil, status.Error(codes.InvalidArgument, "amount is required")
	}
	investmentId, err := s.handler.service.Create(ctx, req.GetAmount(), req.GetCurrency())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create failed in handler:%s", err)
	}
	return &gen.CreateResponse{InvestmentId: investmentId}, nil
}
func (s *InvestmentServer) Get(ctx context.Context, req *gen.GetRequest) (*gen.GetResponse, error) {
	investment, err := s.handler.service.Get(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get failed in handler:%s", err)
	}
	var protoInvestment []*gen.Investments
	for _, investments := range investment {
		protoInvestment = append(protoInvestment, &gen.Investments{
			ID:       investments.ID,
			Amount:   investments.Amount,
			Currency: investments.Currency,
		})
	}
	return &gen.GetResponse{Investment: protoInvestment}, nil
}
func (s *InvestmentServer) Delete(ctx context.Context, req *gen.DeleteRequest) (*gen.DeleteResponse, error) {
	if req.GetInvestmentId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "investment id is required")
	}
	if err := s.handler.service.Delete(ctx, req.GetInvestmentId()); err != nil {
		return nil, status.Errorf(codes.Internal, "delete failed in handler:%s", err)
	}
	return &gen.DeleteResponse{}, nil
}
