package datastore

import (
	"fmt"
	"microtecture/infrastructure/config"
	"time"

	"github.com/couchbase/gocb/v2"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/pkg/errors"
)

// Session sql and nosql databases session
type Session struct {
	SQLSession       sqlSession
	CouchbaseSession couchbaseSession
}

// NewSession creates and returns session
func NewSession() (*Session, error) {
	sqlSession, err := newSQLSession()
	if err != nil {
		return nil, err
	}

	couchbaseSession, err := newCouchbaseSession()
	if err != nil {
		return nil, err
	}

	return &Session{*sqlSession, *couchbaseSession}, nil
}

// NewTestSession creates and returns session for test goals
func NewTestSession() (*Session, error) {
	sqlSession, err := newSQLTestSession()
	if err != nil {
		return nil, err
	}

	couchbaseSession, err := newCouchbaseSession()
	if err != nil {
		return nil, err
	}

	return &Session{*sqlSession, *couchbaseSession}, nil
}

type sqlSession struct {
	*gorm.DB
}

type couchbaseSession struct {
	*gocb.Cluster
}

func getSQLSession(driverName, url string) (*gorm.DB, error) {
	session, err := gorm.Open(driverName, url)
	if err != nil {
		return nil, errors.New(
			fmt.Sprintf(
				"%s %s %s %v", "gorm unable connect to", driverName, "database", err,
			),
		)
	}

	return session, nil
}

func getDBConfig() (*config.DataStoreConfig, error) {
	c, err := config.ConfigFactory(config.DATASTORE_CONFIG)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	conf := c.(*config.DataStoreConfig)

	return conf, nil
}

func newSQLSession() (*sqlSession, error) {
	conf, err := getDBConfig()
	if err != nil {
		return nil, err
	}

	session, err := getSQLSession(conf.Databases.Postgres.Driver, conf.Databases.Postgres.URL)
	if err != nil {
		return nil, err
	}

	return &sqlSession{session}, nil
}

func newSQLTestSession() (*sqlSession, error) {
	conf, err := getDBConfig()
	if err != nil {
		return nil, err
	}

	session, err := getSQLSession(conf.Databases.Postgres.Driver, conf.Databases.Postgres.Test)
	if err != nil {
		return nil, err
	}

	return &sqlSession{session}, nil
}

func newCouchbaseSession() (*couchbaseSession, error) {
	c, err := config.ConfigFactory(config.DATASTORE_CONFIG)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	conf := c.(*config.DataStoreConfig)

	cluster, err := gocb.Connect(
		conf.Databases.Couchbase.URL,
		gocb.ClusterOptions{
			Username: conf.Databases.Couchbase.Username,
			Password: conf.Databases.Couchbase.Password,
		},
	)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	err = cluster.WaitUntilReady(
		time.Second, &gocb.WaitUntilReadyOptions{
			DesiredState: gocb.ClusterStateOnline,
			ServiceTypes: []gocb.ServiceType{gocb.ServiceTypeQuery},
		},
	)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	_, err = cluster.Ping(&gocb.PingOptions{ServiceTypes: []gocb.ServiceType{gocb.ServiceTypeQuery}})
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return &couchbaseSession{cluster}, nil
}
