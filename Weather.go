package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type HourlyWeather struct {
	Context  []interface{} `json:"@context"`
	Type     string        `json:"type"`
	Geometry struct {
		Type        string        `json:"type"`
		Coordinates [][][]float64 `json:"coordinates"`
	} `json:"geometry"`
	Properties struct {
		Updated           time.Time `json:"updated"`
		Units             string    `json:"units"`
		ForecastGenerator string    `json:"forecastGenerator"`
		GeneratedAt       time.Time `json:"generatedAt"`
		UpdateTime        time.Time `json:"updateTime"`
		ValidTimes        time.Time `json:"validTimes"`
		Elevation         struct {
			Value    float64 `json:"value"`
			UnitCode string  `json:"unitCode"`
		} `json:"elevation"`
		Periods []struct {
			Number           int         `json:"number"`
			Name             string      `json:"name"`
			StartTime        string      `json:"startTime"`
			EndTime          string      `json:"endTime"`
			IsDaytime        bool        `json:"isDaytime"`
			Temperature      int         `json:"temperature"`
			TemperatureUnit  string      `json:"temperatureUnit"`
			TemperatureTrend interface{} `json:"temperatureTrend"`
			WindSpeed        string      `json:"windSpeed"`
			WindDirection    string      `json:"windDirection"`
			Icon             string      `json:"icon"`
			ShortForecast    string      `json:"shortForecast"`
			DetailedForecast string      `json:"detailedForecast"`
		} `json:"periods"`
	} `json:"properties"`
}

var hourlyWeather HourlyWeather

func main() {
	fmt.Println("Calling API...")
	client := &http.Client{}
	response, err := http.Get("https://api.weather.gov/gridpoints/TOP/31,80/forecast/hourly")
	//req, err := http.NewRequest("GET", "https://icanhazdadjoke.com/", nil)
	if err != nil {
		fmt.Print(err.Error())
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Print(err.Error())
	}

	json.Unmarshal(bodyBytes, &hourlyWeather)
	fmt.Printf(" Hourly Weather Type %+v\n", hourlyWeather.Type)
	fmt.Printf(" Temperature on Sundaye %+v\n", hourlyWeather.Properties.Periods[0].Temperature)
}
