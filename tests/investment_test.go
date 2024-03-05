package tests

import (
	"GRPC/gen"
	"GRPC/internal/service"
	"GRPC/tests/suite"
	"context"
	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/rand"
	"google.golang.org/grpc/metadata"
	"testing"
	"time"
)

func TestInvestment_OK(t *testing.T) {
	ctx, st, err := suite.New(t)
	require.NoError(t, err)

	accessToken := GetAccessToken(ctx, st, t)
	require.NotEmpty(t, accessToken)
	ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", "Bearer "+accessToken)

	amount := randomAmount()
	currency := randomCurrency(7)

	createResp, err := st.InvestmentClient.Create(ctx, &gen.CreateRequest{
		Amount:   amount,
		Currency: currency,
	})
	require.NoError(t, err)

	id := createResp.GetInvestmentId()
	require.NotEmpty(t, id)

	getResp, err := st.InvestmentClient.Get(ctx, &gen.GetRequest{})
	require.NoError(t, err)

	investment := getResp.GetInvestment()
	require.NotEmpty(t, investment)

	delResp, err := st.InvestmentClient.Delete(ctx, &gen.DeleteRequest{
		InvestmentId: id,
	})
	require.NoError(t, err)
	require.NotEmpty(t, delResp)
}

// TestCreate_Fail run via the TestInvestment_OK func
func TestCreate_Fail(t *testing.T) {
	ctx, st, err := suite.New(t)
	require.NoError(t, err)

	accessToken := GetAccessToken(ctx, st, t)
	require.NotEmpty(t, accessToken)
	ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", "Bearer "+accessToken)

	tests := []struct {
		name        string
		amount      int64
		currency    string
		expectedErr string
	}{
		{
			name:        "empty fields",
			amount:      0,
			currency:    "",
			expectedErr: "currency in required",
		},
		{
			name:        "empty amount",
			amount:      0,
			currency:    randomCurrency(7),
			expectedErr: "amount is required",
		},
		{
			name:        "empty currency",
			amount:      randomAmount(),
			currency:    "",
			expectedErr: "currency in required",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err = st.InvestmentClient.Create(ctx, &gen.CreateRequest{
				Amount:   tt.amount,
				Currency: tt.currency,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func TestDelete_Fail(t *testing.T) {
	ctx, st, err := suite.New(t)
	require.NoError(t, err)

	accessToken := GetAccessToken(ctx, st, t)
	require.NotEmpty(t, accessToken)
	ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", "Bearer "+accessToken)

	tests := []struct {
		name        string
		id          int64
		expectedErr string
	}{
		{
			name:        "empty id",
			id:          0,
			expectedErr: "investment id is required",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err = st.InvestmentClient.Delete(ctx, &gen.DeleteRequest{
				InvestmentId: tt.id,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

// For test Investment`s methods by access token
func GetAccessToken(ctx context.Context, st *suite.Suite, t *testing.T) string {
	//logic for access token
	email := gofakeit.Email()
	password := fakePass()
	secretKey := st.Cfg.AUTH.SecretKey

	registerResp, err := st.AuthClient.Register(ctx, &gen.RegisterRequest{
		Email:    email,
		Password: password,
	})
	id := registerResp.GetUserId()
	require.NotEmpty(t, id)
	require.NoError(t, err)

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
	return accessToken
}

func randomCurrency(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	rand.Seed(uint64(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func randomAmount() int64 {
	x := rand.Intn(2147483647) + 1 //max int in PostgreSql
	return int64(x)
}
