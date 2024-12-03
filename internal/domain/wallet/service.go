package wallet

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"strconv"
	"time"
)

type ServiceWallet struct {
	logger           *slog.Logger
	storage          Storage
	cache            Cache
	exchangerService ExchangerService
}

func NewService(storage Storage, logger *slog.Logger, cache Cache, es ExchangerService) Service {
	return &ServiceWallet{
		storage:          storage,
		logger:           logger,
		cache:            cache,
		exchangerService: es,
	}
}

func (s *ServiceWallet) GetExchangeRates(ctx context.Context) (ExchangeRateResponse, error) {
	const op = "wallet.GetExchangeRates"
	log := s.logger.With(slog.String("op", op))

	rates, _ := s.cache.GetValue(context.Background(), "exchange_rates")
	if rates != "" {
		var result ExchangeRateResponse
		err := json.Unmarshal([]byte(rates), &result)
		if err != nil {
			log.Error(err.Error())
		}
		return result, nil
	}
	res, err := s.exchangerService.GetExchangeRates(ctx)
	if err != nil {
		return ExchangeRateResponse{}, err
	}
	jsonData, err := json.Marshal(res)
	if err != nil {
		log.Error(err.Error())
	}
	_ = s.cache.SetValue(ctx, "exchange_rates", string(jsonData), 30*time.Second)
	return res, nil
}

func (s *ServiceWallet) CreateUserWallet(ctx context.Context, userID string) error {
	const op = "wallet.CreateUserWallet"
	log := s.logger.With("op", op)

	if err := s.storage.CreateWallet(ctx, userID); err != nil {
		log.Error(err.Error())
		return ErrSmtWentWrong
	}
	return nil

}

func (s *ServiceWallet) GetBalance(ctx context.Context, userID string) (Wallet, error) {
	const op = "wallet.GetBalance"
	log := s.logger.With("op", op)

	w, err := s.storage.GetWalletByUserID(ctx, userID)
	if err != nil {
		log.Error(err.Error())
		return Wallet{}, ErrSmtWentWrong
	}
	return w, nil
}

func (s *ServiceWallet) WalletDeposit(ctx context.Context, userID string, amount float32, currency string) (Wallet, error) {
	const op = "wallet.WalletDeposit"
	log := s.logger.With("op", op)

	w, err := s.storage.GetWalletByUserID(ctx, userID)
	if err != nil {
		log.Error(err.Error())
		return Wallet{}, ErrSmtWentWrong
	}
	switch currency {
	case "EUR":
		w.BalanceEUR += amount
	case "USD":
		w.BalanceUSD += amount
	case "RUB":
		w.BalanceRUB += amount
	default:
		return Wallet{}, ErrInvalidAmountOrCurrency
	}
	err = s.storage.UpdateWallet(ctx, userID, w)
	if err != nil {
		log.Error(err.Error())
		return Wallet{}, ErrSmtWentWrong
	}

	return w, nil
}

func (s *ServiceWallet) WalletWithdraw(ctx context.Context, userID string, amount float32, currency string) (Wallet, error) {
	const op = "wallet.WalletWithdraw"
	log := s.logger.With("op", op)

	w, err := s.storage.GetWalletByUserID(ctx, userID)
	if err != nil {
		log.Error(err.Error())
		return Wallet{}, err
	}
	switch currency {
	case "EUR":
		if w.BalanceEUR-amount < 0 {
			return Wallet{}, ErrInvalidAmountOrCurrency
		}
		w.BalanceEUR -= amount
	case "USD":
		if w.BalanceUSD-amount < 0 {
			return Wallet{}, ErrInvalidAmountOrCurrency
		}
		w.BalanceUSD -= amount
	case "RUB":
		if w.BalanceRUB-amount < 0 {
			return Wallet{}, ErrInvalidAmountOrCurrency
		}
		w.BalanceRUB -= amount
	default:
		return Wallet{}, ErrInvalidAmountOrCurrency
	}
	err = s.storage.UpdateWallet(ctx, userID, w)
	if err != nil {
		log.Error(err.Error())
		return Wallet{}, ErrSmtWentWrong
	}

	return w, nil
}

func (s *ServiceWallet) ExchangeCurrency(ctx context.Context, userID string, amount float32, fromCurrency, toCurrency string) (ExchangeResponse, error) {
	const op = "wallet.ExchangeCurrency"
	log := s.logger.With("op", op)

	w, err := s.storage.GetWalletByUserID(ctx, userID)
	if err != nil {
		log.Error(err.Error())
		return ExchangeResponse{}, ErrSmtWentWrong
	}
	rate, err := s.getRate(ctx, fromCurrency, toCurrency)
	if err != nil {
		log.Error(err.Error())
		return ExchangeResponse{}, ErrSmtWentWrong
	}
	fromCur, err := getBalanceByCurrency(w, fromCurrency)
	if err != nil {
		log.Error(err.Error())
		return ExchangeResponse{}, ErrSmtWentWrong
	}
	if amount > fromCur {
		return ExchangeResponse{}, ErrNotEnoughFunds
	}
	fromCur -= amount
	toCur, err := getBalanceByCurrency(w, toCurrency)
	if err != nil {
		log.Error(err.Error())
		return ExchangeResponse{}, ErrSmtWentWrong
	}
	toCur = toCur + amount/rate
	if err = updateBalanceByCurrency(&w, fromCurrency, fromCur); err != nil {
		log.Error(err.Error())
		return ExchangeResponse{}, ErrSmtWentWrong
	}
	if err = updateBalanceByCurrency(&w, toCurrency, toCur); err != nil {
		log.Error(err.Error())
		return ExchangeResponse{}, ErrSmtWentWrong
	}
	if err = s.storage.UpdateWallet(ctx, userID, w); err != nil {
		log.Error(err.Error())
		return ExchangeResponse{}, ErrSmtWentWrong
	}
	res := ExchangeResponse{
		Message:         "Exchange successful",
		ExchangedAmount: toCur,
		NewBalance: map[string]float32{
			fromCurrency: fromCur,
			toCurrency:   toCur,
		},
	}
	return res, nil

}

func (s *ServiceWallet) getRate(ctx context.Context, fromCurrency, toCurrency string) (float32, error) {
	const op = "wallet.getRate"
	log := s.logger.With(slog.String("op", op))

	res, _ := s.cache.GetValue(ctx, fmt.Sprintf("exchange_rate:%s:%s", fromCurrency, toCurrency))
	if res != "" {
		parsedRate, err := strconv.ParseFloat(res, 32)
		if err != nil {
			log.Error(err.Error())
		}
		return float32(parsedRate), nil
	}
	rate, err := s.exchangerService.GetExchangeRateForCurrency(ctx, fromCurrency, toCurrency)
	if err != nil {
		log.Error(err.Error())
		return 0, ErrSmtWentWrong
	}
	_ = s.cache.SetValue(ctx, fmt.Sprintf("exchange_rate:%s:%s", fromCurrency, toCurrency), rate, 30*time.Second)
	return rate, nil
}

func getBalanceByCurrency(wallet Wallet, currency string) (float32, error) {
	fieldName := "Balance" + currency
	v := reflect.ValueOf(wallet)
	field := v.FieldByName(fieldName)

	if !field.IsValid() {
		return 0, errors.New("invalid currency: " + currency)
	}

	if field.Kind() != reflect.Float32 {
		return 0, errors.New("field is not of type float32")
	}

	return float32(field.Float()), nil
}

func updateBalanceByCurrency(wallet *Wallet, currency string, newValue float32) error {
	fieldName := "Balance" + currency
	v := reflect.ValueOf(wallet).Elem()
	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		return errors.New("invalid currency: " + currency)
	}

	if !field.CanSet() {
		return errors.New("field cannot be set")
	}

	if field.Kind() != reflect.Float32 {
		return errors.New("field is not of type float32")
	}

	field.SetFloat(float64(newValue))
	return nil
}
