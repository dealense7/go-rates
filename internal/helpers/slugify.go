package helpers

import (
	"golang.org/x/text/unicode/norm"
	"regexp"
	"strings"
	"unicode"
)

func Slugify(s string) string {
	// 1. normalize and strip accents
	t := norm.NFD.String(s)
	b := make([]rune, 0, len(t))
	for _, r := range t {
		if unicode.Is(unicode.Mn, r) { // Mn = nonspacing mark
			continue
		}
		b = append(b, r)
	}
	// 2. to lower
	out := strings.ToLower(string(b))

	// 3. replace non-alphanum runs with dash
	//    pattern: anything except a–z, 0–9
	re := regexp.MustCompile(`[^a-z0-9]+`)
	out = re.ReplaceAllString(out, "-")

	// 4. trim leading/trailing dashes
	out = strings.Trim(out, "-")

	return out
}
