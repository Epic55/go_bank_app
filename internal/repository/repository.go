package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/epic55/BankApp/internal/mocks"
	_ "github.com/lib/pq"
)

type Repository struct {
	Db *sql.DB
}

func NewRepository(ConnectionString string) *Repository {
	host := os.Getenv("db_host")
	port := os.Getenv("db_port")
	user := os.Getenv("db_user")
	password := os.Getenv("db_password")
	dbname := os.Getenv("db_name")
	ConnectionString1 := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", ConnectionString1)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	db.SetMaxOpenConns(39)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(3 * time.Minute)

	if err := db.Ping(); err != nil {
		db.Close()
		return nil
	}

	CreateTable(db)

	return &Repository{
		Db: db,
	}
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
		_, err := db.Query("CREATE TABLE accounts (id serial PRIMARY KEY, name VARCHAR(20) NOT NULL, account VARCHAR(20) NOT NULL, balance NUMERIC(10, 2) NOT NULL, currency VARCHAR(3) NOT NULL, date VARCHAR(30) NOT NULL, blocked BOOLEAN NOT NULL, defaultaccount BOOLEAN NOT NULL);")
		if err != nil {
			log.Println("failed to execute query", err)
			return
		} else {
			log.Println("Table accounts created successfully")
		}

		for _, account := range mocks.Accounts {
			queryStmt := `INSERT INTO accounts (name,account,balance,currency,date,blocked,defaultaccount) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`
			date := time.Now()
			date1 := date.Format("2006-01-02 15:04:05")
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
		_, err := db.Query("CREATE TABLE payments (id serial PRIMARY KEY, username VARCHAR(20) NOT NULL, date VARCHAR(20) NOT NULL, service VARCHAR(20) NOT NULL, price NUMERIC(10, 2) NOT NULL, currency VARCHAR(3) NOT NULL);")
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
		_, err := db.Query("CREATE TABLE history (id serial PRIMARY KEY, username VARCHAR(20) NOT NULL, date VARCHAR(20) NOT NULL, quantity NUMERIC(10, 2) NOT NULL, currency VARCHAR(3) NOT NULL, typeofoperation VARCHAR(50) NOT NULL);")
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

func (r *Repository) UpdateAccount(w http.ResponseWriter, updatedBalance, changesToAccountBalance float64, id, AccountCurrency, typeofoperation2 string, date1 string) {
	queryStmt := `UPDATE accounts SET balance = $2, currency = $3, date = $4 WHERE id = $1 RETURNING id;`
	err := r.Db.QueryRow(queryStmt, &id, &updatedBalance, &AccountCurrency, date1).Scan(&id)
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

func (r *Repository) UpdateHistory(typeofoperation,
	accountName,
	accountCurrency string,
	changesToAccountBalance float64,
	date string) {
	queryStmt := `INSERT INTO history (username, date, quantity, currency, typeofoperation) VALUES ($1, $2, $3, $4, $5);`
	_, err := r.Db.Exec(queryStmt, accountName, date, changesToAccountBalance, accountCurrency, typeofoperation) //USE Exec FOR INSERT
	if err != nil {
		log.Println("failed to execute query - update history:", err)
		return
	} else {
		fmt.Println("History is updated")
	}
}

func (r *Repository) UpdateAccountTransferLocal(w http.ResponseWriter,
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
	changesToAccountReceiverBalance,
	updatedBalanceSender,
	updatedBalanceReceiver float64,
	date string) {

	updatedBalanceSender = accountSenderBalance - changesToAccountSenderBalance

	queryStmt := `UPDATE accounts SET balance = $2, date = $3  WHERE account = $1 RETURNING id;`
	err := r.Db.QueryRow(queryStmt, &id, &updatedBalanceSender, date).Scan(&id)
	fmt.Printf("Sender account is withdrawed on %.2f Result: %.2f\n", changesToAccountSenderBalance, updatedBalanceSender)
	if err != nil {
		log.Println("failed to execute query - update accounts withdraw", err)
		w.WriteHeader(500)
		return
	}

	//updatedBalanceReceiver := accountReceiverBalance + changesToAccountReceiverBalance

	queryStmt1 := `UPDATE accounts SET balance = $2, date = $3 WHERE account = $1 RETURNING id;`
	err = r.Db.QueryRow(queryStmt1, &id2, &updatedBalanceReceiver, date).Scan(&id2)
	fmt.Printf("Receiver account is topped up on %.2f Result: %.2f\n", changesToAccountReceiverBalance, updatedBalanceReceiver)
	if err != nil {
		log.Println("failed to execute query - update accounts topup", err)
		w.WriteHeader(500)
		return
	}

}

func (r *Repository) UpdateHistoryTransferLocal(typeofoperation,
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

	queryStmt := `INSERT INTO history (username, date, quantity, currency, typeofoperation) VALUES ($1, $2, $3, $4, $5);`
	_, err := r.Db.Exec(queryStmt, accountSenderName, date, changesToAccountSenderBalance, accountSenderCurrency, typeofoperation+accountSenderAccount) //USE Exec FOR INSERT
	if err != nil {
		log.Println("failed to execute query - update history sender:", err)
		return
	} else {
		fmt.Println("History is updated")
	}

	queryStmt1 := `INSERT INTO history (username, date, quantity, currency, typeofoperation) VALUES ($1, $2, $3, $4, $5);`
	_, err = r.Db.Exec(queryStmt1, accountReceiverName, date, changesToAccountReceiverBalance, accountReceiverCurrency, typeofoperation2+accountReceiverAccount) //USE Exec FOR INSERT
	if err != nil {
		log.Println("failed to execute query - update history receiver:", err)
		return
	} else {
		fmt.Println("History is updated")
	}
}

func (r *Repository) UpdateAccountPayment(w http.ResponseWriter,
	updatedBalance, changesToAccountBalance float64,
	id, accountCurrency, date1 string,
) {
	queryStmt := `UPDATE accounts SET balance = $2, currency = $3, date = $4 WHERE id = $1 RETURNING id;`
	err := r.Db.QueryRow(queryStmt, &id, &updatedBalance, &accountCurrency, date1).Scan(&id)
	fmt.Println("Balance is substracted on", changesToAccountBalance, "Result:", updatedBalance)
	if err != nil {
		log.Println("failed to execute query", err)
		w.WriteHeader(500)
		return
	}

}

func (r *Repository) UpdateHistoryPayment(
	accountName, accountCurrency, date1, changesToPaymentsService string,
	changesToAccountBalance float64) {
	typeofoperation := "payment to "
	queryStmt := `INSERT INTO history (username, date, quantity, currency, typeofoperation) VALUES ($1, $2, $3, $4, $5);`
	_, err := r.Db.Exec(queryStmt, accountName, date1, changesToAccountBalance, accountCurrency, typeofoperation+changesToPaymentsService) //USE Exec FOR INSERT
	if err != nil {
		log.Println("failed to execute query - update history:", err)
		return
	} else {
		fmt.Println("History is updated")
	}
}

func (r *Repository) UpdatePayments(
	accountName, date1, changesToPaymentsService string,
	changesToPaymentsPrice float64,
	changesToPaymentsCurrency string,
) {
	queryStmt := `INSERT INTO payments (username, date, service, price, currency) VALUES ($1, $2, $3, $4, $5);`
	_, err := r.Db.Exec(queryStmt, accountName, date1, changesToPaymentsService, changesToPaymentsPrice, changesToPaymentsCurrency) //USE Exec FOR INSERT
	if err != nil {
		log.Println("failed to execute query - update payments:", err)
		return
	} else {
		fmt.Println("Payments is updated")
	}
}

func (r *Repository) UpdateAccountTransfer(w http.ResponseWriter,
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
	changesToAccountReceiverBalance,
	updatedBalanceSender,
	updatedBalanceReceiver float64,
	date string) {
	queryStmt := `UPDATE accounts SET balance = $2, date = $3  WHERE id = $1 RETURNING id;`
	err := r.Db.QueryRow(queryStmt, &id, &updatedBalanceSender, date).Scan(&id)
	fmt.Printf("Sender Balance is withdrawed on %.2f Result: %.2f\n", changesToAccountSenderBalance, updatedBalanceSender)
	if err != nil {
		log.Println("failed to execute query", err)
		w.WriteHeader(500)
		return
	}
	queryStmt = `UPDATE accounts SET balance = $2, date = $3 WHERE id = $1 RETURNING id;`
	err = r.Db.QueryRow(queryStmt, &id2, &updatedBalanceReceiver, date).Scan(&id2)
	fmt.Printf("Receiver Balance is topped on %.2f Result: %.2f\n", changesToAccountReceiverBalance, updatedBalanceReceiver)
	if err != nil {
		log.Println("failed to execute query", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Balances is updated")
}

func (r *Repository) UpdateHistoryTransfer(
	accountSenderName,
	accountSenderCurrency,
	accountSenderAccount,
	accountReceiverName,
	accountReceiverCurrency,
	accountReceiverAccount string,
	changesToAccountSenderBalance,
	changesToAccountReceiverBalance float64,
	date string) {

	typeofoperation := "transfer to "
	queryStmt := `INSERT INTO history (username, date, quantity, currency, typeofoperation) VALUES ($1, $2, $3, $4, $5);`
	_, err := r.Db.Exec(queryStmt, accountSenderName, date, changesToAccountSenderBalance, accountSenderCurrency, typeofoperation+accountReceiverName) //USE Exec FOR INSERT
	if err != nil {
		log.Println("failed to execute query - update history:", err)
		return
	} else {
		fmt.Println("History is updated")
	}

	typeofoperation2 := "topup from user "
	queryStmt = `INSERT INTO history (username, date, quantity, currency, typeofoperation) VALUES ($1, $2, $3, $4, $5);`
	_, err = r.Db.Exec(queryStmt, accountReceiverName, date, changesToAccountReceiverBalance, accountReceiverCurrency, typeofoperation2+accountSenderName) //USE Exec FOR INSERT
	if err != nil {
		log.Println("failed to execute query - update history:", err)
		return
	} else {
		fmt.Println("History is updated")
	}
}
