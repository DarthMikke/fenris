package main

import (
	"fmt"
	"net/http"
	"os"

	"millim.no/fenris/frost"
	"millim.no/fenris/handlers"
	"millim.no/fenris/store"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func indexHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "Hello, world!")
}

var frostApi *frost.Api

type SrvWrapper struct {
	frostApi 	*frost.Api
	obsStore	store.ObservationsStore
}

func (s SrvWrapper) HandleFunc(
	pattern string,
	handler func(*frost.Api, *store.ObservationsStore, http.ResponseWriter, *http.Request),
) {

	http.HandleFunc(pattern, func (w http.ResponseWriter, r *http.Request) {
		handler(s.frostApi, &s.obsStore, w, r)
	})
}

func main() {
	godotenv.Load()

	// Ctx = context.Background()
	redisClient := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		Password:	"",
		DB:	0,
	})

	frostApi = &frost.Api {}
	frostApi.Setup(os.Getenv("client_id"),os.Getenv("client_secret"))
	obsStore := store.NewRedisObservationsStore(frostApi, redisClient)

	server := SrvWrapper{
		frostApi: frostApi,
		obsStore: *obsStore,
	}

	server.HandleFunc("/api/s/{stationId}/from/{fromYear}/to/{toYear}", handlers.StatsHandler)
	server.HandleFunc("/api/s/{stationId}", handlers.StationHandler)
	// http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":5000", nil)
}
