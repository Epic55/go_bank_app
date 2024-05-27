package models

type User struct {
	Id   int    `json:"Id"`
	Name string `json:"name"`
}

type Account struct {
	Id       int    `json:"Id"`
	Name     string `json:"name"`
	Account  string `json:"account"`
	Balance  int    `json:"balance"`
	Currency string `json:"currency"`
	Date     string `json:"date"`
	Blocked  bool   `json:"blocked"`
}

type History struct {
	Id              int    `json:"Id,omitempty"`
	Username        string `json:"username,omitempty"`
	Date            string `json:"date"`
	Quantity        int    `json:"quantity"`
	Currency        string `json:"currency"`
	Typeofoperation string `json:"typeofoperation"`
}

type Payments struct {
	Id       int    `json:"Id,omitempty"`
	Username string `json:"username,omitempty"`
	Date     string `json:"date"`
	Service  string `json:"service"`
	Quantity int    `json:"quantity"`
	Currency string `json:"currency"`
}
