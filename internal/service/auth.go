package service

import (
	"GRPC/gen"
	"GRPC/internal/storage"
	"context"
)

type AuthService struct {
	storage storage.Auth
}

func NewAuthService(storage storage.Auth) *AuthService {
	return &AuthService{storage: storage}
}

func (s *AuthService) Register(ctx context.Context, req *gen.RegisterRequest) (*gen.RegisterResponse, error) {
	panic("implement me")
}
func (s *AuthService) Login(ctx context.Context, req *gen.LoginRequest) (*gen.LoginResponse, error) {
	panic("implement me")
}
