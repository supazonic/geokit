# geokit

A Go geocoding library with a provider-agnostic interface. Import one backend or swap between them without changing your application code.

## Installation

```sh
go get github.com/supazonic/geokit@v0.1.0
```

## What it does

geokit defines a `Geocoder` interface and a unified `Place` type. Provider-specific packages (Google, HERE) implement the interface and handle parsing each API's response into the same `Place` struct — including country-aware address formatting for Japan, India, and Western addresses.

## Providers

| Package | Provider |
|---|---|
| `geocoders/google` | Google Geocoding API |
| `geocoders/here` | HERE Geocoding & Search API v7 |

## Quick start

```go
import (
    "context"

    "github.com/supazonic/geokit"
    "github.com/supazonic/geokit/geocoders/google"
)

client := geokit.NewGeocoder(google.NewGoogleGeocoder("YOUR_API_KEY"))

// Forward geocoding — address to coordinates
places, err := client.Geocode(ctx, "6 Chome-29-62 Nishikicho, Tachikawa, Tokyo 190-0022, Japan")

// Reverse geocoding — coordinates to address
places, err := client.ReverseGeocode(ctx, 35.6897, 139.6922)
```

Swap to HERE by changing one line:

```go
import "github.com/supazonic/geokit/geocoders/here"

client := geokit.NewGeocoder(here.NewHereGeocoder("YOUR_API_KEY"))
```

## The Place type

Both `Geocode` and `ReverseGeocode` return `[]geokit.Place`.

```go
type Place struct {
    // Coordinates
    Lat          float64
    Lng          float64
    LocationType string  // precision: ROOFTOP, RANGE_INTERPOLATED, GEOMETRIC_CENTER, APPROXIMATE

    // Address
    Line1      string  // street-level: "6 Chome-29-62 Nishikicho" (JP) / "1600 Amphitheatre Pkwy" (US)
    Line2      string  // building name or room number, if any
    City       string
    District   string  // county/gun level, e.g. "Iruma District" (Japan)
    State      string  // prefecture in Japan, state/province elsewhere
    PostalCode string
    Country    string

    // Identifiers
    FormattedAddress string  // full human-readable address from the provider
    PlaceID          string  // provider-specific place identifier

    // Confidence
    PartialMatch bool     // true if the provider could not fully resolve the input
    Types        []string // result type tags, e.g. ["street_address"]
}
```

### Line1 and Line2 by country

Address formatting differs significantly between countries. Each provider package includes country-specific parsers:

| Country | Line1 | Line2 |
|---|---|---|
| Western (default) | house number + street name | building name |
| Japan | chome-ban-go + neighborhood (cho) | building name |
| India | building + house number + street | neighborhood/colony |

### Reading confidence

```go
for _, p := range places {
    fmt.Println(p.LocationType) // "ROOFTOP" = exact building coordinate
    fmt.Println(p.PartialMatch) // true = input was not fully matched
    fmt.Println(p.Types)        // ["street_address"] = specific address result
}

// Multiple results returned = provider was uncertain which place you meant
if len(places) > 1 {
    // ask the user to pick
}
```

## Implementing a custom provider

Implement the `geokit.Geocoder` interface and pass it to `NewGeocoder`:

```go
type MyGeocoder struct{}

func (m *MyGeocoder) Geocode(ctx context.Context, address string) ([]geokit.Place, error) {
    // call your API, map response to []geokit.Place
}

func (m *MyGeocoder) ReverseGeocode(ctx context.Context, lat, lng float64) ([]geokit.Place, error) {
    // call your API, map response to []geokit.Place
}

client := geokit.NewGeocoder(&MyGeocoder{})
```
