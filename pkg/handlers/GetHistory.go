package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/epic55/BankApp/pkg/models"
	"github.com/gorilla/mux"
)

func (h handler) GetHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["username"]

	queryStmt := `SELECT date, quantity, currency, typeofoperation FROM history WHERE username = $1;`
	results, err := h.DB.Query(queryStmt, id)
	//fmt.Println(results)
	if err != nil {
		log.Println("failed to execute query - get history", err)
		w.WriteHeader(500)
		return
	}

	var history2 = make([]models.History, 0)
	for results.Next() {
		var history models.History
		err = results.Scan(&history.Date, &history.Quantity, &history.Currency, &history.Typeofoperation)
		if err != nil {
			log.Println("failed to scan", err)
			w.WriteHeader(500)
			return
		}
		history2 = append(history2, history)
	}
	//fmt.Println(history2)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(history2)
}
