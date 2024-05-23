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

func (h handler) Transfer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	id2 := vars["id2"]

	// Read request body
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)

	if err != nil {
		log.Fatalln(err)
	}

	date1 := time.Now()

	var updatedAccountSender models.Account
	json.Unmarshal(body, &updatedAccountSender)

	queryStmt := `SELECT * FROM accounts WHERE id = $1 ;`
	results, err := h.DB.Query(queryStmt, id)
	if err != nil {
		log.Println("failed to execute query", err)
		w.WriteHeader(500)
		return
	}

	var accountSender models.Account
	for results.Next() {
		err = results.Scan(&accountSender.Id, &accountSender.Name, &accountSender.Balance, &account.Currency, &accountSender.Date)
		if err != nil {
			log.Println("failed to scan", err)
			w.WriteHeader(500)
			return
		}
	}

	if accountSender.Balance >= updatedAccountSender.Balance {
		updatedBalanceSender := accountSender.Balance - updatedAccountSender.Balance

		queryStmt2 := `UPDATE accounts SET balance = $2, currency = $4, date = $3  WHERE id = $1 RETURNING id;`
		err = h.DB.QueryRow(queryStmt2, &id, &updatedBalanceSender, &account.Currency, date1).Scan(&id)
		fmt.Println("Sender balance is substracted on", updatedAccountSender.Balance, "Result:", updatedBalanceSender)
		if err != nil {
			log.Println("failed to execute query", err)
			w.WriteHeader(500)
			return
		}

		//2ND USER
		var updatedAccountReceiver models.Account
		json.Unmarshal(body, &updatedAccountReceiver)

		queryStmt3 := `SELECT * FROM accounts WHERE id = $1 ;`
		results2, err := h.DB.Query(queryStmt3, id2)
		if err != nil {
			log.Println("failed to execute query", err)
			w.WriteHeader(500)
			return
		}

		var accountReceiver models.Account
		for results2.Next() {
			err = results2.Scan(&accountReceiver.Id, &accountReceiver.Name, &accountReceiver.Balance, &account.Currency, &accountReceiver.Date)
			if err != nil {
				log.Println("failed to scan", err)
				w.WriteHeader(500)
				return
			}
		}

		updatedBalanceReceiver := accountReceiver.Balance + updatedAccountReceiver.Balance

		queryStmt4 := `UPDATE accounts SET balance = $2, currency = $4, date = $3 WHERE id = $1 RETURNING id;`
		err = h.DB.QueryRow(queryStmt4, &id2, &updatedBalanceReceiver, &account.Currency, date1).Scan(&id2)
		fmt.Println("Receiver balance is topped up on", updatedAccountReceiver.Balance, "Result:", updatedBalanceReceiver)
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
