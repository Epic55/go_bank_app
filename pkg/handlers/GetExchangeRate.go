package handlers

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/epic55/AccountRestApi/pkg/models"
)

func (h handler) GetExchangeRate(w http.ResponseWriter, r *http.Request) {
	date1 := time.Now()
	date := date1.Format("02.01.2006")

	response, err := http.Get("https://nationalbank.kz/rss/get_rates.cfm?fdate=" + date)
	if err != nil {
		log.Println(err.Error())
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}

	var rate1 models.Rate
	err = xml.Unmarshal([]byte(responseData), &rate1)
	if err != nil {
		log.Println("Error - ", err)
	}
	for i, v := range rate1 {
		fmt.Println(i, v)
	}

}
