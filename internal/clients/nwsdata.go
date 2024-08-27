package clients

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"
	wamiddleware "weatherApp/internal/middleware"
)

// National Weather Service API

var Vendor = "nws"

// NWSDataService use an interface for dependency injection.
// Allows switching to other clients w/ minimal upstream code changes
type NWSDataService interface {
	Forecast(context.Context) ([]byte, error)
}

type nwsData struct {
	url string
}

func NewNWSDataService(url string) NWSDataService {
	nd := &nwsData{url: url}
	return nd
}

func (nd *nwsData) Forecast(ctx context.Context) ([]byte, error) {
	// get user given lat and long from query
	lat := ctx.Value(wamiddleware.LatKey).(string)
	lng := ctx.Value(wamiddleware.LngKey).(string)

	// 2-step process to get forecast data from NWS.
	// First, given latitude and longitude, make a call to get grid coordinates
	// Second, using grid coordinates, make second call to get forecast

	gridRequestUrl := fmt.Sprintf("%s/points/%s,%s", nd.url, lat, lng)
	slog.Info("nws grid request", "url", gridRequestUrl)

	// get forecast URL
	gridData, err := fetchGridData(gridRequestUrl)
	if err != nil {
		return nil, err
	}

	forecastRequestUrl := gridData.Properties.Forecast
	slog.Info("nws forecast request", "url", forecastRequestUrl)
	// get forecast data
	forecastData, err := fetchForecastData(forecastRequestUrl)
	if err != nil {
		return nil, err
	}
	if forecastData.Properties.Periods == nil || len(forecastData.Properties.Periods) == 0 {
		return nil, errors.New("forecast period is empty")
	}
	// we only care about the first period, e.g., morning, afternoon
	period := forecastData.Properties.Periods[1]

	feels := generateFeels(period.Temperature)
	slog.Info("returning nws forecast",
		"period", period.Name,
		"temp", period.Temperature,
		"feels", feels,
		"forecast", period.ShortForecast)

	res := map[string]string{
		"period":   period.Name,
		"temp":     strconv.Itoa(period.Temperature),
		"feels":    feels,
		"forecast": period.ShortForecast,
	}
	jsonResp, _ := json.Marshal(res)
	return jsonResp, nil
}

func generateFeels(temp int) string {
	switch {
	case temp < 33:
		return "freezing"
	case temp < 60:
		return "cold"
	case temp < 70:
		return "cool"
	case temp < 80:
		return "moderate"
	case temp < 95:
		return "warm"
	case temp < 105:
		return "hot"
	default:
		return "boiling"
	}
}

// fetchGridData get the grid data which contains the forecast URL
func fetchGridData(url string) (grid, error) {
	var result grid
	resBody, err := fetchClientData(url)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(resBody, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

// fetchForecastData get the forecast data
func fetchForecastData(url string) (forecast, error) {
	var result forecast
	resBody, err := fetchClientData(url)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(resBody, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

// fetchClientData handle the actual client request
func fetchClientData(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating nwsData request: %w", err)
	}
	req.Header.Set("User-Agent", "WeatherApp/v1")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending nwsData request: %w", err)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading nwsData body: %w", err)
	}
	return resBody, nil
}

type forecast struct {
	Context  []interface{} `json:"@context"`
	Type     string        `json:"type"`
	Geometry struct {
		Type        string        `json:"type"`
		Coordinates [][][]float64 `json:"coordinates"`
	} `json:"geometry"`
	Properties struct {
		Units             string    `json:"units"`
		ForecastGenerator string    `json:"forecastGenerator"`
		GeneratedAt       time.Time `json:"generatedAt"`
		UpdateTime        time.Time `json:"updateTime"`
		ValidTimes        string    `json:"validTimes"`
		Elevation         struct {
			UnitCode string  `json:"unitCode"`
			Value    float64 `json:"value"`
		} `json:"elevation"`
		Periods []Period `json:"periods"`
	} `json:"properties"`
}

type Period struct {
	Number                     int       `json:"number"`
	Name                       string    `json:"name"`
	StartTime                  time.Time `json:"startTime"`
	EndTime                    time.Time `json:"endTime"`
	IsDaytime                  bool      `json:"isDaytime"`
	Temperature                int       `json:"temperature"`
	TemperatureUnit            string    `json:"temperatureUnit"`
	TemperatureTrend           string    `json:"temperatureTrend"`
	ProbabilityOfPrecipitation struct {
		UnitCode string `json:"unitCode"`
		Value    *int   `json:"value"`
	} `json:"probabilityOfPrecipitation"`
	WindSpeed        string `json:"windSpeed"`
	WindDirection    string `json:"windDirection"`
	Icon             string `json:"icon"`
	ShortForecast    string `json:"shortForecast"`
	DetailedForecast string `json:"detailedForecast"`
}

type grid struct {
	Context  []interface{} `json:"@context"`
	Id       string        `json:"id"`
	Type     string        `json:"type"`
	Geometry struct {
		Type        string    `json:"type"`
		Coordinates []float64 `json:"coordinates"`
	} `json:"geometry"`
	Properties struct {
		Id                  string `json:"@id"`
		Type                string `json:"@type"`
		Cwa                 string `json:"cwa"`
		ForecastOffice      string `json:"forecastOffice"`
		GridId              string `json:"gridId"`
		GridX               int    `json:"gridX"`
		GridY               int    `json:"gridY"`
		Forecast            string `json:"forecast"`
		ForecastHourly      string `json:"forecastHourly"`
		ForecastGridData    string `json:"forecastGridData"`
		ObservationStations string `json:"observationStations"`
		RelativeLocation    struct {
			Type     string `json:"type"`
			Geometry struct {
				Type        string    `json:"type"`
				Coordinates []float64 `json:"coordinates"`
			} `json:"geometry"`
			Properties struct {
				City     string `json:"city"`
				State    string `json:"state"`
				Distance struct {
					UnitCode string  `json:"unitCode"`
					Value    float64 `json:"value"`
				} `json:"distance"`
				Bearing struct {
					UnitCode string `json:"unitCode"`
					Value    int    `json:"value"`
				} `json:"bearing"`
			} `json:"properties"`
		} `json:"relativeLocation"`
		ForecastZone    string `json:"forecastZone"`
		County          string `json:"county"`
		FireWeatherZone string `json:"fireWeatherZone"`
		TimeZone        string `json:"timeZone"`
		RadarStation    string `json:"radarStation"`
	} `json:"properties"`
}
