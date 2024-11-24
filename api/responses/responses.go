package responses

type HalResponse struct {
	Links map[string](map[string]string)	`json:"_links"`
}

type StationResponse struct {
	StationId	string	`json:"stationId"`
	// AvailableFrom	string	`json:"availableFrom"`
	// AvailableTo	string	`json:"availableTo"`
	// Municipality	string	`json:"municipality"`
	// County	string	`json:"county"`
}
