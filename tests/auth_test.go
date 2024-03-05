package tests

//заменяют нам postman

import (
	"GRPC/gen"
	"GRPC/internal/service"
	"GRPC/tests/suite"
	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"
	"testing"
)

/*
	require используется для проверки условий и завершает тест немедленно, если проверка не удалась.

assert также используется для проверки условий, но в отличие от require, он продолжает выполнение теста, даже если проверка не удалась.
*/

func TestAuth_Ok(t *testing.T) {
	ctx, st, err := suite.New(t)
	require.NoError(t, err)

	email := gofakeit.Email()
	password := fakePass()
	secretKey := st.Cfg.AUTH.SecretKey

	registerResp, err := st.AuthClient.Register(ctx, &gen.RegisterRequest{
		Email:    email,
		Password: password,
	})

	id := registerResp.GetUserId()
	require.NotEmpty(t, id)

	require.NoError(t, err) //проверяем,что err=nil.если нет-тест дропнется

	loginResp, err := st.AuthClient.Login(ctx, &gen.LoginRequest{
		Email:    email,
		Password: password,
	})

	require.NoError(t, err)

	accessToken := loginResp.GetAccessToken()
	require.NotEmpty(t, accessToken)
	ok, err := service.ValidateAccessToken(accessToken, secretKey)
	require.True(t, ok)
	require.NoError(t, err)

	refreshToken := loginResp.GetRefreshToken()
	require.NotEmpty(t, refreshToken)

	refreshTokenResp, err := st.AuthClient.RefreshToken(ctx, &gen.RefreshTokenRequest{
		RefreshToken: refreshToken,
	})
	require.NoError(t, err)

	newAccessToken := refreshTokenResp.GetAccessToken()
	require.NotEmpty(t, newAccessToken)

	newRefreshToken := refreshTokenResp.GetRefreshToken()
	require.NotEmpty(t, newRefreshToken)
}

func TestRegister_Fail(t *testing.T) {
	ctx, st, err := suite.New(t)
	require.NoError(t, err)

	tests := []struct {
		name        string
		email       string
		password    string
		expectedErr string
	}{
		{
			name:        "empty fields",
			email:       "",
			password:    "",
			expectedErr: "email is required",
		},
		{
			name:        "empty email",
			email:       "",
			password:    fakePass(),
			expectedErr: "email is required",
		},
		{
			name:        "empty password",
			email:       gofakeit.Email(),
			password:    "",
			expectedErr: "password is required",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err = st.AuthClient.Register(ctx, &gen.RegisterRequest{
				Email:    tt.email,
				Password: tt.password,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr) //есть ли ошибка в err.error в tt.expectedErr(выводятся также ненужные данные-rpc method...)
		})
	}
}

func TestLogin_Fail(t *testing.T) {
	ctx, st, err := suite.New(t)
	require.NoError(t, err)

	tests := []struct {
		name        string
		email       string
		password    string
		expectedErr string
	}{
		{
			name:        "empty fields",
			email:       "",
			password:    "",
			expectedErr: "email is required",
		},
		{
			name:        "empty email",
			email:       "",
			password:    fakePass(),
			expectedErr: "email is required",
		},
		{
			name:        "empty password",
			email:       gofakeit.Email(),
			password:    "",
			expectedErr: "password is required",
		},
		{
			name:        "incorrect data",
			email:       gofakeit.Email(),
			password:    fakePass(),
			expectedErr: "user not found with email",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err = st.AuthClient.Login(ctx, &gen.LoginRequest{
				Email:    tt.email,
				Password: tt.password,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func TestRefreshToken_Fail(t *testing.T) {
	ctx, st, err := suite.New(t)
	require.NoError(t, err)

	tests := []struct {
		name         string
		refreshToken string
		expectedErr  string
	}{
		{
			name:         "empty refreshToken",
			refreshToken: "",
			expectedErr:  "refresh token is required",
		},
		{
			name:         "incorrect number of characters in RefreshToken",
			refreshToken: "qwerty123",
			expectedErr:  "token contains an invalid number of segments",
		},
		{
			name:         "invalid RefreshToken",
			refreshToken: "000000000000000",
			expectedErr:  "parse refreshTOken failed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err = st.AuthClient.RefreshToken(ctx, &gen.RefreshTokenRequest{
				RefreshToken: tt.refreshToken,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func fakePass() string {
	return gofakeit.Password(true, true, true, true, false, 15)
}

/*lower-символы нижнего регистра (a-z)
upper-символы верхнего регистра (A-Z).
numeric-цифры (0-9).
special-специальные символы.
space-пробелы
num-кол-во символов в пароле*/
