package geokit

// Location is the result of a forward geocoding lookup.
type Location struct {
	Lat        float64
	Lng        float64
	FormattedAddress string
	// PlaceID is a provider-specific identifier for the place, if available.
	PlaceID string
}

// Address is the result of a reverse geocoding lookup.
type Address struct {
	FormattedAddress string
	Street           string
	City             string
	State            string
	PostalCode       string
	Country          string
	// PlaceID is a provider-specific identifier for the place, if available.
	PlaceID string
}
