package currency

import (
	"errors"
	"fmt"
	"os"
	"time"

	"currency/config"
	"currency/core/mongodb"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/sirupsen/logrus"
)

var ErrNotFound = errors.New("not found")

var (
	collection = "currencies"
	db         *mgo.Database
)

func init() {
	db = mongodb.Db

	appEnv := os.Getenv("APP_ENV")
	if config.IsTestingAppEnv(appEnv) {
		return
	}

	ensureIndexes()
}

type Currency struct {
	Id        bson.ObjectId `bson:"_id" json:"id"`
	Name      string        `bson:"name" json:"name"`
	Code      string        `bson:"code" json:"code"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updated_at"`
}

func New(name, code string) *Currency {
	now := time.Now().UTC()

	currency := &Currency{
		Id:        bson.NewObjectId(),
		Name:      name,
		Code:      code,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return currency
}

func FindAll() ([]Currency, error) {
	log := logrus.WithFields(logrus.Fields{"module": "model", "method": "FindAll"})

	var currencies []Currency
	err := db.C(collection).Find(bson.M{}).All(&currencies)
	if err != nil {
		log.WithField("err", err).Error("Failed to find all docs")
		return nil, err
	}

	log.WithField("count", len(currencies)).Debug("Successfully found all docs")

	return currencies, nil
}

func FindById(id string) (*Currency, error) {
	log := logrus.WithFields(logrus.Fields{"module": "model", "method": "FindById", "id": id})

	var currency *Currency
	err := db.C(collection).FindId(bson.ObjectIdHex(id)).One(&currency)
	if err == mgo.ErrNotFound {
		log.WithField("err", err).Warn("Failed to find doc by id")
		return nil, ErrNotFound
	} else if err != nil {
		log.WithField("err", err).Error("Failed to find doc by id")
		return nil, err
	}

	log.Debug("Successfully got doc by id")

	return currency, nil
}

func FindByIds(ids []string) ([]Currency, error) {
	log := logrus.WithFields(logrus.Fields{"module": "model", "method": "FindByIds", "ids": ids})

	objectIds := make([]bson.ObjectId, 0, len(ids))
	for _, id := range ids {
		objectIds = append(objectIds, bson.ObjectIdHex(id))
	}

	var currencies []Currency
	err := db.C(collection).Find(bson.M{"_id": bson.M{"$in": objectIds}}).All(&currencies)
	if err == mgo.ErrNotFound {
		log.WithField("err", err).Warn("Failed to find docs by ids")
		return nil, ErrNotFound
	} else if err != nil {
		log.WithField("err", err).Error("Failed to find docs by ids")
		return nil, err
	}

	log.Debug("Successfully got docs by ids")

	return currencies, nil
}

func Insert(currency *Currency) error {
	log := logrus.WithFields(logrus.Fields{"module": "model", "method": "Insert", "id": currency.Id.String(), "name": currency.Name})

	err := db.C(collection).Insert(currency)
	if err != nil {
		log.WithField("err", err).Error("Failed to insert doc")
		return err
	}

	log.Debug("Successfully inserted doc")

	return nil
}

func Delete(currency *Currency) error {
	log := logrus.WithFields(logrus.Fields{"module": "model", "method": "Delete", "id": currency.Id.String()})

	err := db.C(collection).Remove(currency)
	if err != nil {
		log.WithField("err", err).Error("Failed to delete doc")
		return err
	}

	log.Debug("Successfully deleted doc")

	return nil
}

func Update(currency *Currency) error {
	log := logrus.WithFields(logrus.Fields{"module": "model", "method": "Update", "id": currency.Id.String()})

	err := db.C(collection).UpdateId(currency.Id, currency)
	if err != nil {
		log.WithField("err", err).Error("Failed to update doc")
		return err
	}

	log.Debug("Successfully updated doc")

	return nil
}

func Search(query map[string]interface{}, skip, limit int, orderBy string) ([]Currency, error) {
	log := logrus.WithFields(logrus.Fields{"module": "model", "method": "Search", "query": fmt.Sprintf("%+v", query), "skip": skip, "limit": limit, "order_by": orderBy})

	var currencies []Currency

	q := db.C(collection).Find(query)
	q.Skip(skip)
	q.Limit(limit)

	// default order by created_at in DESC order
	if orderBy == "" {
		orderBy = "-created_at"
	}
	q.Sort(orderBy)

	err := q.All(&currencies)
	if err == mgo.ErrNotFound {
		log.WithField("err", err).Warn("Failed to find docs by query")
		return nil, ErrNotFound
	} else if err != nil {
		log.WithField("err", err).Error("Failed to find docs by query")
		return nil, err
	}

	log.Debug("Successfully got docs by query")

	return currencies, nil
}

func Count(query map[string]interface{}) (int, error) {
	log := logrus.WithFields(logrus.Fields{"module": "model", "method": "Count", "query": fmt.Sprintf("%+v", query)})

	q := db.C(collection).Find(query)
	count, err := q.Count()
	if err != nil {
		log.WithField("err", err).Error("Failed to count docs by query")
		return 0, err
	}

	log.Debug("Successfully counted docs by query")

	return count, nil
}

func ensureIndexes() {
	log := logrus.WithFields(logrus.Fields{"module": "model", "method": "ensureIndexes"})

	currenciesByCode := mongodb.NewIndex("by_code", []string{"code"}, true, false, false)

	indexes := []mgo.Index{}
	indexes = append(indexes, *currenciesByCode)

	for _, index := range indexes {
		log = log.WithFields(logrus.Fields{"index_name": index.Name, "collection_name": collection})

		_, err := mongodb.CreateCollectionIndex(collection, &index)
		if err != nil {
			log.WithFields(logrus.Fields{"err": err}).Error("Failed to create collection index")
			continue
		}

		log.Debug("Successfully ensured collection index")
	}

	currencyRateByBaseDateCreatedAt := mongodb.NewIndex("by_base_date_created_at", []string{"base", "date", "-created_at"}, false, false, false)

	currencyRatesIndexes := []mgo.Index{}
	currencyRatesIndexes = append(currencyRatesIndexes, *currencyRateByBaseDateCreatedAt)

	for _, index := range currencyRatesIndexes {
		log = log.WithFields(logrus.Fields{"index_name": index.Name, "collection_name": rateCollection})

		_, err := mongodb.CreateCollectionIndex(rateCollection, &index)
		if err != nil {
			log.WithFields(logrus.Fields{"err": err}).Error("Failed to create collection index")
			continue
		}

		log.Debug("Successfully ensured collection index")
	}
}
