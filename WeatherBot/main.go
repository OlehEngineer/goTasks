package main

import (
	logger "github.com/OlehEngineer/goTask/WeatherBot/logger"
	log "github.com/sirupsen/logrus"

	"github.com/OlehEngineer/goTask/WeatherBot/bot"
	"github.com/OlehEngineer/goTask/WeatherBot/config/config"
	DB "github.com/OlehEngineer/goTask/WeatherBot/database"
)

func main() {

	//get general Bot configuration
	BotConfig, cfgErr := config.LoadBotConfiguration()
	if cfgErr != nil {
		log.Fatalf("Cannot read configuration file. Error - %v", cfgErr)
	}

	//start logging
	logger.StartLogging(BotConfig.LogLvl)

	//connect to the MongoDB
	clientDB, errDBconnect := DB.ConnectDataBase()
	if errDBconnect != nil {
		log.Fatalf("Cannot connect to the MongoDB. Error - %s", errDBconnect)
	}
	log.Info("Connected to the database")

	//start BOT itself
	bot.StartWeatherBot(clientDB, BotConfig)
}
