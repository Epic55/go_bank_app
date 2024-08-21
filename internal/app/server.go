package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/epic55/BankApp/internal/handlers"
	"github.com/epic55/BankApp/internal/initconfig"
	"github.com/epic55/BankApp/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/epic55/BankApp/internal/repository"
	"github.com/gorilla/mux"
)

type Application struct {
}

func NewApplication() *Application {
	return &Application{}
}

var (
	Repo *repository.Repository
	Hand *handlers.Handler
	Cnfg *models.Config
)

func init() {
	var err error
	Cnfg, err = initconfig.InitConfig("config.json")
	if err != nil {
		fmt.Println("Failed to initialize the config:", err)
		return
	}
	checkConnectionDB()
	checkConnectionMinio()
	Repo = repository.NewRepository(Cnfg.ConnectionString)
	Hand = handlers.NewHandler(Repo, Cnfg)

}

func checkConnectionDB() {
	connString := "postgres://postgres:1@localhost:5432/postgres"

	// Create a connection pool
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		log.Fatalf("Unable to parse connection string: %v\n", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}
	defer pool.Close()

	// Ping the database to check connection
	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatalf("Unable to connect to DB") //: %v\n", err)
	}

	fmt.Println("Successfully connected to PostgreSQL DB!")
}

func checkConnectionMinio() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file 1")
	}

	minio_url := os.Getenv("minio_url")
	accessKeyID := os.Getenv("minio_access_key")
	secretAccessKey := os.Getenv("minio_secret_Key")
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(minio_url, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln("Couldnt connect to MinIO - ", err)
	}

	//Check if the connection is successful by listing buckets
	ctx := context.Background()
	buckets, err := minioClient.ListBuckets(ctx)
	if err != nil {
		log.Fatalln("Couldnt connect to MinIO - ", err)
	}
	fmt.Print("Successfully connected to MinIO - Bucket: ")
	for _, bucket := range buckets {
		fmt.Println(bucket.Name)
	}
}

func (a *Application) StartServer() {
	r := mux.NewRouter()

	r.HandleFunc("/accounts/topup/{id}", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithDeadline(r.Context(), time.Now().Add(30*time.Second))
		defer cancel()
		Hand.Topup(w, r.WithContext(ctx), ctx)

	})

	r.HandleFunc("/accounts", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithDeadline(r.Context(), time.Now().Add(30*time.Second))
		defer cancel()
		Hand.GetAllAccounts(w, r.WithContext(ctx), ctx)

	})

	r.HandleFunc("/accounts/{id}", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithDeadline(r.Context(), time.Now().Add(30*time.Second))
		defer cancel()
		Hand.GetAccount(w, r.WithContext(ctx), ctx)

	})

	r.HandleFunc("/history/{username}", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithDeadline(r.Context(), time.Now().Add(30*time.Second))
		defer cancel()
		Hand.GetHistory(w, r.WithContext(ctx), ctx)

	})

	r.HandleFunc("/statement/{username}", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithDeadline(r.Context(), time.Now().Add(30*time.Second))
		defer cancel()
		Hand.GetStatement(w, r.WithContext(ctx), ctx)

	})

	r.HandleFunc("/history/payments/{username}", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithDeadline(r.Context(), time.Now().Add(30*time.Second))
		defer cancel()
		Hand.GetHistoryPayments(w, r.WithContext(ctx), ctx)

	})

	r.HandleFunc("/history/transfers/{username}", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithDeadline(r.Context(), time.Now().Add(30*time.Second))
		defer cancel()
		Hand.GetHistoryTransfers(w, r.WithContext(ctx), ctx)

	})

	r.HandleFunc("/accounts/delete/{id}", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithDeadline(r.Context(), time.Now().Add(30*time.Second))
		defer cancel()
		Hand.DeleteAccount(w, r.WithContext(ctx), ctx)

	})

	r.HandleFunc("/accounts/withdraw/{id}", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithDeadline(r.Context(), time.Now().Add(30*time.Second))
		defer cancel()
		Hand.Withdraw(w, r.WithContext(ctx), ctx)

	})

	r.HandleFunc("/accounts/transfer/{id}/{id2}", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithDeadline(r.Context(), time.Now().Add(30*time.Second))
		defer cancel()
		Hand.Transfer(w, r.WithContext(ctx), ctx)

	})

	r.HandleFunc("/accounts/transferlocal/{account1}/{account2}", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithDeadline(r.Context(), time.Now().Add(30*time.Second))
		defer cancel()
		Hand.TransferLocal(w, r.WithContext(ctx), ctx)

	})

	r.HandleFunc("/accounts/blocking/{id}", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithDeadline(r.Context(), time.Now().Add(30*time.Second))
		defer cancel()
		Hand.BlockAccount(w, r.WithContext(ctx), ctx)

	})

	r.HandleFunc("/payments/{id}", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithDeadline(r.Context(), time.Now().Add(30*time.Second))
		defer cancel()
		Hand.Payments(w, r.WithContext(ctx), ctx)

	})

	server := &http.Server{
		Addr:         "localhost:" + Cnfg.ListenPort,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      r,
	}
	quit := make(chan os.Signal, 1)
	go shutdown(quit)
	fmt.Println("Listening on port", Cnfg.ListenPort, "...")
	fmt.Println(server.ListenAndServe())

}

func shutdown(quit chan os.Signal) {
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	s := <-quit
	fmt.Println("caught signal", "signal", s.String())
	os.Exit(0)
}
