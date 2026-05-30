package here

import (
	"strings"

	"github.com/supazonic/geokit"
)

// indiaParser handles Indian addresses.
//
// Indian addresses typically look like:
//   "Flat 101, Sunshine Apartments, MG Road, Koramangala, Bangalore 560095"
//
// HERE maps these as:
//   building    → complex/building name  e.g. "Sunshine Apartments"
//   houseNumber → flat/unit number      e.g. "101"
//   street      → street name           e.g. "MG Road"
//   district    → colony/neighborhood   e.g. "Koramangala"
//
// Line1 = building + houseNumber + street
// Line2 = district (colony/neighborhood)
type indiaParser struct{}

func (i indiaParser) Parse(addr hereAddress, p *geokit.Place) {
	parseCommon(addr, p)

	parts := make([]string, 0, 3)
	if addr.Building != "" {
		parts = append(parts, addr.Building)
	}
	if addr.HouseNumber != "" {
		parts = append(parts, addr.HouseNumber)
	}
	if addr.Street != "" {
		parts = append(parts, addr.Street)
	}
	p.Line1 = strings.Join(parts, ", ")
	p.Line2 = addr.District
}
