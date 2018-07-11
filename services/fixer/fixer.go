package fixer

import (
	"github.com/TinkLabs/go-webservices/fixer"
	"github.com/TinkLabs/go-webservices/fixer/latest"
	"github.com/TinkLabs/go-webservices/fixer/symbols"
	"github.com/TinkLabs/go-webservices/fixer/time-series"

	"github.com/sirupsen/logrus"
)

func ListCurrencies() (*fixer.SymbolsResp, error) {
	log := logrus.WithFields(logrus.Fields{"module": "service", "method": "ListCurrencies"})

	resp, err := symbols.List()
	if err != nil {
		log.WithField("err", err).Error("Failed to list currencies")
		return nil, err
	}

	log.Debug("Successfully listed fixer currencies")

	return resp, nil
}

func ListRates(base string, toCurrencies []string) (*fixer.LatestResp, error) {
	log := logrus.WithFields(logrus.Fields{"module": "service", "method": "ListCurrencies"})

	resp, err := latest.Get(base, toCurrencies)
	if err != nil {
		log.WithField("err", err).Error("Failed to list currency rates")
		return nil, err
	}

	log.Debug("Successfully listed currency rates")

	return resp, nil
}

func ListTimeSeriesRates(start_date, end_date, base string, toCurrencies []string) (*fixer.TimeSeriesResp, error) {
	log := logrus.WithFields(logrus.Fields{"module": "service", "method": "ListTimeSeriesRates"})

	resp, err := time_series.Get(start_date, end_date, base, toCurrencies)
	if err != nil {
		log.WithField("err", err).Error("Failed to list time series currency rates")
		return nil, err
	}

	log.Debug("Successfully listed time series currency rates")

	return resp, nil

}
