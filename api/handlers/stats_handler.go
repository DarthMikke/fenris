package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"millim.no/fenris/frost"
	"millim.no/fenris/stats"
)

type DailySummary struct {
	Min 	float64 `json:"min"`
	Avg 	float64 `json:"avg"`
	Max 	float64 `json:"max"`
}

func StatsHandler(f *frost.Api, w http.ResponseWriter, r *http.Request) {
	stationId := r.PathValue("stationId")
	fromYear, err := strconv.Atoi(r.PathValue("fromYear"))
	if (err != nil) {
		panic(err)
	}
	toYear, err := strconv.Atoi(r.PathValue("toYear"))
	if (err != nil) {
		panic(err)
	}

	upstreamResponse, cached, err := f.Observations(
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

	var series []stats.Measurement[float64]
	for _, v := range decodedData.Data {
		/**
		 * Multiple observations can be done at the observation time. If this is
		 * the case, average them and pass a notice to the log/console.
		 */
		switch len(v.Observations) {
		case 1:
			series = append(series, stats.Measurement[float64]{
				Timestamp: v.ReferenceTime,
				Data: float64(v.Observations[0].Value),
			})
		default:
			fmt.Printf("Multiple observations at %s\n", v.ReferenceTime)
			series = append(series, stats.Measurement[float64]{
				Timestamp: v.ReferenceTime,
				Data: stats.AverageWithAccessor(
					v.Observations,
					func(o frost.Observation) float64 {
						return float64(o.Value)
					},
				),
			})
		}
	}

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

	encoder := json.NewEncoder(w)
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
