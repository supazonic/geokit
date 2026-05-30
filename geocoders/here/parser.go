package here

import "github.com/supazonic/geokit"

// AddressParser builds the address fields of a Place from a HERE address object.
// Implement this interface to add support for a new country's address format.
type AddressParser interface {
	Parse(addr hereAddress, p *geokit.Place)
}

// parserForCountry returns the appropriate parser for a given country name.
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

// parseCommon fills the fields that map the same way across all countries.
// HERE's address object already has structured fields — no component walking needed.
func parseCommon(addr hereAddress, p *geokit.Place) {
	p.City = addr.City
	p.District = addr.County // HERE county = administrative_area_level_2
	p.State = addr.State
	p.PostalCode = addr.PostalCode
	p.Country = addr.CountryName
}
