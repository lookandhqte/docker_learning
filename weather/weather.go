package weather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func GetGeocoding(city, apiKey string) (*GeocodingAPI, error) {
	urlForGeo := fmt.Sprintf("http://api.openweathermap.org/geo/1.0/direct?q=%s&appid=%s", city, apiKey)

	resp, err := http.Get(urlForGeo)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var dataGeo GeocodingAPI
	err = json.Unmarshal(body, &dataGeo)
	if err != nil {
		return nil, err
	}

	if len(dataGeo) == 0 {
		return nil, fmt.Errorf("город не найден")
	}

	return &dataGeo, nil
}

func GetWeather(lat, lon float64, apiKey string, amount int) (*WeatherResponse, error) {
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/forecast?lat=%.6f&lon=%.6f&appid=%s", lat, lon, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data WeatherResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	if amount > 0 && len(data.List) > amount {
		data.List = data.List[:amount]
	}

	return &data, nil
}

func ReturnHTML(city string, weather *WeatherResponse) string {
	html := `
    <html>
    <head>
        <title>Прогноз погоды для ` + city + `</title>
        <style>
            .forecast-day {
                border: 1px solid #ddd;
                border-radius: 8px;
                padding: 15px;
                margin-bottom: 20px;
                background-color: #f9f9f9;
            }
            .forecast-header {
                font-size: 1.2em;
                font-weight: bold;
                margin-bottom: 10px;
                color: #333;
            }
            .forecast-item {
                margin-bottom: 8px;
            }
            img.weather-icon {
                vertical-align: middle;
                margin-right: 10px;
            }
        </style>
    </head>
    <body>
        <h1>Прогноз погоды для ` + weather.City.Name + `</h1>
    `

	currentDay := ""

	for i := 0; i < len(weather.List); i++ {
		forecastTime := time.Unix(weather.List[i].Dt, 0)
		dayStr := forecastTime.Format("02.01.2006")

		// Если день сменился, добавляем заголовок
		if dayStr != currentDay {
			if currentDay != "" {
				html += `</div>` // Закрываем предыдущий день
			}
			html += fmt.Sprintf(`
            <div class="forecast-day">
                <div class="forecast-header">%s</div>
            `, forecastTime.Format("Monday, 02 January 2006"))
			currentDay = dayStr
		}

		// Добавляем прогноз на конкретное время
		html += fmt.Sprintf(`
            <div class="forecast-item">
                <strong>%s</strong><br>
                <img class="weather-icon" src="https://openweathermap.org/img/wn/%s@2x.png" alt="Weather icon">
                Температура: %.1f°C <br>
                Влажность: %d%%<br>
                %s<br>
                Ветер: %.1f м/с
            </div>
        `,
			forecastTime.Format("15:04"),
			weather.List[i].Weather[0].Icon,
			weather.List[i].Main.Temp-273.15,
			weather.List[i].Main.Humidity,
			weather.List[i].Weather[0].Description,
			weather.List[i].Wind.Speed,
		)
	}

	html += `
        </div>
    </body>
    </html>
    `
	return html
}

func WeatherHandler(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if city == "" {
		http.Error(w, "Укажите город: /weather?city=Москва", http.StatusBadRequest)
		return
	}

	apiKey := "c6b1ff1b5bdba6853fba34421f6149c3"
	if apiKey == "" {
		http.Error(w, "API ключ не настроен", http.StatusInternalServerError)
		return
	}

	// Получаем координаты города
	geocoding, err := GetGeocoding(city, apiKey)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка геокодинга: %v", err), http.StatusInternalServerError)
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, fmt.Sprintf("Ошибка id: %s", "айди пустое"), http.StatusInternalServerError)
		return
	}
	intId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка STRid->int: %v", err), http.StatusInternalServerError)
		return
	}

	firstLocation := (*geocoding)[0]
	lat, lon := firstLocation.Lat, firstLocation.Lon

	weather, err := GetWeather(lat, lon, apiKey, intId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка запроса погоды: %v", err), http.StatusInternalServerError)
		return
	}

	if len(weather.List) == 0 {
		http.Error(w, "Нет данных о погоде", http.StatusInternalServerError)
		return
	}

	html := ReturnHTML(city, weather)

	// Группируем прогнозы по дням

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}
