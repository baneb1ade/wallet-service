package config

import (
	"github.com/joho/godotenv"
	"os"
	"time"
)

type Config struct {
	Server  ServerConfig
	Storage StorageConfig
	Cache   CacheConfig
	Clients Clients
	Secret  string
}

type CacheConfig struct {
	Addr     string
	Username string
	Password string
	DB       string
}

type StorageConfig struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
}

type ServerConfig struct {
	Address string
	Port    string
}

type Clients struct {
	Auth     AuthClientConfig
	Exchange ExchangeClientConfig
}

type AuthClientConfig struct {
	Address string
	Timeout time.Duration
	Retries uint
}

type ExchangeClientConfig struct {
	Address string
	Timeout time.Duration
	Retries uint
}

func getEnvWithDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func MustLoad(cfgPath string) *Config {
	if cfgPath != "" {
		err := godotenv.Load(cfgPath)
		if err != nil {
			panic(err)
		}
	}

	config := Config{
		Secret: getEnvWithDefault("SECRET", "secret"),
		Server: ServerConfig{
			Address: getEnvWithDefault("SERVER_ADDRESS", "0.0.0.0"),
			Port:    getEnvWithDefault("SERVER_PORT", "8080"),
		},
		Storage: StorageConfig{
			DBUser:     getEnvWithDefault("POSTGRES_USER", "postgres"),
			DBPassword: getEnvWithDefault("POSTGRES_PASSWORD", "password"),
			DBHost:     getEnvWithDefault("POSTGRES_HOST", "localhost"),
			DBPort:     getEnvWithDefault("POSTGRES_PORT", "5432"),
			DBName:     getEnvWithDefault("POSTGRES_DB", "mydatabase"),
		},
		Cache: CacheConfig{
			Addr:     getEnvWithDefault("REDIS_ADDRESS", "redis:6379"),
			Username: getEnvWithDefault("REDIS_USER", "redis"),
			Password: getEnvWithDefault("REDIS_PASSWORD", "redis"),
			DB:       getEnvWithDefault("REDIS_DB", "0"),
		},
		Clients: Clients{
			Auth: AuthClientConfig{
				Address: getEnvWithDefault("AUTH_GRPC_ADDR", "auth:44045"),
				Timeout: 5 * time.Second,
				Retries: 5,
			},
			Exchange: ExchangeClientConfig{
				Address: getEnvWithDefault("EXCHANGE_GRPC_ADDR", "exchanger:44044"),
				Timeout: 5 * time.Second,
				Retries: 5,
			},
		},
	}

	return &config
}
