package service

import (
	"GRPC/internal/domain/models"
	"GRPC/internal/storage"
	"context"
	"github.com/sirupsen/logrus"
)

type InvestmentService struct {
	storage storage.Investment
	logger  *logrus.Logger
}

func NewInvestmentService(storage storage.Investment, logger *logrus.Logger) *InvestmentService {
	return &InvestmentService{storage: storage, logger: logger}
}

func (s *InvestmentService) Create(ctx context.Context, amount int64, currency string) (investmentId int64, err error) {
	log := s.logger.WithFields(logrus.Fields{
		"amount":   amount,
		"currency": currency,
	})
	log.Info("received create req")

	investmentId, err = s.storage.Create(ctx, amount, currency)
	if err != nil {
		log.Error("failed to create")
		return 0, err
	}

	log.Info("investment created")
	return investmentId, nil
}

func (s *InvestmentService) Get(ctx context.Context) (investment models.Investment, err error) {
	s.logger.Info("received get req")

	investment, err = s.storage.Get(ctx)
	if err != nil {
		s.logger.Error("get failed")
		return models.Investment{}, err
	}

	s.logger.Info("investment got")
	return investment, nil
}

func (s *InvestmentService) Delete(ctx context.Context, investmentId int64) error {
	log := s.logger.WithField("investmentId", investmentId)
	log.Info("received delete req")

	if err := s.storage.Delete(ctx, investmentId); err != nil {
		log.Error("delete failed")
		return err
	}

	log.Info("investment deleted")
	return nil
}
