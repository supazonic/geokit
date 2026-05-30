package here

import (
	"strings"

	"github.com/supazonic/geokit"
)

// defaultParser handles Western-style addresses.
// Line1 = houseNumber + street   e.g. "1600 Amphitheatre Pkwy"
// Line2 = building               e.g. "Googleplex"
type defaultParser struct{}

func (d defaultParser) Parse(addr hereAddress, p *geokit.Place) {
	parseCommon(addr, p)

	parts := make([]string, 0, 2)
	if addr.HouseNumber != "" {
		parts = append(parts, addr.HouseNumber)
	}
	if addr.Street != "" {
		parts = append(parts, addr.Street)
	}
	p.Line1 = strings.Join(parts, " ")
	p.Line2 = addr.Building
}
