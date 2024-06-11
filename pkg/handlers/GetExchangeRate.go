package handlers

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/epic55/AccountRestApi/pkg/models"
	"github.com/labstack/gommon/log"
)

func (h handler) GetExchangeRate(w http.ResponseWriter, r *http.Request) {
	date1 := time.Now()
	date1.Format("02.01.2006")
	fmt.Println(date1)

	response, err := http.Get("https://nationalbank.kz/rss/get_rates.cfm?fdate=" + date1)
	if err != nil {
		log.Error(err.Error())
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Error(err)
	}

	var rate1 models.Rate
	err = xml.Unmarshal([]byte(responseData), &rate1)
	if err != nil {
		log.Error("Error - ", err)
	}
	fmt.Println(rate1)

}
