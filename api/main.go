package main

import (
	"fmt"
	"net/http"
	"os"

	"millim.no/fenris/frost"
	"millim.no/fenris/handlers"

	"github.com/joho/godotenv"
)

func indexHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "Hello, world!")
}

var frostApi *frost.Api

func main() {
	godotenv.Load()

	frostApi = &frost.Api {}
	frostApi.Setup(os.Getenv("client_id"),os.Getenv("client_secret"))

	http.HandleFunc("/api/s/{stationId}/from/{fromYear}/to/{toYear}", handlers.StatsHandler(frostApi))
	http.HandleFunc("/api/s/{stationId}", handlers.StationHandler(frostApi))
	// http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":5000", nil)
}
