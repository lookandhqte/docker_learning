package humidity

import (
	"fmt"
	"gelzh/weather" // Импорт своего модуля
	"log"
	"net/http"
	"strconv"
	"time"
)

func humidityHandler(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	daysStr := r.URL.Query().Get("id")

	days := 5 // По умолчанию 5 дней
	if daysStr != "" {
		var err error
		days, err = strconv.Atoi(daysStr)
		if err != nil {
			http.Error(w, "id должен быть числом", http.StatusBadRequest)
			return
		}
	}
	apiKey := "c6b1ff1b5bdba6853fba34421f6149c3"

	// Используем функции из модуля meow
	geo, err := weather.GetGeocoding(city, apiKey)
	if err != nil {
		http.Error(w, fmt.Sprintf("Геокодинг ошибка: %v", err), http.StatusInternalServerError)
		return
	}
	weather, err := weather.GetWeather(geo.Lat, geo.Lon, apiKey, days)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка погоды: %v", err), http.StatusInternalServerError)
		return
	}

	// Генерация HTML с акцентом на влажность
	html := fmt.Sprintf(`<h1>Влажность в %s</h1>`, city)
	for _, forecast := range weather.List {
		date := time.Unix(forecast.Dt, 0).Format("02.01")
		html += fmt.Sprintf(`
			<div>%s: %d%% влажности</div>
		`, date, forecast.Main.Humidity)
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func main() {
	http.HandleFunc("/humidity", humidityHandler)
	log.Println("Сервер влажности запущен на :1213")
	log.Fatal(http.ListenAndServe(":1213", nil))
}
