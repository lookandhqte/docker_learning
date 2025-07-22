package meow

type WeatherData struct {
	Dt   int64 `json:"dt"`
	Main struct {
		Temp     float64 `json:"temp"`
		Humidity int     `json:"humidity"`
	} `json:"main"`
	Weather []struct {
		Icon        string `json:icon`
		Description string `json:"description"`
	} `json:"weather"`
	Wind struct {
		Speed float64 `json:speed`
	} `json:"wind"`
	Name string `json:"name"`
}

type WeatherResponse struct {
	Cod     string        `json:"cod"`
	Message int           `json:"message"`
	List    []WeatherData `json: "list"`
	City    struct {
		Name string `json:"name"`
	} `json:"city"`
}

type GeocodingAPI []struct {
	Name string  `json:"name"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
}
