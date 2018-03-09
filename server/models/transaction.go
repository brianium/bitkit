package models

import (
	"github.com/lib/pq"
)

// Transaction represents a transaction in the mempool
type Transaction struct {
	ID      string  `json:"txid"`
	FeeRate float32 `json:"fee_rate"`
	Weight  int     `json:"weight"`
}

// InsertTransactions does a batch insert of all the given transactions
func (db *DB) InsertTransactions(transactions []*Transaction) (err error) {
	txn, err := db.Begin()
	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			txn.Rollback()
			return
		}
		err = txn.Commit()
	}()

	stmt, err := txn.Prepare(pq.CopyInSchema("bitkit", "transactions", "id", "fee_rate", "weight"))
	if err != nil {
		return
	}

	for _, transaction := range transactions {
		_, err = stmt.Exec(transaction.ID, transaction.FeeRate, transaction.Weight)
		if err != nil {
			return
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return
	}

	err = stmt.Close()
	return
}
