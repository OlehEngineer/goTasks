package config

import "github.com/caarlos0/env/v8"

//general Bot configuration constants
type BotConfiguration struct {
	BotToken               string `env:"TGBOTTOKEN"`
	WeatherGeoToken        string `env:"WEATHERAPITOKEN"`
	AboutEN                string `env:"ABOUTEN"`
	AboutUA                string `env:"ABOUTUA"`
	AskUserCityEN          string `env:"CITYNAMEREQUESTEN"`
	AskUserCityUA          string `env:"CITYNAMEREQUESTUA"`
	BotDebugMode           bool   `env:"BOTDEBUG" default:"false"`
	GeoBodyLink            string `env:"GEOBODYLINK"`
	WeatherBodyLink        string `env:"WEATHERBODYLINK"`
	Geolimit               int    `env:"GEOLIMIT" default:"5"`
	DefaultUA              string `env:"DEFAULTUK"`
	DefaultUS              string `env:"DEFAULTUS"`
	LogLvl                 string `env:"LOGLEVEL"`
	GeoAPIErrorMessageUS   string `env:"GeoAPIErrorMessageUS"`
	GeoAPIErrorMessageUA   string `env:"GeoAPIErrorMessageUA"`
	WeatherForecastErrorUA string `env:"WeatherForecastErrorUA"`
	WeatherForecastErrorUS string `env:"WeatherForecastErrorUS"`
	OldButtonClickUA       string `env:"OldButtonClickUA"`
	OldButtonClickUS       string `env:"OldButtonClickUS"`
	MongoDB                string `env:"MONGODB"`
	DBname                 string `env:"DBNAME"`
	DBcollectionName       string `env:"DBCOLLECTIONNAME"`
	MailingTimeUA          string `env:"MAILINGTIMEUA"`
	MailingTimeUS          string `env:"MAILINGTIMEUS"`
	YouSubcribeUA          string `env:"YOUSUBSCRIBENUA"`
	YouSubcribeUS          string `env:"YOUSUBSCRIBENUS"`
	AlreadySubcribenUA     string `env:"ALREADYSUBSCRIBENUA"`
	AlreadySubcribenUS     string `env:"ALREADYSUBSCRIBENUS"`
}

//get environment variables from .ENV file
func LoadBotConfiguration() (*BotConfiguration, error) {
	var cfg BotConfiguration
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
