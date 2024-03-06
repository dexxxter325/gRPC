package suite

import (
	"GRPC/gen"
	"GRPC/internal/app"
	"GRPC/internal/config"
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // used by migrator
	_ "github.com/golang-migrate/migrate/v4/source/file"       // used by migrator
	"github.com/sirupsen/logrus"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"golang.org/x/exp/rand"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

const (
	migrationPath = "../migrations"
)

type Suite struct {
	*testing.T
	Cfg              *config.Config
	AuthClient       gen.AuthClient
	InvestmentClient gen.InvestmentClient
}

func New(t *testing.T) (context.Context, *Suite, error) {
	t.Parallel()
	t.Helper() //ошибка будет лучше отображена в отладке
	ctx := context.Background()
	port, err := createPostgresContainer(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("createPostgresContainer failed:%s\n", err)
	}

	if err = updateDBPortInConfig(port); err != nil {
		t.Fatalf("update DBPort failed:%s", err)
	}
	if err = updateGRPCPortInConfig(); err != nil {
		t.Fatalf("update GRPCPort failed:%s", err)
	}

	logger := logrus.New()
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	cfg, err := config.InitByPath("../config/local_tests.yml")
	if err != nil {
		logger.Fatalf("init config failed:%s", err)
	}

	go func() {
		app.Run(logger, cfg)
	}()

	time.Sleep(time.Second * 1) //for stop

	cc, err := grpc.DialContext(ctx, net.JoinHostPort("localhost", cfg.GRPC.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //небезопасное соед.для тестов
	if err != nil {
		return nil, nil, fmt.Errorf("create client failed:%s", err)
	}

	return ctx, &Suite{
		T:                t,
		Cfg:              cfg,
		AuthClient:       gen.NewAuthClient(cc),
		InvestmentClient: gen.NewInvestmentClient(cc),
	}, nil
}

func createPostgresContainer(ctx context.Context) (string, error) {
	pgContainer, err := postgres.RunContainer(
		ctx,
		testcontainers.WithImage("postgres:16.2"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("qwerty"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return "", fmt.Errorf("run Container failed:%s\n", err)
	}
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return "", fmt.Errorf("conn to container failed:%s", err)
	}

	portWithTcp, err := pgContainer.MappedPort(ctx, "5432")
	if err != nil {
		return "", fmt.Errorf("get mapped port failed: %s", err)
	}
	port := strings.TrimSuffix(string(portWithTcp), "/tcp")

	if err = applyMigrations(connStr, migrationPath); err != nil {
		return "", fmt.Errorf("applyMigrations failed:%s", err)
	}
	return port, nil
}

func applyMigrations(connStr, migrationPath string) error {
	migrations, err := migrate.New(
		"file://"+migrationPath,
		connStr,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %v", err)
	}
	defer func() {
		sourceErr, dbErr := migrations.Close() //для освобождения рес-ов и предотвращение их утечки
		if sourceErr != nil {
			logrus.Errorf("close migrations failed:%s", sourceErr)
		}
		if dbErr != nil {
			logrus.Errorf("close migrations failed:%s", dbErr)
		}
	}()

	if err = migrations.Up(); err != nil {
		return fmt.Errorf("failed to apply migrations: %v", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrations re already applied:%s", err.Error())
	}
	return nil
}

func updateDBPortInConfig(DBPort string) error {
	cfg, err := config.InitByPath("../config/local_tests.yml")
	if err != nil {
		return err
	}

	// Чтение файла config.yml
	data, err := os.ReadFile("../config/local_tests.yml")
	if err != nil {
		return err
	}
	// Замена на сгенерированный порт
	updatedDBPort := strings.Replace(string(data), cfg.DB.Postgres.Port, DBPort, 1)

	// Запись обновленных данных в файл
	err = os.WriteFile("../config/local_tests.yml", []byte(updatedDBPort), 0644)
	if err != nil {
		return err
	}
	return nil
}

// updateGRPCPortInConfig создана для того ,чтобы тесты запускались на разных портах и избегали проблемы прослушивания 1 и того же порта
func updateGRPCPortInConfig() error {
	cfg, err := config.InitByPath("../config/local_tests.yml")
	if err != nil {
		return err
	}

	randGRPCPortInt := rand.Intn(8081-8000) + 8000
	randGRPCPort := strconv.Itoa(randGRPCPortInt)

	data, err := os.ReadFile("../config/local_tests.yml")
	if err != nil {
		return err
	}

	updatedGRPCPort := strings.Replace(string(data), cfg.GRPC.Port, randGRPCPort, 1)

	err = os.WriteFile("../config/local_tests.yml", []byte(updatedGRPCPort), 0644)
	if err != nil {
		return err
	}
	return nil
}
