package service

import (
	"GRPC/internal/storage"
	"context"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

type AuthService struct {
	storage storage.User
	logger  *logrus.Logger
}

func NewAuthService(storage storage.User, logger *logrus.Logger) *AuthService {
	return &AuthService{storage: storage, logger: logger}
}

const (
	accessTokenTTL  = 24 * 30 * time.Hour //1 month
	refreshTokenTTL = 24 * 30 * time.Hour
)

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
func (s *AuthService) Login(ctx context.Context, email, password string) (token string, err error) {
	log := s.logger.WithFields(logrus.Fields{
		"email":    email,
		"password": password,
	})
	log.Info("Received login request")

	user, err := s.storage.GetUserByEmail(ctx, email)
	if err != nil {
		log.Error("failed to getUser")
		return "", err
	}
	if err = bcrypt.CompareHashAndPassword(user.Password, []byte(password)); err != nil {
		log.Error("password wrong")
		return "", err
	}
	secretKey := os.Getenv("SECRETKEY")
	token, err = GenerateNewAccessToken(user, accessTokenTTL, secretKey)
	if err != nil {
		log.Errorf("failed in GenerateNewAccessToken:%s", err)
		return "", err
	}

	log.Info("user logged in")
	return token, nil
}
