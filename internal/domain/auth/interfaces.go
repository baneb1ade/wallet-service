package auth

import "context"

type ServiceAuth interface {
	Register(ctx context.Context, wc WalletCreator, email, username, password string) (string, error)
	Login(ctx context.Context, username, password string) (string, error)
}

type WalletCreator interface {
	CreateUserWallet(ctx context.Context, userID string) error
}
