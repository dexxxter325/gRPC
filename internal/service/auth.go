package service

import (
	"GRPC/internal/config"
	"GRPC/internal/storage"
	"context"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AuthService struct {
	storage storage.User
	logger  *logrus.Logger
	cfg     *config.Config
}

func NewAuthService(storage storage.User, logger *logrus.Logger, config *config.Config) *AuthService {
	return &AuthService{storage: storage, logger: logger, cfg: config}
}

func (s *AuthService) Register(ctx context.Context, email, password string) (userId int64, err error) {
	log := s.logger.WithFields(logrus.Fields{
		"email":    email,
		"password": password,
	})
	log.Info("Received registration request")

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) //salt already inside.Cost-hash security lvl
	if err != nil {
		log.Error("failed to generate hash password")
		return 0, err
	}

	userId, err = s.storage.SaveUser(ctx, email, hashPassword)
	if err != nil {
		log.Error("failed to saveUser")
		return 0, err
	}

	log.Info("user registered")
	return userId, nil

}
func (s *AuthService) Login(ctx context.Context, email, password string) (accessToken, refreshToken string, err error) {
	log := s.logger.WithFields(logrus.Fields{
		"email":    email,
		"password": password,
	})
	log.Info("Received login request")

	user, err := s.storage.GetUserByEmail(ctx, email)
	if err != nil {
		log.Error("failed to getUser")
		return "", "", err
	}
	if err = bcrypt.CompareHashAndPassword(user.Password, []byte(password)); err != nil {
		log.Error("password wrong")
		return "", "", err
	}

	secretKey := s.cfg.AUTH.SecretKey

	accessTokenTTLStr := s.cfg.AUTH.AccessTokenTTl
	accessTokenTTL, err := time.ParseDuration(accessTokenTTLStr)
	if err != nil {
		log.Error("convert accessTokenTTL to type time failed")
		return "", "", err
	}
	accessToken, err = GenerateNewAccessToken(user, accessTokenTTL, secretKey)
	if err != nil {
		log.Errorf("failed in GenerateNewAccessToken:%s", err)
		return "", "", err
	}

	refreshTokenTTLStr := s.cfg.AUTH.RefreshTokenTTl
	refreshTokenTTL, err := time.ParseDuration(refreshTokenTTLStr)
	if err != nil {
		log.Error("convert refreshTokenTTL to type time failed")
		return "", "", err
	}
	refreshToken, err = GenerateNewRefreshToken(user.ID, refreshTokenTTL, secretKey)
	if err != nil {
		log.Errorf("failed in GenerateNewRefreshToken:%s", err)
		return "", "", err
	}

	log.Info("user logged in")
	return accessToken, refreshToken, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (NewAccessToken, NewRefreshToken string, err error) {
	log := s.logger.WithField("refreshToken", refreshToken)
	log.Info("Received RefreshToken request")

	secretKey := s.cfg.AUTH.SecretKey

	userId, err := ParseRefreshToken(refreshToken, secretKey)
	if err != nil {
		log.Errorf("parse Refresh Token failed:%s", err)
		return "", "", err
	}

	user, err := s.storage.GetUserById(ctx, userId)
	if err != nil {
		log.Errorf("Get user by id failed:%s", err)
		return "", "", err
	}

	accessTokenTTLStr := s.cfg.AUTH.AccessTokenTTl
	accessTokenTTL, err := time.ParseDuration(accessTokenTTLStr)
	if err != nil {
		log.Errorf("parse accessTOkenTTL to type time failed:%s", err)
		return "", "", nil
	}
	NewAccessToken, err = GenerateNewAccessToken(user, accessTokenTTL, secretKey)
	if err != nil {
		log.Errorf("generateNewTokenPair failed:%s", err)
		return "", "", err
	}

	refreshTTLStr := s.cfg.AUTH.RefreshTokenTTl
	refreshTokenTTL, err := time.ParseDuration(refreshTTLStr)
	if err != nil {
		log.Errorf("parse accessTOkenTTL to type time failed:%s", err)
		return "", "", nil
	}
	NewRefreshToken, err = GenerateNewRefreshToken(user.ID, refreshTokenTTL, secretKey)
	if err != nil {
		log.Errorf("generateNewTokenPair failed:%s", err)
		return "", "", err
	}

	log.Info("tokens generated")
	return NewAccessToken, NewRefreshToken, err
}
