package fixer

import (
	"currency/config"

	"github.com/TinkLabs/go-webservices/fixer"
)

var (
	url    string
	apiKey string
)

func init() {
	url := config.Config.Fixer.Url
	apiKey := config.Config.Fixer.ApiKey

	fixer.Setup(url, apiKey)
	fixer.SetDebug(false)
}
