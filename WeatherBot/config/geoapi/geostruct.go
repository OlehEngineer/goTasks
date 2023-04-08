package geoapi

//Structure of GEO API response
type Geolocation struct {
	CityName  string     `json:"name"`
	LocalName LocalNames `json:"local_names"`
	Latitude  float32    `json:"lat"`
	Longitude float32    `json:"lon"`
	Country   string     `json:"country"`
	State     string     `json:"state"`
}
type LocalNames struct {
	FeatureName string `json:"feature_name"`
	EN          string `json:"en"`
	UA          string `json:"uk"`
}
