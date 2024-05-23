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

func (h handler) Subtraction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Read request body
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)

	if err != nil {
		log.Fatalln(err)
	}

	var updatedAccount models.Account
	json.Unmarshal(body, &updatedAccount)

	queryStmt := `SELECT * FROM accounts WHERE id = $1 ;`
	results, err := h.DB.Query(queryStmt, id)
	if err != nil {
		log.Println("failed to execute query", err)
		w.WriteHeader(500)
		return
	}

	var account models.Account
	for results.Next() {
		err = results.Scan(&account.Id, &account.Name, &account.Balance, &account.Date)
		if err != nil {
			log.Println("failed to scan", err)
			w.WriteHeader(500)
			return
		}
	}

	date1 := time.Now()

	if account.Balance >= updatedAccount.Balance {
		updatedBalance := account.Balance - updatedAccount.Balance

		queryStmt2 := `UPDATE accounts SET balance = $2, date = $3 WHERE id = $1 RETURNING id;`
		err = h.DB.QueryRow(queryStmt2, &id, &updatedBalance, date1).Scan(&id)
		fmt.Println("Balance is substracted on", updatedAccount.Balance, "Result:", updatedBalance)
		if err != nil {
			log.Println("failed to execute query", err)
			w.WriteHeader(500)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("Balances is updated")
	} else {
		fmt.Println("Not enough money")

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("Not enough money")
	}

}
