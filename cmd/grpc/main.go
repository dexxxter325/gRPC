package main

import (
	"GRPC/gen"
	"GRPC/internal/config"
	"GRPC/internal/service"
	"GRPC/internal/storage"
	"GRPC/internal/storage/postgres"
	"GRPC/internal/transport/grpc/handler"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := logrus.New()
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	cfg, err := config.Init()
	if err != nil {
		logger.Fatalf("init config failed:%s", err)
	}

	db, err := postgres.ConnToPostgres(cfg)
	if err != nil {
		logger.Fatalf("connect to postgres failed:%s", err)
	}
	defer db.Close() //для избежания утечки рес-ов

	storages := storage.NewStorage(db)
	services := service.NewService(storages, logger, cfg)
	handlers := handler.NewHandler(services)

	server := grpc.NewServer(grpc.UnaryInterceptor(handler.UnaryInterceptor(cfg)))

	investmentRegistrar := handler.NewInvestmentServer(handlers)
	gen.RegisterInvestmentServer(server, investmentRegistrar)

	authRegistrar := handler.NewAuthServer(handlers)
	gen.RegisterAuthServer(server, authRegistrar)
	//запускаем сервер в параллельно ,чтобы дальше в коде ждать от нее сигнала.Если бы запускали без горутины(не параллельно),то код не пошел бы дальше и мы не получили опр.сигналов
	go func() {
		listener, err := net.Listen("tcp", cfg.GRPC.Port)
		if err != nil {
			logger.Fatalf("listen failed:%s", err)
		}
		logger.Info("server started on port 8000!")
		if err = server.Serve(listener); err != nil {
			logger.Fatalf("serve failed:%s", err)
		}
	}()
	stop := make(chan os.Signal, 1)                      //канал для передачи информации о сигналах
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT) //когда прийдет один из сигналов-запишет это в канал

	<-stop //будет в режиме ожидания получения значения.После получения-код идет дальше

	server.GracefulStop() //1)Остановка приема новых запросов.2)ожидает завершения обработки всех текущих запросов.3)стопает сервер

	logger.Info("application stopped")
}
