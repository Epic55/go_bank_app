package initconfig

import (
	"encoding/json"
	"os"

	"github.com/epic55/BankApp/internal/models"
)

func InitConfig(filename string) (*models.Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var Config models.Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&Config); err != nil {
		return nil, err
	}
	return &Config, nil
}
