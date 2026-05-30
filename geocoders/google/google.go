package google

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/supazonic/geokit"
)

const baseURL = "https://maps.googleapis.com/maps/api/geocode/json"

type google struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

// NewGoogleGeocoder returns a geokit.Geocoder backed by the Google Geocoding API.
func NewGoogleGeocoder(apiKey string) *google {
	return &google{
		apiKey:     apiKey,
		httpClient: &http.Client{},
		baseURL:    baseURL,
	}
}

// --- internal response types ---

type apiResponse struct {
	Status       string      `json:"status"`
	ErrorMessage string      `json:"error_message"`
	Results      []apiResult `json:"results"`
}

type apiResult struct {
	FormattedAddress  string                `json:"formatted_address"`
	PlaceID           string                `json:"place_id"`
	Geometry          apiGeometry           `json:"geometry"`
	AddressComponents []apiAddressComponent `json:"address_components"`
	Types             []string              `json:"types"`
	PartialMatch      bool                  `json:"partial_match"`
}

type apiGeometry struct {
	Location     apiLatLng `json:"location"`
	LocationType string    `json:"location_type"`
	Viewport     apiBounds `json:"viewport"`
	Bounds       apiBounds `json:"bounds"`
}

type apiLatLng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type apiBounds struct {
	Northeast apiLatLng `json:"northeast"`
	Southwest apiLatLng `json:"southwest"`
}

type apiAddressComponent struct {
	LongName  string   `json:"long_name"`
	ShortName string   `json:"short_name"`
	Types     []string `json:"types"`
}

// --- shared HTTP helper ---

func (g *google) fetch(ctx context.Context, params url.Values) (*apiResponse, error) {
	params.Set("key", g.apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, g.baseURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	switch apiResp.Status {
	case "OK", "ZERO_RESULTS":
		return &apiResp, nil
	default:
		msg := apiResp.Status
		if apiResp.ErrorMessage != "" {
			msg += ": " + apiResp.ErrorMessage
		}
		return nil, fmt.Errorf("google geocoder: %s", msg)
	}
}

// --- geokit.Geocoder implementation ---

// Geocode converts an address string into Locations via the Google Geocoding API.
func (g *google) Geocode(ctx context.Context, address string) ([]geokit.Location, error) {
	params := url.Values{"address": {address}}
	apiResp, err := g.fetch(ctx, params)
	if err != nil {
		return nil, err
	}

	locations := make([]geokit.Location, len(apiResp.Results))
	for i, r := range apiResp.Results {
		locations[i] = geokit.Location{
			Lat:              r.Geometry.Location.Lat,
			Lng:              r.Geometry.Location.Lng,
			FormattedAddress: r.FormattedAddress,
			PlaceID:          r.PlaceID,
		}
	}
	return locations, nil
}

// ReverseGeocode converts coordinates into Addresses via the Google Geocoding API.
func (g *google) ReverseGeocode(ctx context.Context, lat, lng float64) ([]geokit.Address, error) {
	latlng := strconv.FormatFloat(lat, 'f', -1, 64) + "," + strconv.FormatFloat(lng, 'f', -1, 64)
	params := url.Values{"latlng": {latlng}}
	apiResp, err := g.fetch(ctx, params)
	if err != nil {
		return nil, err
	}

	addresses := make([]geokit.Address, len(apiResp.Results))
	for i, r := range apiResp.Results {
		addresses[i] = toAddress(r)
	}
	return addresses, nil
}

// toAddress maps a Google API result to a geokit.Address by extracting
// typed address_components.
func toAddress(r apiResult) geokit.Address {
	a := geokit.Address{
		FormattedAddress: r.FormattedAddress,
		PlaceID:          r.PlaceID,
	}

	var streetNumber, route string
	for _, c := range r.AddressComponents {
		switch {
		case hasType(c.Types, "street_number"):
			streetNumber = c.LongName
		case hasType(c.Types, "route"):
			route = c.LongName
		case hasType(c.Types, "locality"):
			a.City = c.LongName
		case hasType(c.Types, "administrative_area_level_1"):
			a.State = c.LongName
		case hasType(c.Types, "postal_code"):
			a.PostalCode = c.LongName
		case hasType(c.Types, "country"):
			a.Country = c.LongName
		}
	}

	parts := make([]string, 0, 2)
	if streetNumber != "" {
		parts = append(parts, streetNumber)
	}
	if route != "" {
		parts = append(parts, route)
	}
	a.Street = strings.Join(parts, " ")

	return a
}

func hasType(types []string, target string) bool {
	for _, t := range types {
		if t == target {
			return true
		}
	}
	return false
}
