package tomorrowio

import (
	"encoding/json"
)

func WeatherForecast(apikey, location, timesteps string) (CoreWeatherForecast, error) {
	forecast := CoreWeatherForecast{}

	// Create an API client and format the request
	apiClient := NewAPIClient(apikey)
	apiClient.FormatReqURL(location, timesteps)

	// Execute the request and return the forecast
	rawJSON, err := apiClient.ExecuteRequest()
	if err != nil {
		return forecast, err
	}

	err = json.Unmarshal(rawJSON, &forecast)
	return forecast, err
}
