package service

import (
	"GRPC/internal/storage"
	"context"
)

type AuthService struct {
	storage storage.Auth
}

func NewAuthService(storage storage.Auth) *AuthService {
	return &AuthService{storage: storage}
}

func (s *AuthService) Register(ctx context.Context, email, password string) (userId int, err error) {
	panic("implement me")
}
func (s *AuthService) Login(ctx context.Context, email, password string) (token string, err error) {
	panic("implement me")
}
