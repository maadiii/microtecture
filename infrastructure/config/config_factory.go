package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/logrusorgru/aurora"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type Config interface {
	Init() error
}

func ConfigFactory(initializer int) (Config, error) {
	switch initializer {
	case DATASTORE_CONFIG:
		config := new(DataStoreConfig)
		if err := config.Init(); err != nil {
			return nil, errors.New(err.Error())
		}
		return config, nil
	case APPLICATION_CONFIG:
		config := new(ApplicationConfig)
		if err := config.Init(); err != nil {
			return nil, errors.New(err.Error())
		}
		return config, nil
	default:
		return nil, errors.New(fmt.Sprintf("Initializer method %d not recogonized.", initializer))
	}
}

func marshal(config interface{}) error {
	var configFilePath string

	if viper.ConfigFileUsed() == "" {
		viper.SetConfigFile(CONFIG_FILE_NAME)
	}

	if os.Getenv(ENVIRONMENT_NAME) == "dev" {
		_, b, _, _ := runtime.Caller(0)
		dir := filepath.Dir(filepath.Dir(filepath.Dir(b)))
		configFilePath = filepath.Join(dir, viper.ConfigFileUsed())
	} else if os.Getenv(ENVIRONMENT_NAME) == "prod" {
		configFilePath = viper.ConfigFileUsed()
	} else {
		panic(aurora.BgRed(ENVIRONMENT_NAME + " environment varialbe not set."))
	}

	configFile, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return errors.New(err.Error())
	}
	err = yaml.Unmarshal(configFile, config)
	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}
