package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/epic55/BankApp/internal/models"
)

func (h *Handler) GetAllAccounts(w http.ResponseWriter, r *http.Request, ctx context.Context) {

	results, err := h.R.Db.Query("SELECT * FROM accounts;")
	if err != nil {
		log.Println("failed to execute query", err)
		w.WriteHeader(500)
		return
	}

	var accounts = make([]models.Account, 0)
	for results.Next() {
		var account models.Account
		err = results.Scan(&account.Id, &account.Name, &account.Account, &account.Balance, &account.Currency, &account.Date, &account.Blocked, &account.Defaultaccount)
		if err != nil {
			log.Println("failed to scan", err)
			w.WriteHeader(500)
			return
		}

		accounts = append(accounts, account)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accounts)
}
