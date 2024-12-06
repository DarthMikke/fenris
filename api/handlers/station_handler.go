package handlers

import (
	"encoding/json"
	"net/http"
	"millim.no/fenris/frost"
)

func StationHandler(f *frost.Api) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		stationId := r.PathValue("stationId")

		/*
		response := responses.StationResponse{
			StationId: stationId,
		}
		*/

		upstreamResponse, cached, err := f.Sources([]string{stationId})

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
}
