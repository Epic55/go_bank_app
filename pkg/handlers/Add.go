package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/epic55/AccountRestApi/pkg/models"
)

func (h handler) Add(w http.ResponseWriter, r *http.Request) {
	// Read to request body
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
		w.WriteHeader(500)
		return
	}

	var account models.Account
	json.Unmarshal(body, &account)
	if err != nil {
		fmt.Println("Error - ", err)
	}

	//account.Id = (uuid.New()).String()
	date1 := time.Now()
	queryStmt := `INSERT INTO accounts (name,balance,date) VALUES ($1, $2, $3) RETURNING id;`
	err = h.DB.QueryRow(queryStmt, &account.Name, &account.Balance, date1).Scan(&account.Id)
	fmt.Println("Record is added")
	if err != nil {
		log.Println("failed to execute query", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("Created")

}
