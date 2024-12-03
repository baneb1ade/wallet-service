package app

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log/slog"
	authClient "wallet/internal/clients/auth"
	exchangeClient "wallet/internal/clients/exchange"
	"wallet/internal/config"
	"wallet/internal/domain/auth"
	wallet2 "wallet/internal/domain/wallet"
	"wallet/internal/domain/wallet/db"
	"wallet/pkg/clients/psql"
	"wallet/pkg/clients/redis"
)

type App struct {
	config *config.Config
	logger *slog.Logger
	router *gin.Engine
}

func NewApp(logger *slog.Logger, cfg *config.Config) (*App, error) {
	logger.Info("Initializing application")
	logger.Info("Connect to Postgresql")
	c, err := psql.NewClient(context.Background(), psql.PostgresConfig{
		Addr:     cfg.Storage.DBHost,
		Port:     cfg.Storage.DBPort,
		Username: cfg.Storage.DBUser,
		Password: cfg.Storage.DBPassword,
		Database: cfg.Storage.DBName,
	})
	logger.Info("Connect to Redis")
	rdb, err := redis.NewClient(redis.ConfigRedis{
		Addr:     cfg.Cache.Addr,
		Username: cfg.Cache.Username,
		Password: cfg.Cache.Password,
		DB:       cfg.Cache.DB,
	})
	if err != nil {
		return nil, err
	}
	repo := db.NewRepository(c, logger)
	cache := db.NewCache(logger, rdb)

	authGRPC, err := authClient.New(
		logger,
		cfg.Clients.Auth.Address,
		cfg.Clients.Auth.Timeout,
		cfg.Clients.Auth.Retries,
	)
	exchangeGRPC, err := exchangeClient.New(
		logger,
		cfg.Clients.Exchange.Address,
		cfg.Clients.Exchange.Timeout,
		cfg.Clients.Exchange.Retries,
	)
	s := wallet2.NewService(repo, logger, cache, exchangeGRPC)
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	v := validator.New()
	apiV1 := r.Group("/api/v1/")
	walletGroup := apiV1.Group("/wallet")
	authGroup := apiV1.Group("/auth")
	exchangeGroup := apiV1.Group("/exchange")

	walletGroup.Use(auth.AuthorizationMiddleware([]byte(cfg.Secret)))

	walletGroup.GET("/balance/", wallet2.GetWalletBalanceHandler(s))
	walletGroup.POST("/deposit/", wallet2.UpdateWalletBalanceDeposit(s, v))
	walletGroup.POST("/withdraw/", wallet2.UpdateWalletBalanceWithdraw(s, v))

	authGroup.POST("/register/", auth.Register(authGRPC, s, v))
	authGroup.POST("/login/", auth.Login(authGRPC, v))

	exchangeGroup.Use(auth.AuthorizationMiddleware([]byte(cfg.Secret)))
	exchangeGroup.POST("/", wallet2.ExchangeRatesForCurrency(s, v))
	exchangeGroup.GET("/rates/", wallet2.GetExchangeRates(s))

	var app = &App{
		config: cfg,
		logger: logger,
		router: r,
	}
	return app, nil
}

func (app *App) Start() {
	app.logger.Info("Start HTTP server")
	errChan := make(chan error)
	serverAddr := fmt.Sprintf("%s:%s", app.config.Server.Address, app.config.Server.Port)

	go func() {
		if err := app.router.Run(serverAddr); err != nil {
			app.logger.Error("Failed to start HTTP server", "error", err)
			errChan <- err
		}
	}()
	select {
	case err := <-errChan:
		if err != nil {
			app.logger.Error("Failed to start HTTP server", "error", err)
		}
	}
}
