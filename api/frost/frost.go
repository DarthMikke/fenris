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

type SourceObject struct {
	Id	string `json:"id"`
	Name	string `json:"name"`
    ShortName	string `json:"shortName"`
    Country	string `json:"country"`
    CountryCode	string  `json:"countryCode"`
    WmoId int `json:"wmoId"`
	ValidFrom 	string  `json:"validFrom"`
	County	string  `json:"county"`
	CountyId	int  `json:"countyId"`
	Municipality	string  `json:"municipality"`
	MunicipalityId	int  `json:"municipalityId"`
	/*
      "@type": "SensorSystem",
      "geometry": {
        "@type": "Point",
        "coordinates": [
          10.4762,
          61.0917
        ],
        "nearest": false
      },
      "masl": 240,
      "ontologyId": 0,
      "stationHolders": [
        "MET.NO"
      ],
      "externalIds": [
        "0-20000-0-01378",
        "10.249.1.155"
      ],
      "wigosId": "0-20000-0-01378"
      */
}

type SourcesResponse struct {
	Context	string 	`json:"@context"`
	Type  	string 	`json:"@type"`
	ApiVersion 	string 	`json:"apiVersion"`
	Data	[]SourceObject 	`json:"data"`
	/*
  "apiVersion": "v0",
  "license": "https://creativecommons.org/licenses/by/3.0/no/",
  "createdAt": "2024-11-24T19:55:42Z",
  "queryTime": 1.37,
  "currentItemCount": 1,
  "itemsPerPage": 1,
  "offset": 0,
  "totalItemCount": 1,
  "currentLink": "https://frost.met.no/sources/v0.jsonld?ids=SN12680",
  */
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
    	Password: 	"",
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

		a.redis.Set(a.Ctx, url, &s, 3_600*1_000_000_000)
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
