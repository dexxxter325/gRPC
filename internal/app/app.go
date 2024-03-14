package app

import (
	"GRPC/gen"
	"GRPC/internal/config"
	"GRPC/internal/metrics"
	"GRPC/internal/service"
	"GRPC/internal/storage"
	"GRPC/internal/storage/postgres"
	"GRPC/internal/traces"
	"GRPC/internal/transport/grpc/handler"
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func RunGRPC(logger *logrus.Logger, cfg *config.Config) {

	shutDownTraces, err := traces.InitTraces("http://172.17.0.1:14268/api/traces", "grpc_service") //D0cker host
	if err != nil {
		logger.Fatalf("Init Jaeger failed:%s", err)
	}
	defer func() {
		shutDownTraces()
		logger.Info("traces stopped.")
	}()

	shutDownMetrics, err := metrics.InitMetrics(cfg.Metrics.Port, logger)
	if err != nil {
		logger.Fatalf("create metrics failed:%s", err)
	}
	defer func() {
		shutDownMetrics()
		logger.Info("metrics stopped")
	}()

	db, err := postgres.ConnToPostgres(cfg)
	if err != nil {
		logger.Fatalf("connect to postgres failed:%s", err)
	}
	defer func() {
		db.Close() //для избежания утечки рес-ов
		logger.Info("Postgres Connection closed")
	}()

	storages := storage.NewStorage(db)
	services := service.NewService(storages, logger, cfg)
	handlers := handler.NewHandler(services)

	server := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()), //auto generating trace spans
		grpc.UnaryInterceptor(handler.UnaryInterceptor(cfg)),
	)

	investmentRegistrar := handler.NewInvestmentServer(handlers)
	gen.RegisterInvestmentServer(server, investmentRegistrar)

	authRegistrar := handler.NewAuthServer(handlers)
	gen.RegisterAuthServer(server, authRegistrar)

	//запускаем сервер в параллельно ,чтобы дальше в коде ждать от нее сигнала.Если бы запускали без горутины(не параллельно),то код не пошел бы дальше и мы не получили опр.сигналов
	go func() {
		listener, err := net.Listen("tcp", ":"+cfg.GRPC.Port)
		if err != nil {
			logger.Fatalf("listen failed:%s", err)
		}
		logger.Infof("server started on port:%s", cfg.GRPC.Port)
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

func RunGrpcGateway(logger *logrus.Logger, cfg *config.Config) {
	// Register gRPC server endpoint
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := gen.RegisterInvestmentHandlerFromEndpoint(context.Background(), mux, "grpc:"+cfg.GRPC.Port, opts); err != nil { ////docker host
		logger.Fatalf("register GRPC gateway failed:%s", err)
	}

	if err := gen.RegisterAuthHandlerFromEndpoint(context.Background(), mux, "grpc:"+cfg.GRPC.Port, opts); err != nil {
		logger.Fatalf("register GRPC gateway failed:%s", err)
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	srv := &http.Server{
		Addr:    ":" + cfg.GrpcGateway.Port,
		Handler: mux,
	}

	go func() {
		logger.Infof("GrpcGateway started on port:%s", cfg.GrpcGateway.Port)
		if err := srv.ListenAndServe(); err != nil {
			logger.Fatalf("run GrpcGateway failed:%s", err)
		}
	}()
	/*defer func() {
		if err = srv.Shutdown(context.Background()); err != nil {
			logger.Fatalf("Error shutting down GrpcGateway:%s", err)
		}
	}()*/
}
