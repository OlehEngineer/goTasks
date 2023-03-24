package holidays

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

func MakeHolidayRequest(userCountry string, token string) string {
	country := userCountry[len(userCountry)-2:] //cut the country's abreviation from user's input
	t := time.Now()
	year := t.Year()
	month := int(t.Month())
	day := t.Day()
	link := ("https://holidays.abstractapi.com/v1/?api_key=" + token + "&country=" + country + "&year=" + strconv.Itoa(year) + "&month=" + strconv.Itoa(month) + "&day=" + strconv.Itoa(day))

	resp, err := http.Get(link)
	if err != nil {
		log.Fatalf("Holiday request error => %s\n", err)
	}

	defer resp.Body.Close()

	var holidays []APIResponse // slice of struct
	err = json.NewDecoder(resp.Body).Decode(&holidays)
	if err != nil {
		log.Error(err)
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
	LocalName   string `json;"name_local"`
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
