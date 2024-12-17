package frost

import (
	"bytes"
	"context"
	"fmt"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/redis/go-redis/v9"
)

type HttpError struct {
	Code	int
	Message	string
}

func (e HttpError) Error() string {
	return fmt.Sprintf("%d %s", e.Code, e.Message)
}

type Api struct {
	Url string
	clientSecret	string
	clientId	string

	httpClient	http.Client
	redis 	*redis.Client
	Ctx	context.Context
}

func (a *Api) Setup (ClientId string, ClientSecret string) {
	a.Url = "frost.met.no"
	a.clientId =  ClientId
	a.clientSecret = ClientSecret

	a.httpClient = http.Client{}

	a.Ctx = context.Background()
	a.redis = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		Password:	"",
		DB:	0,
	})
}

func (a *Api) get (url string) (*string, bool, error) {
	value, err := a.redis.Get(a.Ctx, url).Result()
	if (err == redis.Nil) {
		resp, err := a._get(url)
		fmt.Println("Caching:", url)

		if (err != nil) {
			return nil, false, err
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		s := buf.String()

		if resp.StatusCode == 200 {
			a.redis.Set(a.Ctx, url, &s, 3_600*1_000_000_000)
		} else if resp.StatusCode >= 400 {
			return &s, false, HttpError{
				Code: resp.StatusCode,
				Message: s,
			}
		}

		return &s, false, nil
	} else if (err != nil) {
		return nil, false, err
	} else {
		fmt.Println("Cache hit:", url)
		return &value, true, nil
	}
}

func (a *Api) _get (url string) (*http.Response, error) {
	request, _ := http.NewRequest("GET", url, nil)
	request.SetBasicAuth(a.clientId, a.clientSecret)
	return a.httpClient.Do(request)
}

func (a *Api) Sources (ids []string) (*SourcesResponse, bool, error) {
	response, cached, err := a.get("https://frost.met.no/sources/v0.jsonld?ids=" + strings.Join(ids, ","))

	if (err != nil) {
		return nil, false, err
	}

	var decoded SourcesResponse;
	reader := strings.NewReader(*response)
	decoder := json.NewDecoder(reader)
	decoder.Decode(&decoded)

	return &decoded, cached, nil
}

/*
Implements Frost API's `observations` endpoint.

https://frost.met.no/api.html#!/observations/observations
 */
func (a *Api) Observations (sources []string, referenceTime string, elements []string) (*ObservationResponse, bool, error) {
	queryArgs := []string{
		"sources=" + strings.Join(sources, ","),
		"referencetime=" + referenceTime,
		"elements=" + strings.Join(elements, ","),
	}

	response, cached, err := a.get("https://frost.met.no/observations/v0.jsonld?" + strings.Join(queryArgs, "&"))

	if (err != nil) {
		return nil, false, err
	}

	var decodedData ObservationResponse
	decoder := json.NewDecoder(strings.NewReader(*response))
	err = decoder.Decode(&decodedData)

	if (err != nil) {
		return nil, false, err
	}

	return &decodedData, cached, nil
}
