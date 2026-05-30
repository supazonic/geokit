package here

import (
	"strings"

	"github.com/supazonic/geokit"
)

// japanParser handles Japanese addresses.
//
// HERE returns Japan's address components as:
//
//	houseNumber → specific lot        e.g. "62"
//	block       → chome              e.g. "6-chome"
//	district    → cho/neighborhood   e.g. "Nishikicho"
//	city        → city               e.g. "Tachikawa"
//	county      → district/gun       e.g. "Iruma District"
//	state       → prefecture         e.g. "Tokyo"
//
// Line1 = houseNumber + block + district  e.g. "62 6-chome Nishikicho"
// Line2 = building, if present
type japanParser struct{}

func (j japanParser) Parse(addr hereAddress, p *geokit.Place) {
	parseCommon(addr, p)

	parts := make([]string, 0, 3)
	if addr.HouseNumber != "" {
		parts = append(parts, addr.HouseNumber)
	}
	if addr.Block != "" {
		parts = append(parts, addr.Block)
	}
	if addr.District != "" {
		parts = append(parts, addr.District)
	}
	p.Line1 = strings.Join(parts, " ")
	p.Line2 = addr.Building
}
