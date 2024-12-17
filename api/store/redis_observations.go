package store

import (
	"fmt"

	"github.com/redis/go-redis/v9"
	"millim.no/fenris/frost"
	"millim.no/fenris/stats"
)

type RedisObservationsStore struct {
	f 	*frost.Api
	r 	*redis.Client
}

func NewRedisObservationsStore(f *frost.Api, r *redis.Client) *RedisObservationsStore {
	return &RedisObservationsStore{
		f: f,
		r: r,
	}
}

func (obss RedisObservationsStore)GetObservations (station string, fromYear int, toYear int)	([]stats.Measurement[float64], error) {
	upstreamResponse, _, err := obss.f.Observations(
		[]string{station},
		fmt.Sprintf("%d-01-01/%d-01-01", fromYear, toYear + 1),
		[]string{"air_temperature"},
	)
	if (err != nil) {
		panic(err)
	}


	var series []stats.Measurement[float64]
	for _, v := range upstreamResponse.Data {
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

	return series, nil
}
