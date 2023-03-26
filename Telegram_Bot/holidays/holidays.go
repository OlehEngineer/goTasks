package holidays

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

func MakeHolidayRequest(userCountry string, token string) string {
	country := userCountry[len(userCountry)-2:] //cut the country's abreviation from user's input
	t := time.Now()
	link, err := url.Parse("https://holidays.abstractapi.com/v1/")
	if err != nil {
		log.Errorf("Cannot build the API request URL. Error - %v", err)
		return fmt.Sprintf("Cannot build the API request URL. Error - %v", err)
	}
	query := link.Query()
	query.Set("api_key", token)
	query.Set("country", country)
	query.Set("year", strconv.Itoa(t.Year()))
	query.Set("month", strconv.Itoa(int(t.Month())))
	query.Set("day", strconv.Itoa(int(t.Day())))
	link.RawQuery = query.Encode()

	resp, err := http.Get(link.String())
	if err != nil {
		log.Errorf("API request error => %v\n", err)
		return fmt.Sprintf("API connection problem. Error - %v. Please try again later", err)
	}

	defer resp.Body.Close()

	var holidays []APIResponse // slice of struct
	err = json.NewDecoder(resp.Body).Decode(&holidays)
	if err != nil {
		log.Errorf("Decoding problem. Error - %v", err)
		return fmt.Sprintf("API response reading problem. Error - %v. Please try again later", err)
	}

	var HolidayList []string //slice of holidays in chosen country for today
	for _, holiday := range holidays {
		HolidayList = append(HolidayList, holiday.Name)
	}
	if len(HolidayList) >= 1 {
		answer := strings.Join(HolidayList, ", ")
		return answer
	} else {
		return "there is no holiday today"
	}
}

type APIResponse struct {
	Name        string `json:"name"`
	LocalName   string `json:"name_local"`
	Language    string `json:"language"`
	Description string `json:"description"`
	Country     string `json:"country"`
	Location    string `json:"location"`
	Type        string `json:"type"`
	Date        string `json:"date"`
	Date_year   string `json:"date_year"`
	Date_month  string `json:"date_month"`
	Date_day    string `json:"date_day"`
	Week_day    string `json:"week_day"`
}
