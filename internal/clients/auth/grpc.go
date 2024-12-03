package auth

import (
	"context"
	"fmt"
	v3 "github.com/baneb1ade/auth-protos/gen/go"
	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"time"
	"wallet/internal/domain/auth"
)

type Client struct {
	api v3.AuthClient
	log *slog.Logger
}

func New(log *slog.Logger, addr string, timeOut time.Duration, retries uint) (*Client, error) {
	const op = "clients.auth.New"

	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(retries),
		grpcretry.WithPerRetryTimeout(timeOut),
	}

	logOpts := []grpclog.Option{grpclog.WithLogOnEvents(grpclog.PayloadReceived)}

	cc, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpclog.UnaryClientInterceptor(InterceptorLogger(log), logOpts...),
			grpcretry.UnaryClientInterceptor(retryOpts...),
		),
	)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Client{api: v3.NewAuthClient(cc), log: log}, nil
}

func InterceptorLogger(log *slog.Logger) grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, level grpclog.Level, msg string, fields ...any) {
		log.Log(ctx, slog.Level(level), msg, fields...)
	})
}

func (c *Client) Register(ctx context.Context, wc auth.WalletCreator, email, username, password string) (string, error) {
	const op = "grpc.Register"
	log := c.log.With(slog.With("op", op))

	res, err := c.api.Register(ctx, &v3.RegisterRequest{
		Email:    email,
		Username: username,
		Password: password,
	})
	if err != nil {
		log.Error(err.Error())
		return "", fmt.Errorf("%s: %w", op, err)
	}
	err = wc.CreateUserWallet(ctx, res.UserId)
	if err != nil {
		log.Error(err.Error())
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return res.UserId, nil
}

func (c *Client) Login(ctx context.Context, username, password string) (string, error) {
	const op = "grpc.Login"
	log := c.log.With(slog.With("op", op))

	res, err := c.api.Login(ctx, &v3.LoginRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		log.Error(err.Error())
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return res.Token, nil
}
