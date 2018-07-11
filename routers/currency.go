package routers

import (
	hd "currency/handlers"
	mid "currency/middlewares"

	"github.com/kataras/iris"
)

func currency(app *iris.Application) {
	app.Post("/currencies", hd.CreateCurrencies)
	app.Get("/currencies", mid.Pagination, hd.ListCurrencies)
	app.Get("/currencies/:code", mid.FetchCurrency, hd.GetCurrency)
	app.Post("/currencies/:code/rates", mid.FetchCurrency, hd.CreateCurrencyRate)
	app.Get("/currencies/:code/rates", mid.FetchCurrency, hd.ListCurrencyRates)
	app.Post("/currencies/rates", hd.CreateCurrenciesRate)
	app.Get("/currency/convert", hd.ConvertCurrencies)
	app.Post("/currencies/:code/timeseries", mid.FetchCurrency, hd.CreateTimeSeriesCurrencyRate)
	app.Post("/currencies/timeseries", hd.CreateTimeSeriesCurrenciesRate)
}