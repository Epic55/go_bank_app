package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/epic55/BankApp/internal/models"
	"github.com/gorilla/mux"
)

func (h *Handler) Payments(w http.ResponseWriter, r *http.Request, ctx context.Context) {
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
	results, err := h.R.Db.Query(queryStmt, id)
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

	date := time.Now()
	date1 := date.Format("2006-01-02 15:04:05")

	services := []string{"tele2", "beeline", "kcell"}

	if checkservice(services, changesToPayments.Service) {
		if !account.Blocked {

			if account.Balance >= changesToAccount.Balance {
				updatedBalance := account.Balance - changesToAccount.Balance

				h.R.UpdateAccountPayment(w, updatedBalance, changesToAccount.Balance, id, account.Currency, date1)

				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode("Payment is done")

				h.R.UpdateHistoryPayment(account.Name, account.Currency, date1, changesToPayments.Service, changesToAccount.Balance)

				h.R.UpdatePayments(account.Name, date1, changesToPayments.Service, changesToAccount.Balance, account.Currency)
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
