package main

import (
	"github.com/caarlos0/env/v7"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

var CurrentConfiguration = Config{} // variable with current Bot configuration

func main() {
	//real .ENV configuration file
	envErr := godotenv.Load("local.env")
	if envErr != nil {
		log.Fatalf(".ENV file is missing. Error > %s < occurs.", envErr)
	}

	// parsing .ENV file with "github.com/caarlos0/env/v7"
	configErr := env.Parse(&CurrentConfiguration)
	if configErr != nil {
		log.Fatalf(".ENV file parsing error => %s.", &configErr)
	}
	//start logging
	LogFile, logErr := StartLogging()
	if logErr != nil {
		log.Fatalf("problem with log file, error => %s", logErr)
	}
	defer LogFile.Close()

	//get Bot Token and Bot API link from current configuration struct Config{}
	botToken := CurrentConfiguration.Token
	botURL := CurrentConfiguration.Botapi + botToken
	offset := 0 // using as counter for UpdateID incresing through the looping

	//↓ infinity cycle for updates checking ↓
	for {
		updates, err := getUpdates(botURL, offset)
		if err != nil {
			log.Errorf("Something went wrong: ", err.Error())
		}
		//iteration throught each element of updates
		for _, update := range updates {
			err = respond(botURL, update)
			if err != nil {
				log.Error(err)
			}
			//new updates id. Increase for each new update
			offset = update.UpdateId + 1
		}
		//write logging parameters
		for _, currentUpd := range updates {
			log.Infof("Update ID - %v, Chat ID - %v, Massage - «%s»\n", currentUpd.UpdateId, currentUpd.Message.Chat.ChatId, currentUpd.Message.Text)
		}
	}
}
