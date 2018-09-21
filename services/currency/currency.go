package currency

import (
	"errors"
	"github.com/globalsign/mgo/bson"

	encurrency "currency/entities/currency"

	"github.com/sirupsen/logrus"

	// TODO: remove explicit init
	_ "currency/core/webservices/fixer"
)

var ErrNotFound = errors.New("not found")

type UpdateCurrencyReq encurrency.UpdateCurrencyReq

func Create(name, code string) (*encurrency.Currency, error) {
	log := logrus.WithFields(logrus.Fields{"module": "services/currency", "method": "Create", "name": name, "code": code})

	enCurrency, err := encurrency.Create(name, code)
	if err != nil {
		log.WithField("err", err).Error("Failed to create currency")
		return nil, err
	}

	log.Debug("Successfully created currency")

	return enCurrency, nil
}

func FindById(id string) (*encurrency.Currency, error) {
	log := logrus.WithFields(logrus.Fields{"module": "services/currency", "method": "FindById", "id": id})

	enCurrency, err := encurrency.FindById(id)
	if err == encurrency.ErrNotFound {
		log.WithField("err", err).Warn("Failed to find by id")
		return nil, ErrNotFound
	} else if err != nil {
		log.WithField("err", err).Error("Failed to find by id")
		return nil, err
	}

	log.Debug("Successfully found currency by id")

	return enCurrency, nil
}

func FindByCode(code string) (*encurrency.Currency, error) {
	log := logrus.WithFields(logrus.Fields{"module": "services/currency", "method": "FindByCode", "code": code})

	query := map[string]interface{}{}
	query["code"] = code

	enCurrencies, _, err := Search(query, 0, 1, "")
	if err != nil {
		log.WithField("err", err).Error("Failed to search currency by code")
		return nil, err
	}

	if len(enCurrencies) == 0 {
		log.WithField("err", ErrNotFound).Warn("Failed to find currency by code")
		return nil, ErrNotFound
	}
	enCurrency := enCurrencies[0]
	log.WithFields(logrus.Fields{"id": enCurrency.Id}).Debug("Successfully find currency by code")
	return &enCurrency, nil
}

func FindByCodes(codes []string) ([]encurrency.Currency, error) {
	log := logrus.WithFields(logrus.Fields{"module": "services/currency", "method": "FindByCodes", "codes": codes})

	query := map[string]interface{}{}
	//query["code"] = codes
	query = bson.M{"code": bson.M{"$in": codes}}

	enCurrencies, _, err := Search(query, 0, 0, "")
	if err != nil {
		log.WithField("err", err).Error("Failed to search currency by code")
		return nil, err
	}

	if len(enCurrencies) == 0 {
		log.WithField("err", ErrNotFound).Warn("Failed to find currency by code")
		return nil, ErrNotFound
	}
	log.WithFields(logrus.Fields{"currencies": enCurrencies}).Debug("Successfully find currency by code")
	return enCurrencies, nil
}

func FindByIds(ids []string) ([]encurrency.Currency, error) {
	log := logrus.WithFields(logrus.Fields{"module": "services/currency", "method": "FindByIds", "ids": ids})

	enCurrencies, err := encurrency.FindByIds(ids)
	if err == encurrency.ErrNotFound {
		log.WithField("err", err).Warn("Failed to find by ids")
		return nil, ErrNotFound
	} else if err != nil {
		log.WithField("err", err).Error("Failed to find by ids")
		return nil, err
	}

	log.Debug("Successfully found currency by ids")

	return enCurrencies, nil
}

// PartialUpdate partial updates existing currency object by request received through API
func PartialUpdate(enCurrency *encurrency.Currency, r *UpdateCurrencyReq) (*encurrency.Currency, error) {
	log := logrus.WithFields(logrus.Fields{"module": "services/currency", "method": "PartialUpdate", "id": enCurrency.Id, "req": r})

	if r.Name != nil {
		enCurrency.Name = *r.Name
	}
	if r.Code != nil {
		enCurrency.Code = *r.Code
	}

	err := encurrency.Update(enCurrency)
	if err != nil {
		log.WithField("err", err).Error("Failed to update currency")
		return nil, err
	}

	updatedFile, err := FindById(enCurrency.Id)
	if err != nil {
		log.WithField("err", err).Error("Failed to get currency")
		return nil, err
	}

	log.Debug("Successfully partial updated currency")

	return updatedFile, nil
}

func Update(enCurrency *encurrency.Currency) error {
	log := logrus.WithFields(logrus.Fields{"module": "services/currency", "method": "Update", "id": enCurrency.Id})

	err := encurrency.Update(enCurrency)
	if err != nil {
		log.WithField("err", err).Error("Failed to update currency")
		return err
	}

	log.Debug("Successfully updated currency")

	return nil
}

func Delete(enCurrency *encurrency.Currency) error {
	log := logrus.WithFields(logrus.Fields{"module": "services/currency", "method": "Delete", "id": enCurrency.Id})

	err := encurrency.Delete(enCurrency)
	if err != nil {
		log.WithField("err", err).Error("Failed to remove currency")
		return err
	}

	log.Debug("Successfully removed currency")

	return nil
}

func FindAll() ([]encurrency.Currency, error) {
	log := logrus.WithFields(logrus.Fields{"module": "services/currency", "method": "Delete"})

	enCurrencies, err := encurrency.FindAll()
	if err != nil {
		log.WithField("err", err).Error("Failed to list currencies")
		return nil, err
	}

	log.Debug("Successfully list currencies")

	return enCurrencies, nil
}

func Search(query map[string]interface{}, skip, limit int, orderBy string) ([]encurrency.Currency, int, error) {
	log := logrus.WithFields(logrus.Fields{"module": "services/currency", "method": "Search", "skip": skip, "limit": limit, "order_by": orderBy})

	enCurrencies, err := encurrency.Search(query, skip, limit, orderBy)
	if err != nil {
		log.WithField("err", err).Error("Failed to search currencies")
		return nil, 0, err
	}

	total, err := Count(query)
	if err != nil {
		log.WithField("err", err).Error("Failed to count currencies")
		return nil, 0, err
	}

	log.WithFields(logrus.Fields{"count": len(enCurrencies), "total": total}).Debug("Successfully searched currencies")

	return enCurrencies, total, nil
}

func Count(query map[string]interface{}) (int, error) {
	log := logrus.WithFields(logrus.Fields{"module": "services/currency", "method": "Count"})

	total, err := encurrency.Count(query)
	if err != nil {
		log.WithField("err", err).Error("Failed to search currencies")
		return 0, err
	}

	log.Debug("Successfully counted currencies")

	return total, nil
}

func Convert(from, to string, amount float64) (interface{}, error) {
	log := logrus.WithFields(logrus.Fields{"module": "services/currency", "method": "Convert", "from": from, "to": to, "amount": amount})

	enBaseRate, err := FindLatestRatesByBase(from)
	if err != nil {
		log.WithField("err", err).Error("Failed to find base currency")
		return 0, err
	}

	toAmount, err := enBaseRate.Convert(to, amount)
	if err != nil {
		log.WithField("err", err).Error("Failed to convert currencies")
		return 0, err
	}

	resp := make(map[string]interface{})
	resp["from"] = from
	resp["to"] = to
	resp["amount"] = amount
	resp["result"] = toAmount
	//resp["rate"] = toAmount

	log.Debug("Successfully converted currencies")
	return resp, nil
}
