package wallet

type Wallet struct {
	UUID       string  `json:"uuid"`
	UserUUID   string  `json:"user_uuid"`
	BalanceEUR float32 `json:"balance_eur"`
	BalanceUSD float32 `json:"balance_usd"`
	BalanceRUB float32 `json:"balance_rub"`
}
