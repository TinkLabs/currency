package middlewares

import (
	currencysrv "currency/services/currency"
	"github.com/kataras/iris"
	"github.com/sirupsen/logrus"
	"strings"
)

func FetchCurrency(ctx iris.Context) {
	log := logrus.WithFields(logrus.Fields{"module": "middleware", "method": "FetchCurrency"})

	language := ctx.URLParam("code");
	log.Debug(language)

	code := ctx.Params().Get("code")
	xRequestId := ctx.Values().GetString("_x_request_id")

	log = log.WithFields(logrus.Fields{"x_request_id": xRequestId, "code": code})

	enCurrency, err := currencysrv.FindByCode(code)
	if err == currencysrv.ErrNotFound {
		log.Warn("Failed to get currency by code")
		ctx.StatusCode(iris.StatusNotFound)
		return
	} else if err != nil {
		log.WithField("err", err).Error("Failed to get currency by code")
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}

	ctx.Values().Set("_encurrency", enCurrency)

	ctx.Next()
}

func FetchCurrencies(ctx iris.Context) {

	log := logrus.WithFields(logrus.Fields{"module": "middleware", "method": "FetchCurrencies"})

	codeString := ctx.URLParam("code");
	codes := strings.Split(codeString, ",");

	xRequestId := ctx.Values().GetString("_x_request_id")

	log = log.WithFields(logrus.Fields{"x_request_id": xRequestId, "code": codes})

	enCurrencies, err := currencysrv.FindByCodes(codes)
	if err == currencysrv.ErrNotFound {
		log.Warn("Failed to get currency by codes")
		ctx.StatusCode(iris.StatusNotFound)
		return
	} else if err != nil {
		log.WithField("err", err).Error("Failed to get currency by codes")
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}

	ctx.Values().Set("_encurrencies", enCurrencies)

	ctx.Next()
}

const currencyCode = "X-Currency-Code"

// FetchCurrencyInHeader fetches X-Currency-Code in header
func FetchCurrencyInHeader(ctx iris.Context) {
	log := logrus.WithFields(logrus.Fields{"module": "middleware", "method": "FetchCurrencyInHeader"})

	currencyCode := ctx.GetHeader(currencyCode)

	if currencyCode == "" {
		// TODO: think about what to do when there is no currency provided
		// default HKD
		currencyCode = "HKD"
	}

	enCurrency, err := currencysrv.FindByCode(currencyCode)
	if err == currencysrv.ErrNotFound {
		log.Warn("Failed to get currency by code")
		ctx.StatusCode(iris.StatusBadRequest)
		return
	} else if err != nil {
		log.WithField("err", err).Error("Failed to get currency by code")
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}

	ctx.Values().Set("_encurrency", enCurrency)

	ctx.Next()
}
