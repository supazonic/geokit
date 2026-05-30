package geokit

import "context"

// Geocoder is the interface that geocoding backends must implement.
// Importers pass their implementation to NewGeocoder to obtain a Client.
type Geocoder interface {
	// Geocode converts a human-readable address into one or more Locations.
	Geocode(ctx context.Context, address string) ([]Location, error)

	// ReverseGeocode converts coordinates into one or more Addresses.
	ReverseGeocode(ctx context.Context, lat, lng float64) ([]Address, error)
}
