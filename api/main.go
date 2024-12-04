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
	var threed [][][]stats.Measurement[[]frost.Observation]
	threed = stats.Wrap(bins, 12)
	threed, err = stats.Transpose(threed)
	if (err != nil) {
		panic(err)
	}
	var flattened [][]stats.Measurement[[]frost.Observation]
	flattened = stats.Flatten3D(threed)

	averages := stats.Reduce(
		flattened,
		func(obs []stats.Measurement[[]frost.Observation]) float64 {
			var numbers []float64
			for _, e := range obs {
				numbers = append(numbers, float64(e.Data[0].Value))
			}
			return stats.Average(numbers)
		},
	)
	maxs := stats.Reduce(
		flattened,
		func(obs []stats.Measurement[[]frost.Observation]) stats.Measurement[float64] {
			var numbers []stats.Measurement[float64]
			for _, e := range obs {
				mapped := stats.Measurement[float64]{
					Timestamp: e.Timestamp,
					Data: float64(e.Data[0].Value),
				}
				numbers = append(numbers, mapped)
			}
			return stats.AnnotatedMax(numbers)
		},
	)
	mins := stats.Reduce(
		flattened,
		func(obs []stats.Measurement[[]frost.Observation]) stats.Measurement[float64] {
			var numbers []stats.Measurement[float64]
			for _, e := range obs {
				mapped := stats.Measurement[float64]{
					Timestamp: e.Timestamp,
					Data: float64(e.Data[0].Value),
				}
				numbers = append(numbers, mapped)
			}
			return stats.AnnotatedMin(numbers)
		},
	)

	encoder := json.NewEncoder(w)
	encoder.Encode(map[string]any{
		"average": averages,
		"max": maxs,
		"min": mins,
	})
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
