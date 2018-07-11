package currency

import (
	"time"

	encurrency "currency/entities/currency"
	fixersrv "currency/services/fixer"

	"github.com/sirupsen/logrus"
)

func CreateCurrencies() error {
	log := logrus.WithFields(logrus.Fields{"module": "services/currency/fixer", "method": "CreateCurrencies"})

	resp, err := fixersrv.ListCurrencies()
	if err != nil {
		log.WithField("err", err).Error("Failed to list currencies")
		return err
	}

	for name, code := range resp.Symbols {
		log = log.WithFields(logrus.Fields{"code": code, "name": name})

		enCurrency, err := Create(string(code), string(name))
		if err != nil {
			log.WithField("err", err).Error("Failed to create currency")
			continue
		}
		log = log.WithFields(logrus.Fields{"id": enCurrency.Id})
		log.Debug("Successfully created currency")
	}

	log.Debug("Successfully created currencies")
	return nil
}

func GetOrCreateCurrencyRate(code string) (*encurrency.Rate, error) {
	log := logrus.WithFields(logrus.Fields{"module": "services/currency/fixer", "method": "GetOrCreateCurrencyRate", "currency_code": code})

	base := code

	date := time.Now().UTC().Format("2006-01-02") // Get current date in the format YYYY-MM-DD

	enRates, _, err := FindRatesByBaseDate(base, date, 0, 1, "")
	if err != nil && err != ErrNotFound {
		log.WithField("err", err).Error("Failed to find currency rate by base and date")
		return nil, err
	}

	if enRates != nil {
		log.Debug("Successfully got currency rate by base and date")
		enRate := enRates[0]
		return &enRate, nil
	}

	toCurrencies := []string{} // default all

	baseRate, err := fixersrv.GetLatestRates(base, toCurrencies)
	if err != nil {
		log.WithField("err", err).Error("Failed to get currency rate")
		return nil, err
	}

	base = baseRate.Base
	date = baseRate.Date

	rates := make(map[string]float64)
	for code, rate := range baseRate.Rates {
		rates[string(code)] = float64(rate)
	}

	enRate, err := CreateRate(base, date, rates)
	if err != nil {
		log.WithField("err", err).Error("Failed to create currency rate")
		return nil, err
	}

	log.Debug("Successfully created currency rate")
	return enRate, nil
}

func CreateCurrenciesRate() {
	log := logrus.WithFields(logrus.Fields{"module": "services/currency/fixer", "method": "CreateCurrenciesRate"})

	enCurrencies, err := FindAll()
	if err != nil {
		log.WithField("err", err).Error("Failed to get currencies")
		return
	}

	for _, enCurrency := range enCurrencies {
		log = log.WithFields(logrus.Fields{"currency_code": enCurrency.Code})
		enRate, err := GetOrCreateCurrencyRate(enCurrency.Code)
		if err != nil {
			log.WithField("err", err).Error("Failed to create currency rate")
			continue
		} else {
			log.WithField("rate_id", enRate.Id).Debug("Successfully created currency rate")
		}
	}
	
	return
}

func CreateTimeSeriesCurrencyRate(code, startDate, endDate string) {
	log := logrus.WithFields(logrus.Fields{"module": "services/currency/fixer", "method": "CreateTimeSeriesCurrencyRate", "start_date": startDate, "end_date": endDate, "currency_code": code})

	base := code
	toCurrencies := []string{} // default all

	baseRate, err := fixersrv.GetRatesByDates(startDate, endDate, base, toCurrencies)
	if err != nil {
		log.WithField("err", err).Error("Failed to get time series currency rates")
		return
	}

	for date, rateByDate := range baseRate.Rates {
		rates := make(map[string]float64)

		for code, rate := range rateByDate {
			rates[string(code)] = float64(rate)
		}

		enRate, err := GetOrCreateRate(base, date, rates)
		if err != nil {
			log.WithField("err", err).Error("Failed to create currency rate")
			continue
		} else {
			log.WithField("rate_id", enRate.Id).Debug("Successfully created currency rate")
		}
	}

	return
}

func CreateTimeSeriesCurrenciesRate(startDate, endDate string) {
	log := logrus.WithFields(logrus.Fields{"module": "services/currency/fixer", "method": "CreateCurrenciesRate"})

	enCurrencies, err := FindAll()
	if err != nil {
		log.WithField("err", err).Error("Failed to get currencies")
		return
	}

	for _, enCurrency := range enCurrencies {
		log = log.WithFields(logrus.Fields{"currency_code": enCurrency.Code})
		go CreateTimeSeriesCurrencyRate(enCurrency.Code, startDate, endDate)

	}

	return
}