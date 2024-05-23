package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/epic55/AccountRestApi/pkg/mocks"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "1"
	dbname   = "postgres"
)

func Connect() *sql.DB {
	connInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connInfo)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	log.Println("Successfully connected to db!")

	return db
}

func CloseConnection(db *sql.DB) {
	defer db.Close()
}

func CreateTable(db *sql.DB) {
	var exists bool
	if err := db.QueryRow("SELECT EXISTS (SELECT FROM pg_tables WHERE  schemaname = 'public' AND tablename = 'accounts' );").Scan(&exists); err != nil {
		log.Println("failed to execute query", err)
		return
	}
	if !exists {
		results, err := db.Query("CREATE TABLE accounts (id serial PRIMARY KEY, name VARCHAR(100) NOT NULL, balance int NOT NULL, currency VARCHAR(100) NOT NULL, date timestamp NOT NULL, blocked BOOLEAN NOT NULL);")
		if err != nil {
			log.Println("failed to execute query", err)
			return
		}
		log.Println("Table created successfully", results)

		for _, account := range mocks.Accounts {
			queryStmt := `INSERT INTO accounts (name,balance,currency,date,blocked) VALUES ($1, $2, $3, $4, $5) RETURNING id;`

			date1 := time.Now()
			err := db.QueryRow(queryStmt, &account.Name, &account.Balance, &account.Currency, date1, &account.Blocked).Scan(&account.Id)
			if err != nil {
				log.Println("failed to execute query", err)
				return
			}
		}
		log.Println("Mock accounts included in Table", results)
	} else {
		log.Println("Table 'account' already exists ")
	}

}
