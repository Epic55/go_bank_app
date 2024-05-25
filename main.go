package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/epic55/AccountRestApi/pkg/db"
	"github.com/epic55/AccountRestApi/pkg/handlers"
	"github.com/gorilla/mux"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the account REST API!")
	fmt.Println("Account REST API")
}

func handleRequests(DB *sql.DB) {
	h := handlers.New(DB)
	// create a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/accounts", h.GetAllAccounts).Methods(http.MethodGet)
	myRouter.HandleFunc("/accounts/{id}", h.GetAccount).Methods(http.MethodGet)
	myRouter.HandleFunc("/history/{username}", h.GetHistory).Methods(http.MethodGet)
	myRouter.HandleFunc("/accounts/delete/{id}", h.DeleteAccount).Methods(http.MethodDelete)
	myRouter.HandleFunc("/accounts/topup/{id}", h.Topup).Methods(http.MethodPut)
	myRouter.HandleFunc("/accounts/withdraw/{id}", h.Withdraw).Methods(http.MethodPut)
	myRouter.HandleFunc("/accounts/transfer/{id}/{id2}", h.Transfer).Methods(http.MethodPut)
	myRouter.HandleFunc("/accounts/blocking/{id}", h.BlockAccount).Methods(http.MethodPut)

	log.Fatal(http.ListenAndServe("localhost:8080", myRouter))
	fmt.Println("Listening in port 8080")
}

func main() {
	DB := db.Connect()
	db.CreateTable(DB)
	handleRequests(DB)
	db.CloseConnection(DB)
}
