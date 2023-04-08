package main

import (
	config "WeatherBot/config/config"
	"WeatherBot/config/geoapi"
	"WeatherBot/config/weather"
	"WeatherBot/logger"

	"fmt"

	"github.com/caarlos0/env/v7"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

var BotConfig = config.BotConfiguration{} //genral bot configuration constants
var CurrentCitylist map[string][]float32  //map with cities for user choosing.

func main() {
	FullCityList := make(map[int64]map[string][]float32) //map which store userID and cities for the user choosing and geo coordintater of cities
	UserLanguages := make(map[int64]string)              //map which store user id and user's language
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
	// START BOT
	bot, errBot := tgbotapi.NewBotAPI(BotConfig.BotToken)
	if errBot != nil {
		log.Panicf("cannot start the Bot. Error - %s\n", errBot)
	}
	//set Bot debug mode. Could be TRUE or FALSE
	bot.Debug = BotConfig.BotDebugMode

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	defer func() {
		if r := recover(); r != nil {
			// Log the error and continue program execution
			log.Printf("Recovered from panic: %v", r)
		}
	}()
	//get updates from Telegram API
	updates, errUpd := bot.GetUpdatesChan(u)
	if errUpd != nil {
		log.Errorf("Cannot get update. Error - %s\n", errUpd)
	}

	for update := range updates {

		if update.Message != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			switch update.Message.Text {
			case "/about":
				msg.Text = fmt.Sprintf("%s\n%s", BotConfig.AboutEN, BotConfig.AboutUA)
			case "/start":
				msg.Text = fmt.Sprintf("%s\n%s", BotConfig.AskUserCityEN, BotConfig.AskUserCityUA)
			default:
				userLanguage := update.Message.From.LanguageCode     //get user's language
				UserLanguages[update.Message.Chat.ID] = userLanguage //save user's language to map

				var geoErr error

				if userLanguage == "uk" || userLanguage == "ua" {
					msg.ReplyMarkup, CurrentCitylist, geoErr = geoapi.GetGeoLocation(update.Message.Text, BotConfig.WeatherGeoToken, BotConfig.Geolimit, BotConfig.GeoBodyLink)
					if geoErr != nil {
						msg.Text = BotConfig.GeoAPIErrorMessageUA
					} else {
						msg.Text = BotConfig.DefaultUA
						FullCityList[update.Message.Chat.ID] = CurrentCitylist //add cities to the general list
					}
				} else {
					msg.ReplyMarkup, CurrentCitylist, geoErr = geoapi.GetGeoLocation(update.Message.Text, BotConfig.WeatherGeoToken, BotConfig.Geolimit, BotConfig.GeoBodyLink)
					if geoErr != nil {
						msg.Text = BotConfig.GeoAPIErrorMessageUS
					} else {
						msg.Text = BotConfig.DefaultUS
						FullCityList[update.Message.Chat.ID] = CurrentCitylist //add cities to the general list
					}
				}
			}
			if _, err := bot.Send(msg); err != nil {
				log.Error(err)
			}
		} else if update.CallbackQuery != nil {

			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.Text)

			userCity := update.CallbackQuery.Data                       //get city which user chosen
			lang := UserLanguages[update.CallbackQuery.Message.Chat.ID] //get user's language
			// if user use old buttons with city name it cause error "index out of range [0] with length 0"
			if _, ok := FullCityList[update.CallbackQuery.Message.Chat.ID][userCity]; ok {

				coordintates := FullCityList[update.CallbackQuery.Message.Chat.ID][userCity] //get chosen city's geo coordinates

				forecast, weatherErr := weather.GetWeather(coordintates[0], coordintates[1], BotConfig.WeatherBodyLink, BotConfig.WeatherGeoToken, lang)

				//check if the weather forecast available
				if weatherErr != nil {
					if lang == "ua" || lang == "uk" {
						msg.Text = BotConfig.WeatherForecastErrorUA
					} else {
						msg.Text = BotConfig.WeatherForecastErrorUS
					}
				} else {
					msg.Text = forecast
				}
				if _, err := bot.Send(msg); err != nil {
					log.Error(err)
				}
			} else {
				// in case there is not available data in FullCityList bot returm message to user
				switch lang {
				case "ua":
					msg.Text = fmt.Sprintf("%s", BotConfig.OldButtonClickUA)
				case "uk":
					msg.Text = fmt.Sprintf("%s", BotConfig.OldButtonClickUA)
				default:
					msg.Text = fmt.Sprintf("%s", BotConfig.OldButtonClickUS)
				}
				if _, err := bot.Send(msg); err != nil {
					log.Error(err)
				}
				continue
			}
		}
	}
}
