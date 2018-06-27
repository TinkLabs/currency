package currency

import (
	"errors"
	"fmt"
	"time"

	mdcurrency "currency/models/currency"

	"github.com/globalsign/mgo/bson"
	"github.com/sirupsen/logrus"
)

var ErrNotFound = errors.New("not found")

type UpdateCurrencyReq struct {
	Name *string `json:"name"`
	Code *string `json:"code"`
}

type Currency struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (f *Currency) ToDoc() *mdcurrency.Currency {
	currency := &mdcurrency.Currency{
		Id:        bson.ObjectIdHex(f.Id),
		Name:      f.Name,
		Code:      f.Code,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
	}
	return currency
}

func NewFromDoc(doc *mdcurrency.Currency) *Currency {
	idBytes, _ := doc.Id.MarshalText()
	currency := &Currency{
		Id:        string(idBytes),
		Name:      doc.Name,
		Code:      doc.Code,
		CreatedAt: doc.CreatedAt,
		UpdatedAt: doc.UpdatedAt,
	}

	return currency
}

func NewFromDocs(docs []mdcurrency.Currency) []Currency {
	currencies := make([]Currency, 0, len(docs))
	for _, doc := range docs {
		currencies = append(currencies, *NewFromDoc(&doc))
	}

	return currencies
}

func Create(name, code string) (*Currency, error) {
	log := logrus.WithFields(logrus.Fields{"module": "entities/currency", "method": "Create", "name": name, "code": code})

	currency := mdcurrency.New(name, code)
	err := mdcurrency.Insert(currency)
	if err != nil {
		log.WithField("err", err).Error("Failed to insert into db")
		return nil, err
	}

	enFile := NewFromDoc(currency)

	log.WithField("id", enFile.Id).Debug("Successfully created currency")

	return enFile, nil
}

func FindById(id string) (*Currency, error) {
	log := logrus.WithFields(logrus.Fields{"module": "entities/currency", "method": "FindById", "id": id})

	currency, err := mdcurrency.FindById(id)
	if err == mdcurrency.ErrNotFound {
		log.WithField("err", err).Warn("Failed to find by id")
		return nil, ErrNotFound
	} else if err != nil {
		log.WithField("err", err).Error("Failed to find by id")
		return nil, err
	}

	enFile := NewFromDoc(currency)

	log.Debug("Successfully found currency by id")

	return enFile, nil
}

func FindByIds(ids []string) ([]Currency, error) {
	log := logrus.WithFields(logrus.Fields{"module": "entities/currency", "method": "FindByIds", "ids": ids})

	currencies, err := mdcurrency.FindByIds(ids)
	if err == mdcurrency.ErrNotFound {
		log.WithField("err", err).Warn("Failed to find by ids")
		return nil, ErrNotFound
	} else if err != nil {
		log.WithField("err", err).Error("Failed to find by ids")
		return nil, err
	}

	enFiles := NewFromDocs(currencies)

	log.Debug("Successfully found currencies by ids")

	return enFiles, nil
}

func Update(enFile *Currency) error {
	log := logrus.WithFields(logrus.Fields{"module": "entities/currency", "method": "Update", "id": enFile.Id})

	now := time.Now().UTC()
	currency := enFile.ToDoc()
	currency.UpdatedAt = now

	err := mdcurrency.Update(currency)
	if err != nil {
		log.WithField("err", err).Error("Failed to update currency")
		return err
	}

	log.Debug("Successfully updated currency")

	return nil
}

func Delete(enFile *Currency) error {
	log := logrus.WithFields(logrus.Fields{"module": "entities/currency", "method": "Delete", "id": enFile.Id})

	currency := enFile.ToDoc()
	err := mdcurrency.Delete(currency)
	if err != nil {
		log.WithField("err", err).Error("Failed to remove currency")
		return err
	}

	log.Debug("Successfully removed currency")

	return nil
}

func FindAll() ([]Currency, error) {
	log := logrus.WithFields(logrus.Fields{"module": "entities/currency", "method": "FindAll"})

	currencies, err := mdcurrency.FindAll()
	if err != nil {
		log.WithField("err", err).Error("Failed to list currencies")
		return nil, err
	}

	enFiles := NewFromDocs(currencies)

	log.WithField("count", len(enFiles)).Debug("Successfully listed currencies")

	return enFiles, nil
}

func Search(query map[string]interface{}, skip, limit int, orderBy string) ([]Currency, error) {
	log := logrus.WithFields(logrus.Fields{"module": "entities/currency", "method": "Search", "query": fmt.Sprintf("%+v", query), "skip": skip, "limit": limit, "order_by": orderBy})

	currencies, err := mdcurrency.Search(query, skip, limit, orderBy)
	if err != nil {
		log.WithField("err", err).Error("Failed to search currencies")
		return nil, err
	}

	enFiles := NewFromDocs(currencies)

	log.WithField("count", len(enFiles)).Debug("Successfully searched currencies")

	return enFiles, nil
}

func Count(query map[string]interface{}) (int, error) {
	log := logrus.WithFields(logrus.Fields{"module": "entities/currency", "method": "Count", "query": fmt.Sprintf("%+v", query)})

	total, err := mdcurrency.Count(query)
	if err != nil {
		log.WithField("err", err).Error("Failed to search currencies")
		return 0, err
	}

	log.WithField("total", total).Debug("Successfully counted currencies")

	return total, nil
}
