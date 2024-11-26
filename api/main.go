package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"millim.no/fenris/frost"
	// "millim.no/fenris/responses"
	"github.com/joho/godotenv"
	"os"
)

func indexHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "Hello, world!")
}

func stationHandler(w http.ResponseWriter, r *http.Request) {
	stationId := r.PathValue("stationId")

	/*
	response := responses.StationResponse{
		StationId: stationId,
	}
	*/

	upstreamResponse, cached, err := frostApi.Sources([]string{stationId})

	if (err != nil) {
		panic(err)
	}
	if (cached) {
		w.Header().Add("X-Cache-Hit", "1")
	}

	/*
	response := new (map[string]any)
	response["id"] = (upstreamResponse.Data[0]).Id
	*/

	encoder := json.NewEncoder(w)
	encoder.Encode(upstreamResponse.Data[0])
}

func statsHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "{}");
}

var frostApi *frost.Api

func main() {
	godotenv.Load()

	frostApi = &frost.Api {}
	frostApi.Setup(os.Getenv("client_id"),os.Getenv("client_secret"))

	http.HandleFunc("/api/s/{stationId}/from/{fromYear}/to/{toYear}", statsHandler)
	http.HandleFunc("/api/s/{stationId}", stationHandler)
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":5000", nil)
}
