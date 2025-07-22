package main

import (
	"log"
	"net/http"

	"github.com/lookandhqte/docker_learning/humidity"
	"github.com/lookandhqte/docker_learning/weather"
)

func main() {

	http.HandleFunc("/humidity", humidity.HumidityHandler)
	http.HandleFunc("/weather", weather.WeatherHandler)
	log.Println("Сервер влажности запущен на :1213")
	log.Fatal(http.ListenAndServe(":1213", nil))
	log.Printf("Сервер погоды запущен на :1212")
	log.Fatal(http.ListenAndServe(":1212", nil))
}
