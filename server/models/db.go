package models

import (
	"database/sql"

	_ "github.com/lib/pq" // models package is where we do db business
)

// Datastore represents the interface for accessing bitkit data
type Datastore interface {
	InsertTransactions(transactions []*Transaction) error
	ReplaceTransactions(transactions []*Transaction) error
	GetTransaction(id string) (*Transaction, error)
	GetRandomTransaction() (*Transaction, error)
}

// DB is the type all database logic is implemented against
type DB struct {
	*sql.DB
}

// NewDB creates a DB type after opening a connection
func NewDB(pguri string) (*DB, error) {
	db, err := sql.Open("postgres", pguri)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}
