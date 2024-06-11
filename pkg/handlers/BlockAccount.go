package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/epic55/BankApp/pkg/models"
	"github.com/gorilla/mux"
)

func (h handler) BlockAccount(w http.ResponseWriter, r *http.Request) {
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

	queryStmt2 := `UPDATE accounts SET blocked = $2 WHERE id = $1 RETURNING id;`
	err = h.DB.QueryRow(queryStmt2, &id, &changesToAccount.Blocked).Scan(&id)
	fmt.Println("Blocking status is changed on ", changesToAccount.Blocked)
	if err != nil {
		log.Println("failed to execute query", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Blocking status is altered")

}
