package mocks

import "github.com/epic55/AccountRestApi/pkg/models"

var Accounts = []models.Account{
	{Id: 1, Name: "Epic", Balance: 100, Currency: "Tg", Date: "2024-01-17"},
	{Id: 2, Name: "Shifty", Balance: 200, Currency: "Tg", Date: "2024-01-17"},
}
