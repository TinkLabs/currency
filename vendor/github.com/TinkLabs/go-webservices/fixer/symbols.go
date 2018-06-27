package fixer

type SymbolsResp struct {
	Success bool `json:"success"`
	Symbols Symbols `json:"symbols"`
}

type Symbols map[CurrencyCode]CurrencyName

type CurrencyCode string
type CurrencyName string