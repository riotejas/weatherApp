package clients

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchClientDataForecast(t *testing.T) {
	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Mock response data
		mockResponse := `{
			"@context": [],
			"type": "Feature",
			"geometry": {
				"type": "Polygon",
				"coordinates": [[[0, 0], [1, 1], [2, 2]]]
			},
			"properties": {
				"units": "us",
				"forecastGenerator": "generator",
				"generatedAt": "2023-08-01T00:00:00Z",
				"updateTime": "2023-08-01T01:00:00Z",
				"validTimes": "2023-08-01T00:00:00Z/2023-08-02T00:00:00Z",
				"elevation": {
					"unitCode": "wmoUnit:m",
					"value": 10
				},
				"periods": []
			}
		}`
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockResponse))
	}))
	defer mockServer.Close()

	// Call the function with the mock server URL
	result, err := fetchClientData[forecast](mockServer.URL)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, "Feature", result.Type)
	assert.Equal(t, "us", result.Properties.Units)
}

func TestFetchClientDataGrid(t *testing.T) {
	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Mock response data
		mockResponse := `{
			"@context": [],
			"id": "gridId",
			"type": "Feature",
			"geometry": {
				"type": "Point",
				"coordinates": [1.0, 2.0]
			},
			"properties": {
				"@id": "gridId",
				"@type": "Grid",
				"cwa": "CWA",
				"forecastOffice": "Office",
				"gridId": "GridId",
				"gridX": 1,
				"gridY": 2,
				"forecast": "forecast",
				"forecastHourly": "forecastHourly",
				"forecastGridData": "forecastGridData",
				"observationStations": "stations",
				"relativeLocation": {
					"type": "Point",
					"geometry": {
						"type": "Point",
						"coordinates": [1.0, 2.0]
					},
					"properties": {
						"city": "City",
						"state": "State",
						"distance": {
							"unitCode": "wmoUnit:km",
							"value": 10.0
						},
						"bearing": {
							"unitCode": "wmoUnit:degree",
							"value": 180
						}
					}
				},
				"forecastZone": "Zone",
				"county": "County",
				"fireWeatherZone": "FireZone",
				"timeZone": "TimeZone",
				"radarStation": "Radar"
			}
		}`
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockResponse))
	}))
	defer mockServer.Close()

	// Call the function with the mock server URL
	result, err := fetchClientData[grid](mockServer.URL)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, "gridId", result.Id)
	assert.Equal(t, "Feature", result.Type)
	assert.Equal(t, "CWA", result.Properties.Cwa)
}
