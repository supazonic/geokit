package geokit

import "context"

// Client wraps a Geocoder implementation and is the primary entry point for
// geocoding operations.
type Client struct {
	backend Geocoder
}

// NewGeocoder returns a Client backed by the provided Geocoder implementation.
// Importers supply their own backend (e.g. a Google Maps or Mapbox adapter)
// that satisfies the Geocoder interface.
func NewGeocoder(g Geocoder) *Client {
	return &Client{backend: g}
}

// Geocode converts an address string into a slice of Locations.
func (c *Client) Geocode(ctx context.Context, address string) ([]Location, error) {
	return c.backend.Geocode(ctx, address)
}

// ReverseGeocode converts coordinates into a slice of Addresses.
func (c *Client) ReverseGeocode(ctx context.Context, lat, lng float64) ([]Address, error) {
	return c.backend.ReverseGeocode(ctx, lat, lng)
}
