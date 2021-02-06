package testutil

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"microtecture/infrastructure/config"
	"microtecture/interface/router"
	"microtecture/registry"
	"microtecture/usecase/controllers"
)

type T struct {
	Controller  controllers.Root
	Testing     *testing.T
	httpRequest *http.Request
	httpRoute   *httprouter.Router
}

func removeTestDB(conf config.DataStoreConfig) error {
	dbs, err := sql.Open(conf.Databases.Postgres.Driver, conf.Databases.Postgres.Admin)
	if err != nil {
		return errors.New(err.Error())
	}
	defer dbs.Close()

	err = dbs.Ping()
	if err != nil {
		return errors.New(err.Error())
	}

	var dbname string
	splited := strings.Split(conf.Databases.Postgres.Test, " ")
	for _, s := range splited {
		tmp := strings.Split(s, "=")
		if tmp[0] == "dbname" {
			dbname = tmp[1]
		}
	}

	statement := fmt.Sprintf("DROP DATABASE %v", dbname)
	_, err = dbs.Exec(statement)
	if err != nil && err.Error() != `pq: database "`+config.NAME+`_test" does not exist` {
		return errors.New(err.Error())
	}

	return nil
}

func createTestDB() error {
	conf, err := config.ConfigFactory(config.DATASTORE_CONFIG)
	if err != nil {
		return err
	}
	dbConfig := conf.(*config.DataStoreConfig)

	if err := removeTestDB(*dbConfig); err != nil {
		return err
	}

	dbs, err := sql.Open(dbConfig.Databases.Postgres.Driver, dbConfig.Databases.Postgres.Admin)
	if err != nil {
		return errors.New(err.Error())
	}
	defer dbs.Close()

	err = dbs.Ping()
	if err != nil {
		return errors.New(err.Error())
	}

	var dbname string
	var dbuser string
	splited := strings.Split(dbConfig.Databases.Postgres.Test, " ")
	for _, s := range splited {
		tmp := strings.Split(s, "=")
		if tmp[0] == "dbname" {
			dbname = tmp[1]
		} else if tmp[0] == "user" {
			dbuser = tmp[1]
		}
	}

	statement := fmt.Sprintf("CREATE DATABASE %v", dbname)
	_, err = dbs.Exec(statement)
	if err != nil {
		return errors.New(err.Error())
	}

	statement = fmt.Sprintf(`grant ALL privileges on database "%v" to "%v"`, dbname, dbuser)
	_, err = dbs.Exec(statement)
	if err != nil {
		return errors.New(err.Error())
	}

	return err
}

func New() *T {
	err := createTestDB()
	if err != nil {
		logrus.Fatalf("%+v", err)
	}

	reg, err := registry.NewTest()
	if err != nil {
		logrus.Fatalf("%+v", err)
	}
	c := reg.NewRootController()

	testctrl := &T{
		Controller: c,
	}

	return testctrl
}

func (t *T) Init(te *testing.T) {
	// TODO: write drop, migrate and base data for database
	// t.Controller.GetBase().Application.DropDB()
	// t.Controller.GetBase().Application.MigrateDB()
	// t.Controller.GetBase().Application.InsertBaseData()
	t.Testing = te
}

func (t *T) Close() {
	err := t.Controller.GetBase().Application.Close()
	if err != nil {
		logrus.Fatalf("%+v", err)
	}

	conf, err := config.ConfigFactory(config.DATASTORE_CONFIG)
	if err != nil {
		logrus.Fatalf("%+v", err)
	}
	dbConfig := conf.(*config.DataStoreConfig)

	if err := removeTestDB(*dbConfig); err != nil {
		logrus.Fatalf("%+v", err)
	}
}

// Login make authenticated request
func (t *T) Login(request *http.Request, id uuid.UUID, firstName, lastName string, roles ...string) (*http.Request, error) {
	app := t.Controller.GetBase().Application
	jwt, err := app.CreateJWT(id, firstName, lastName, false, roles...)
	if err != nil {
		return nil, err
	}

	if err == nil {
		cookie := &http.Cookie{
			Name:  config.ACCESS_TOKEN_NAME,
			Value: jwt,
		}

		request.AddCookie(cookie)
	}

	return request, err
}

func (t *T) SendRestRequest(data interface{}) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()

	if data != nil {
		switch ty := data.(type) {
		case string:
			ty, _ = data.(string)
			t.httpRequest.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(ty)))
		default:
			d, _ := json.Marshal(&data)
			t.httpRequest.Body = ioutil.NopCloser(bytes.NewBuffer(d))
		}
	}

	t.httpRoute.ServeHTTP(rr, t.httpRequest)

	return rr
}

func (t *T) SendBadJson(req *http.Request, router *httprouter.Router) {
	t.Testing.Run("when send bad json", func(te *testing.T) {
		req.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`"{badjson}`)))
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(te, 400, rr.Code)
	})
}

func (t *T) SetHTTPRequest(method, url string) {
	var err error
	t.httpRequest, err = http.NewRequest(strings.ToUpper(method), url, nil)
	if err != nil {
		panic(err)
	}

	t.httpRoute = httprouter.New()
	router.Route(t.httpRoute, t.Controller)
}

func Fatal(err error, t *testing.T) {
	if err != nil {
		t.Fatal(err)
	}
}
