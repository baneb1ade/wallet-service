package wallet

type UpdateWalletRequest struct {
	Amount   float64 `json:"amount" validate:"required"`
	Currency float64 `json:"currency" validate:"required"`
}

type CurrenciesResponse struct {
	EUR float32 `json:"EUR"`
	USD float32 `json:"USD"`
	RUB float32 `json:"RUB"`
}

type ChangeBalanceRequest struct {
	Amount   float32 `json:"amount"`
	Currency string  `json:"currency" validate:"required,oneof=USD EUR RUB"`
}

type ExchangeRateResponse struct {
	Rates map[string]float32 `json:"rates"`
}

type ExchangeRequest struct {
	FromCurrency string  `json:"from_currency" validate:"required,oneof=USD EUR RUB"`
	ToCurrency   string  `json:"to_currency" validate:"required,oneof=USD EUR RUB"`
	Amount       float32 `json:"amount"`
}

type ExchangeResponse struct {
	Message         string             `json:"message"`
	ExchangedAmount float32            `json:"exchanged_amount"`
	NewBalance      map[string]float32 `json:"new_balance"`
}
