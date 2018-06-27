package mongodb

import (
	"os"
	"time"

	"currency/config"

	"github.com/globalsign/mgo"
	"github.com/sirupsen/logrus"
)

var (
	sess *mgo.Session
	Db   *mgo.Database
)

type Index mgo.Index

func init() {
	// skip to connect to db if it's testing env
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "testing" {
		return
	}

	user := config.Config.Mongo.User
	password := config.Config.Mongo.Password
	urls := config.Config.Mongo.Urls
	database := config.Config.Mongo.Database
	port := config.Config.Mongo.Port
	replicaSet := config.Config.Mongo.ReplicaSet

	log := logrus.WithFields(logrus.Fields{"user": user, "urls": urls, "database": database, "port": port, "replica_set": replicaSet})

	dialInfo := &mgo.DialInfo{
		Addrs:    urls,
		Database: database,
		Username: user,
		Password: password,

		// required for MongoDB Atlas
		//ReplicaSetName: replicaSet,

		// required for MongoDB Atlas
		//Source:         "admin",
	}

	// below dialInfo will also works, keep here for a reference for future reference
	//dialInfo, err := mgo.ParseURL("mongodb://admin:<PASSWORD>@cluster0-shard-00-00-quudj.mongodb.net:27017,cluster0-shard-00-01-quudj.mongodb.net:27017,cluster0-shard-00-02-quudj.mongodb.net:27017/currency?replicaSet=Cluster0-shard-0&authSource=admin")
	//if err != nil {
	//	log.WithFields(logrus.Fields{"err": err}).Error("Failed to parse urls")
	//	panic(err)
	//}

	// disable ssl connection due to mlab free try version does not support ssl
	// to connect to an SSL-enabled MongoDB,
	// just define the DialServer function in mgo.DialInfo,
	// and make use of the tls.Dial to perform the connection.
	// https://github.com/go-mgo/mgo/issues/84
	//dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
	//	return tls.Dial("tcp", addr.String(), &tls.Config{})
	//}

	dialInfo.Timeout = time.Second * 30

	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		log.WithFields(logrus.Fields{"err": err}).Error("Failed to open db connection")
		panic(err)
	}

	sess = session
	Db = sess.DB(database)

	log.Info("Successfully connected to db")
}

func GetSession() *mgo.Session {
	return sess.Copy()
}

func GetCollectionNames() ([]string, error) {
	return Db.CollectionNames()
}

func GetCollections() ([]interface{}, error) {
	names, err := GetCollectionNames()
	if err != nil {
		return nil, err
	}

	collections := make([]interface{}, 0, len(names))
	for _, name := range names {
		c := GetCollection(name)
		collections = append(collections, c)
	}

	return collections, nil
}

func GetCollection(collectionName string) interface{} {
	c := Db.C(collectionName)
	docCount, _ := c.Count()
	indexes, _ := c.Indexes()

	collectionInfo := map[string]interface{}{
		"name":           collectionName,
		"document_count": docCount,
		"indexes":        indexes,
	}

	return collectionInfo
}

func GetCollectionIndexes(collectionName string) []interface{} {
	c := Db.C(collectionName)
	cIndexes, _ := c.Indexes()

	indexes := make([]interface{}, 0, len(cIndexes))
	for _, index := range cIndexes {
		indexes = append(indexes, index)
	}

	return indexes
}

func CreateCollectionIndex(collectionName string, index *mgo.Index) (*mgo.Index, error) {
	c := Db.C(collectionName)

	err := c.EnsureIndex(*index)
	if err != nil {
		return nil, err
	}

	return index, nil
}

func NewIndex(name string, keys []string, isUnique, isDropDups, isProcessInBackground bool) *mgo.Index {
	index := &mgo.Index{
		Name:       name,
		Key:        keys,
		Unique:     isUnique,
		DropDups:   isDropDups,
		Background: isProcessInBackground,
	}
	return index
}
