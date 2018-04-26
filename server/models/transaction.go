package models

import (
	"fmt"
	"strings"
)

// Transaction represents a transaction in the mempool
type Transaction struct {
	ID               string  `json:"txid"`
	FeeRate          float32 `json:"fee_rate"`
	Weight           int     `json:"weight"`
	Fee				 int     `json:"fee"`
	TransactionCount int     `json:"transaction_count"`
	TotalWeight      int     `json:"total_weight"`
}

func insertTransactions(db *DB, transactions []*Transaction) error {
	length := len(transactions)
	valueStrings := make([]string, 0, length)
	valueArgs := make([]interface{}, 0, length*3)
	i := 0
	for _, transaction := range transactions {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d)", i*3+1, i*3+2, i*3+3, i*3+4))
		valueArgs = append(valueArgs, transaction.ID)
		valueArgs = append(valueArgs, transaction.FeeRate)
		valueArgs = append(valueArgs, transaction.Weight)
		valueArgs = append(valueArgs, transaction.Fee)
		i++
	}

	sql := `
	INSERT INTO bitkit.transactions (id, fee_rate, weight, fee)
	VALUES %s
	ON CONFLICT (id)
	DO UPDATE SET
    	fee_rate = EXCLUDED.fee_rate,
		weight = EXCLUDED.weight,
		fee = EXCLUDED.fee
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

// GetTransaction returns a transaction record with a count of transactions ahead of it
// and the total weight of all transactions being queried
func (db *DB) GetTransaction(id string) (*Transaction, error) {
	sql := `
	WITH fr as (
		SELECT id, fee_rate, weight, fee
		FROM bitkit.transactions
		WHERE id = $1
	)
	SELECT fr.id, fr.fee_rate, fr.weight, fr.fee, COUNT(tx.id) - 1 as transaction_count, SUM(tx.weight) - fr.weight as total_weight
	FROM bitkit.transactions as tx
	JOIN fr ON tx.fee_rate >= fr.fee_rate
	GROUP BY fr.id, fr.fee_rate, fr.weight
	`
	stmt, err := db.Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var (
		txID             string
		feeRate          float32
		weight           int
		fee				 int
		transactionCount int
		totalWeight      int
	)
	err = stmt.QueryRow(id).Scan(&txID, &feeRate, &weight, &fee, &transactionCount, &totalWeight)
	if err != nil {
		return nil, err
	}
	return &Transaction{txID, feeRate, weight, fee, transactionCount, totalWeight}, nil
}

// GetRandomTransaction returns a random transaction record with a count of transactions ahead of it
// and the total weight of all transactions being queried
func (db *DB) GetRandomTransaction() (*Transaction, error) {
	sql := `
	WITH fr as (
		SELECT id, fee_rate, weight, fee
		FROM bitkit.transactions
		OFFSET floor(random() * (select count(*)-1 from bitkit.transactions))
		LIMIT 1
	)
	SELECT fr.id, fr.fee_rate, fr.weight, fr.fee, COUNT(tx.id) - 1 as transaction_count, SUM(tx.weight) - fr.weight as total_weight
	FROM bitkit.transactions as tx
	JOIN fr ON tx.fee_rate >= fr.fee_rate
	GROUP BY fr.id, fr.fee_rate, fr.weight
	`
	var (
		txID             string
		feeRate          float32
		weight           int
		fee				 int
		transactionCount int
		totalWeight      int
	)
	err := db.QueryRow(sql).Scan(&txID, &feeRate, &weight, &fee, &transactionCount, &totalWeight)
	if err != nil {
		return nil, err
	}
	return &Transaction{txID, feeRate, weight, fee, transactionCount, totalWeight}, nil
}
