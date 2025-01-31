package config

import (
	logging "github.com/ipfs/go-log/v2"
	"github.com/spf13/viper"
	"strconv"
)

var log = logging.Logger("config")

type AppConfig struct {
	DB_USERNAME string
	DB_PASSWORD string
	DB_HOSTNAME string
	DB_PORT     int
	DB_NAME     string

	GOOGLE_APPLICATION_PROJECT_ID string
	GOOGLE_APPLICATION_BUCKET     string

	LOCAL_PATH          string
	REPLACE_PREFIX_PATH string
	WORK_COUNT          int
}

func InitConfig() *AppConfig {
	app := AppConfig{}

	viper.AddConfigPath(".")
	viper.SetConfigFile("local.yml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("read local.yml config error:", err)
		return nil
	}
	app.DB_USERNAME = viper.Get("DB_USERNAME").(string)
	app.DB_PASSWORD = viper.Get("DB_PASSWORD").(string)
	app.DB_HOSTNAME = viper.Get("DB_HOSTNAME").(string)
	app.DB_PORT, _ = strconv.Atoi(viper.Get("DB_PORT").(string))
	app.DB_NAME = viper.Get("DB_NAME").(string)
	app.GOOGLE_APPLICATION_PROJECT_ID = viper.Get("GOOGLE_APPLICATION_PROJECT_ID").(string)
	app.GOOGLE_APPLICATION_BUCKET = viper.Get("GOOGLE_APPLICATION_BUCKET").(string)
	app.LOCAL_PATH = viper.Get("LOCAL_PATH").(string)
	app.REPLACE_PREFIX_PATH = viper.Get("REPLACE_PREFIX_PATH").(string)
	app.WORK_COUNT, _ = strconv.Atoi(viper.Get("WORK_COUNT").(string))
	if app.WORK_COUNT <= 0 {
		app.WORK_COUNT = 10
	}

	log.Infof("successful read config: %+v\n", app)

	return &app
}
