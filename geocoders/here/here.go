package here

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/supazonic/geokit"
)

const (
	defaultGeocodeURL    = "https://geocode.search.hereapi.com/v1/geocode"
	defaultRevgeocodeURL = "https://revgeocode.search.hereapi.com/v1/revgeocode"
)

type here struct {
	apiKey        string
	httpClient    *http.Client
	geocodeURL    string
	revgeocodeURL string
}

// NewHereGeocoder returns a geokit.Geocoder backed by the HERE Geocoding & Search API v7.
func NewHereGeocoder(apiKey string) *here {
	return &here{
		apiKey:        apiKey,
		httpClient:    &http.Client{},
		geocodeURL:    defaultGeocodeURL,
		revgeocodeURL: defaultRevgeocodeURL,
	}
}

// --- internal response types ---

type apiResponse struct {
	Items []apiItem `json:"items"`
}

type apiItem struct {
	Title           string      `json:"title"`
	ID              string      `json:"id"`
	ResultType      string      `json:"resultType"`
	HouseNumberType string      `json:"houseNumberType"`
	Address         hereAddress `json:"address"`
	Position        hereLatLng  `json:"position"`
	Scoring         hereScoring `json:"scoring"`
}

type hereAddress struct {
	Label       string `json:"label"`
	CountryCode string `json:"countryCode"`
	CountryName string `json:"countryName"`
	StateCode   string `json:"stateCode"`
	State       string `json:"state"`
	County      string `json:"county"`
	City        string `json:"city"`
	District    string `json:"district"`
	Subdistrict string `json:"subdistrict"`
	Street      string `json:"street"`
	Block       string `json:"block"`
	Subblock    string `json:"subblock"`
	PostalCode  string `json:"postalCode"`
	HouseNumber string `json:"houseNumber"`
	Building    string `json:"building"`
	Unit        string `json:"unit"`
}

type hereLatLng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type hereScoring struct {
	QueryScore float64 `json:"queryScore"`
}

// apiError mirrors HERE's error response body.
type apiError struct {
	Status int    `json:"status"`
	Title  string `json:"title"`
	Cause  string `json:"cause"`
}

// --- shared HTTP helper ---

func (h *here) fetch(ctx context.Context, endpoint string, params url.Values) (*apiResponse, error) {
	params.Set("apiKey", h.apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var e apiError
		json.NewDecoder(resp.Body).Decode(&e)
		if e.Title != "" {
			return nil, fmt.Errorf("here geocoder: %s — %s", e.Title, e.Cause)
		}
		return nil, fmt.Errorf("here geocoder: HTTP %d", resp.StatusCode)
	}

	var apiResp apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	return &apiResp, nil
}

// --- geokit.Geocoder implementation ---

// Geocode converts an address string into Places via the HERE Geocoding API.
func (h *here) Geocode(ctx context.Context, address string) ([]geokit.Place, error) {
	params := url.Values{"q": {address}}
	apiResp, err := h.fetch(ctx, h.geocodeURL, params)
	if err != nil {
		return nil, err
	}

	places := make([]geokit.Place, len(apiResp.Items))
	for i, item := range apiResp.Items {
		places[i] = toPlace(item)
	}
	return places, nil
}

// ReverseGeocode converts coordinates into Places via the HERE Reverse Geocoding API.
func (h *here) ReverseGeocode(ctx context.Context, lat, lng float64) ([]geokit.Place, error) {
	at := strconv.FormatFloat(lat, 'f', -1, 64) + "," + strconv.FormatFloat(lng, 'f', -1, 64)
	params := url.Values{"at": {at}}
	apiResp, err := h.fetch(ctx, h.revgeocodeURL, params)
	if err != nil {
		return nil, err
	}

	places := make([]geokit.Place, len(apiResp.Items))
	for i, item := range apiResp.Items {
		places[i] = toPlace(item)
	}
	return places, nil
}

// toPlace maps a HERE API item to a geokit.Place.
// LocationType is set from houseNumberType when available (more specific),
// falling back to resultType.
func toPlace(item apiItem) geokit.Place {
	locType := item.ResultType
	if item.HouseNumberType != "" {
		locType = item.HouseNumberType
	}

	types := []string{item.ResultType}
	if item.HouseNumberType != "" {
		types = append(types, item.HouseNumberType)
	}

	p := geokit.Place{
		Lat:              item.Position.Lat,
		Lng:              item.Position.Lng,
		LocationType:     locType,
		FormattedAddress: item.Address.Label,
		PlaceID:          item.ID,
		Types:            types,
	}

	parserForCountry(item.Address.CountryName).Parse(item.Address, &p)

	return p
}
