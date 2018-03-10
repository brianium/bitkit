package models

import (
	"fmt"
	"strings"
)

// Transaction represents a transaction in the mempool
type Transaction struct {
	ID      string  `json:"txid"`
	FeeRate float32 `json:"fee_rate"`
	Weight  int     `json:"weight"`
}

func insertTransactions(db *DB, transactions []*Transaction) error {
	length := len(transactions)
	valueStrings := make([]string, 0, length)
	valueArgs := make([]interface{}, 0, length*3)
	i := 0
	for _, transaction := range transactions {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d)", i*3+1, i*3+2, i*3+3))
		valueArgs = append(valueArgs, transaction.ID)
		valueArgs = append(valueArgs, transaction.FeeRate)
		valueArgs = append(valueArgs, transaction.Weight)
		i++
	}

	sql := `
	INSERT INTO bitkit.transactions (id, fee_rate, weight)
	VALUES %s
	ON CONFLICT (id)
	DO UPDATE SET
    	fee_rate = EXCLUDED.fee_rate,
    	weight = EXCLUDED.weight
	`
	stmt := fmt.Sprintf(sql, strings.Join(valueStrings, ","))
	_, err := db.Exec(stmt, valueArgs...)
	return err
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

	err = insertTransactions(db, transactions)

	return
}

// ReplaceTransactions drops all transaction records and then does a batch insert
func (db *DB) ReplaceTransactions(transactions []*Transaction) (err error) {
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

	// drop records first
	_, err = db.Exec("DELETE FROM bitkit.transactions")
	if err != nil {
		return
	}

	err = insertTransactions(db, transactions)

	return
}
