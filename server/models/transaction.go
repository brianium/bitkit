package models

// Transaction represents a transaction in the mempool
type Transaction struct {
	ID      string  `json:"txid"`
	FeeRate float32 `json:"fee_rate"`
	Weght   int     `json:"weight"`
}
