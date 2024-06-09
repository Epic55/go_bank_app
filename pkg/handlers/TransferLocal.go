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

func (h handler) TransferLocal(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["account1"]
	id2 := vars["account2"]

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body) // Read request body
	if err != nil {
		log.Fatalln(err)
	}

	date1 := time.Now()

	var changesToAccountSender models.Account
	json.Unmarshal(body, &changesToAccountSender)

	queryStmt := `SELECT * FROM accounts WHERE account = $1 ;`
	results, err := h.DB.Query(queryStmt, id)
	if err != nil {
		log.Println("failed to execute query 1", err)
		w.WriteHeader(500)
		return
	}

	var accountSender models.Account
	for results.Next() {
		err = results.Scan(&accountSender.Id, &accountSender.Name, &accountSender.Account, &accountSender.Balance, &accountSender.Currency, &accountSender.Date, &accountSender.Blocked, &accountSender.Defaultaccount)
		if err != nil {
			log.Println("failed to scan", err)
			w.WriteHeader(500)
			return
		}
	}

	//RECEIVER USER
	var changesToAccountReceiver models.Account
	json.Unmarshal(body, &changesToAccountReceiver)

	queryStmt3 := `SELECT * FROM accounts WHERE account = $1 ;`
	results2, err := h.DB.Query(queryStmt3, id2)
	if err != nil {
		log.Println("failed to execute query 2", err)
		w.WriteHeader(500)
		return
	}

	var accountReceiver models.Account
	for results2.Next() {
		err = results2.Scan(&accountReceiver.Id, &accountReceiver.Name, &accountReceiver.Account, &accountReceiver.Balance, &accountReceiver.Currency, &accountReceiver.Date, &accountReceiver.Blocked, &accountReceiver.Defaultaccount)
		if err != nil {
			log.Println("failed to scan", err)
			w.WriteHeader(500)
			return
		}
	}
	if accountSender.Blocked || accountReceiver.Blocked {

		fmt.Println("Operation is not permitted. Account is blocked -")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("Operation is not permitted. Account is blocked")

	} else {
		if accountReceiver.Currency == "usd" && accountSender.Currency == "tg" && accountSender.Balance >= changesToAccountSender.Balance*450 {
			changesToAccountSender.Balance = changesToAccountSender.Balance * 450
			if accountSender.Balance >= changesToAccountSender.Balance { //CHECK BALANCE OF SENDER, CAN HE AFFORD TO SEND MONEY

				h.UpdateAccounts(w, id, id2,
					accountSender.Name,
					accountSender.Currency,
					accountSender.Account,
					accountReceiver.Name,
					accountReceiver.Currency,
					accountReceiver.Account,
					accountReceiver.Balance,
					accountSender.Balance,
					changesToAccountSender.Balance,
					changesToAccountReceiver.Balance,
					date1)

				typeofoperation := "transfer btwn my acccounts from "
				typeofoperation2 := "transfer btwn my acccounts to "

				h.UpdateHistory(typeofoperation,
					typeofoperation2,
					accountSender.Name,
					accountSender.Currency,
					accountSender.Account,
					accountReceiver.Name,
					accountReceiver.Currency,
					accountReceiver.Account,
					changesToAccountSender.Balance,
					changesToAccountReceiver.Balance,
					date1)
			} else {
				NotEnoughMoney(w)
			}

		} else if accountReceiver.Currency == "tg" && accountSender.Currency == "usd" && accountSender.Balance >= changesToAccountSender.Balance/450 && changesToAccountSender.Balance >= 450 {
			changesToAccountSender.Balance = changesToAccountSender.Balance / 450
			if accountSender.Balance >= changesToAccountSender.Balance { //CHECK BALANCE OF SENDER, CAN HE AFFORD TO SEND MONEY

				h.UpdateAccounts(w, id, id2,
					accountSender.Name,
					accountSender.Currency,
					accountSender.Account,
					accountReceiver.Name,
					accountReceiver.Currency,
					accountReceiver.Account,
					accountReceiver.Balance,
					accountSender.Balance,
					changesToAccountSender.Balance,
					changesToAccountReceiver.Balance,
					date1)

				typeofoperation := "transfer btwn my acccounts from "
				typeofoperation2 := "transfer btwn my acccounts to "

				h.UpdateHistory(typeofoperation,
					typeofoperation2,
					accountSender.Name,
					accountSender.Currency,
					accountSender.Account,
					accountReceiver.Name,
					accountReceiver.Currency,
					accountReceiver.Account,
					changesToAccountSender.Balance,
					changesToAccountReceiver.Balance,
					date1)
			} else {
				NotEnoughMoney(w)
			}
		} else if accountReceiver.Currency == "tg" && accountSender.Currency == "tg" && accountSender.Balance >= changesToAccountSender.Balance {

			if accountSender.Balance >= changesToAccountSender.Balance { //CHECK BALANCE OF SENDER, CAN HE AFFORD TO SEND MONEY

				h.UpdateAccounts(w, id, id2,
					accountSender.Name,
					accountSender.Currency,
					accountSender.Account,
					accountReceiver.Name,
					accountReceiver.Currency,
					accountReceiver.Account,
					accountReceiver.Balance,
					accountSender.Balance,
					changesToAccountSender.Balance,
					changesToAccountReceiver.Balance,
					date1)

				typeofoperation := "transfer btwn my acccounts from "
				typeofoperation2 := "transfer btwn my acccounts to "

				h.UpdateHistory(typeofoperation,
					typeofoperation2,
					accountSender.Name,
					accountSender.Currency,
					accountSender.Account,
					accountReceiver.Name,
					accountReceiver.Currency,
					accountReceiver.Account,
					changesToAccountSender.Balance,
					changesToAccountReceiver.Balance,
					date1)

			} else {
				NotEnoughMoney(w)
			}
		} else {
			fmt.Println("Not enough money for convertation")

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode("Not enough money for convertation")
		}
	}
}

func (h handler) UpdateAccounts(w http.ResponseWriter,
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
	changesToAccountReceiverBalance int,
	date time.Time) {

	updatedBalanceSender := accountSenderBalance - changesToAccountSenderBalance

	queryStmt2 := `UPDATE accounts SET balance = $2, date = $3  WHERE account = $1 RETURNING id;`
	err := h.DB.QueryRow(queryStmt2, &id, &updatedBalanceSender, date).Scan(&id)
	fmt.Println("Sender account is withdrawed on", changesToAccountSenderBalance, "Result:", updatedBalanceSender)
	if err != nil {
		log.Println("failed to execute query - update accounts withdraw", err)
		w.WriteHeader(500)
		return
	}

	updatedBalanceReceiver := accountReceiverBalance + changesToAccountReceiverBalance

	queryStmt4 := `UPDATE accounts SET balance = $2, date = $3 WHERE account = $1 RETURNING id;`
	err = h.DB.QueryRow(queryStmt4, &id2, &updatedBalanceReceiver, date).Scan(&id2)
	fmt.Println("Receiver account is topped up on", changesToAccountReceiverBalance, "Result:", updatedBalanceReceiver)
	if err != nil {
		log.Println("failed to execute query - update accounts topup", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Balances is updated on " + strconv.Itoa(changesToAccountReceiverBalance) + ". Result: " + strconv.Itoa(updatedBalanceReceiver))
}

func (h handler) UpdateHistory(typeofoperation,
	typeofoperation2,
	accountSenderName,
	accountSenderCurrency,
	accountSenderAccount,
	accountReceiverName,
	accountReceiverCurrency,
	accountReceiverAccount string,
	changesToAccountSenderBalance,
	changesToAccountReceiverBalance int,
	date time.Time) {

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

func NotEnoughMoney(w http.ResponseWriter) {
	fmt.Println("Not enough money")

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Not enough money")
}
