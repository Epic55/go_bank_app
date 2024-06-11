package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/epic55/BankApp/pkg/models"
	"github.com/gorilla/mux"
)

func (h handler) Payments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)

	if err != nil {
		log.Fatalln(err)
	}

	var changesToAccount models.Account
	json.Unmarshal(body, &changesToAccount)

	var changesToPayments models.Payments
	json.Unmarshal(body, &changesToPayments)

	queryStmt := `SELECT * FROM accounts WHERE id = $1 ;`
	results, err := h.DB.Query(queryStmt, id)
	if err != nil {
		log.Println("failed to execute query", err)
		w.WriteHeader(500)
		return
	}

	var account models.Account
	for results.Next() {
		err = results.Scan(&account.Id, &account.Name, &account.Account, &account.Balance, &account.Currency, &account.Date, &account.Blocked, &account.Defaultaccount)
		if err != nil {
			log.Println("failed to scan", err)
			w.WriteHeader(500)
			return
		}
	}

	date1 := time.Now()

	services := []string{"tele2", "beeline", "kcell"}

	if checkservice(services, changesToPayments.Service) {
		if !account.Blocked {

			if account.Balance >= changesToAccount.Balance {
				updatedBalance := account.Balance - changesToAccount.Balance

				queryStmt2 := `UPDATE accounts SET balance = $2, currency = $3, date = $4 WHERE id = $1 RETURNING id;`
				err = h.DB.QueryRow(queryStmt2, &id, &updatedBalance, &account.Currency, date1).Scan(&id)
				fmt.Println("Balance is substracted on", changesToAccount.Balance, "Result:", updatedBalance)
				if err != nil {
					log.Println("failed to execute query", err)
					w.WriteHeader(500)
					return
				}

				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode("Payment is done")

				typeofoperation := "payment to "
				queryStmt3 := `INSERT INTO history (username, date, quantity, currency, typeofoperation) VALUES ($1, $2, $3, $4, $5);`
				_, err = h.DB.Exec(queryStmt3, account.Name, date1, changesToAccount.Balance, account.Currency, typeofoperation+changesToPayments.Service) //USE Exec FOR INSERT
				if err != nil {
					log.Println("failed to execute query - update history:", err)
					return
				} else {
					fmt.Println("History is updated")
				}

				queryStmt4 := `INSERT INTO payments (username, date, service, quantity, currency) VALUES ($1, $2, $3, $4, $5);`
				_, err = h.DB.Exec(queryStmt4, account.Name, date1, changesToPayments.Service, changesToAccount.Balance, account.Currency) //USE Exec FOR INSERT
				if err != nil {
					log.Println("failed to execute query - update payments:", err)
					return
				} else {
					fmt.Println("Payments is updated")
				}

			} else {
				NotEnoughMoney(w)
			}

		} else {
			AccountIsBlocked(w, account.Name, account.Id)
		}

	} else {
		fmt.Println("No such service in application")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("No such service in application")
	}

}

func checkservice(services []string, item string) bool { //CHECK FOR EXISTING SERVICE
	for _, v := range services {
		if v == item {
			return true
		}
	}
	return false
}
