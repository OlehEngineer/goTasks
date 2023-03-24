package main

import (
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// .ENV file structure
type Config struct {
	Token       string `env:"BOTTOKEN"`
	Myinfo      string `env:"MYINFO" default:"no personal info"`
	MuGitHub    string `env:"MYGITHUB" default:"https://github.com/"`
	Myemail     string `env:"MYMAIL" default:"nomail@server.com"`
	Botapi      string `env:"BOTAPI" default:"https://api.telegram.org/bot"`
	LogrusLevel string `env:"LOGRUSLEVEL" default:"info"`
	Apitoken    string `env:"APITOKEN"`
	Botdebug    bool   `env:"BOTDEBUG" default:"false"`
	ABOUT       string `env:"ABOUT"`
	HELP        string `env:"HELP"`
	DefaultText string `env:"DEFAULTTEXT"`

	Country_1 string `env:"Country1"`
	Country_2 string `env:"Country2"`
	Country_3 string `env:"Country3"`
	Country_4 string `env:"Country4"`
	Country_5 string `env:"Country5"`
	Country_6 string `env:"Country6"`
	Country_7 string `env:"Country7"`
	Country_8 string `env:"Country8"`

	Country_1_Flag string `env:"FLAG1"`
	Country_2_Flag string `env:"FLAG2"`
	Country_3_Flag string `env:"FLAG3"`
	Country_4_Flag string `env:"FLAG4"`
	Country_5_Flag string `env:"FLAG5"`
	Country_6_Flag string `env:"FLAG6"`
	Country_7_Flag string `env:"FLAG7"`
	Country_8_Flag string `env:"FLAG8"`
}

func CountryFlagUnicode(Ucode string) string {
	flag, err := strconv.Unquote(`"` + Ucode + `"`)
	if err != nil {
		return ""
	}
	return flag
}

func SendKeyboardToUser(CurrentConfiguration Config) tgbotapi.ReplyKeyboardMarkup {
	//keyboard with list of countries for choosing by user
	var CountryList = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(CountryFlagUnicode(CurrentConfiguration.Country_1_Flag)+" "+CurrentConfiguration.Country_1),
			tgbotapi.NewKeyboardButton(CountryFlagUnicode(CurrentConfiguration.Country_2_Flag)+" "+CurrentConfiguration.Country_2),
			tgbotapi.NewKeyboardButton(CountryFlagUnicode(CurrentConfiguration.Country_3_Flag)+" "+CurrentConfiguration.Country_3),
			tgbotapi.NewKeyboardButton(CountryFlagUnicode(CurrentConfiguration.Country_4_Flag)+" "+CurrentConfiguration.Country_4),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(CountryFlagUnicode(CurrentConfiguration.Country_5_Flag)+" "+CurrentConfiguration.Country_5),
			tgbotapi.NewKeyboardButton(CountryFlagUnicode(CurrentConfiguration.Country_6_Flag)+" "+CurrentConfiguration.Country_6),
			tgbotapi.NewKeyboardButton(CountryFlagUnicode(CurrentConfiguration.Country_7_Flag)+" "+CurrentConfiguration.Country_7),
			tgbotapi.NewKeyboardButton(CountryFlagUnicode(CurrentConfiguration.Country_8_Flag)+" "+CurrentConfiguration.Country_8),
		),
	)
	return CountryList
}
