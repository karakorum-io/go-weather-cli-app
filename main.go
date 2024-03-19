package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
	model "github.com/karakorum.io/mausam/model"
)

func main() {
	color.Blue("Karakorum Weather Software")

	q := "agra"

	if len(os.Args) >= 2 {
		q = os.Args[1]
	}

	res, err := http.Get("http://api.weatherapi.com/v1/forecast.json?key=503ee91321924099894192614241903&q=" + q + "&days=1&aqi=no&alerts=no")

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		panic("Weather API not available")
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		panic(err)
	}

	var weather model.Weather
	er := json.Unmarshal(body, &weather)

	if er != nil {
		panic("Weather API not available")
	}

	location, current, hours := weather.Location, weather.Current, weather.Forecast.ForecastDay[0].Hour

	color.Cyan("%s, %s: %.0fC, %s\n", location.Name, location.Country, current.TempC, current.Condition.Text)

	for _, hour := range hours {

		date := time.Unix(hour.TimeEpoc, 0)

		if date.Before(time.Now()) {
			continue
		}

		message := fmt.Sprintf(
			"Time: %s - Temp: %.0fC, Humidity: %.0f , Rain: %.0f%%, Condition: %s",
			date.Format("15:05"),
			hour.TempC,
			hour.Humidity,
			hour.ChanceOfRain,
			hour.Condition.Text,
		)

		if hour.ChanceOfRain == 0 {
			color.Green(message)
		} else {
			color.Red(message)
		}
	}
}
