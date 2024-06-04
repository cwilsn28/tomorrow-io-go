package tomorrowio

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const TESTLOC = "42.3584, -71.0597"
const TIMESTEP = "1h"

func TestEnvLoad(t *testing.T) {
	// Make sure we can load a .env file in the cwd.
	LoadEnv()

	// Make sure the key=val pairs are read and set in the environment.
	apikey := os.Getenv("APIKEY")
	assert.NotEqual(t, apikey, "")

	// Cleanup
	os.Unsetenv("APIKEY")
}

func TestWeatherForecastClient(t *testing.T) {
	LoadEnv()
	apikey := os.Getenv("APIKEY")

	// Create API client.
	// Make sure the client can execute a request and no error is returned.
	apiClient := NewAPIClient(apikey)

	// // Create and execute a test API reqVuest
	apiClient.FormatReqURL(TESTLOC, TIMESTEP)
	resp, err := apiClient.ExecuteRequest()

	// Errors should be nil.
	assert.Equal(t, err, nil)

	// Response should not be an empty string.
	assert.NotEqual(t, resp, "")

	// Cleanup
	os.Unsetenv("APIKEY")
}

func TestCoreWeatherForecast(t *testing.T) {
	var err error
	var output_filename = "write_test.csv"

	LoadEnv()
	apikey := os.Getenv("APIKEY")

	// Make sure we can retrive a forecast and no error is returned.
	t.Log("Instantiating CoreWeatherForecast class...")
	fx, err := WeatherForecast(apikey, TESTLOC, TIMESTEP)
	assert.Equal(t, err, nil)

	// Make sure the forecast object is not empty.
	t.Log("Checking forecast is not empty...")
	assert.NotEqual(t, len(fx.Timelines.Forecasts), 0)

	// Make sure daily forecasts have 24 recordsgg
	t.Log("Checking length of daily forecasts...")
	forecasts := fx.DailyForecasts()
	for _, vals := range forecasts {
		assert.Equal(t, len(vals), 24)
	}

	// Test forecast val methods.
	// Ensure each forcast has positive length.
	t.Log("Checking length of forecast vals...")
	for fxDate, _ := range forecasts {
		cloudcover := fx.CloudCover(fxDate)
		assert.Equal(t, len(cloudcover), 24)
		temperature := fx.Temperature(fxDate)
		assert.Equal(t, len(temperature), 24)
		windgust := fx.WindGust(fxDate)
		assert.Equal(t, len(windgust), 24)
		windspeed := fx.WindSpeed(fxDate)
		assert.Equal(t, len(windspeed), 24)
	}

	// Test ability to write forecast to csv
	t.Log("Checking forecast export to csv...")
	err = fx.CommonForecastToCSV("2024-05-14", output_filename)
	assert.Equal(t, err, nil)

	fileBytes, err := os.ReadFile(output_filename)
	assert.Equal(t, err, nil)
	assert.Greater(t, len(fileBytes), 0)

	// Cleanup
	os.Remove(output_filename)
	os.Unsetenv("APIKEY")
}
