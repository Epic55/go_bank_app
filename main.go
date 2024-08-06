package main

import "github.com/epic55/BankApp/internal/app"

func main() {
	app := app.NewApplication()
	app.StartServer()
}
