package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/sys/unix"

	"microtecture/infrastructure/config"
)

var rootCli = &cobra.Command{
	Use:   config.NAME,
	Short: config.DESCRIPTION,
	PersistentPreRun: func(cli *cobra.Command, args []string) {
		if !terminal.IsTerminal(unix.Stdout) {
			logrus.SetFormatter(&logrus.JSONFormatter{})
		} else {
			logrus.SetFormatter(&logrus.TextFormatter{
				FullTimestamp:   true,
				TimestampFormat: time.RFC3339Nano,
			})
		}
	},
}

func Execute() {
	_ = rootCli.Execute()
}

var configFile string

func init() {
	if os.Getenv(config.ENVIRONMENT_NAME) == "" {
		fmt.Println(
			`Set`,
			config.ENVIRONMENT_NAME,
			`environment variable to "dev" or "prod" for development or production.`,
		)
		os.Exit(0)
	}

	var configFileName string
	if os.Getenv(config.ENVIRONMENT_NAME) == "dev" {
		configFileName = config.CONFIG_FILE_NAME
	}

	cobra.OnInitialize(initConfig)
	rootCli.PersistentFlags().StringVarP(
		&configFile,
		"config",
		"c",
		configFileName,
		"set config file in yml format.",
	)
}

func initConfig() {
	viper.SetConfigFile(configFile)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if err.Error() == `Config File "config" Not Found in "[]"` {
			fmt.Println(`App is in production mode, use --config to set config file.`)
			os.Exit(0)
		}

		fmt.Println(err)
		os.Exit(0)
	}
}
