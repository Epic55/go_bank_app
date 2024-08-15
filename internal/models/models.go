package models

type User struct {
	Id   int    `json:"Id"`
	Name string `json:"name"`
}

type Account struct {
	Id             int     `json:"Id,omitempty"`
	Name           string  `json:"name"`
	Account        string  `json:"account"`
	Balance        float64 `json:"balance"`
	Currency       string  `json:"currency"`
	Date           string  `json:"date"`
	Blocked        bool    `json:"blocked"`
	Defaultaccount bool    `json:"defaultaccount"`
}

type History struct {
	Id              int     `json:"Id,omitempty"`
	Username        string  `json:"username,omitempty"`
	Date            string  `json:"date"`
	Quantity        float64 `json:"quantity"`
	Currency        string  `json:"currency"`
	Typeofoperation string  `json:"typeofoperation"`
}

type Payments struct {
	Id       int     `json:"Id,omitempty"`
	Username string  `json:"username,omitempty"`
	Date     string  `json:"date"`
	Service  string  `json:"service"`
	Price    float64 `json:"price"`
	Currency string  `json:"currency"`
}

type Rate struct {
	Items []Item `xml:"item"`
}

type Item struct {
	Code  string  `xml:"title"`
	Value float64 `xml:"description"`
}
