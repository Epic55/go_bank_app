package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/epic55/BankApp/pkg/models"
	"github.com/gorilla/mux"
)

func (h handler) GetAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	queryStmt := `SELECT name, account, balance, currency, date, blocked, defaultaccount FROM accounts WHERE name = $1 ;`
	results, err := h.DB.Query(queryStmt, id)
	if err != nil {
		log.Println("failed to execute query", err)
		w.WriteHeader(500)
		return
	}

	var accounts = make([]models.Account, 0)
	for results.Next() {
		var account models.Account
		err = results.Scan(&account.Name, &account.Account, &account.Balance, &account.Currency, &account.Date, &account.Blocked, &account.Defaultaccount)
		if err != nil {
			log.Println("failed to scan", err)
			w.WriteHeader(500)
			return
		}

		accounts = append(accounts, account)
	}

	// var account models.Account
	// for results.Next() {
	// 	err = results.Scan(&account.Id, &account.Name, &account.Account, &account.Balance, &account.Currency, &account.Date, &account.Blocked, &account.Defaultaccount)
	// 	//fmt.Println(results)
	// 	if err != nil {
	// 		log.Println("failed to scan", err)
	// 		w.WriteHeader(500)
	// 		return
	// 	}
	// }

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accounts)
}
