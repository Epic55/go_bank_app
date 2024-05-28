package mocks

import "github.com/epic55/AccountRestApi/pkg/models"

var Users = []models.User{
	{Id: 1, Name: "Epic"},
	{Id: 2, Name: "Shifty"},
}

var Accounts = []models.Account{
	{Id: 1, Name: "Epic", Account: "qwe1", Balance: 100, Currency: "tg", Date: "2024-01-17", Blocked: false, Defaultaccount: true},
	{Id: 2, Name: "Shifty", Account: "asd2", Balance: 100, Currency: "tg", Date: "2024-01-17", Blocked: false, Defaultaccount: true},
	{Id: 3, Name: "Epic", Account: "rty3", Balance: 200, Currency: "usd", Date: "2024-01-17", Blocked: false, Defaultaccount: false},
	{Id: 4, Name: "Shifty", Account: "fgh4", Balance: 200, Currency: "usd", Date: "2024-01-17", Blocked: false, Defaultaccount: false},
}
