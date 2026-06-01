package geokit

import "context"

// Geocoder is the interface that geocoding backends must implement.
// Importers pass their implementation to NewGeocoder to obtain a Client.
type Geocoder interface {
	// Geocode converts a human-readable address into one or more Places.
	Geocode(ctx context.Context, address string) ([]Place, error)

	// ReverseGeocode converts coordinates into one or more Places.
	ReverseGeocode(ctx context.Context, lat, lng float64) ([]Place, error)
}
