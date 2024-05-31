package types

type WeatherResponse struct {
	Data []Data `json:"data"`
}

type Data struct {
	Datetime      string  `json:"datetime"`
	Temperature   float64 `json:"temp"`
	Humidity      float64 `json:"rh"`
	Precipitation float64 `json:"precip"`
	Weather       Weather `json:"weather"`
	WindSpeed     float64 `json:"wind_spd"`
}
type Weather struct {
	Description string `json:"description"`
}
