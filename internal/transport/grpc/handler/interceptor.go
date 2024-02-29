package handler

import (
	"GRPC/internal/config"
	"GRPC/internal/service"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

/*UnaryInterceptor-когда клиент отправляет запрос и сервер отвечает на него одним ответом.
exml:позволит добавить проверку на то,что запрос к моему серверу проверяется на аутентификацию,прежде чем он будет обработан.
StreamInterceptor-когда данные передаются как поток в обе стороны, в отличие от одиночных запросов и ответов.
exml:приложение чата,где клиент и сервер могут передавать сообщения в реальном времени.StreamInterceptor поможет добавить логику проверки сообщений в реальном времени*/

func UnaryInterceptor(cfg *config.Config) grpc.UnaryServerInterceptor {
	// Возвращаем функцию-перехватчик
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if info.FullMethod != "/Auth/Register" && info.FullMethod != "/Auth/Login" && info.FullMethod != "/Auth/RefreshToken" {
			authCtx, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				return nil, status.Error(codes.Unauthenticated, "failed get ctx in interceptor")
			}
			header := authCtx.Get("Authorization")
			if len(header) == 0 {
				return nil, status.Errorf(codes.Unauthenticated, "authorization header is empty")
			}

			bearerAndToken := header[0] //наглядно указываем ,что будет в запросе(для дальнейшего разделения)
			headerParts := strings.Split(bearerAndToken, " ")

			if len(headerParts) != 2 && headerParts[0] != "Bearer" {
				return nil, status.Errorf(codes.Unauthenticated, "invalid auth header")
			}
			if len(headerParts[1]) == 0 {
				return nil, status.Error(codes.Unauthenticated, "empty auth token")
			}

			ok, err := service.ValidateAccessToken(headerParts[1], cfg.AUTH.SecretKey)
			if !ok {
				return nil, status.Errorf(codes.Unauthenticated, "ValidateAccessToken failed:%s", err)
			}
		}

		return handler(ctx, req)
	}
}

/*настроить ci/cd,покрыть функц.тестами,мб из main  перенести в internal/app .shutdown*/
