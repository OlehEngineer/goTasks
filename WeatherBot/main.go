package main

import (
	config "github.com/OlehEngineer/goTask/WeatherBot/config/config"

	logger "github.com/OlehEngineer/goTask/WeatherBot/logger"

	"github.com/OlehEngineer/goTask/WeatherBot/bot"

	"github.com/caarlos0/env/v8"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

var BotConfig = config.BotConfiguration{} //genral bot configuration constants

func main() {

	//strat logging
	logger.StartLogging(BotConfig.LogLvl)

	//read .ENV file for general Bot configuration
	errEnv := godotenv.Load("set.env")
	if errEnv != nil {
		log.Fatalf("Cannot read .ENV file. Error - %s\n", &errEnv)
	}
	// set configuration variable and parse .ENV file
	errConfg := env.Parse(&BotConfig)
	if errConfg != nil {
		log.Fatalf("cannot parse .ENV file. Error - %s\n", errConfg)
	}

	bot.StartWeatherBot(BotConfig)
}
