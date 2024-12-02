package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"millim.no/fenris/frost"
	"millim.no/fenris/stats"

	"github.com/joho/godotenv"
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

func statsHandler(w http.ResponseWriter, r *http.Request) {
	stationId := r.PathValue("stationId")
	fromYear, err := strconv.Atoi(r.PathValue("fromYear"))
	if (err != nil) {
		panic(err)
	}
	toYear, err := strconv.Atoi(r.PathValue("toYear"))
	if (err != nil) {
		panic(err)
	}

	upstreamResponse, cached, err := frostApi.Observations(
		[]string{stationId},
		fmt.Sprintf("%d-01-01/%d-01-01", fromYear, toYear + 1),
		[]string{"air_temperature"},
	)
	if (err != nil) {
		panic(err)
	}
	if (cached) {
		w.Header().Add("X-Cache-Hit", "1")
	}

	var decodedData frost.ObservationResponse
	decoder := json.NewDecoder(strings.NewReader(*upstreamResponse))
	decoder.Decode(&decodedData)

	// Sort into bins with different YYYY-MM component in
	// the referenceTime field.
	var series []stats.Measurement[[]frost.Observation]
	for _, v := range decodedData.Data {
		series = append(series, stats.Measurement[[]frost.Observation]{
			Timestamp: v.ReferenceTime,
			Data: v.Observations,
		})
	}
	var bins [][]stats.Measurement[[]frost.Observation]
	bins = stats.Periodise(series)
	var wrapped [][][]stats.Measurement[[]frost.Observation]
	wrapped = stats.Wrap(bins, 12)

	encoder := json.NewEncoder(w)
	encoder.Encode(wrapped)
	// fmt.Fprintf(w, *upstreamResponse);
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
