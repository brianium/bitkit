package main

import (
	"app/models"
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"golang.org/x/crypto/acme/autocert"
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
	mux.HandleFunc("/transactions", secured(env.transactions))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world"))
	})
	if os.Getenv("ENV") == "production" {
		certManager := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist("api.bitkit.live"), //Your domain here
			Cache:      autocert.DirCache("certs"),                //Folder for storing certificates
		}

		server := &http.Server{
			Addr: ":https",
			TLSConfig: &tls.Config{
				GetCertificate: certManager.GetCertificate,
			},
			Handler: mux,
		}

		go http.ListenAndServe(":http", certManager.HTTPHandler(nil))

		log.Fatal(server.ListenAndServeTLS("", ""))
	} else {
		log.Fatal(http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", mux))
	}
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
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
