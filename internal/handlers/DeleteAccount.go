package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (h *Handler) DeleteAccount(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	vars := mux.Vars(r)
	id := vars["id"]

	queryStmt := `DELETE FROM accounts WHERE id = $1;`
	_, err := h.R.Db.Query(queryStmt, &id)
	if err != nil {
		log.Println("failed to execute query", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Deleted")
}
