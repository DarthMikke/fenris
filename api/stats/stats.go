package stats

import (
	"fmt"
	"time"
)

type Measurement[V any] struct {
    Timestamp string;
    Data V;
}

func Transpose[V any](matrix [][]V) (transposed [][]V, err error) {
    defer func() {
        if r := recover(); r != nil {
            transposed = nil
            if iserr, ok := r.(error); ok {
            	err = iserr
            }
        }
    }()

	rows := len(matrix)
	cols := len(matrix[0])

	for i := 0; i < cols; i++ {
		var newrow []V
		for j := 0; j < rows; j++ {
			newrow = append(newrow, matrix[j][i])
		}
		transposed = append(transposed, newrow)
	}
	err = nil
	return
}
func Wrap[V any](series []V, columns int) (rows [][]V) {
	for i := 0; i < len(series); i += columns {
		j := min(i + columns, len(series))
		rows = append(rows, series[i:j])
	}
	return rows
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
			j := 0
			for !( date1.Before(date0) ) {
				if (j > 90) {
					fmt.Printf(
						"Too large gap in time series between %s and %s.\n",
						e.Timestamp,
						series[i - 1].Timestamp,
					)
					break
				}
				bins = append(bins, *new ([]Measurement[V]))
				date0 = time.Date(date0.Year(), date0.Month() + 1, date0.Day(), 0, 0, 0, 0, date0.Location())
				j++
			}
		}
		bins[len(bins) - 1] = append(bins[len(bins) - 1], e)
	}
	return bins
}
