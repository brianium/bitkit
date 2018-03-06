package main

import (
	"app/models"
	"encoding/json"
	"log"
	"net/http"
)

// Higher order handler for ensuring the request method is POST
func postHandler(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			f(w, r)
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	}
}

// TransactionsRequest models a request with multiple transaction from the mempool
type TransactionsRequest struct {
	Data []models.Transaction `json:"data"`
}

func transactions(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var txns TransactionsRequest
	err := decoder.Decode(&txns)

	if err != nil {
		panic(err)
	}

	txnsJSON, err := json.Marshal(txns)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(txnsJSON)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/transactions", postHandler(transactions))
	log.Fatal(http.ListenAndServe(":8080", mux))
}
