package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"server/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/handlerfunc"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Handles a production environment
	if os.Getenv("ENV") != "development" {
		lambda.Start(Handler)
	} else { // Handles the dev environment
		router := createRouter()
		handler := corsHandler().Handler(router)
		log.Fatal(http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", handler))
	}
}

var initialized = false
var handlerLambda *handlerfunc.HandlerFuncAdapter

// Handler serves as the endpoint for the aws lambda function
func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if !initialized {
		router := createRouter()
		handler := corsHandler().Handler(router)
		handlerLambda = handlerfunc.New(handler.ServeHTTP)
		initialized = true
	}
	return handlerLambda.Proxy(req)
}

func createRouter() *mux.Router {
	db, err := models.NewDB(os.Getenv("POSTGRES_URI") + "?sslmode=disable")
	if err != nil {
		log.Panic(err)
	}

	env := &Env{db}

	// Define routes
	mux := mux.NewRouter()
	mux.HandleFunc("/transactions", secured(env.transactions))
	mux.HandleFunc("/transactions/random", env.randomTransaction).Methods("GET")
	mux.HandleFunc("/transactions/{id}", env.transaction).Methods("GET")
	return mux
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

func (env *Env) randomTransaction(w http.ResponseWriter, r *http.Request) {
	tx, err := env.db.GetRandomTransaction()
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
	return cors.New(cors.Options{
		AllowedOrigins: []string{allowed},
		AllowedMethods: []string{"GET"},
	})
}
