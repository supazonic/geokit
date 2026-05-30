package google

import (
	"strings"

	"github.com/supazonic/geokit"
)

// defaultParser handles Western-style addresses where Line1 is built from
// street_number + route (e.g. "1600 Amphitheatre Pkwy").
type defaultParser struct{}

func (d defaultParser) Parse(components []apiAddressComponent, p *geokit.Place) {
	parseCommon(components, p)

	var streetNumber, route string
	for _, c := range components {
		switch {
		case hasType(c.Types, "street_number"):
			streetNumber = c.LongName
		case hasType(c.Types, "route"):
			route = c.LongName
		case hasType(c.Types, "premise"):
			p.Line2 = c.LongName
		}
	}

	parts := make([]string, 0, 2)
	if streetNumber != "" {
		parts = append(parts, streetNumber)
	}
	if route != "" {
		parts = append(parts, route)
	}
	p.Line1 = strings.Join(parts, " ")
}
