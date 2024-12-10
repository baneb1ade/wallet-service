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

// GetWalletBalanceHandler godoc
// @Summary      Get wallet balance
// @Description  Retrieve the balance of the user's wallet
// @Tags         wallet
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                  true  "Bearer Token"  default(Bearer <token>)
// @Success      200  {object}  CurrenciesResponse
// @Router       /api/v1/wallet/balance/ [get]
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

// UpdateWalletBalanceDeposit godoc
// @Summary      Deposit money into wallet
// @Description  Add a specified amount to the user's wallet
// @Tags         wallet
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                  true  "Bearer Token"  default(Bearer <token>)
// @Param        request  body      ChangeBalanceRequest  true  "Deposit request"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}  "Validation failed"
// @Failure      401      {object}  map[string]string       "userID not found in context"
// @Failure      500      {object}  map[string]string       "internal server error"
// @Router       /api/v1/wallet/deposit/ [post]
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

// UpdateWalletBalanceWithdraw godoc
// @Summary      Withdraw money from wallet
// @Description  Deduct a specified amount from the user's wallet
// @Tags         wallet
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                  true  "Bearer Token"  default(Bearer <token>)
// @Param        request  body      ChangeBalanceRequest  true  "Withdraw request"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}  "Validation failed"
// @Failure      401      {object}  map[string]string       "userID not found in context"
// @Failure      500      {object}  map[string]string       "internal server error"
// @Router       /api/v1/wallet/withdraw/ [post]
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

// GetExchangeRates godoc
// @Summary      Get exchange rates
// @Description  Retrieve the latest exchange rates for supported currencies
// @Tags         exchange
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                  true  "Bearer Token"  default(Bearer <token>)
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]string  "internal server error"
// @Router       /api/v1/exchange/rates/ [get]
func GetExchangeRates(s Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		res, err := s.GetExchangeRates(context.Background())
		if err != nil {
			writeJSONError(c, err)
		}
		c.JSON(http.StatusOK, res)
	}
}

// ExchangeRatesForCurrency godoc
// @Summary      Exchange currency
// @Description  Exchange a specified amount from one currency to another
// @Tags         exchange
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                  true  "Bearer Token"  default(Bearer <token>)
// @Param        request  body      ExchangeRequest  true  "Exchange request"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}  "Validation failed or invalid request"
// @Failure      401      {object}  map[string]string       "user not found"
// @Failure      500      {object}  map[string]string       "internal server error"
// @Router       /api/v1/exchange/ [post]
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
