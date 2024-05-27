package types

type WeatherResponse struct {
	Data []Data `json:"data"`
}

type Data struct {
	Datetime string  `json:"datetime"`
	Weather  Weather `json:"weather"`
}

type Weather struct {
	Description string `json:"description"`
}
