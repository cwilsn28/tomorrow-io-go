package tomorrowio

import (
	"io"
	"net/http"
	"net/url"
	"time"
)

const BASE_URL = "https://api.tomorrow.io"
const API_PATH = "/v4/weather/forecast"
const DefaultHTTPTimeout = 80 * time.Second

type APIClient struct {
	HTTPClient *http.Client
	APIKey     string
	fxURL      string
}

/* Return an instance of an API client object to the caller */
func NewAPIClient(apikey string) *APIClient {
	client := APIClient{
		HTTPClient: NewHTTPClient(),
		APIKey:     apikey,
	}

	return &client
}

func (c *APIClient) FormatReqURL(location, timesteps string) {
	// Create a base URL object.
	reqURL, err := url.Parse(BASE_URL)
	if err != nil {
		panic(err)
	}

	// Add the API path.
	reqURL.Path += API_PATH

	// Add the query parameters.
	queryParams := url.Values{}
	queryParams.Add("location", location)
	queryParams.Add("timesteps", timesteps)
	queryParams.Add("apikey", c.APIKey)

	// Encode and set the query string.
	reqURL.RawQuery = queryParams.Encode()
	c.fxURL = reqURL.String()
}

func (c *APIClient) ExecuteRequest() ([]byte, error) {
	var bodyBytes []byte
	var err error

	// Create the forecast request.
	req, err := http.NewRequest("GET", c.fxURL, nil)
	if err != nil {
		return bodyBytes, err
	}

	// Execute the request.
	req.Header.Set("accept", "application/json")
	resp, err := c.HTTPClient.Do(req)
	bodyBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		return bodyBytes, err
	}

	return bodyBytes, nil
}
