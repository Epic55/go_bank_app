package models

type Account struct {
	Id   int    `json:"Id"`
	Name string `json:"name"`
	//Account  int    `json:"account"`
	Balance  int    `json:"balance"`
	Currency string `json:"currency"`
	Date     string `json:"date"`
	Blocked  bool   `json:"blocked"`
}

type History struct {
	Id              int    `json:"Id"`
	Username        string `json:"username"`
	Typeofoperation string `json:"typeofoperation"`
	Quantity        int    `json:"quantity"`
	Currency        string `json:"currency"`
	Date            string `json:"date"`
}
