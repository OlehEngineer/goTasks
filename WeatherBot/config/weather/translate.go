package weather

// translation structure for respond to user
type Translation struct {
	VisibilityUS    string `env:"VisibilityUS"`
	CloudsUS        string `env:"CloudsUS"`
	RainUS          string `env:"RainUS"`
	SnowUS          string `env:"SnowUS"`
	TempUS          string `env:"TempUS"`
	FeelsLikeUS     string `env:"FeelsLikeUS"`
	TempMinUS       string `env:"TempMinUS"`
	TempMaxUS       string `env:"TempMaxUS"`
	PressureUS      string `env:"PressureUS"`
	HumidityUS      string `env:"HumidityUS"`
	WindSpeedUS     string `env:"WindSpeedUS"`
	SunRiseUS       string `env:"SunRiseUS"`
	SunSetUS        string `env:"SunSetUS"`
	PrecipitationUS string `env:"PrecipitationUS"`
	VisibilityUA    string `env:"VisibilityUA"`
	CloudsUA        string `env:"CloudsUA"`
	RainUA          string `env:"RainUA"`
	SnowUA          string `env:"SnowUA"`
	TempUA          string `env:"TempUA"`
	FeelsLikeUA     string `env:"FeelsLikeUA"`
	TempMinUA       string `env:"TempMinUA"`
	TempMaxUA       string `env:"TempMaxUA"`
	PressureUA      string `env:"PressureUA"`
	HumidityUA      string `env:"HumidityUA"`
	WindSpeedUA     string `env:"WindSpeedUA"`
	SunRiseUA       string `env:"SunRiseUA"`
	SunSetUA        string `env:"SunSetUA"`
	PrecipitationUA string `env:"PrecipitationUA"`
}
