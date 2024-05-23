package models

type Account struct {
	Id       int    `json:"Id"`
	Name     string `json:"name"`
	Balance  int    `json:"balance"`
	Currency string `json:"currency"`
	Date     string `json:"date"`
}
