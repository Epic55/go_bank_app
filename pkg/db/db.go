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
	var exists1 bool
	var exists2 bool
	var exists3 bool
	var exists4 bool
	if err := db.QueryRow("SELECT EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'users' );").Scan(&exists1); err != nil {
		log.Println("failed to execute query", err)
		return
	}

	if err := db.QueryRow("SELECT EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'accounts' );").Scan(&exists2); err != nil {
		log.Println("failed to execute query", err)
		return
	}

	if err := db.QueryRow("SELECT EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'payments' );").Scan(&exists3); err != nil {
		log.Println("failed to execute query", err)
		return
	}

	if err := db.QueryRow("SELECT EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'history' );").Scan(&exists4); err != nil {
		log.Println("failed to execute query", err)
		return
	}

	if !exists1 {
		_, err := db.Query("CREATE TABLE users (id serial PRIMARY KEY, name VARCHAR(20) NOT NULL);")
		if err != nil {
			log.Println("failed to execute query", err)
			return
		} else {
			log.Println("Table Users created successfully")
		}

		for _, user := range mocks.Users {
			queryStmt := `INSERT INTO users (name) VALUES ($1);`

			_, err := db.Exec(queryStmt, &user.Name)
			//err := db.QueryRow(queryStmt, &user.Name).Scan(&user.Id)
			if err != nil {
				log.Println("failed to execute query", err)
				return
			}
		}
		log.Println("Mock users included in Table")
	} else {
		log.Println("Table 'users' already exists ")
	}

	if !exists2 {
		_, err := db.Query("CREATE TABLE accounts (id serial PRIMARY KEY, name VARCHAR(20) NOT NULL, account VARCHAR(20) NOT NULL, balance int NOT NULL, currency VARCHAR(3) NOT NULL, date timestamp NOT NULL, blocked BOOLEAN NOT NULL, defaultaccount BOOLEAN NOT NULL);")
		if err != nil {
			log.Println("failed to execute query", err)
			return
		} else {
			log.Println("Table accounts created successfully")
		}

		for _, account := range mocks.Accounts {
			queryStmt := `INSERT INTO accounts (name,account,balance,currency,date,blocked,defaultaccount) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`
			date1 := time.Now()
			err := db.QueryRow(queryStmt, &account.Name, &account.Account, &account.Balance, &account.Currency, date1, &account.Blocked, &account.Defaultaccount).Scan(&account.Id)
			if err != nil {
				log.Println("failed to execute query", err)
				return
			}
		}
		log.Println("Mock accounts included in Table")
	} else {
		log.Println("Table 'accounts' already exists ")
	}

	if !exists3 {
		_, err := db.Query("CREATE TABLE payments (id serial PRIMARY KEY, username VARCHAR(20) NOT NULL, date timestamp NOT NULL, service VARCHAR(20) NOT NULL, quantity int NOT NULL, currency VARCHAR(3) NOT NULL);")
		if err != nil {
			log.Println("failed to execute query", err)
			return

		} else {
			log.Println("Table payments created successfully")
		}

	} else {
		log.Println("Table 'payments' already exists ")
	}

	if !exists4 {
		_, err := db.Query("CREATE TABLE history (id serial PRIMARY KEY, username VARCHAR(20) NOT NULL, date timestamp NOT NULL, quantity int NOT NULL, currency VARCHAR(3) NOT NULL, typeofoperation VARCHAR(50) NOT NULL);")
		if err != nil {
			log.Println("failed to execute query", err)
			return

		} else {
			log.Println("Table history created successfully")
		}

	} else {
		log.Println("Table 'history' already exists ")
	}

}
