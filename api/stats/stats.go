package stats

import (
	"fmt"
	"time"
)

type numeric interface {
	int32|int64|float32|float64
}

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

func Flatten3D[V any](matrix [][][]V) [][]V {
	var flattened [][]V
	for _, plane := range matrix {
		var newrow []V
		for _, row := range plane {
			for _, e := range row {
				newrow = append(newrow, e)
			}
		}
		flattened = append(flattened, newrow)
	}

	return flattened
}

// Convert a vector to a matrix by starting a new row
// every `columns` number of columns. The last row might have
// a different number of columns.
func Wrap[V any](series []V, columns int) (rows [][]V) {
	for i := 0; i < len(series); i += columns {
		j := min(i + columns, len(series))
		rows = append(rows, series[i:j])
	}
	return rows
}

// Convert time series data to a matrix with one column per calendar month.
func Periodise[V any](series []Measurement[V], period string) [][]Measurement[V] {
	var bins [][]Measurement[V]
	for i, e := range series {
		if (i == 0) {
			bins = append(bins, *new ([]Measurement[V]))
			bins[len(bins) - 1] = append(bins[len(bins) - 1], e)
			continue
		}

		samebin := false
		if (period == "P1M") {
			samebin = e.Timestamp[:7] == series[i - 1].Timestamp[:7]
		} else if period == "P1D" {
			samebin = e.Timestamp[:10] == series[i - 1].Timestamp[:10]
		} else {
			panic("Invalid time period.")
		}

		if !samebin {
			date1, err := time.Parse(time.RFC3339, e.Timestamp)
			if (err != nil) {
				panic(err)
			}
			date0, err := time.Parse(time.RFC3339, series[i - 1].Timestamp)
			if (err != nil) {
				panic(err)
			}

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

func Reduce[V any, U any](vector [][]V, reducer func([]V) U) (reduced []U) {
	for _, row := range vector {
		reduced = append(reduced, reducer(row))
	}
	return
}

func Average[V numeric](numbers []V) float64 {
	var sum V
	sum = 0
	for _, e := range numbers {
		sum += e
	}
	return float64(sum)/float64(len(numbers));
}

func AnnotatedMax[V numeric](numbers []Measurement[V]) Measurement[V] {
	max := numbers[0]
	for _, e := range numbers {
		if (e.Data > max.Data) {
			max = e
		}
	}
	return max;
}

func AnnotatedMin[V numeric](numbers []Measurement[V]) Measurement[V] {
	min := numbers[0]
	for _, e := range numbers {
		if (e.Data < min.Data) {
			min = e
		}
	}
	return min;
}

// Return average of a single field in a slice of structs of type V.
func AverageWithAccessor[V any](series []V, accessor func (V) float64) float64 {
	var sum float64
	sum = 0
	for _, v := range series {
		sum += accessor(v)
	}
	return sum/float64(len(series))
}
