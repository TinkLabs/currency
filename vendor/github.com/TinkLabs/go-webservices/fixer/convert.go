package fixer

type ConvertResp struct {
	Success bool `json:"success"`
	Query   Query `json:"query"`
	Info   Info `json:"info"`
	Date   string  `json:"date"`
	Result float64 `json:"result"`
}

type Query   struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount int    `json:"amount"`
}

type Info struct {
	Timestamp int     `json:"timestamp"`
	Rate      float64 `json:"rate"`
}