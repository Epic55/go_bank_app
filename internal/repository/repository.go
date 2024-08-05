package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Repository struct {
	DB *sql.DB
}

func New(db *sql.DB) Repository {
	return Repository{db}
}

func (h *Repository) UpdateAccount(w http.ResponseWriter, updatedBalance, changesToAccountBalance float64, id, AccountCurrency, typeofoperation2 string, date1 string) {
	queryStmt2 := `UPDATE accounts SET balance = $2, currency = $3, date = $4 WHERE id = $1 RETURNING id;`
	err := h.DB.QueryRow(queryStmt2, &id, &updatedBalance, &AccountCurrency, date1).Scan(&id)
	if err != nil {
		log.Println("failed to execute query:", err)
		w.WriteHeader(500)
		return
	} else {
		fmt.Printf("Balance is %s on %.2f Result: %.2f\n", typeofoperation2, changesToAccountBalance, updatedBalance)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Balance is " + string(typeofoperation2) + " on " + strconv.FormatFloat(changesToAccountBalance, 'f', 2, 64))
}

func (h *Repository) UpdateHistory2(typeofoperation,
	accountName,
	accountCurrency string,
	changesToAccountBalance float64,
	date string) {
	queryStmt3 := `INSERT INTO history (username, date, quantity, currency, typeofoperation) VALUES ($1, $2, $3, $4, $5);`
	_, err := h.DB.Exec(queryStmt3, accountName, date, changesToAccountBalance, accountCurrency, typeofoperation) //USE Exec FOR INSERT
	if err != nil {
		log.Println("failed to execute query - update history:", err)
		return
	} else {
		fmt.Println("History is updated")
	}
}

func (h *Repository) UpdateAccounts(w http.ResponseWriter,
	id, id2,
	accountSenderName,
	accountSenderCurrency,
	accountSenderAccount,
	accountReceiverName,
	accountReceiverCurrency,
	accountReceiverAccount string,
	accountReceiverBalance,
	accountSenderBalance,
	changesToAccountSenderBalance,
	changesToAccountReceiverBalance float64,
	date string) {

	updatedBalanceSender := accountSenderBalance - changesToAccountSenderBalance

	queryStmt2 := `UPDATE accounts SET balance = $2, date = $3  WHERE account = $1 RETURNING id;`
	err := h.DB.QueryRow(queryStmt2, &id, &updatedBalanceSender, date).Scan(&id)
	fmt.Printf("Sender account is withdrawed on %.2f Result: %.2f\n", changesToAccountSenderBalance, updatedBalanceSender)
	if err != nil {
		log.Println("failed to execute query - update accounts withdraw", err)
		w.WriteHeader(500)
		return
	}

	updatedBalanceReceiver := accountReceiverBalance + changesToAccountReceiverBalance

	queryStmt4 := `UPDATE accounts SET balance = $2, date = $3 WHERE account = $1 RETURNING id;`
	err = h.DB.QueryRow(queryStmt4, &id2, &updatedBalanceReceiver, date).Scan(&id2)
	fmt.Printf("Receiver account is topped up on %.2f Result: %.2f\n", changesToAccountReceiverBalance, updatedBalanceReceiver)
	if err != nil {
		log.Println("failed to execute query - update accounts topup", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Balances is updated on " + strconv.FormatFloat(changesToAccountReceiverBalance, 'f', 2, 64) + ". Result: " + strconv.FormatFloat(updatedBalanceReceiver, 'f', 2, 64))
}

func (h *Repository) UpdateHistory(typeofoperation,
	typeofoperation2,
	accountSenderName,
	accountSenderCurrency,
	accountSenderAccount,
	accountReceiverName,
	accountReceiverCurrency,
	accountReceiverAccount string,
	changesToAccountSenderBalance,
	changesToAccountReceiverBalance float64,
	date string) {

	queryStmt3 := `INSERT INTO history (username, date, quantity, currency, typeofoperation) VALUES ($1, $2, $3, $4, $5);`
	_, err := h.DB.Exec(queryStmt3, accountSenderName, date, changesToAccountSenderBalance, accountSenderCurrency, typeofoperation+accountSenderAccount) //USE Exec FOR INSERT
	if err != nil {
		log.Println("failed to execute query - update history sender:", err)
		return
	} else {
		fmt.Println("History is updated")
	}

	queryStmt3 = `INSERT INTO history (username, date, quantity, currency, typeofoperation) VALUES ($1, $2, $3, $4, $5);`
	_, err = h.DB.Exec(queryStmt3, accountReceiverName, date, changesToAccountReceiverBalance, accountReceiverCurrency, typeofoperation2+accountReceiverAccount) //USE Exec FOR INSERT
	if err != nil {
		log.Println("failed to execute query - update history receiver:", err)
		return
	} else {
		fmt.Println("History is updated")
	}
}
