package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/janirefdez/ArticleRestApi/pkg/models"
)

func (h handler) GetAllAccounts(w http.ResponseWriter, r *http.Request) {

	results, err := h.DB.Query("SELECT * FROM accounts;")
	if err != nil {
		log.Println("failed to execute query", err)
		w.WriteHeader(500)
		return
	}

	var accounts = make([]models.Account, 0)
	for results.Next() {
		var account models.Account
		err = results.Scan(&account.Id, &account.Name, &account.Balance, &account.Date)
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
