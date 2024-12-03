package wallet

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type ExchangeRateGetter interface {
	GetExchangeRateForCurrency(ctx context.Context, fromCurrency, toCurrency string) (float32, error)
}

func GetWalletBalanceHandler(s Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			return
		}
		userIDStr, ok := userID.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": ErrSmtWentWrong.Error()})
			return
		}
		w, err := s.GetBalance(context.Background(), userIDStr)
		if err != nil {
			writeJSONError(c, err)
			return
		}

		balance := CurrenciesResponse{
			EUR: w.BalanceEUR,
			USD: w.BalanceUSD,
			RUB: w.BalanceRUB,
		}
		res := map[string]CurrenciesResponse{
			"balance": balance,
		}
		c.JSON(http.StatusOK, res)
	}
}

func UpdateWalletBalanceDeposit(s Service, v *validator.Validate) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req ChangeBalanceRequest
		if err := c.ShouldBind(&req); err != nil {
			writeJSONError(c, err)
			return
		}

		if err := v.Struct(req); err != nil {
			var validationErrors validator.ValidationErrors
			errors.As(err, &validationErrors)
			invalidFields := make([]string, len(validationErrors))

			for i, fieldError := range validationErrors {
				invalidFields[i] = fieldError.Field()
			}

			c.JSON(http.StatusBadRequest, gin.H{
				"error":  "Validation failed",
				"fields": invalidFields,
			})
			return
		}
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "userID not found in context"})
			return
		}
		userIDStr, ok := userID.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "userID is not a valid string"})
			return
		}
		w, err := s.WalletDeposit(context.Background(), userIDStr, req.Amount, req.Currency)
		if err != nil {
			writeJSONError(c, err)
			return
		}
		balance := CurrenciesResponse{
			EUR: w.BalanceEUR,
			USD: w.BalanceUSD,
			RUB: w.BalanceRUB,
		}
		res := map[string]interface{}{
			"message":     "Account topped up successfully",
			"new_balance": balance,
		}
		c.JSON(http.StatusOK, res)
	}
}

func UpdateWalletBalanceWithdraw(s Service, v *validator.Validate) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req ChangeBalanceRequest
		if err := c.ShouldBind(&req); err != nil {
			writeJSONError(c, err)
			return
		}

		if err := v.Struct(req); err != nil {
			var validationErrors validator.ValidationErrors
			errors.As(err, &validationErrors)
			invalidFields := make([]string, len(validationErrors))

			for i, fieldError := range validationErrors {
				invalidFields[i] = fieldError.Field()
			}

			c.JSON(http.StatusBadRequest, gin.H{
				"error":  "Validation failed",
				"fields": invalidFields,
			})
			return
		}
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "userID not found in context"})
			return
		}
		userIDStr, ok := userID.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "userID is not a valid string"})
			return
		}
		w, err := s.WalletWithdraw(context.Background(), userIDStr, req.Amount, req.Currency)
		if err != nil {
			writeJSONError(c, err)
			return
		}
		balance := CurrenciesResponse{
			EUR: w.BalanceEUR,
			USD: w.BalanceUSD,
			RUB: w.BalanceRUB,
		}
		res := map[string]interface{}{
			"message":     "Withdrawal successful",
			"new_balance": balance,
		}
		c.JSON(http.StatusOK, res)
	}
}

func GetExchangeRates(s Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		res, err := s.GetExchangeRates(context.Background())
		if err != nil {
			writeJSONError(c, err)
		}
		c.JSON(http.StatusOK, res)
	}
}

func ExchangeRatesForCurrency(s Service, v *validator.Validate) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req ExchangeRequest
		if err := c.ShouldBind(&req); err != nil {
			writeJSONError(c, err)
			return
		}

		if err := v.Struct(req); err != nil {
			var validationErrors validator.ValidationErrors
			errors.As(err, &validationErrors)
			invalidFields := make([]string, len(validationErrors))

			for i, fieldError := range validationErrors {
				invalidFields[i] = fieldError.Field()
			}

			c.JSON(http.StatusBadRequest, gin.H{
				"error":  "Validation failed",
				"fields": invalidFields,
			})
			return
		}
		if req.ToCurrency == req.FromCurrency {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			return
		}
		userIDStr, ok := userID.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": ErrSmtWentWrong.Error()})
			return
		}

		r, err := s.ExchangeCurrency(
			context.Background(),
			userIDStr,
			req.Amount,
			req.FromCurrency,
			req.ToCurrency,
		)
		if err != nil {
			writeJSONError(c, err)
			return
		}
		c.JSON(http.StatusOK, r)
	}
}

func writeJSONError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrSmtWentWrong):
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	case errors.Is(err, ErrInvalidAmountOrCurrency):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, ErrNotEnoughFunds):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
