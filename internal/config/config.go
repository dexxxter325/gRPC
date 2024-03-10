package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	GRPC    GRPC
	Metrics Metrics
	DB      DB
	AUTH    AUTH
}

type GRPC struct {
	Port string
}

type Metrics struct {
	Port string
}

type DB struct {
	Postgres Postgres
}

type Postgres struct {
	Host     string
	Port     string
	User     string
	DbName   string
	Password string
	Sslmode  string
}

type AUTH struct {
	SecretKey       string
	RefreshTokenTTl string
	AccessTokenTTl  string
}

func Init() (*Config, error) {
	return InitByPath("config/local.yml")
}

func InitByPath(configPath string) (*Config, error) {
	var cfg Config

	// Указываем путь к файлу конфигурации
	viper.SetConfigFile(configPath)

	// Читаем конфигурационный файл
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config file failed:%s", err)
	}

	// Заполняем структуру конфигурации значениями из файла
	cfg = Config{
		GRPC: GRPC{
			Port: viper.GetString("grpc.port"),
		},
		Metrics: Metrics{
			Port: viper.GetString("metrics.port"),
		},
		DB: DB{
			Postgres: Postgres{
				Host:     viper.GetString("db.postgres.host"),
				Port:     viper.GetString("db.postgres.port"),
				User:     viper.GetString("db.postgres.user"),
				DbName:   viper.GetString("db.postgres.dbName"),
				Password: viper.GetString("db.postgres.password"),
				Sslmode:  viper.GetString("db.postgres.sslmode"),
			},
		},
		AUTH: AUTH{
			SecretKey:       viper.GetString("auth.secretKey"),
			RefreshTokenTTl: viper.GetString("auth.refreshTokenTTl"),
			AccessTokenTTl:  viper.GetString("auth.accessTokenTTl"),
		},
	}

	return &cfg, nil
}
