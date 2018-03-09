package main

import (
	"app/models"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

// Env serves as the context the app runs in
// All route handlers are implemented against this struct
type Env struct {
	db models.Datastore
}

func main() {
	db, err := models.NewDB(os.Getenv("POSTGRES_URI") + "?sslmode=disable")
	if err != nil {
		log.Panic(err)
	}

	env := &Env{db}

	mux := http.NewServeMux()
	mux.HandleFunc("/transactions", env.transactions)
	log.Fatal(http.ListenAndServe(":8080", mux))
}

// TransactionsRequest models a request with multiple transaction from the mempool
type TransactionsRequest struct {
	Data []*models.Transaction `json:"data"`
}

func (env *Env) transactions(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	decoder := json.NewDecoder(r.Body)

	var txns TransactionsRequest
	err := decoder.Decode(&txns)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	err = env.db.InsertTransactions(txns.Data)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
