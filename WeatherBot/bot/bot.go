package bot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"

	config "github.com/OlehEngineer/goTask/WeatherBot/config/config"
	"github.com/OlehEngineer/goTask/WeatherBot/config/geoapi"
	"github.com/OlehEngineer/goTask/WeatherBot/config/weather"
	DB "github.com/OlehEngineer/goTask/WeatherBot/database"
)

var CurrentCitylist map[string][]float32 //map with cities for user choosing.

func StartWeatherBot(client *mongo.Client, BotConfig *config.BotConfiguration) {

	FullCityList := make(map[int64]map[string][]float32) //map which store userID and cities for the user choosing and geo coordintater of cities
	UserLanguages := make(map[int64]string)              //map which store user id and user's language
	TempUserData := make(map[int64]DB.UserSubData)       //temporary user's data for subscription
	DBcollection := client.Database(BotConfig.DBname).Collection(BotConfig.DBcollectionName)

	// START BOT
	Bot, errBot := tgbotapi.NewBotAPI(BotConfig.BotToken)
	if errBot != nil {
		log.Panicf("cannot start the Bot. Error - %s\n", errBot)
	}
	//set Bot debug mode. Could be TRUE or FALSE
	Bot.Debug = BotConfig.BotDebugMode
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// ! start checking mailing time for subscribed users
	go Ticker(DBcollection, Bot, BotConfig.WeatherBodyLink, BotConfig.WeatherGeoToken)

	//get updates from Telegram API
	updates, errUpd := Bot.GetUpdatesChan(u)
	if errUpd != nil {
		log.Errorf("Cannot get update. Error - %s\n", errUpd)
	}

	for update := range updates {

		if update.Message != nil {
			userID := update.Message.Chat.ID
			userLanguage := update.Message.From.LanguageCode //get user's language
			msg := tgbotapi.NewMessage(userID, update.Message.Text)

			switch update.Message.Text {
			case "/about":
				msg.Text = fmt.Sprintf("%s\n%s", BotConfig.AboutEN, BotConfig.AboutUA)
			case "/start":
				msg.Text = fmt.Sprintf("%s\n%s", BotConfig.AskUserCityEN, BotConfig.AskUserCityUA)
			case "/subscribe":

				isUserSubs, err := DB.IsUserAlreadySubscribed(DBcollection, userID)
				if err != nil {
					log.Errorf("Cannot check if the user is subscribed. Error - %s", err)
				}
				if isUserSubs == true {
					if userLanguage == "uk" || userLanguage == "ua" {
						msg.Text = BotConfig.AlreadySubcribenUA
					} else {
						msg.Text = BotConfig.AlreadySubcribenUS
					}
				} else {
					TempUserData[userID] = DB.UserSubData{User_id: userID, Language: userLanguage, Code: 1}
					if userLanguage == "uk" || userLanguage == "ua" {
						msg.Text = BotConfig.AskUserCityUA
					} else {
						msg.Text = BotConfig.AskUserCityEN
					}
					log.Info(TempUserData)
				}

			case "/unsubscribe":
				isUserSubs, err := DB.IsUserAlreadySubscribed(DBcollection, userID)
				if err != nil {
					log.Errorf("Cannot check if the user is subscribed. Error - %s", err)
				}
				if isUserSubs == true {
					err := DB.DeleteUsersSubscription(DBcollection, userID)
					if err != nil {
						msg.Text = "Cannot unsubscribe."
					} else {
						msg.Text = "You successfully unsubscribed"
					}
					Bot.Send(msg)
				}

			default:
				UserLanguages[userID] = userLanguage //save user's language to map

				if TempUserData[userID].Code == 2 {
					isTime, userTime := DB.IsitTime(update.Message.Text)
					if isTime == true {

						if t, ok := TempUserData[userID]; ok {
							t.MailingTime = userTime
							TempUserData[userID] = t
						}

						isSubsc, err := DB.SubscribeUser(DBcollection, TempUserData[userID])
						if err != nil {
							log.Errorf("Cannot subscribe the user. Error - %s", err)
						}
						if isSubsc == true {
							log.Info(TempUserData)
							if userLanguage == "uk" || userLanguage == "ua" {
								msg.Text = BotConfig.YouSubcribeUA
							} else {
								msg.Text = BotConfig.YouSubcribeUS
							}
							Bot.Send(msg)
							continue
						} else {
							msg.Text = "Cannot subscribe you"
						}

					} else {
						msg.Text = "You enter incorrect time. Please try one more time"
					}

				}
				var geoErr error

				if userLanguage == "uk" || userLanguage == "ua" {
					msg.ReplyMarkup, CurrentCitylist, geoErr = geoapi.GetGeoLocation(update.Message.Text, BotConfig.WeatherGeoToken, BotConfig.Geolimit, BotConfig.GeoBodyLink)
					if geoErr != nil {
						msg.Text = BotConfig.GeoAPIErrorMessageUA
					} else {
						msg.Text = BotConfig.DefaultUA
						FullCityList[userID] = CurrentCitylist //add cities to the general list
					}
				} else {
					msg.ReplyMarkup, CurrentCitylist, geoErr = geoapi.GetGeoLocation(update.Message.Text, BotConfig.WeatherGeoToken, BotConfig.Geolimit, BotConfig.GeoBodyLink)
					if geoErr != nil {
						msg.Text = BotConfig.GeoAPIErrorMessageUS
					} else {
						msg.Text = BotConfig.DefaultUS
						FullCityList[userID] = CurrentCitylist //add cities to the general list
					}
				}
			}
			if _, err := Bot.Send(msg); err != nil {
				log.Error(err)
			}
		} else if update.CallbackQuery != nil {
			userID := update.CallbackQuery.Message.Chat.ID
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.Text)
			userCity := update.CallbackQuery.Data //get city which user chosen
			lang := UserLanguages[userID]         //get user's language

			// if user use old buttons with city name it cause error "index out of range [0] with length 0"
			if _, ok := FullCityList[userID][userCity]; ok {

				coordinates := FullCityList[userID][userCity] //get chosen city's geo coordinates

				// check user who just chose the city is under subscribe procedure
				if TempUserData[userID].Code == 1 {
					if _, ok := TempUserData[userID]; ok {

						TempUserData[userID] = DB.UserSubData{
							User_id: userID,
							Location: DB.GeoCoordinates{
								Latitude:  coordinates[0],
								Longitude: coordinates[1]},
							Language: UserLanguages[userID],
							Code:     2}
						if lang == "ua" || lang == "uk" {
							msg.Text = BotConfig.MailingTimeUA
						} else {
							msg.Text = BotConfig.MailingTimeUS
						}
						Bot.Send(msg)
					}
					log.Info(TempUserData)
				}

				forecast, weatherErr := weather.GetWeather(coordinates[0], coordinates[1], BotConfig.WeatherBodyLink, BotConfig.WeatherGeoToken, lang)

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
				if _, err := Bot.Send(msg); err != nil {
					log.Error(err)
				}
			} else {
				// in case there is not available data in FullCityList bot return message to user
				switch lang {
				case "ua":
					msg.Text = fmt.Sprintf("%s", BotConfig.OldButtonClickUA)
				case "uk":
					msg.Text = fmt.Sprintf("%s", BotConfig.OldButtonClickUA)
				default:
					msg.Text = fmt.Sprintf("%s", BotConfig.OldButtonClickUS)
				}
				if _, err := Bot.Send(msg); err != nil {
					log.Error(err)
				}
				continue
			}
		}
	}
}
