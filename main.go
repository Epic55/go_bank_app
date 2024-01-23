package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/janirefdez/ArticleRestApi/pkg/db"
	"github.com/janirefdez/ArticleRestApi/pkg/handlers"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the ccount REST API!")
	fmt.Println("ccount REST API")
}

func handleRequests(DB *sql.DB) {
	h := handlers.New(DB)
	// create a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/accounts", h.GetAllAccounts).Methods(http.MethodGet)
	myRouter.HandleFunc("/accounts/{id}", h.GetAccount).Methods(http.MethodGet)
	myRouter.HandleFunc("/accounts/{id}", h.DeleteAccount).Methods(http.MethodDelete)
	myRouter.HandleFunc("/accounts/a/{id}", h.Addition).Methods(http.MethodPut)
	myRouter.HandleFunc("/accounts/s/{id}", h.Subtraction).Methods(http.MethodPut)
	myRouter.HandleFunc("/accounts/t/{id}/{id2}", h.Transfer).Methods(http.MethodPut)
	myRouter.HandleFunc("/accounts/t2/{id}/{id2}", h.Transfer2).Methods(http.MethodPut)

	log.Fatal(http.ListenAndServe("localhost:8080", myRouter))
	fmt.Println("Listening in port 8080")
}

func main() {
	DB := db.Connect()
	db.CreateTable(DB)
	handleRequests(DB)
	db.CloseConnection(DB)
}
