package mocks

import "github.com/epic55/BankApp/internal/models"

var Users = []models.User{
	{Id: 1, Name: "Epic"},
	{Id: 2, Name: "Shifty"},
}

var Accounts = []models.Account{
	{Id: 1, Name: "Epic", Account: "q1", Balance: 1000, Currency: "tg", Date: "2024-01-17", Blocked: false, Defaultaccount: true},
	{Id: 3, Name: "Epic", Account: "e3", Balance: 200, Currency: "usd", Date: "2024-01-17", Blocked: false, Defaultaccount: false},
	{Id: 5, Name: "Epic", Account: "t5", Balance: 200, Currency: "tg", Date: "2024-01-17", Blocked: false, Defaultaccount: false},
	{Id: 2, Name: "Shifty", Account: "w2", Balance: 100, Currency: "tg", Date: "2024-01-17", Blocked: false, Defaultaccount: true},
	{Id: 4, Name: "Shifty", Account: "r4", Balance: 200, Currency: "usd", Date: "2024-01-17", Blocked: false, Defaultaccount: false},
}
