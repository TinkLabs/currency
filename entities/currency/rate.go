package currency

import (
	"errors"
	"fmt"
	"time"

	mdcurrency "currency/models/currency"

	"github.com/globalsign/mgo/bson"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

var (
	ErrInvalidParams = errors.New("currency: invalid params")
)

type Rate struct {
	Id        string             `json:"id"`
	Base      string             `json:"base"`
	Date      string             `json:"date"`
	Rates     map[string]float64 `json:"rates"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

// Convert calculates how much to currency is needed given amount of from currency
func (r *Rate) Convert(from, to string, amount float64) (float64, error) {
	if from != r.Base {
		return 0, ErrInvalidParams
	}

	toRate, present := r.Rates[to]
	if !present {
		return 0, ErrInvalidParams
	}

	return toRate * amount, nil
}

// Convert2 uses EUR as rate before we pay for API provider
func (r *Rate) Convert2(from, to string, amount float64) (float64, error) {
	if r.Base != "EUR" {
		return 0, ErrInvalidParams
	}

	fromRate, present := r.Rates[from]
	if !present {
		return 0, ErrInvalidParams
	}

	toRate, present := r.Rates[to]
	if !present {
		return 0, ErrInvalidParams
	}

	toRate2 := decimal.NewFromFloat(toRate)
	fromRate2 := decimal.NewFromFloat(fromRate)
	amount2 := decimal.NewFromFloat(amount)

	result, _ := toRate2.Div(fromRate2).Mul(amount2).Round(3).Float64()
	return result, nil
}

func (r *Rate) ToDoc() *mdcurrency.Rate {
	rate := &mdcurrency.Rate{
		Id:        bson.ObjectIdHex(r.Id),
		Base:      r.Base,
		Date:      r.Date,
		Rates:     r.Rates,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
	return rate
}

func NewRateFromDoc(doc *mdcurrency.Rate) *Rate {
	idBytes, _ := doc.Id.MarshalText()
	rate := &Rate{
		Id:        string(idBytes),
		Base:      doc.Base,
		Date:      doc.Date,
		Rates:     doc.Rates,
		CreatedAt: doc.CreatedAt,
		UpdatedAt: doc.UpdatedAt,
	}

	return rate
}

func NewRatesFromDocs(docs []mdcurrency.Rate) []Rate {
	rates := make([]Rate, 0, len(docs))
	for _, doc := range docs {
		rates = append(rates, *NewRateFromDoc(&doc))
	}

	return rates
}

func CreateRate(base, date string, rates map[string]float64) (*Rate, error) {
	log := logrus.WithFields(logrus.Fields{"module": "entities/rate", "method": "CreateRate", "base": base, "date": date, "rates": rates})

	rate := mdcurrency.NewRate(base, date, rates)
	err := mdcurrency.InsertRate(rate)
	if err != nil {
		log.WithField("err", err).Error("Failed to insert into db")
		return nil, err
	}

	enRate := NewRateFromDoc(rate)

	log.WithField("id", enRate.Id).Debug("Successfully created rate")

	return enRate, nil
}

func FindRateById(id string) (*Rate, error) {
	log := logrus.WithFields(logrus.Fields{"module": "entities/rate", "method": "FindRateById", "id": id})

	rate, err := mdcurrency.FindRateById(id)
	if err == mdcurrency.ErrNotFound {
		log.WithField("err", err).Warn("Failed to find by id")
		return nil, ErrNotFound
	} else if err != nil {
		log.WithField("err", err).Error("Failed to find by id")
		return nil, err
	}

	enRate := NewRateFromDoc(rate)

	log.Debug("Successfully found rate by id")

	return enRate, nil
}

func UpdateRate(enRate *Rate) error {
	log := logrus.WithFields(logrus.Fields{"module": "entities/rate", "method": "UpdateRate", "id": enRate.Id})

	now := time.Now().UTC()
	rate := enRate.ToDoc()
	rate.UpdatedAt = now

	err := mdcurrency.UpdateRate(rate)
	if err != nil {
		log.WithField("err", err).Error("Failed to update rate")
		return err
	}

	log.Debug("Successfully updated rate")

	return nil
}

func DeleteRate(enRate *Rate) error {
	log := logrus.WithFields(logrus.Fields{"module": "entities/rate", "method": "DeleteRate", "id": enRate.Id})

	rate := enRate.ToDoc()
	err := mdcurrency.DeleteRate(rate)
	if err != nil {
		log.WithField("err", err).Error("Failed to remove rate")
		return err
	}

	log.Debug("Successfully removed rate")

	return nil
}

func FindAllRates() ([]Rate, error) {
	log := logrus.WithFields(logrus.Fields{"module": "entities/rate", "method": "FindAllRates"})

	rates, err := mdcurrency.FindAllRates()
	if err != nil {
		log.WithField("err", err).Error("Failed to list rates")
		return nil, err
	}

	enRates := NewRatesFromDocs(rates)

	log.WithField("count", len(enRates)).Debug("Successfully listed rates")

	return enRates, nil
}

func SearchRates(query map[string]interface{}, skip, limit int, orderBy string) ([]Rate, error) {
	log := logrus.WithFields(logrus.Fields{"module": "entities/mail", "method": "SearchRates", "query": fmt.Sprintf("%+v", query), "skip": skip, "limit": limit, "order_by": orderBy})

	rates, err := mdcurrency.SearchRates(query, skip, limit, orderBy)
	if err != nil {
		log.WithField("err", err).Error("Failed to search rates")
		return nil, err
	}

	enRates := NewRatesFromDocs(rates)

	log.WithField("count", len(enRates)).Debug("Successfully searched rates")

	return enRates, nil
}

func CountRates(query map[string]interface{}) (int, error) {
	log := logrus.WithFields(logrus.Fields{"module": "entities/mail", "method": "CountRates", "query": fmt.Sprintf("%+v", query)})

	total, err := mdcurrency.CountRates(query)
	if err != nil {
		log.WithField("err", err).Error("Failed to search rates")
		return 0, err
	}

	log.WithField("total", total).Debug("Successfully counted rates")

	return total, nil
}
