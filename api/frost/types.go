package frost

type JsonLdResponse struct {
	Context	string 	`json:"@context"`
	Type  	string 	`json:"@type"`
}

type ResponseBase struct {
	JsonLdResponse
	ApiVersion string `json:"apiVersion"` // The version of the API that generated this response. ,
	License string `json:"license"` // The license that applies to this content. ,
	CreatedAt string `json:"createdAt"` // The time at which this document was created (RFC 3339). ,
	QueryTime float64 `json:"queryTime"` // The time, in seconds, that this document took to generate. ,
	CurrentLink string `json:"currentLink"` // The current link indicates the URI that was used to generate the current API response ,
}

type Source struct {
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
	ResponseBase
	Data	[]Source 	`json:"data"`
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


type ObservationResponse struct {
ResponseBase
CurrentItemCount int `json:"currentItemCount"` // The current number of items in this result set. ,
ItemsPerPage int `json:"itemsPerPage"` // The maximum number of items in a result set. ,
Offset int `json:"offset"` // The offset of the first item in the result set. The Frost API uses a zero-base index. ,
TotalItemCount int `json:"totalItemCount"` // The total number of items in this specific result set. ,
NextLink string `json:"nextLink"` // The next link indicates how more data can be retrieved. It points to the URI to load the next set of data. ,
PreviousLink string `json:"previousLink"` // The previous link indicates how more data can be retrieved. It points to the URI to load the previous set of data. ,
Data []ObservationsAtRefTime `json:"data"` // Container for all the data from the response.
}
type ObservationsAtRefTime struct {
SourceId string `json:"sourceId"` // The sourceId at which values were observed. ,
Geometry Point `json:"geometry"` // The latitude and longitude at which values were observed (if known). ,
ReferenceTime string `json:"referenceTime"` // The base time at which values were observed (any timeOffset must be added). ,
Observations []Observation `json:"observations"` // The observed values. This is a map of the form [ElementId (as a String), Observation]
}
type Point struct {
Type string `json:"@type"` // The type of the geometry object ,
Coordinates []float32 `json:"coordinates"` // Coordinates of the geometry object
}
type Observation struct {
ElementId string `json:"elementId"` // The ID of the element being observed. ,
Value float32 `json:"value"` // The observed value (either a number or a UTC datetime of the format YYYY-MM-DD hh:mm:ss.sss). ,
OrigValue float32 `json:"origValue"` // The original observed value (either a number or a UTC datetime of the format YYYY-MM-DD hh:mm:ss.sss). ,
Unit string `json:"unit"` // The unit of measurement of the observed value. ,
CodeTable string `json:"codeTable"` // If the unit is a code, the codetable that describes the codes used. ,
Level Level `json:"level"` // The vertical level at which the value was observed (if known). ,
TimeOffset string `json:"timeOffset"` // The offset from referenceTime at which the observation applies. ,
TimeResolution string `json:"timeResolution"` // The time between consecutive observations in the time series to which the observation belongs. ,
// TimeSeriesId object `json:"timeSeriesId"` // The internal ID of the time series to which the observation belongs. ,
PerformanceCategory string `json:"performanceCategory"` // The performance category of the source when the value was observed. ,
ExposureCategory string `json:"exposureCategory"` // The exposure category of the source when the value was observed. ,
// QualityCode object `json:"qualityCode"` // The quality control flag of the observed value. ,
ControlInfo string `json:"controlInfo"` // The control info of the observed value. ,
// DataVersion object `json:"dataVersion"` // The data version of the data value, if one exists (**Note: Currently not available for any observation data).
}
type Level struct {
LevelType string `json:"levelType"` // The reference type of the level value. ,
Unit string `json:"unit"` // The unit of measurement of the level value. ,
Value int `json:"value"` // The level value.
}

type ErrorResponse struct {
	ResponseBase
	Error 	ResponseError `json:"error"`
}

type ResponseError struct {
	Code	int		`json:"code"`
	Message	string	`json:"message"`
	Reason	string	`json:"reason"`
}
