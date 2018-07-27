package currency

import (
	encurrency "currency/entities/currency"

	"github.com/sirupsen/logrus"
)

func CreateRate(base, date string, rates map[string]float64) (*encurrency.Rate, error) {
	log := logrus.WithFields(logrus.Fields{"module": "services/currency/rate", "method": "CreateRate", "base": base, "date": date, "rates": rates})

	enRate, err := encurrency.CreateRate(base, date, rates)
	if err != nil {
		log.WithField("err", err).Error("Failed to create currency rate")
		return nil, err
	}

	log = log.WithFields(logrus.Fields{"id": enRate.Id})
	log.Debug("Successfully created currency rate")

	return enRate, nil
}

func GetOrCreateRate(code, date string, rates map[string]float64) (*encurrency.Rate, error) {
	log := logrus.WithFields(logrus.Fields{"module": "services/currency/fixer", "method": "GetOrCreateCurrencyRate", "currency_code": code})

	base := code

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

	enRate, err := CreateRate(base, date, rates)
	if err != nil {
		log.WithField("err", err).Error("Failed to create currency rate")
		return nil, err
	}

	log.Debug("Successfully created currency rate")
	return enRate, nil
}

func FindRateById(id string) (*encurrency.Rate, error) {
	log := logrus.WithFields(logrus.Fields{"module": "services/currency/rate", "method": "FindRateById", "id": id})

	enRate, err := encurrency.FindRateById(id)
	if err == encurrency.ErrNotFound {
		log.WithField("err", err).Warn("Failed to find by id")
		return nil, ErrNotFound
	} else if err != nil {
		log.WithField("err", err).Error("Failed to find by id")
		return nil, err
	}

	log.Debug("Successfully found rate by id")

	return enRate, nil
}

func FindLatestRatesByBase(base string) (*encurrency.Rate, error) {
	log := logrus.WithFields(logrus.Fields{"module": "services/currency/rate", "method": "FindLatestRatesByBase", "base": base})

	query := map[string]interface{}{
		"base": base,
	}

	enRates, _, err := SearchRates(query, 0, 1, "-created_at")
	if err != nil {
		log.WithField("err", err).Error("Failed to find rates by base")
		return nil, err
	}

	if len(enRates) == 0 {
		log.Warn("Failed to find latest rate by code")
		return nil, ErrNotFound
	}

	enRate := &enRates[0]

	log.WithFields(logrus.Fields{"id": enRate.Id}).Debug("Successfully found latest rate by code")

	return enRate, nil
}

func FindRatesByBase(base string, skip, limit int, orderBy string) ([]encurrency.Rate, int, error) {
	log := logrus.WithFields(logrus.Fields{"module": "services/currency/rate", "method": "FindRateById", "base": base, "skip": skip, "limit": limit, "order_by": orderBy})

	query := map[string]interface{}{
		"base": base,
	}

	enRates, total, err := SearchRates(query, skip, limit, orderBy)
	if err != nil {
		log.WithField("err", err).Error("Failed to find rates by base")
		return nil, 0, err
	}

	log.WithFields(logrus.Fields{"count": len(enRates), "total": total}).Debug("Successfully found rates by base")
	return enRates, total, nil
}

func FindRatesByBaseDate(base, date string, skip, limit int, orderBy string) ([]encurrency.Rate, int, error) {
	log := logrus.WithFields(logrus.Fields{"module": "services/currency/rate", "method": "FindRatesByBaseDate", "base": base, "date": date, "skip": skip, "limit": limit, "order_by": orderBy})

	query := map[string]interface{}{
		"base": base,
		"date": date,
	}

	enRates, total, err := SearchRates(query, skip, limit, orderBy)
	if err != nil {
		log.WithField("err", err).Error("Failed to find rates by base")
		return nil, 0, err
	}

	log.WithFields(logrus.Fields{"count": len(enRates), "total": total}).Debug("Successfully found rates by base and date")
	return enRates, total, nil
}

func UpdateRate(enRate *encurrency.Rate) error {
	log := logrus.WithFields(logrus.Fields{"module": "services/currency/rate", "method": "UpdateRate", "id": enRate.Id})

	err := encurrency.UpdateRate(enRate)
	if err != nil {
		log.WithField("err", err).Error("Failed to update rate")
		return err
	}

	log.Debug("Successfully updated rate")

	return nil
}

func DeleteRate(enRate *encurrency.Rate) error {
	log := logrus.WithFields(logrus.Fields{"module": "services/currency/rate", "method": "DeleteRate", "id": enRate.Id})

	err := encurrency.DeleteRate(enRate)
	if err != nil {
		log.WithField("err", err).Error("Failed to remove rate")
		return err
	}

	log.Debug("Successfully removed rate")

	return nil
}

func SearchRates(query map[string]interface{}, skip, limit int, orderBy string) ([]encurrency.Rate, int, error) {
	log := logrus.WithFields(logrus.Fields{"module": "services/currency/rate", "method": "SearchRates", "skip": skip, "limit": limit, "order_by": orderBy})

	enRates, err := encurrency.SearchRates(query, skip, limit, orderBy)
	if err != nil {
		log.WithField("err", err).Error("Failed to search rates")
		return nil, 0, err
	}

	total, err := CountRates(query)
	if err != nil {
		log.WithField("err", err).Error("Failed to count rates")
		return nil, 0, err
	}

	log.WithFields(logrus.Fields{"count": len(enRates), "total": total}).Debug("Successfully searched rates")

	return enRates, total, nil
}

func CountRates(query map[string]interface{}) (int, error) {
	log := logrus.WithFields(logrus.Fields{"module": "services/currency/rate", "method": "CountRates"})

	total, err := encurrency.CountRates(query)
	if err != nil {
		log.WithField("err", err).Error("Failed to search rates")
		return 0, err
	}

	log.Debug("Successfully counted rates")

	return total, nil
}
