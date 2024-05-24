package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/epic55/AccountRestApi/pkg/models"
	"github.com/gorilla/mux"
)

func (h handler) Topup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Read request body
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)

	if err != nil {
		log.Fatalln(err)
	}

	var changesToAccount models.Account
	json.Unmarshal(body, &changesToAccount)

	queryStmt := `SELECT * FROM accounts WHERE id = $1 ;`
	results, err := h.DB.Query(queryStmt, id)
	//fmt.Println(results)
	if err != nil {
		log.Println("failed to execute query", err)
		w.WriteHeader(500)
		return
	}

	var account models.Account //CURRENT INFO ABT ACCOUNT
	for results.Next() {
		err = results.Scan(&account.Id, &account.Name, &account.Balance, &account.Currency, &account.Date, &account.Blocked)
		if err != nil {
			log.Println("failed to scan", err)
			w.WriteHeader(500)
			return
		}
	}
	fmt.Println("Account - ", account)
	fmt.Println("changesToAccount - ", changesToAccount)

	date1 := time.Now()
	if account.Blocked == false {
		updatedBalance := account.Balance + changesToAccount.Balance

		queryStmt2 := `UPDATE accounts SET balance = $2, currency = $3, date = $4 WHERE id = $1 RETURNING id;`
		err = h.DB.QueryRow(queryStmt2, &id, &updatedBalance, &account.Currency, date1).Scan(&id)
		fmt.Println("Balance is added on ", changesToAccount.Balance, "Result: ", updatedBalance)
		if err != nil {
			log.Println("failed to execute query", err)
			w.WriteHeader(500)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("Balance is added on")

	} else {
		fmt.Println("Operation is not permitted. Account is blocked. Name -", account.Name, "ID -", account.Id)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("Operation is not permitted. Account is blocked")
	}

}