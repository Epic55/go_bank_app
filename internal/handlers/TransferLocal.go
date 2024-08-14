package handlers

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/epic55/BankApp/internal/models"
	"github.com/gorilla/mux"
)

var ExchangeRate float64

func (h *Handler) TransferLocal(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	vars := mux.Vars(r)
	id := vars["account1"]
	id2 := vars["account2"]

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body) // Read request body
	if err != nil {
		log.Fatalln(err)
	}

	date := time.Now()
	date1 := date.Format("2006-01-02 15:04:05")

	var changesToAccountSender models.Account
	json.Unmarshal(body, &changesToAccountSender)

	queryStmt := `SELECT * FROM accounts WHERE account = $1 ;`
	results, err := h.R.Db.Query(queryStmt, id)
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
	results2, err := h.R.Db.Query(queryStmt3, id2)
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

	GetExchangeRate(w, r)

	if accountSender.Blocked || accountReceiver.Blocked {

		fmt.Println("Operation is not permitted. Account is blocked -")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("Operation is not permitted. Account is blocked")

	} else {
		if accountReceiver.Currency == "usd" && accountSender.Currency == "tg" && accountSender.Balance >= changesToAccountSender.Balance*ExchangeRate {
			changesToAccountSender.Balance = changesToAccountSender.Balance * ExchangeRate

			if accountSender.Balance >= changesToAccountSender.Balance { //CHECK BALANCE OF SENDER, CAN HE AFFORD TO SEND MONEY

				h.R.UpdateAccountTransferLocal(w, id, id2,
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

				h.R.UpdateHistoryTransferLocal(typeofoperation,
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

		} else if accountReceiver.Currency == "tg" && accountSender.Currency == "usd" && accountSender.Balance >= changesToAccountSender.Balance/ExchangeRate && changesToAccountSender.Balance >= ExchangeRate {
			changesToAccountSender.Balance = changesToAccountSender.Balance / ExchangeRate
			if accountSender.Balance >= changesToAccountSender.Balance { //CHECK BALANCE OF SENDER, CAN HE AFFORD TO SEND MONEY

				h.R.UpdateAccountTransferLocal(w, id, id2,
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

				h.R.UpdateHistoryTransferLocal(typeofoperation,
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

				h.R.UpdateAccountTransferLocal(w, id, id2,
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

				h.R.UpdateHistoryTransferLocal(typeofoperation,
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

func NotEnoughMoney(w http.ResponseWriter) {
	fmt.Println("Not enough money")
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Not enough money")
}

func GetExchangeRate(w http.ResponseWriter, r *http.Request) float64 {
	date1 := time.Now()
	date := date1.Format("02.01.2006")

	response, err := http.Get("https://nationalbank.kz/rss/get_rates.cfm?fdate=" + date)
	if err != nil {
		log.Println("Error1 - ", err.Error())
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("Error2 - ", err)
	}

	var rate1 models.Rate
	err = xml.Unmarshal([]byte(responseData), &rate1)
	if err != nil {
		log.Println("Error3 - ", err)
	}

	for _, item := range rate1.Items {
		if item.Code == "USD" {
			ExchangeRate = item.Value
			fmt.Println("USD Exchange Rate is ", ExchangeRate)
		}
	}
	return ExchangeRate
}
