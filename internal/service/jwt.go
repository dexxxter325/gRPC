package service

import (
	"GRPC/internal/domain/models"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func GenerateNewAccessToken(user models.User, accessTokenTTL time.Duration, secretKey string) (string, error) {
	accessToken := jwt.New(jwt.SigningMethodHS256)
	claims := accessToken.Claims.(jwt.MapClaims)
	claims["userId"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(accessTokenTTL).Unix()
	accessTokenString, err := accessToken.SignedString([]byte(secretKey)) //подписываем токен секретным ключом
	if err != nil {
		return "", fmt.Errorf("convert access token to string failed:%s", err)
	}
	return accessTokenString, nil
}

func ParseAccessToken(accessToken string) (bool, error) {
	panic("скорей всего будем юзать в middleware")
}

func GenerateNewRefreshToken() (refreshToken string, err error) {
	panic("imppppp me")
}

func ParseRefreshToken(refreshToken string) (bool, error) {
	panic("скорей всего будем юзать в middleware")
}
