package main

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/mark3labs/mcp-go/mcp"
	"io"
	"net/http"
)

const (
	wfApiBae           = "https://api.weather.gov"
	userAgent          = "weather-app/1.0"
	maxForecastPeriods = 5
)

func makeCallRequest(url string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func callForecast(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid arguments: %v", request.Params.Arguments)
	}
	latitude, ok := args[WeiDu].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid latitude argument: %v", args[WeiDu])
	}
	longitude, ok := args[JingDu].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid longitude argument: %v", args[JingDu])
	}

	pointsURL := fmt.Sprintf("%s/points/%.4f,%.4f", wfApiBae, latitude, longitude)
	points, err := makeCallRequest(pointsURL)
	if err != nil {
		return nil, fmt.Errorf("call points API failed: %v", err)
	}
	pointsData := &WeatherPoints{}
	if err := json.Unmarshal(points, pointsData); err != nil {
		return nil, fmt.Errorf("parse points call response failed: %v", err)
	}
	if pointsData == nil || pointsData.Properties.Forecast == "" {
		return nil, fmt.Errorf("no forecast url found")
	}

	//第二次请求weather request API
	forecastResp, err := makeCallRequest(pointsData.Properties.Forecast)
	if err != nil {
		return nil, fmt.Errorf("get weather forecast data failed: %v", err)
	}
	forecastData := &Forecast{}
	if err := json.Unmarshal(forecastResp, forecastData); err != nil {
		return nil, fmt.Errorf("parse forecast call response failed: %v", err)
	}

	type ForecastResult struct {
		Name        string `json:"Name"`
		Temperature string `json:"Temperature"`
		Wind        string `json:"Wind"`
		Forecast    string `json:"Forecast"`
	}

	// 整理前5个预报的结果
	forecast := make([]ForecastResult, maxForecastPeriods)
	for i, period := range forecastData.Properties.Periods {
		if i == maxForecastPeriods {
			break
		}
		forecast[i] = ForecastResult{
			Name:        period.Name,
			Temperature: fmt.Sprintf("%d°%s", period.Temperature, period.TemperatureUnit),
			Wind:        fmt.Sprintf("%s %s", period.WindSpeed, period.WindDirection),
			Forecast:    period.DetailedForecast,
		}
	}

	bytes, err := json.Marshal(forecast)
	if err != nil {
		return nil, fmt.Errorf("marshal forecast call response failed: %v", err)
	}
	return mcp.NewToolResultText(string(bytes)), nil
}
