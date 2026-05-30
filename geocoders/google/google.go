package google

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

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
	Results      []ApiResult `json:"results"`
}

type ApiResult struct {
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

type locationType string

func (lt locationType) String() string {
	return string(lt)
}

const (
	//	Exact coordinates of the specific building/address
	LocationType_ROOFTOP locationType = "ROOFTOP"
	// Estimated point between two known rooftop addresses on a street
	LocationType_RANGE_INTERPOLATED locationType = "RANGE_INTERPOLATED"
	// Center point of a region (road, neighborhood, city, etc.)
	LocationType_GEOMETRIC_CENTER locationType = "GEOMETRIC_CENTER"
	// Rough location, only loosely associated with the address
	LocationType_APPROXIMATE locationType = "APPROXIMATE"
)

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

// Geocode converts an address string into Places via the Google Geocoding API.
func (g *google) Geocode(ctx context.Context, address string) ([]geokit.Place, error) {
	params := url.Values{"address": {address}}
	apiResp, err := g.fetch(ctx, params)
	if err != nil {
		return nil, err
	}

	places := make([]geokit.Place, len(apiResp.Results))
	for i, r := range apiResp.Results {
		places[i] = toPlace(r)
	}
	return places, nil
}

// ReverseGeocode converts coordinates into Places via the Google Geocoding API.
func (g *google) ReverseGeocode(ctx context.Context, lat, lng float64) ([]geokit.Place, error) {
	latlng := strconv.FormatFloat(lat, 'f', -1, 64) + "," + strconv.FormatFloat(lng, 'f', -1, 64)
	params := url.Values{"latlng": {latlng}}
	apiResp, err := g.fetch(ctx, params)
	if err != nil {
		return nil, err
	}

	places := make([]geokit.Place, len(apiResp.Results))
	for i, r := range apiResp.Results {
		places[i] = toPlace(r)
	}
	return places, nil
}

// toPlace maps a Google API result to a geokit.Place, populating both
// coordinates and structured address fields from address_components.
func toPlace(r ApiResult) geokit.Place {
	p := geokit.Place{
		Lat:              r.Geometry.Location.Lat,
		Lng:              r.Geometry.Location.Lng,
		LocationType:     r.Geometry.LocationType,
		FormattedAddress: r.FormattedAddress,
		PlaceID:          r.PlaceID,
	}

	country := countryFromComponents(r.AddressComponents)
	parserForCountry(country).Parse(r.AddressComponents, &p)

	return p
}

func hasType(types []string, target string) bool {
	for _, t := range types {
		if t == target {
			return true
		}
	}
	return false
}
