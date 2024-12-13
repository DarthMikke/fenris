package store

import (
	"millim.no/fenris/stats"
)

type ObservationsStore interface {
	// Get observations for a given time period
	GetObservations(string, int, int)	([]stats.Measurement[float64], error)
}

	// FetchObservations(station string, fromYear int, toYear int)	([]stats.Measurement[float64], error)
