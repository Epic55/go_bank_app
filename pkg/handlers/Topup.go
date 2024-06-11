package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/epic55/AccountRestApi/pkg/models"
	"github.com/gorilla/mux"
)

func (h handler) Topup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	date1 := time.Now()

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
		h.UpdateAccount(w, updatedBalance, changesToAccount.Balance, id, account.Currency, typeofoperation2, date1)

		typeofoperation := "topup"
		h.UpdateHistory2(typeofoperation, account.Name, account.Currency, changesToAccount.Balance, date1)

	} else {
		AccountIsBlocked(w, account.Name, account.Id)
	}
}

func (h handler) UpdateAccount(w http.ResponseWriter, updatedBalance, changesToAccountBalance float64, id, AccountCurrency, typeofoperation2 string, date1 time.Time) {
	queryStmt2 := `UPDATE accounts SET balance = $2, currency = $3, date = $4 WHERE id = $1 RETURNING id;`
	err := h.DB.QueryRow(queryStmt2, &id, &updatedBalance, &AccountCurrency, date1).Scan(&id)
	if err != nil {
		log.Println("failed to execute query:", err)
		w.WriteHeader(500)
		return
	} else {
		fmt.Println("Balance is ", typeofoperation2, " on ", changesToAccountBalance, "Result: ", updatedBalance)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Balance is " + string(typeofoperation2) + " on " + strconv.FormatFloat(changesToAccountBalance, 'f', 2, 64))
}

func (h handler) UpdateHistory2(typeofoperation,
	accountName,
	accountCurrency string,
	changesToAccountBalance float64,
	date time.Time) {
	queryStmt3 := `INSERT INTO history (username, date, quantity, currency, typeofoperation) VALUES ($1, $2, $3, $4, $5);`
	_, err := h.DB.Exec(queryStmt3, accountName, date, changesToAccountBalance, accountCurrency, typeofoperation) //USE Exec FOR INSERT
	if err != nil {
		log.Println("failed to execute query - update history:", err)
		return
	} else {
		fmt.Println("History is updated")
	}
}

func AccountIsBlocked(w http.ResponseWriter, accountName string, accountId int) {
	fmt.Println("Operation is not permitted. Account is blocked. Name -", accountName, "ID -", accountId)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Operation is not permitted. Account is blocked")
}

// func (h handler) GetInfoAboutAccount(w http.ResponseWriter, accountId, accountName, accountAccount, accountCurrency, accountDate string, accountBalance int, accountBlocked, accountDefaultaccount bool) Account {
// 	queryStmt := `SELECT * FROM accounts WHERE id = $1 ;`
// 	results, err := h.DB.Query(queryStmt, accountId)
// 	if err != nil {
// 		log.Println("failed to execute query", err)
// 		w.WriteHeader(500)
// 		return err
// 	}

// 	for results.Next() {
// 		err = results.Scan(&accountId, &accountName, &accountAccount, &accountBalance, &accountCurrency, &accountDate, &accountBlocked, &accountDefaultaccount)
// 		if err != nil {
// 			log.Println("failed to scan", err)
// 			w.WriteHeader(500)
// 			return err
// 		}
// 	}
// 	return Account{Name: accountName, Account: accountAccount, Balance: accountBalance, Currency: accountCurrency, Date: accountDate, Blocked: accountBlocked}
// }
