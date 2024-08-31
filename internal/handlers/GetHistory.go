package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/epic55/BankApp/internal/models"
	"github.com/gorilla/mux"
)

func (h *Handler) GetHistory(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	vars := mux.Vars(r)
	id := vars["username"]

	queryStmt := `SELECT date, quantity, currency, typeofoperation FROM history WHERE username = $1 ORDER BY date DESC;`
	results, err := h.R.Db.Query(queryStmt, id)

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

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(history2)
}
