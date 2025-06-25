package main

// WeatherPoints 表示 /points/latitude,longitude 的响应结构
// 主要关注 Properties.Forecast 字段
type WeatherPoints struct {
	Properties struct {
		Forecast string `json:"forecast"`
	} `json:"properties"`
}

// Forecast 表示 /gridpoints/.../forecast 的响应结构
type Forecast struct {
	Properties struct {
		Periods []struct {
			Name             string `json:"name,omitempty"`
			Temperature      int    `json:"temperature,omitempty"`
			TemperatureUnit  string `json:"temperatureUnit,omitempty"`
			WindSpeed        string `json:"windSpeed,omitempty"`
			WindDirection    string `json:"windDirection,omitempty"`
			ShortForecast    string `json:"shortForecast,omitempty"`
			DetailedForecast string `json:"detailedForecast,omitempty"`
		} `json:"periods,omitempty"`
	} `json:"properties,omitempty"`
}
