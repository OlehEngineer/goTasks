package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

var RespTransl = Translation{} //struct with translation for respond to user

// send GET request to Weather API
func GetWeather(lat float32, lon float32, weatherBodyLink string, apiToken string, lang string) (string, error) {

	apiLink := fmt.Sprintf("%s%f&lon=%f&appid=%s&units=metric&lang=%s", weatherBodyLink, lat, lon, apiToken, lang)

	resp, err := http.Get(apiLink)
	if err != nil {
		log.Errorf("Cannot GET feedback from Weather API. Error - %v", err)
		return "Cannot get weather forecast. Please try once more time.", err
	}
	defer resp.Body.Close()
	return WeatherAPIResponseParsing(resp.Body, lang)

}

// parse Weather API response and build response to user
func WeatherAPIResponseParsing(respBody io.ReadCloser, lang string) (string, error) {

	bytes, err := ioutil.ReadAll(respBody)
	if err != nil {
		log.Errorf("JSON answer parsing problem. Error - %s", err)
		return "No data available", err
	}

	currentWeather := WeatherForecast{} // current weather forecast for chosen city
	errWeather := json.Unmarshal(bytes, &currentWeather)
	if errWeather != nil {
		log.Errorf("Cannot parse the JSON weather response. Error - %v", errWeather)
		return "No data available", errWeather
	}

	// check if bot got correct response from weather server
	if currentWeather.Cod != 200 {
		return "cannot reach weather server", errors.New("Cannot reach weather server")
	}

	//load .ENV file with response translation on two languages
	errTrans := godotenv.Load("config/weather/lang.env")
	if errTrans != nil {
		log.Errorf("Cannot load translate.go file. Error - %v", errTrans)
		return "Cannot send weather forecast on your language", errTrans
	}

	errTransConf := env.Parse(&RespTransl)
	if errTransConf != nil {
		log.Errorf("Cannot read translate.go file. Error - %v", errTransConf)
		return "Cannot send weather forecast on your language", errTransConf
	}

	//general weather data
	weather := currentWeather.Cloudsinfo.All
	weatherDesc := currentWeather.WeatherInfo[0].Description
	precipitation := precipitationCheck(currentWeather, lang)
	pressure := pressureConverter(currentWeather)
	temperature := currentWeather.Forecast.Temp
	feelsLikeTemp := currentWeather.Forecast.FeelsLike
	visibility := currentWeather.Visibility
	windSpeed := currentWeather.Windinfo.Speed
	sunRise := timeConverter(currentWeather.SunInfo.SunRise, currentWeather.TimeZone)
	sunSet := timeConverter(currentWeather.SunInfo.SunSet, currentWeather.TimeZone)

	//final response to user
	weatherResponse := createWeatherResponse(weather, visibility, weatherDesc, precipitation, sunRise, sunSet, lang, pressure, temperature, feelsLikeTemp, windSpeed)

	return weatherResponse, nil
}

// convert parsed time value into time string for message
func timeConverter(UTCtime, offset int) string {
	utcTime := time.Unix(int64(UTCtime), 0).UTC()
	localTime := utcTime.Add(time.Duration(offset) * time.Second)

	return localTime.Format("15:04:05")
}

// convert pressure value from hPa to millimeters of mercury
func pressureConverter(weather WeatherForecast) float32 {
	k := 1.333
	hPaPressure := weather.Forecast.Pressure
	if hPaPressure > 0 {
		return float32(hPaPressure) / float32(k)
	} else {
		return float32(0)
	}
}

// check is there any precipitation information inside weather API response
func precipitationCheck(weather WeatherForecast, language string) string {
	//return information about precipitation
	rainValue := weather.RainInfo.Rain1H
	snowValue := weather.SnowInfo.Snow1H
	rainUA, rainUS := "Дощ", "Rain"
	snowUA, snowUS := "Сніг", "Snow"

	if language == "ua" || language == "uk" {
		if rainValue > 0 && snowValue > 0 {
			return fmt.Sprintf("%s: %v мм\t%s: %v мм", rainUA, rainValue, snowUA, snowValue)
		} else if rainValue > 0 && snowValue == 0 {
			return fmt.Sprintf("%s: %v мм", rainUA, rainValue)
		} else if rainValue == 0 && snowValue > 0 {
			return fmt.Sprintf("%s: %v мм", snowUA, snowValue)
		} else {
			return "Без опадів"
		}
	} else {
		if rainValue > 0 && snowValue > 0 {
			return fmt.Sprintf("%s: %v mm\t%s: %v mm", rainUS, rainValue, snowUS, snowValue)
		} else if rainValue > 0 && snowValue == 0 {
			return fmt.Sprintf("%s: %v mm", rainUS, rainValue)
		} else if rainValue == 0 && snowValue > 0 {
			return fmt.Sprintf("%s: %v mm", snowUS, snowValue)
		} else {
			return "There is no precipitation"
		}
	}
}

// build message to the user based on weather API response
func createWeatherResponse(weather, visibility int, weatherDesc, precipitation, sunRise, sunSet, lang string, pressure, temperature, feelsLikeTemp, windSpeed float32) string {
	var weatherResponse string
	if lang == "ua" || lang == "uk" {
		weatherResponse = fmt.Sprintf("%s: %s - %v%%\n%s\n%s: %v℃, %s: %v℃\n%s: %v м\n%s: %v м/с\n%s: %0.2f ммРс\n%s: %v;\t%s: %v",
			RespTransl.CloudsUA, weatherDesc, weather,
			precipitation,
			RespTransl.TempMaxUA, temperature, RespTransl.FeelsLikeUA, feelsLikeTemp,
			RespTransl.VisibilityUA, visibility,
			RespTransl.WindSpeedUA, windSpeed,
			RespTransl.PressureUA, pressure,
			RespTransl.SunRiseUA, sunRise, RespTransl.SunSetUA, sunSet)
	} else {
		weatherResponse = fmt.Sprintf("%s: %s - %v%%\n%s\n%s: %v℃, %s: %v℃\n%s: %v m\n%s: %v m/s\n%s: %0.2f mmHg\n%s: %v;\t%s: %v",
			RespTransl.CloudsUS, weatherDesc, weather,
			precipitation,
			RespTransl.TempMaxUS, temperature, RespTransl.FeelsLikeUS, feelsLikeTemp,
			RespTransl.VisibilityUS, visibility,
			RespTransl.WindSpeedUS, windSpeed,
			RespTransl.PressureUS, pressure,
			RespTransl.SunRiseUS, sunRise, RespTransl.SunSetUS, sunSet)
	}
	return weatherResponse
}
