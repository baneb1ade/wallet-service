package wallet

import "errors"

var ErrInvalidAmountOrCurrency = errors.New("invalid amount or currency")
var ErrSmtWentWrong = errors.New("something went wrong")
var ErrNotEnoughFunds = errors.New("insufficient funds or invalid currencies")
