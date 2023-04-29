package geoapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

// return inline keyboard markup list of available cities with name user inputted.
func GetGeoLocation(userCity string, token string, limit int, linkBody string) (tgbotapi.InlineKeyboardMarkup, map[string][]float32, error) {
	//userCity - name of city got from user
	// token - token for API access
	// limit - const of quantity of returned cities by API
	//linkBody - the main body of API link
	//http://api.openweathermap.org/geo/1.0/direct?q=London&limit=5&appid={API key}
	city := strings.TrimSpace(userCity)
	link := fmt.Sprintf("%s%s&limit=%v&appid=%s", linkBody, city, limit, token)

	//create pseudo returned data
	pseudoKeyboard, pseudoMap := pseudoData()

	resp, err := http.Get(link)
	if err != nil {
		log.Errorf("Cannot GET feedback from Geo Api. Error - %v", err)
		return pseudoKeyboard, pseudoMap, err
	}
	defer resp.Body.Close()

	keyboard, cities, pasErr := GeoAPiResponseParsing(resp.Body)

	if pasErr != nil {
		return pseudoKeyboard, pseudoMap, pasErr
	}
	return keyboard, cities, nil

}

// parse GEO API response and create map with cities variants and geo coordinates
func GeoAPiResponseParsing(respBody io.ReadCloser) (tgbotapi.InlineKeyboardMarkup, map[string][]float32, error) {

	//create pseudo returned data
	pseudoKeyboard, pseudoMap := pseudoData()

	bytes, err := ioutil.ReadAll(respBody)
	if err != nil {
		log.Errorf("JSON answer parsing problem. Error - %s", err)
		return pseudoKeyboard, pseudoMap, err
	}

	if len(bytes) < 10 {
		return pseudoKeyboard, pseudoMap, errors.New("no data/немає даних")
	}

	var cities []Geolocation
	err = json.Unmarshal(bytes, &cities)
	if err != nil {
		log.Errorf("Problem during parsed data reading. Error - %s", err)
		return pseudoKeyboard, pseudoMap, err
	}
	listOfVariants := make(map[string][]float32)

	for _, city := range cities {
		coordinates := []float32{city.Latitude, city.Longitude}
		listOfVariants[city.CityName+"_"+city.Country] = append(listOfVariants[city.CityName], coordinates...)
	}
	return generateInlineKeyboardMarkup(listOfVariants), listOfVariants, err

}

// return inline keyboard markup for user choose based on cities list
func generateInlineKeyboardMarkup(data map[string][]float32) tgbotapi.InlineKeyboardMarkup {

	var buttons []tgbotapi.InlineKeyboardButton
	var rows [][]tgbotapi.InlineKeyboardButton

	for k, _ := range data {
		button := tgbotapi.NewInlineKeyboardButtonData(k, k)
		buttons = append(buttons, button)

		// Create a new row after every 2 buttons
		if len(buttons)%2 == 0 {
			rows = append(rows, buttons)
			buttons = []tgbotapi.InlineKeyboardButton{}
		}
	}
	// Add any remaining buttons to the last row
	if len(buttons) > 0 {
		rows = append(rows, buttons)
	}
	markup := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return markup
}

// generate pseudo data (empty) inline keyboard markup
func pseudoData() (tgbotapi.InlineKeyboardMarkup, map[string][]float32) {

	pseudoKeyboard := tgbotapi.InlineKeyboardMarkup{InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{}}

	pseudoMap := make(map[string][]float32)

	return pseudoKeyboard, pseudoMap
}
