package routers

import (
	hd "currency/handlers"
	mid "currency/middlewares"

	"github.com/kataras/iris"
	"github.com/kataras/iris/cache"

	"time"
)

func currency(app *iris.Application) {
	cacheMd := cache.Handler(10 * time.Minute)

	app.Post("/currencies", hd.CreateCurrencies)
	app.Get("/currencies", cacheMd, mid.Pagination, hd.ListCurrencies)
	app.Get("/currencies/:code", cacheMd, mid.FetchCurrency, hd.GetCurrency)
	app.Post("/currencies/:code/rates/latest", cacheMd, mid.FetchCurrency, hd.CreateCurrencyRate)
	app.Get("/currencies/:code/rates", cacheMd, mid.FetchCurrency, hd.ListCurrencyRates)
	app.Get("/currencies_rates", cacheMd, hd.ListCurrenciesRates)
	app.Post("/currencies/rates/latest", hd.CreateCurrenciesRate)
	app.Get("/currency/convert", hd.ConvertCurrencies)
	app.Post("/currencies/:code/rates", mid.FetchCurrency, hd.CreateTimeSeriesCurrencyRate)
	app.Post("/currencies/rates", hd.CreateTimeSeriesCurrenciesRate)
}
