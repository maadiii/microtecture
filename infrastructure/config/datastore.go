package config

import (
	"fmt"

	"github.com/pkg/errors"
)

type postgres struct {
	Driver string `yaml:"driver"`
	URL    string `yaml:"main_url"`
	Test   string `yaml:"test_url"`
	Admin  string `yaml:"admin_url"`
}

type couchbase struct {
	URL      string `yaml:"url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type databases struct {
	Postgres  postgres  `yaml:"postgres"`
	Couchbase couchbase `yaml:"couchbase"`
}

type DataStoreConfig struct {
	Databases databases `yaml:"databases"`
}

func (self *DataStoreConfig) Init() error {
	if err := marshal(self); err != nil {
		return err
	}

	const errMsg = "is not set in config file."

	if self.Databases.Postgres.Driver == "" {
		return errors.New(fmt.Sprintf("%s %s", "databases.postgres.driver", errMsg))
	}
	if self.Databases.Postgres.URL == "" {
		return errors.New(fmt.Sprintf("%s %s", "databases.postgres.main_url", errMsg))
	}
	if self.Databases.Postgres.Test == "" {
		return errors.New(fmt.Sprintf("%s %s", "databases.postgres.test_url", errMsg))
	}
	if self.Databases.Postgres.Admin == "" {
		return errors.New(fmt.Sprintf("%s %s", "databases.postgres.admin_url", errMsg))
	}

	if self.Databases.Couchbase.URL == "" {
		return errors.New(fmt.Sprintf("%s %s", "databases.couchbase.url", errMsg))
	}

	return nil
}
