package tomorrowio

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strconv"
	"strings"
)

type CoreWeatherForecast struct {
	Timelines HourlyForecasts `json:"timelines,omitempty`
}

type HourlyForecasts struct {
	Forecasts []HourlyForecast `json:"hourly,omitempty"`
}

type HourlyForecast struct {
	Time   string       `json:"time,omitempty"`
	Values ForecastVals `json:"values,omitempty"`
}

type ForecastVals struct {
	CloudBase                float64 `json:"cloudBase,omitempty"`
	CloudCeiling             float64 `json:"cloudCeiling,omitempty"`
	CloudCover               float64 `json:"cloudCover,omitempty"`
	DewPoint                 float64 `json:"dewPoint,omitempty"`
	Evapotranspiration       float64 `json:"evapotranspiration,omitempty"`
	FreezingRainIntensity    float64 `json:"freezingRainIntensity,omitempty"`
	Humidity                 float64 `json:"humidity,omitempty"`
	IceAccumulation          float64 `json:"iceAccumulation,omitempty"`
	IceAccumulationLwe       float64 `json:"iceAccumulationLwe,omitempty"`
	PrecipitationProbability float64 `json:"precipitationProbability,omitempty"`
	PressureSurfaceLevel     float64 `json:"pressureSurfaceLevel,omitempty"`
	RainAccumulation         float64 `json:"rainAccumulation,omitempty"`
	RainAccumulationLwe      float64 `json:"rainAccumulationLwe,omitempty"`
	RainIntensity            float64 `json:"rainIntensity,omitempty"`
	SleetAccumulation        float64 `json:"sleetAccumulation,omitempty"`
	SleetAccumulationLwe     float64 `json:"sleetAccumulationLwe,omitempty"`
	SleetIntensity           float64 `json:"sleetIntensity,omitempty"`
	SnowAccumulation         float64 `json:"snowAccumulation,omitempty"`
	SnowAccumulationLwe      float64 `json:"snowAccumulationLwe,omitempty"`
	SnowDepth                float64 `json:"snowDepth,omitempty"`
	SnowIntensity            float64 `json:"snowIntensity,omitempty"`
	Temperature              float64 `json:"temperature,omitempty"`
	TemperatureApparent      float64 `json:"temperatureApparent,omitempty"`
	UvHealthConcern          float64 `json:"uvHealthConcern,omitempty"`
	UvIndex                  float64 `json:"uvIndex,omitempty"`
	Visibility               float64 `json:"visibility,omitempty"`
	WeatherCode              float64 `json:"weatherCode,omitempty"`
	WindDirection            float64 `json:"windDirection,omitempty"`
	WindGust                 float64 `json:"windGust,omitempty"`
	WindSpeed                float64 `json:"windSpeed,omitempty"`
}

/* ---
Forecast methods
--- */

func (f *CoreWeatherForecast) DailyForecasts() map[string][]HourlyForecast {
	dailyForecasts := make(map[string][]HourlyForecast)
	for _, fx := range f.Timelines.Forecasts {
		fxDate := strings.Split(fx.Time, "T")[0]

		if _, inMap := dailyForecasts[fxDate]; inMap {
			dailyForecasts[fxDate] = append(dailyForecasts[fxDate], fx)
		}
	}
	return dailyForecasts
}

func (f *CoreWeatherForecast) ForecastForDay(fxDate string) ([]HourlyForecast, error) {
	dailyForecasts := f.DailyForecasts()

	if fx, exists := dailyForecasts[fxDate]; exists {
		return fx, nil
	}

	err := errors.New(fmt.Sprintf("No forecast for date %s", fxDate))
	return make([]HourlyForecast, 0), err
}

func (f *CoreWeatherForecast) ToJSON(filename string) error {
	jsonBYTES, err := json.Marshal(f)
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, jsonBYTES, fs.FileMode(os.O_RDWR))
	return err
}

/* ---
Access to some common forecast vals
--- */

func (f *CoreWeatherForecast) CloudCover(forecastDate string) [][]interface{} {
	cloudCover := make([][]interface{}, 0)
	for _, fx := range f.Timelines.Forecasts {
		if compareDates(forecastDate, fx.Time) {
			cloudCover = append(cloudCover, []interface{}{fx.Time, fx.Values.CloudCover})
		}
	}
	return cloudCover
}

func (f *CoreWeatherForecast) Temperature(forecastDate string) [][]interface{} {
	temperature := make([][]interface{}, 0)
	for _, fx := range f.Timelines.Forecasts {
		if compareDates(forecastDate, fx.Time) {
			temperature = append(temperature, []interface{}{fx.Time, fx.Values.Temperature})
		}
	}
	return temperature
}

func (f *CoreWeatherForecast) WindGust(forecastDate string) [][]interface{} {
	windGust := make([][]interface{}, 0)
	for _, fx := range f.Timelines.Forecasts {
		if compareDates(forecastDate, fx.Time) {
			windGust = append(windGust, []interface{}{fx.Time, fx.Values.WindGust})
		}
	}
	return windGust
}

func (f *CoreWeatherForecast) WindSpeed(forecastDate string) [][]interface{} {
	windSpeed := make([][]interface{}, 0)
	for _, fx := range f.Timelines.Forecasts {
		if compareDates(forecastDate, fx.Time) {
			windSpeed = append(windSpeed, []interface{}{fx.Time, fx.Values.WindSpeed})
		}
	}
	return windSpeed
}

func (f *CoreWeatherForecast) CommonForecastToCSV(forecastDate, filename string) error {
	cloudCover := f.CloudCover(forecastDate)
	temperature := f.Temperature(forecastDate)
	windGust := f.WindGust(forecastDate)
	windSpeed := f.WindGust(forecastDate)

	// Start with the header row
	rows := CommonForecastHeader

	for i := 0; i < len(cloudCover); i++ {
		timeChunks := strings.Split(cloudCover[i][0].(string), "T")
		row := []string{
			timeChunks[0],
			timeChunks[1],
			strconv.FormatFloat(cloudCover[i][1].(float64), 'f', -1, 64),
			strconv.FormatFloat(temperature[i][1].(float64), 'f', -1, 64),
			strconv.FormatFloat(windGust[i][1].(float64), 'f', -1, 64),
			strconv.FormatFloat(windSpeed[i][1].(float64), 'f', -1, 64),
		}
		rows = append(rows, row)
	}

	outfile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	w := csv.NewWriter(outfile)
	err = w.WriteAll(rows)
	return err
}

/* ---
Local helpers
--- */

func compareDates(dateString1, dateString2 string) bool {
	d1 := strings.Split(dateString1, "T")[0]
	d2 := strings.Split(dateString2, "T")[0]

	if d1 == d2 {
		return true
	}
	return false
}
