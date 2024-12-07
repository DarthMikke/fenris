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

type SrvWrapper struct {
	frostApi 	*frost.Api
}

func (s SrvWrapper) HandleFunc(
	pattern string,
	handler func(*frost.Api, http.ResponseWriter, *http.Request),
) {

	http.HandleFunc(pattern, func (w http.ResponseWriter, r *http.Request) {
		handler(s.frostApi, w, r)
	})
}

func main() {
	godotenv.Load()

	frostApi = &frost.Api {}
	frostApi.Setup(os.Getenv("client_id"),os.Getenv("client_secret"))

	server := SrvWrapper{frostApi: frostApi}

	server.HandleFunc("/api/s/{stationId}/from/{fromYear}/to/{toYear}", handlers.StatsHandler)
	server.HandleFunc("/api/s/{stationId}", handlers.StationHandler)
	// http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":5000", nil)
}
