package handler

import (
	"GRPC/gen"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *AuthServer) Register(ctx context.Context, req *gen.RegisterRequest) (*gen.RegisterResponse, error) {
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password in required")
	}
	userId, err := s.handler.service.Register(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Register failed in handler:%s", err)
	}
	return &gen.RegisterResponse{UserId: userId}, nil

}

func (s *AuthServer) Login(ctx context.Context, req *gen.LoginRequest) (*gen.LoginResponse, error) {
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password in required")
	}
	accessToken, err := s.handler.service.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "login failed in handler:%s", err) //internal-сломалось что-то на стороне сервера /,а не клиента
	}
	return &gen.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: "refresh",
	}, nil
}
func (s *AuthServer) RefreshToken(ctx context.Context, req *gen.RefreshTokenRequest) (*gen.RefreshTokenResponse, error) {
	panic("imp me!")
}
