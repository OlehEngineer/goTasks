package weather

// Weather API response structure
type WeatherForecast struct {
	Coordinates GeoCoordinates `json:"coord"`
	WeatherInfo []Weather      `json:"weather"`
	Base        string         `json:"base"` //Internal parameter
	Forecast    Main           `json:"main"`
	Visibility  int            `json:"visibility"` //Visibility, meter. The maximum value of the visibility is 10km
	Windinfo    Wind           `json:"wind"`
	Cloudsinfo  Clouds         `json:"clouds"`
	RainInfo    Rain           `json:"rain"`
	SnowInfo    Snow           `json:"snow"`
	Dt          int            `json:"dt"` //Time of data calculation, unix, UTC
	SunInfo     Sys            `json:"sys"`
	TimeZone    int            `json:"timezone"` //Shift in seconds from UT
	Id          int            `json:"id"`       //City ID. Please note that built-in geocoder functionality has been deprecated.
	Cityname    string         `json:"name"`     //City name
	Cod         int            `json:"cod"`      //Internal parameter
}
type GeoCoordinates struct {
	Longitude float32 `json:"lon"` //City geo location, longitude
	Latitude  float32 `json:"lat"` //City geo location, latitude
}
type Weather struct {
	ID          int    `json:"id"`          //Weather condition id
	Main        string `json:"main"`        //Group of weather parameters (Rain, Snow, Extreme etc.)
	Description string `json:"description"` //Weather condition within the group. You can get the output in your language.
	Icon        string `json:"icon"`        //Weather icon id
}
type Main struct {
	Temp        float32 `json:"temp"`       //Temperature. Unit Default: Kelvin, Metric: Celsius, Imperial: Fahrenheit.
	FeelsLike   float32 `json:"feels_like"` //Temperature. This temperature parameter accounts for the human perception of weather. Unit Default: Kelvin, Metric: Celsius, Imperial: Fahrenhei
	TempMin     float32 `json:"temp_min"`   //Minimum temperature at the moment. This is minimal currently observed temperature (within large megalopolises and urban areas). Unit Default: Kelvin, Metric: Celsius, Imperial: Fahrenheit.
	TempMax     float32 `json:"temp_max"`   //Maximum temperature at the moment. This is maximal currently observed temperature (within large megalopolises and urban areas). Unit Default: Kelvin, Metric: Celsius, Imperial: Fahrenheit.
	Pressure    int     `json:"pressure"`   //Atmospheric pressure (on the sea level, if there is no sea_level or grnd_level data), hPa
	Humidity    int     `json:"humidity"`   //Humidity, %
	SeaLevel    int     `json:"sea_level"`  //Atmospheric pressure on the sea level, hPa
	GroundLevel int     `json:"grnd_level"` //Atmospheric pressure on the ground level, hPa
}
type Wind struct {
	Speed  float32 `json:"speed"` //Wind speed. Unit Default: meter/sec, Metric: meter/sec, Imperial: miles/hour.
	Degree int     `json:"deg"`   //Wind direction, degrees (meteorological)
	Gust   float32 `json:"gust"`  //Wind gust. Unit Default: meter/sec, Metric: meter/sec, Imperial: miles/hour
}
type Rain struct {
	Rain1H float32 `json:"1h"` //Rain volume for the last 1 hour, mm
	Rain3H float32 `json:"3h"` //Rain volume for the last 3 hours, mm
}
type Snow struct {
	Snow1H float32 `json:"1h"` //Snow volume for the last 1 hour, mm
	Snow3H float32 `json:"3h"` //Snow volume for the last 3 hours, mm
}
type Clouds struct {
	All int `json:"all"` //Cloudiness, %
}
type Sys struct {
	Type    int    `json:"type"`    //Internal parameter
	Id      int    `json:"id"`      //Internal parameter
	Message string `json:"message"` //Internal parameter
	Country string `json:"country"` //Country code (GB, JP etc.)
	SunRise int    `json:"sunrise"` //Sunrise time, unix, UTC
	SunSet  int    `json:"sunset"`  //Sunset time, unix, UTC
}
