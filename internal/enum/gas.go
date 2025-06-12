package enum

import (
	"fmt"
	"github.com/dealense7/go-rate-app/internal/helpers"
)

type GasProvider int

const (
	SOCAR     GasProvider = 1
	WISSOL    GasProvider = 2
	PORTAL    GasProvider = 3
	GULF      GasProvider = 4
	ROMPETROL GasProvider = 5
	LUKOILI   GasProvider = 7
	CONNECT   GasProvider = 9
)

var providerNames = map[GasProvider]string{
	SOCAR:     "Socar",
	WISSOL:    "Wissol",
	PORTAL:    "Portal",
	GULF:      "Gulf",
	ROMPETROL: "Rompetrol",
	LUKOILI:   "Lukoili",
	CONNECT:   "Connect",
}

func (p GasProvider) String() string {
	if s, ok := providerNames[p]; ok {
		return s
	}
	return fmt.Sprintf("GasProvider(%d)", int(p))
}

func (p GasProvider) Logo() string {
	var imagePath = map[GasProvider]string{
		SOCAR:     "static/img/logos/gas/socar.webp",
		WISSOL:    "static/img/logos/gas/wissol.webp",
		PORTAL:    "static/img/logos/gas/portal.webp",
		GULF:      "static/img/logos/gas/gulf.webp",
		ROMPETROL: "static/img/logos/gas/rompetrol.webp",
		LUKOILI:   "static/img/logos/gas/lukoili.webp",
		CONNECT:   "static/img/logos/gas/connect.webp",
	}

	if s, ok := imagePath[p]; ok {
		return s
	}
	return fmt.Sprintf("GasProvider(%d)", int(p))
}

func (p GasProvider) Slug() string {
	if s, ok := providerNames[p]; ok {
		return helpers.Slugify(s)
	}
	return fmt.Sprintf("GasProvider(%d)", int(p))
}
