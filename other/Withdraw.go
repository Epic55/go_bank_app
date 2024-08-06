package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/epic55/BankApp/internal/models"
	"github.com/gorilla/mux"
)

func (h *Handler) Withdraw(w http.ResponseWriter, r *http.Request) {
	var m sync.Mutex

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
	results, err := h.R.DB.Query(queryStmt, id)
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
	if !account.Blocked {

		if account.Balance >= changesToAccount.Balance {
			updatedBalance := account.Balance - changesToAccount.Balance

			typeofoperation2 := "withdrawed"
			m.Lock()
			h.R.UpdateAccount(w, updatedBalance, changesToAccount.Balance, id, account.Currency, typeofoperation2, date1)
			m.Unlock()

			typeofoperation := "withdraw"
			h.R.UpdateHistory2(typeofoperation, account.Name, account.Currency, changesToAccount.Balance, date1)

		} else {
			NotEnoughMoney(w)
		}

	} else {
		AccountIsBlocked(w, account.Name, account.Id)
	}

}
