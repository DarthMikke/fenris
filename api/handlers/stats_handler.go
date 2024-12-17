package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"millim.no/fenris/frost"
	"millim.no/fenris/stats"
	"millim.no/fenris/store"
)

type DailySummary struct {
	Min 	float64 `json:"min"`
	Avg 	float64 `json:"avg"`
	Max 	float64 `json:"max"`
}

func StatsHandler(_ *frost.Api, s *store.ObservationsStore, w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	defer func() {
		if r := recover(); r != nil {
			if iserr, ok := r.(error); ok {
				encoder.Encode(iserr.Error())
			} else {
				encoder.Encode(r)
			}
		}
	}()

	stationId := r.PathValue("stationId")
	fromYear, err := strconv.Atoi(r.PathValue("fromYear"))
	if (err != nil) {
		panic(err)
	}
	toYear, err := strconv.Atoi(r.PathValue("toYear"))
	if (err != nil) {
		panic(err)
	}

	series, err := (*s).GetObservations(stationId, fromYear, toYear)
	if err != nil {
		if iserr, ok := err.(store.ObservationsStoreError); ok {
			switch iserr.Details {
			case store.NoData:
				// Fetch data
			case store.FetchingData:
				// Respond accordingly
			case store.OutsideOfRange:
				// Respond accordingly
			default:
				panic(err)
			}
		} else {
			panic(err)
		}
	}

	w.Header().Add("X-Cache-Hit", "1")

	// Sort into bins with different YYYY-MM-DD component in
	// the referenceTime field.
	var binsByDay []stats.Measurement[DailySummary]
	for _, v := range stats.Periodise(series, "P1D") {
		min := v[0]
		max := v[0]
		for _, w := range v {
			if w.Data < min.Data {
				min = w
			}
			if w.Data > max.Data {
				max = w
			}
		}
		newSummary := stats.Measurement[DailySummary] {
			Timestamp:	v[0].Timestamp,
			Data:	DailySummary {
				Min: min.Data,
				Max: max.Data,
				Avg: stats.AverageWithAccessor(
					v, func(o stats.Measurement[float64]) float64 {
						return o.Data
					},
				),
			},
		}
		fmt.Println(v, newSummary)
		binsByDay = append(binsByDay, newSummary)
	}

	binsByMonth := stats.Periodise(binsByDay, "P1M")
	threed := stats.Wrap(binsByMonth, 12)
	threed, err = stats.Transpose(threed)
	if (err != nil) {
		panic(err)
	}
	flattened := stats.Flatten3D(threed)

	averages := stats.Reduce(
		flattened,
		func(obs []stats.Measurement[DailySummary]) float64 {
			var numbers []float64
			for _, e := range obs {
				numbers = append(numbers, float64(e.Data.Avg))
			}
			return stats.Average(numbers)
		},
	)
	avgmaxs := stats.Reduce(
		flattened,
		func(obs []stats.Measurement[DailySummary]) float64 {
			var numbers []float64
			for _, e := range obs {
				numbers = append(numbers, float64(e.Data.Max))
			}
			return stats.Average(numbers)
		},
	)
	avgmins := stats.Reduce(
		flattened,
		func(obs []stats.Measurement[DailySummary]) float64 {
			var numbers []float64
			for _, e := range obs {
				numbers = append(numbers, float64(e.Data.Min))
			}
			return stats.Average(numbers)
		},
	)
	/*
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
	*/

	encoder.Encode(map[string]any{
		"average": averages,
		"avgmax": avgmaxs,
		"avgmin": avgmins,
		/*
		"max": maxs,
		"min": mins,
		*/
	})
	// fmt.Fprintf(w, *upstreamResponse);
}
