package google

import (
	"strings"

	"github.com/supazonic/geokit"
)

// indiaParser handles Indian addresses.
//
// Indian addresses typically look like:
//   "Flat 101, Sunshine Apartments, MG Road, Koramangala, Bangalore 560095"
//
// Google maps these as:
//   premise            → building/complex name  e.g. "Sunshine Apartments"
//   street_number      → unit/flat number       e.g. "Flat 101"
//   route              → street name            e.g. "MG Road"
//   sublocality_level_1→ colony/area/neighborhood e.g. "Koramangala"
//
// Line1 = premise + street_number + route  → "Sunshine Apartments, Flat 101, MG Road"
// Line2 = sublocality_level_1              → "Koramangala"
type indiaParser struct{}

func (i indiaParser) Parse(components []apiAddressComponent, p *geokit.Place) {
	parseCommon(components, p)

	var premise, streetNumber, route, sublocalityL1 string
	for _, c := range components {
		switch {
		case hasType(c.Types, "premise"):
			premise = c.LongName
		case hasType(c.Types, "street_number"):
			streetNumber = c.LongName
		case hasType(c.Types, "route"):
			route = c.LongName
		case hasType(c.Types, "sublocality_level_1"):
			sublocalityL1 = c.LongName
		}
	}

	parts := make([]string, 0, 3)
	if premise != "" {
		parts = append(parts, premise)
	}
	if streetNumber != "" {
		parts = append(parts, streetNumber)
	}
	if route != "" {
		parts = append(parts, route)
	}
	p.Line1 = strings.Join(parts, ", ")
	p.Line2 = sublocalityL1
}
