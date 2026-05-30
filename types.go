package geokit

// Place is the unified result of both forward and reverse geocoding.
// Forward geocoding fills Lat, Lng, LocationType, and basic address fields.
// Reverse geocoding fills all address fields alongside the coordinates.
type Place struct {
	// Coordinates
	Lat          float64
	Lng          float64
	LocationType string // ROOFTOP, RANGE_INTERPOLATED, GEOMETRIC_CENTER, APPROXIMATE

	// Address
	Line1      string // chome-ban-go + neighborhood (JP) / street number + street name (Western)
	Line2      string // building name + room number, if any
	City       string
	District   string // administrative_area_level_2, e.g. Iruma District (Japan)
	State      string // prefecture in Japan, state/province elsewhere
	PostalCode string
	Country    string

	// Shared identifiers
	FormattedAddress string
	PlaceID          string

	// Confidence signals
	PartialMatch bool     // true if Google couldn't fully resolve the input address
	Types        []string // result type tags, e.g. ["street_address", "premise"]
}
