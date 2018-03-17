package main

import (
	"app/models"
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	db, err := models.NewDB(os.Getenv("POSTGRES_URI") + "?sslmode=disable")
	if err != nil {
		log.Panic(err)
	}

	env := &Env{db}

	// Define routes
	mux := mux.NewRouter()
	mux.HandleFunc("/transactions", secured(env.transactions))
	mux.HandleFunc("/transaction/{id}", env.transaction)

	handler := corsHandler().Handler(mux)

	// Handles a production environment
	if os.Getenv("ENV") == "production" {
		certManager := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist("api.bitkit.live"),
			Cache:      autocert.DirCache("certs"),
		}

		server := &http.Server{
			Addr: ":https",
			TLSConfig: &tls.Config{
				GetCertificate: certManager.GetCertificate,
			},
			Handler: handler,
		}

		go http.ListenAndServe(":http", certManager.HTTPHandler(nil))

		log.Fatal(server.ListenAndServeTLS("", ""))
	} else { // Handles the dev environment
		log.Fatal(http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", handler))
	}
}

// Env serves as the context the app runs in
// All route handlers are implemented against this struct
type Env struct {
	db models.Datastore
}

// ********** API Handlers ********** //

// TransactionsRequest models a request with multiple transaction from the mempool
type TransactionsRequest struct {
	Data   []*models.Transaction `json:"data"`
	Method string                `json:"method"`
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

	if txns.Method == "reset" {
		err = env.db.ReplaceTransactions(txns.Data)
	} else {
		err = env.db.InsertTransactions(txns.Data)
	}

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// TransactionResponse represents a single found transaction
type TransactionResponse struct {
	Data *models.Transaction `json:"data"`
}

func (env *Env) transaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	txID := vars["id"]
	tx, err := env.db.GetTransaction(txID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&TransactionResponse{tx})
}

// ********** Helper Functions ********** //
func secured(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok && username == os.Getenv("AUTH_USER") && password == os.Getenv("AUTH_PASSWORD") {
			handler(w, r)
		} else {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
	}
}

func corsHandler() *cors.Cors {
	allowed := "*"
	if os.Getenv("ENV") == "production" {
		allowed = "https://bitkit.live"
	}
	return cors.New(cors.Options{
		AllowedOrigins: []string{allowed},
		AllowedMethods: []string{"GET"},
	})
}
