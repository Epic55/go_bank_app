package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/epic55/BankApp/internal/models"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func (h *Handler) GetStatement(w http.ResponseWriter, r *http.Request, ctx context.Context) {
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

	err = godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file 2")
	}

	minio_url := os.Getenv("minio_url")
	accessKeyID := os.Getenv("minio_access_key")
	secretAccessKey := os.Getenv("minio_secret_Key")
	bucketName := os.Getenv("minio_bucket_name")
	objectName := os.Getenv("minio_object_name")
	filePath := "C:\\Users\\alibe\\Desktop\\statement.txt"
	useSSL := false

	minioClient, err := minio.New(minio_url, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	jsonData, err := json.Marshal(history2)
	if err != nil {
		panic(err)
	}

	//Write the data to the file
	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		log.Fatal(err)
	}

	//Open the file
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	// Get the file size
	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalln(err)
	}
	fileSize := fileInfo.Size()

	// Upload the file
	_, err = minioClient.PutObject(context.Background(), bucketName, objectName, file, fileSize, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Successfully uploaded %s to %s\n", objectName, bucketName)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Statement saved to a file")

}
