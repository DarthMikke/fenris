package stats

import (
	"time"
)

type Measurement[V any] struct {
    Timestamp string;
    Data V;
}

func Periodise[V any](series []Measurement[V]) [][]Measurement[V] {
	var bins [][]Measurement[V]
	for i, e := range series {
		if (i == 0) {
			bins = append(bins, *new ([]Measurement[V]))
			bins[len(bins) - 1] = append(bins[len(bins) - 1], e)
			continue
		}
		date1, err := time.Parse(time.RFC3339, e.Timestamp)
		if (err != nil) {
			panic(err)
		}
		date0, err := time.Parse(time.RFC3339, series[i - 1].Timestamp)
		if (err != nil) {
			panic(err)
		}
		if !(date1.Year() == date0.Year() && date0.Month() == date1.Month()) {
			bins = append(bins, *new ([]Measurement[V]))
		}
		bins[len(bins) - 1] = append(bins[len(bins) - 1], e)
	}
	return bins
}
