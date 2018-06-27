package currency

import (
	"fmt"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/sirupsen/logrus"
)

var (
	// TODO: single collection or separate collection?
	rateCollection = "currency_rates"
)

type Rate struct {
	Id        bson.ObjectId `bson:"_id" json:"id"`
	Base      string        `bson:"base" json:"base"`
	Date      string        `bson:"date" json:"date"`
	Rates     Rates         `bson:"rates" json:"rates"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updated_at"`
}

type Rates map[string]float64

func NewRate(base, date string, rates map[string]float64) *Rate {
	now := time.Now().UTC()

	rate := &Rate{
		Id:        bson.NewObjectId(),
		Base:      base,
		Date:      date,
		Rates:     rates,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return rate
}

func FindAllRates() ([]Rate, error) {
	log := logrus.WithFields(logrus.Fields{"module": "model", "method": "FindAllRates"})

	var rates []Rate
	err := db.C(rateCollection).Find(bson.M{}).All(&rates)
	if err != nil {
		log.WithField("err", err).Error("Failed to find all docs")
		return nil, err
	}

	log.WithField("count", len(rates)).Debug("Successfully found all docs")

	return rates, nil
}

func FindRateById(id string) (*Rate, error) {
	log := logrus.WithFields(logrus.Fields{"module": "model", "method": "FindRateById", "id": id})

	var rate *Rate
	err := db.C(rateCollection).FindId(bson.ObjectIdHex(id)).One(&rate)
	if err == mgo.ErrNotFound {
		log.WithField("err", err).Warn("Failed to find doc by id")
		return nil, ErrNotFound
	} else if err != nil {
		log.WithField("err", err).Error("Failed to find doc by id")
		return nil, err
	}

	log.Debug("Successfully got doc by id")

	return rate, nil
}

func InsertRate(rate *Rate) error {
	log := logrus.WithFields(logrus.Fields{"module": "model", "method": "InsertRate", "id": rate.Id.String(), "base": rate.Base})

	err := db.C(rateCollection).Insert(rate)
	if err != nil {
		log.WithField("err", err).Error("Failed to insert doc")
		return err
	}

	log.Debug("Successfully inserted doc")

	return nil
}

func DeleteRate(rate *Rate) error {
	log := logrus.WithFields(logrus.Fields{"module": "model", "method": "DeleteRate", "id": rate.Id.String()})

	err := db.C(rateCollection).Remove(rate)
	if err != nil {
		log.WithField("err", err).Error("Failed to delete doc")
		return err
	}

	log.Debug("Successfully deleted doc")

	return nil
}

func UpdateRate(rate *Rate) error {
	log := logrus.WithFields(logrus.Fields{"module": "model", "method": "UpdateRate", "id": rate.Id.String()})

	err := db.C(rateCollection).UpdateId(rate.Id, rate)
	if err != nil {
		log.WithField("err", err).Error("Failed to update doc")
		return err
	}

	log.Debug("Successfully updated doc")

	return nil
}

func SearchRates(query map[string]interface{}, skip, limit int, orderBy string) ([]Rate, error) {
	log := logrus.WithFields(logrus.Fields{"module": "model", "method": "SearchRates", "query": fmt.Sprintf("%+v", query), "skip": skip, "limit": limit, "order_by": orderBy})

	var rates []Rate

	q := db.C(rateCollection).Find(query)
	q.Skip(skip)
	q.Limit(limit)

	// default order by created_at in DESC order
	if orderBy == "" {
		orderBy = "-created_at"
	}
	q.Sort(orderBy)

	err := q.All(&rates)
	if err == mgo.ErrNotFound {
		log.WithField("err", err).Warn("Failed to find docs by query")
		return nil, ErrNotFound
	} else if err != nil {
		log.WithField("err", err).Error("Failed to find docs by query")
		return nil, err
	}

	log.Debug("Successfully got docs by query")

	return rates, nil
}

func CountRates(query map[string]interface{}) (int, error) {
	log := logrus.WithFields(logrus.Fields{"module": "model", "method": "CountRates", "query": fmt.Sprintf("%+v", query)})

	q := db.C(rateCollection).Find(query)
	count, err := q.Count()
	if err != nil {
		log.WithField("err", err).Error("Failed to count docs by query")
		return 0, err
	}

	log.Debug("Successfully counted docs by query")

	return count, nil
}
