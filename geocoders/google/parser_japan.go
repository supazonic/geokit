package google

import (
	"strings"

	"github.com/supazonic/geokit"
)

// japanParser handles Japanese addresses.
//
// Google returns Japan's block/lot (chome-ban-go) and neighborhood (cho) as
// sublocality levels rather than street_number + route:
//
//	sublocality_level_2 → chome-ban-go  e.g. "6 Chome-29-62"
//	sublocality_level_1 → cho name      e.g. "Nishikicho"
//
// Line1 = sublocality_level_2 + sublocality_level_1  → "6 Chome-29-62 Nishikicho"
// Line2 = premise (building name), if present
type japanParser struct{}

func (j japanParser) Parse(components []apiAddressComponent, p *geokit.Place) {
	parseCommon(components, p)

	var sublocalityL1, sublocalityL2, premise string
	for _, c := range components {
		switch {
		case hasType(c.Types, "sublocality_level_1"):
			sublocalityL1 = c.LongName
		case hasType(c.Types, "sublocality_level_2"):
			sublocalityL2 = c.LongName
		case hasType(c.Types, "premise"):
			premise = c.LongName
		}
	}

	parts := make([]string, 0, 2)
	if sublocalityL2 != "" {
		parts = append(parts, sublocalityL2)
	}
	if sublocalityL1 != "" {
		parts = append(parts, sublocalityL1)
	}
	p.Line1 = strings.Join(parts, " ")
	p.Line2 = premise
}
