# go-alphavantage

A partial Go wrapper for the tomorrow.io weather forecast API.

**Please Note:** This project is a work in progress.

## Introduction

This package provides access to the tomorrow.io weather forecast API. For a 
complete treatment of tomorrow.io APIs, please visit (tomorrow.io)[https://www.tomorrow.io/].

## Usage

```bash
go get github.com/cwilsn28/tomorrow-io-go
```


```go

    ...

	LoadEnv()
	apikey := os.Getenv("APIKEY")

    // Please refer to https://docs.tomorrow.io/reference/weather-forecast 
    // for a complete discussion on location and timestep values.
    location := "42.3584, -71.0597"
    timestep := "1h"

	fx, err := WeatherForecast(apikey, location, timestep)

    ...

```

From the tomorrow.io API documentation:

Using the weather forecast API you can access up-to-date weather information for your location, including minute-by-minute forecasts (for premium users) for the next hour, hourly forecasts for the next 120 hours, and daily forecasts for the next 5 days.

### Accessing daily forecasts

```go

    ...

    daily := fx.DailyForecasts()

```

### Accessing common forecasts for a forecast date.

```go

    ...

    temps := fx.Temperature("2024-02-01")
    windspeed := fx.Windspeed("2024-02-01")
    windgust := fx.WindGust("2024-02-01")
    cloudcover := fx.CloudCover("2024-02-01")

    ...
```

### Exporting forecast data

#### Dumping to json

```go

    ...
    
    err := fx.ToJSON()

    ...

```

#### Exporting common forecasts to csv

```go

    ...

    err = fx.CommonForecastToCSV("2024-02-01", output_filename)

    ...
```