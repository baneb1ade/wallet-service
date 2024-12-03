package wallet

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"time"
)

type PsqlClient interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

type Storage interface {
	CreateWallet(ctx context.Context, userID string) error
	UpdateWallet(ctx context.Context, userID string, wallet Wallet) error
	GetWalletByUserID(ctx context.Context, UserID string) (Wallet, error)
}

type Cache interface {
	GetValue(ctx context.Context, key string) (string, error)
	SetValue(ctx context.Context, key string, value interface{}, ttl time.Duration) error
}

type Service interface {
	GetBalance(ctx context.Context, userID string) (Wallet, error)
	WalletDeposit(ctx context.Context, userID string, amount float32, currency string) (Wallet, error)
	WalletWithdraw(ctx context.Context, userID string, amount float32, currency string) (Wallet, error)
	CreateUserWallet(ctx context.Context, userID string) error
	ExchangeCurrency(ctx context.Context, userID string, amount float32, fromCurrency, toCurrency string) (ExchangeResponse, error)
	GetExchangeRates(ctx context.Context) (ExchangeRateResponse, error)
}

type ExchangerService interface {
	GetExchangeRateForCurrency(ctx context.Context, fromCurrency, toCurrency string) (float32, error)
	GetExchangeRates(ctx context.Context) (ExchangeRateResponse, error)
}
