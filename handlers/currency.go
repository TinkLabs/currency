package handlers

import (
	encurrency "currency/entities/currency"
	currencysrv "currency/services/currency"
	"currency/services/pagination"

	"github.com/kataras/iris"
	"github.com/sirupsen/logrus"
)

func GetCurrency(ctx iris.Context) {
	log := logrus.WithFields(logrus.Fields{"module": "handler", "method": "GetCurrency", "http_request": ctx.Request()})

	id := ctx.Params().Get("id")
	xRequestId := ctx.Values().GetString("_x_request_id")
	enCurrency := ctx.Values().Get("_encurrency").(*encurrency.Currency)

	log = log.WithFields(logrus.Fields{"x_request_id": xRequestId, "id": id})

	log.Debug("Successfully got currency by id")
	ctx.JSON(enCurrency)
}

func ListCurrencies(ctx iris.Context) {
	log := logrus.WithFields(logrus.Fields{"module": "handler", "method": "ListCurrencies", "http_request": ctx.Request()})

	xRequestId := ctx.Values().GetString("_x_request_id")

	limit := ctx.Values().GetIntDefault("_limit", 10)
	skip := ctx.Values().GetIntDefault("_skip", 0)
	orderBy := ctx.Values().GetStringDefault("_order_by", "")

	log = log.WithFields(logrus.Fields{"x_request_id": xRequestId, "limit": limit, "skip": skip, "order_by": orderBy})

	enCurrencies, total, err := currencysrv.Search(nil, skip, limit, orderBy)
	if err != nil {
		log.WithField("err", err).Error("Failed to list currencies")
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}

	count := len(enCurrencies)

	log = log.WithFields(logrus.Fields{"count": count, "total": total})

	result := pagination.New()
	result.SetCount(count)
	result.SetLimit(limit)
	result.SetSkip(skip)
	result.SetTotal(total)
	result.SetData(enCurrencies)

	log.Debug("Successfully listed currencies")

	ctx.JSON(result.BuildScrollable())
}

func CreateCurrencies(ctx iris.Context) {
	log := logrus.WithFields(logrus.Fields{"module": "handler", "method": "CreateCurrencies", "http_request": ctx.Request()})

	xRequestId := ctx.Values().GetString("_x_request_id")

	log = log.WithFields(logrus.Fields{"x_request_id": xRequestId})

	err := currencysrv.CreateCurrencies()
	if err != nil {
		log.WithField("err", err).Error("Failed to list currencies")
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}

	log.Debug("Successfully created currencies")
	ctx.JSON(nil)
}

func CreateCurrencyRate(ctx iris.Context) {
	log := logrus.WithFields(logrus.Fields{"module": "handler", "method": "CreateCurrencyRate", "http_request": ctx.Request()})

	id := ctx.Params().Get("id")
	xRequestId := ctx.Values().GetString("_x_request_id")
	enCurrency := ctx.Values().Get("_encurrency").(*encurrency.Currency)

	log = log.WithFields(logrus.Fields{"x_request_id": xRequestId, "id": id})

	enCurrencyRate, err := currencysrv.CreateCurrencyRate(enCurrency.Code)
	if err != nil {
		log.WithField("err", err).Error("Failed to create currency rate")
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}

	log.Debug("Successfully created currency rate")
	ctx.JSON(enCurrencyRate)
}

func CreateCurrenciesRate(ctx iris.Context) {
	log := logrus.WithFields(logrus.Fields{"module": "handler", "method": "CreateCurrenciesRate", "http_request": ctx.Request()})

	id := ctx.Params().Get("id")
	xRequestId := ctx.Values().GetString("_x_request_id")

	log = log.WithFields(logrus.Fields{"x_request_id": xRequestId, "id": id})

	go currencysrv.CreateCurrenciesRate()

	log.Debug("Successfully accepted create currencies rate request")
	ctx.StatusCode(202)
}

func ListCurrencyRates(ctx iris.Context) {
	log := logrus.WithFields(logrus.Fields{"module": "handler", "method": "ListCurrencyRates", "http_request": ctx.Request()})

	xRequestId := ctx.Values().GetString("_x_request_id")

	limit := ctx.Values().GetIntDefault("_limit", 10)
	skip := ctx.Values().GetIntDefault("_skip", 0)
	orderBy := ctx.Values().GetStringDefault("_order_by", "")
	enCurrency := ctx.Values().Get("_encurrency").(*encurrency.Currency)

	log = log.WithFields(logrus.Fields{"x_request_id": xRequestId, "limit": limit, "skip": skip, "order_by": orderBy})

	enCurrencyRates, total, err := currencysrv.FindRatesByBase(enCurrency.Code, skip, limit, orderBy)
	if err != nil {
		log.WithField("err", err).Error("Failed to list currencies")
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}

	count := len(enCurrencyRates)

	log = log.WithFields(logrus.Fields{"count": count, "total": total})

	result := pagination.New()
	result.SetCount(count)
	result.SetLimit(limit)
	result.SetSkip(skip)
	result.SetTotal(total)
	result.SetData(enCurrencyRates)

	log.Debug("Successfully listed currency rates")

	ctx.JSON(result.BuildScrollable())
}

func ConvertCurrencies(ctx iris.Context) {
	log := logrus.WithFields(logrus.Fields{"module": "handler", "method": "ConvertCurrencies", "http_request": ctx.Request()})

	xRequestId := ctx.Values().GetString("_x_request_id")

	from := ctx.FormValue("from")
	to := ctx.FormValue("to")
	amount, err := ctx.URLParamFloat64("amount")
	if err != nil {
		log.WithField("err", err).Error("Failed to parse amount")
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}

	log = log.WithFields(logrus.Fields{"x_request_id": xRequestId, "from": from, "to": to, "amount": amount})

	enCurrencyConversion, err := currencysrv.Convert(from, to, amount)
	if err != nil {
		log.WithField("err", err).Error("Failed to convert currencies")
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}

	log.Debug("Successfully created currency rate")
	ctx.JSON(enCurrencyConversion)
}
