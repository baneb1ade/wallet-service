package exchange

import (
	"context"
	"fmt"
	proto "github.com/baneb1ade/exchanger-protos/gen/go"
	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"time"
	"wallet/internal/domain/wallet"
)

type Client struct {
	api proto.ExchangeServiceClient
	log *slog.Logger
}

func New(log *slog.Logger, addr string, timeOut time.Duration, retries uint) (*Client, error) {
	const op = "clients.auth.New"
	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(retries),
		grpcretry.WithPerRetryTimeout(timeOut),
	}

	logOpts := []grpclog.Option{grpclog.WithLogOnEvents(grpclog.PayloadReceived, grpclog.PayloadSent)}

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
	return &Client{api: proto.NewExchangeServiceClient(cc), log: log}, nil
}

func InterceptorLogger(log *slog.Logger) grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, level grpclog.Level, msg string, fields ...any) {
		log.Log(ctx, slog.Level(level), msg, fields...)
	})
}

func (c *Client) GetExchangeRates(ctx context.Context) (wallet.ExchangeRateResponse, error) {
	const op = "grpc.exchange.GetExchangeRates"
	logger := c.log.With(slog.String("op", op))

	res, err := c.api.GetExchangeRates(ctx, &proto.Empty{})
	if err != nil {
		logger.Error(err.Error())
		return wallet.ExchangeRateResponse{}, fmt.Errorf("%s: %w", op, err)
	}
	return wallet.ExchangeRateResponse{
		Rates: res.Rates,
	}, nil
}

func (c *Client) GetExchangeRateForCurrency(ctx context.Context, fromCurrency, toCurrency string) (float32, error) {
	const op = "grpc.exchange.GetExchangeRateForCurrency"
	logger := c.log.With(slog.String("op", op))

	res, err := c.api.GetExchangeRateForCurrency(ctx, &proto.CurrencyRequest{
		FromCurrency: fromCurrency,
		ToCurrency:   toCurrency,
	})
	if err != nil {
		logger.Error(err.Error())
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return res.Rate, nil

}
