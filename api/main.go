package main

import (
	"fmt"
	"encoding/json"
	"net/http"
)

func indexHandler(w http.ResponseWriter, _ *http.Request) {
  fmt.Fprintf(w, "Hello, world!")
}

func stationHandler(w http.ResponseWriter, r *http.Request) {
  encoder := json.NewEncoder(w)
  encoder.Encode(r.PathValue("stationId"))
}

func main() {
  http.HandleFunc("/api/s/{stationId}", stationHandler)
  http.HandleFunc("/api/s/{stationId}/from/{fromYear}/to/{toYear}", stationHandler)
  http.HandleFunc("/", indexHandler)
  http.ListenAndServe(":5000", nil)
}
