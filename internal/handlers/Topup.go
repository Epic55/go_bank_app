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

func (h *Handler) Topup(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	vars := mux.Vars(r)
	id := vars["id"]
	date := time.Now()
	date1 := date.Format("2006-01-02 15:04:05")

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var changesToAccount models.Account
	json.Unmarshal(body, &changesToAccount)

	queryStmt := `SELECT * FROM accounts WHERE id = $1 ;`
	results, err := h.R.Db.Query(queryStmt, id)
	if err != nil {
		log.Println("failed to execute query", err)
		w.WriteHeader(500)
		return
	}

	var account models.Account //CURRENT INFO ABT ACCOUNT
	for results.Next() {
		err = results.Scan(&account.Id, &account.Name, &account.Account, &account.Balance, &account.Currency, &account.Date, &account.Blocked, &account.Defaultaccount)
		if err != nil {
			log.Println("failed to scan", err)
			w.WriteHeader(500)
			return
		}
	}

	// h.GetInfoAboutAccount(w, id, Account.Name, Account.Account, Account.Currency, Account.Date, Account.Balance, Account.Blocked, Account.Defaultaccount)

	if !account.Blocked {
		updatedBalance := account.Balance + changesToAccount.Balance

		typeofoperation2 := "topuped"
		h.R.UpdateAccount(w, updatedBalance, changesToAccount.Balance, id, account.Currency, typeofoperation2, date1)

		typeofoperation := "topup"
		h.R.UpdateHistory(typeofoperation, account.Name, account.Currency, changesToAccount.Balance, date1)

	} else {
		AccountIsBlocked(w, account.Name, account.Id)
	}
}

func AccountIsBlocked(w http.ResponseWriter, accountName string, accountId int) {
	fmt.Println("Operation is not permitted. Account is blocked. Name -", accountName, "ID -", accountId)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Operation is not permitted. Account is blocked")
}
