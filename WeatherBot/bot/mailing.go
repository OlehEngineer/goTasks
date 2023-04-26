package bot

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"

	"github.com/OlehEngineer/goTask/WeatherBot/config/weather"
	DB "github.com/OlehEngineer/goTask/WeatherBot/database"
)

// check mailing time for each user in the database and send message to the user in time.
func ScheduleMailingSender(collection *mongo.Collection, bot *tgbotapi.BotAPI, weatherBodyLink string, apiToken string) error {

	client, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Errorf("Cannot find mailing time in the DB. Error - %v", err)
		return err
	}
	defer client.Close(context.Background())

	for client.Next(context.Background()) {
		var user DB.UserSubData
		if err := client.Decode(&user); err != nil {
			log.Errorf("Cannot decode user data from the database. Error - %v", err)

		}
		now := time.Now()
		if now.Hour() == user.MailingTime.Hour() && now.Minute() == user.MailingTime.Minute() {
			ID := user.User_id
			latitude := user.Location.Latitude
			longitude := user.Location.Longitude
			language := user.Language

			userWeather, err := weather.GetWeather(latitude, longitude, weatherBodyLink, apiToken, language)
			if err != nil {
				log.Errorf("Cannot get weather forecast. Please try once more time. Error - %v", err)
			}

			message := tgbotapi.NewMessage(ID, userWeather)
			bot.Send(message)
		}
	}
	return nil
}

// Ticker  generate cycle in one minute for mailing time checking
func Ticker(collection *mongo.Collection, bot *tgbotapi.BotAPI, weatherBodyLink string, apiToken string) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			err := ScheduleMailingSender(collection, bot, weatherBodyLink, apiToken)
			if err != nil {
				log.Errorf("Cannot send scheduled message to the user. Error - %v", err)
			}
		}
	}
}
