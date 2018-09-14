package handlers

import (
	_ "github.com/TinkLabs/go-webservices/fixer"
	_ "github.com/TinkLabs/go-webservices/fixer/convert"
	_ "github.com/TinkLabs/go-webservices/fixer/latest"
	_ "github.com/TinkLabs/go-webservices/fixer/symbols"
	_ "github.com/TinkLabs/go-webservices/fixer/time-series"

	_ "gopkg.in/go-playground/validator.v9"
	_ "github.com/shopspring/decimal"
	_ "github.com/Rhymond/go-money"
)
