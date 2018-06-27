package fixer


type LatestResp struct {
	Success   bool   `json:"success"`
	Timestamp int    `json:"timestamp"`
	Base      string `json:"base"`
	Date      string `json:"date"`
	Rates     Rates `json:"rates"`
}

type Rates map[CurrencyCode]Rate

type Rate float64