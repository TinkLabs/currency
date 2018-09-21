package handlers

import (
	"time"

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

	enCurrencies, err := currencysrv.FindAll()
	if err != nil {
		log.WithField("err", err).Error("Failed to list currencies")
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}

	total := len(enCurrencies)

	log = log.WithFields(logrus.Fields{"count": total, "total": total})

	result := pagination.New()
	result.SetCount(total)
	result.SetLimit(total)
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

	enCurrencyRate, err := currencysrv.GetOrCreateCurrencyRate(enCurrency.Code)
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

func CreateTimeSeriesCurrencyRate(ctx iris.Context) {
	log := logrus.WithFields(logrus.Fields{"module": "handler", "method": "CreateTimeSeriesCurrenciesRate", "http_request": ctx.Request()})

	xRequestId := ctx.Values().GetString("_x_request_id")
	enCurrency := ctx.Values().Get("_encurrency").(*encurrency.Currency)
	startDate := ctx.URLParam("start_date")
	endDate := ctx.URLParam("end_date")

	if !isValidDates(startDate, endDate) {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}

	log = log.WithFields(logrus.Fields{"x_request_id": xRequestId, "start_date": startDate, "end_date": endDate})

	go currencysrv.CreateTimeSeriesCurrencyRate(enCurrency.Code, startDate, endDate)

	log.Debug("Successfully accepted create time series currencies rate request")
	ctx.StatusCode(202)
}

func CreateTimeSeriesCurrenciesRate(ctx iris.Context) {
	log := logrus.WithFields(logrus.Fields{"module": "handler", "method": "CreateTimeSeriesCurrenciesRate", "http_request": ctx.Request()})

	xRequestId := ctx.Values().GetString("_x_request_id")

	startDate := ctx.URLParam("start_date")
	endDate := ctx.URLParam("end_date")

	if !isValidDates(startDate, endDate) {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}

	log = log.WithFields(logrus.Fields{"x_request_id": xRequestId, "start_date": startDate, "end_date": endDate})

	go currencysrv.CreateTimeSeriesCurrenciesRate(startDate, endDate)

	log.Debug("Successfully accepted create time series currencies rate request")
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

func ListCurrenciesRates(ctx iris.Context) {
	log := logrus.WithFields(logrus.Fields{"module": "handler", "method": "ListCurrenciesRates", "http_request": ctx.Request()})

	xRequestId := ctx.Values().GetString("_x_request_id")

	limit := ctx.Values().GetIntDefault("_limit", 1)
	skip := ctx.Values().GetIntDefault("_skip", 0)
	orderBy := ctx.Values().GetStringDefault("_order_by", "")
	enCurrencies := ctx.Values().Get("_encurrencies").([]encurrency.Currency)
	enCurrenciesRates := make([]*encurrency.Rate, len(enCurrencies))
	for index, item := range enCurrencies {
		enRate, err := currencysrv.FindLatestRatesByBase(item.Code)
		if err != nil {
			log.WithField("err", err).Error("Failed to list currencies")
			ctx.StatusCode(iris.StatusInternalServerError)
			return
		}
		log = log.WithFields(logrus.Fields{"rate": enRate})

		enCurrenciesRates[index] = enRate
	}

	log = log.WithFields(logrus.Fields{"x_request_id": xRequestId, "enCurrencies": enCurrencies, "limit": limit, "skip": skip, "order_by": orderBy})

	count := len(enCurrenciesRates)

	result := pagination.New()
	result.SetCount(count)
	result.SetLimit(count)
	result.SetTotal(count)
	result.SetSkip(skip)
	result.SetData(enCurrenciesRates)

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

func isValidDates(startDate, endDate string) (bool) {
	endTime, endErr := time.Parse("2006-01-02", endDate)
	startTime, startErr := time.Parse("2006-01-02", startDate)
	if startErr != nil || endErr != nil {
		return false
	}

	isValidDates := endTime.After(startTime)
	if !isValidDates {
		return false
	}

	numHours := endTime.Sub(startTime).Hours()
	numDays := numHours / 24.0

	return numDays <= 365
}
