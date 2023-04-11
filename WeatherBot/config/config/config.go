package config

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
}
