package application

import (
	"microtecture/infrastructure/config"
	"microtecture/infrastructure/datastore"

	"github.com/sirupsen/logrus"
)

type application struct {
	Config    config.ApplicationConfig
	DBSession datastore.Session
	Logger    logrus.FieldLogger
}

// New creates and returns Application
func New() (application, error) {
	app := application{Logger: logrus.StandardLogger()}

	conf, err := config.ConfigFactory(config.APPLICATION_CONFIG)
	if err != nil {
		return app, err
	}
	appConfig := conf.(*config.ApplicationConfig)
	app.Config = *appConfig

	dbSession, err := datastore.NewSession()
	if err != nil {
		return app, err
	}

	app.DBSession = *dbSession

	return app, nil
}

// Close close sql database of application and some things else
func (self application) Close() error {
	return self.DBSession.SQLSession.Close()
}
