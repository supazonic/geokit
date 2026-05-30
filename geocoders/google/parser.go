package google

import "github.com/supazonic/geokit"

// AddressParser builds the address fields of a Place from Google address_components.
// Implement this interface to add support for a new country's address format.
type AddressParser interface {
	Parse(components []apiAddressComponent, p *geokit.Place)
}

// parserForCountry returns the appropriate parser for a given country long_name.
// Falls back to defaultParser for countries without a custom implementation.
func parserForCountry(country string) AddressParser {
	switch country {
	case "Japan":
		return japanParser{}
	case "India":
		return indiaParser{}
	default:
		return defaultParser{}
	}
}

// countryFromComponents extracts the country long_name from address_components.
func countryFromComponents(components []apiAddressComponent) string {
	for _, c := range components {
		if hasType(c.Types, "country") {
			return c.LongName
		}
	}
	return ""
}

// parseCommon fills the fields that are structured the same across all countries:
// City, District, State, PostalCode, Country.
func parseCommon(components []apiAddressComponent, p *geokit.Place) {
	for _, c := range components {
		switch {
		case hasType(c.Types, "locality"):
			p.City = c.LongName
		case hasType(c.Types, "administrative_area_level_2"):
			p.District = c.LongName
		case hasType(c.Types, "administrative_area_level_1"):
			p.State = c.LongName
		case hasType(c.Types, "postal_code"):
			p.PostalCode = c.LongName
		case hasType(c.Types, "country"):
			p.Country = c.LongName
		}
	}
}
