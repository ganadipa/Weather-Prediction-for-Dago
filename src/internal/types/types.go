package types

type WeatherResponse struct {
	Data []Data `json:"data"`
}

type DescriptionOnlyData struct {
	Data []DescriptionOnly `json:"data"`
}

type DescriptionOnly struct {
	Datetime string  `json:"datetime"`
	Weather  Weather `json:"weather"`
}

type Data struct {
	Datetime      string  `json:"datetime"`
	Weather       Weather `json:"weather"`
}
type Weather struct {
	Description string `json:"description"`
}

type HistoricalDatum struct {
	Datetime           string
	WeatherDescription string
}
