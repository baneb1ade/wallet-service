package tests

import (
	"context"
	"github.com/joho/godotenv"
	"os"
	"testing"
	"wallet/internal/domain/wallet"
	"wallet/internal/domain/wallet/db"
	"wallet/pkg/clients/psql"
	"wallet/pkg/logger"
)

type TestStorageConfig struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
}

func loadTestConfig(t *testing.T) *TestStorageConfig {
	err := godotenv.Load("../../../../config.env")
	if err != nil {
		t.Fatal("Error loading .env file")
	}
	requiredEnvVars := []string{
		"POSTGRES_USER",
		"POSTGRES_PASSWORD",
		"POSTGRES_HOST",
		"POSTGRES_PORT",
		"POSTGRES_DB",
		"SERVER_ADDRESS",
		"SERVER_PORT",
	}

	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			t.Fatalf("%s must be set", envVar)
		}
	}
	return &TestStorageConfig{
		DBUser:     os.Getenv("POSTGRES_USER"),
		DBPassword: os.Getenv("POSTGRES_PASSWORD"),
		DBHost:     os.Getenv("POSTGRES_HOST"),
		DBPort:     os.Getenv("POSTGRES_PORT"),
		DBName:     os.Getenv("POSTGRES_DB"),
	}
}

func TestService(t *testing.T) {
	cfg := loadTestConfig(t)
	if cfg.DBHost == "postgres" {
		cfg.DBHost = "localhost"
	}
	psqlClient, err := psql.NewClient(context.Background(), psql.PostgresConfig{
		Addr:     cfg.DBHost,
		Port:     cfg.DBPort,
		Username: cfg.DBUser,
		Password: cfg.DBPassword,
		Database: cfg.DBName,
	})
	if err != nil {
		t.Fatal(err)
	}
	const userID = "fcdf2063-6ea1-4897-b664-6082af67e2f1"
	ctx := context.Background()
	qi := `INSERT INTO "user"(id, email, username, password) VALUES ($1, $2, $3, $4)`
	qd := `DELETE FROM "user" WHERE id = $1`
	_, err = psqlClient.Exec(ctx, qi, userID, "test@gmail.com", "username", "password")
	if err != nil {
		t.Fatal(err)
	}
	log := logger.SetupLogger(logger.Local, "")
	storage := db.NewRepository(psqlClient, log)
	newWallet := wallet.Wallet{
		UserUUID:   userID,
		BalanceEUR: 100.12,
		BalanceUSD: 42.12,
		BalanceRUB: 123.1,
	}

	t.Run("Create One", func(t *testing.T) {
		err := storage.CreateWallet(ctx, userID)
		if err != nil {
			t.Fatal(err)
		}
	})
	t.Run("Update One", func(t *testing.T) {
		err := storage.UpdateWallet(ctx, userID, newWallet)
		if err != nil {
			t.Fatal(err)
		}
	})
	t.Run("Get One", func(t *testing.T) {
		w, err := storage.GetWalletByUserID(ctx, userID)
		if err != nil {
			t.Fatal(err)
		}

		if w.UserUUID != newWallet.UserUUID {
			t.Fatalf("want %s, got %s", userID, w.UserUUID)
		}
		if w.BalanceEUR != newWallet.BalanceEUR {
			t.Fatalf("want %f, got %f", newWallet.BalanceEUR, w.BalanceEUR)
		}
		if w.BalanceUSD != newWallet.BalanceUSD {
			t.Fatalf("want %f, got %f", newWallet.BalanceUSD, w.BalanceUSD)
		}
		if w.BalanceRUB != newWallet.BalanceRUB {
			t.Fatalf("want %f, got %f", newWallet.BalanceRUB, w.BalanceRUB)
		}
	})
	_, err = psqlClient.Exec(ctx, qd, userID)
	if err != nil {
		t.Fatal(err)
	}
}
