package nominatim

type OSMResponse struct {
	PlaceID     int        `json:"place_id"`
	Licence     string     `json:"licence"`
	OsmType     string     `json:"osm_type"`
	OsmID       int        `json:"osm_id"`
	Lat         string     `json:"lat"`
	Lon         string     `json:"lon"`
	PlaceRank   float64    `json:"place_rank"`
	Category    string     `json:"category"`
	Type        string     `json:"type"`
	Importance  float64    `json:"importance"`
	AddressType string     `json:"addresstype"`
	Name        *string    `json:"name"`
	DisplayName string     `json:"display_name"`
	Address     OSMAddress `json:"address"`
	BoundingBox []string   `json:"boundingbox"`
}

type OSMAddress struct {
	HouseNumber string `json:"house_number"`
	Road        string `json:"road"`
	Suburb      string `json:"suburb"`
	City        string `json:"city"`
	County      string `json:"county"`
	State       string `json:"state"`
	Postcode    string `json:"postcode"`
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
}
