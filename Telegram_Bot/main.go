package main

import (
	"telegramBot/holidays"

	"github.com/caarlos0/env/v7"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

var CurrentConfiguration = Config{} // variable with current Bot configuration

func main() {
	//real .ENV configuration file
	envErr := godotenv.Load("local.env")
	if envErr != nil {
		log.Errorf(".ENV file is missing. Error > %v < occurs.", envErr)
	}

	// parsing .ENV file with "github.com/caarlos0/env/v7"
	configErr := env.Parse(&CurrentConfiguration)
	if configErr != nil {
		log.Errorf(".ENV file parsing error => %v.", &configErr)
	}
	//start logging
	LogFile, logErr := StartLogging()
	if logErr != nil {
		log.Errorf("problem with log file, error => %v", logErr)
	}
	defer LogFile.Close()

	//get Bot Token and Bot API link from current configuration struct Config{}
	botToken := CurrentConfiguration.Token

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = CurrentConfiguration.Botdebug //current Bot setting. Could be TRUE or FALSE

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	//check each update, any feedback from user
	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			// Extract the command from the Message.
			switch update.Message.Command() {
			case "help":
				msg.Text = CurrentConfiguration.HELP
			case "about":
				msg.Text = CurrentConfiguration.ABOUT
			case "holidays":
				msg.ReplyMarkup = SendKeyboardToUser(CurrentConfiguration) // send to user keyboard with countries list
			case "stop":
				msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true) // remove keyboard with countries list
			default:
				msg.Text = CurrentConfiguration.DefaultText // default answer to unknown command by User. Defined in the .ENV file
			}
			if _, err := bot.Send(msg); err != nil {
				log.Errorf("Bot send message error - %v\n", err)
			}

		} else {
			// check text message from User
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

			switch update.Message.Text {
			case "🇺🇦 UA", "🇵🇱 PL", "🇩🇪 DE", "🇫🇷 FR", "🇯🇵 JP", "🇬🇧 GB", "🇨🇦 CA", "🇺🇸 US":
				msg.Text = holidays.MakeHolidayRequest(update.Message.Text, CurrentConfiguration.Apitoken)
			default:
				msg.Text = "not known country"
			}
			if _, err := bot.Send(msg); err != nil {
				log.Errorf("Bot send message error - %v\n", err)
			}

		}
	}
}
