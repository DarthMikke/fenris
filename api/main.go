package main

import (
	"fmt"
	"encoding/json"
	"net/http"
  "millim.no/fenris/frost"
  // "millim.no/fenris/responses"
  "github.com/joho/godotenv"
  "os"
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

  response, _ := frostApi.Sources([]string{stationId})

  encoder := json.NewEncoder(w)
  encoder.Encode(response.Data[0])
}

var frostApi *frost.Api

func main() {
  godotenv.Load()

  frostApi = &frost.Api {}
  frostApi.Setup(os.Getenv("client_id"),os.Getenv("client_secret"))

  http.HandleFunc("/api/s/{stationId}/from/{fromYear}/to/{toYear}", stationHandler)
  http.HandleFunc("/api/s/{stationId}", stationHandler)
  http.HandleFunc("/", indexHandler)
  http.ListenAndServe(":5000", nil)
}
