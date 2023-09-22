package account

import (
	"strings"
)

func joinPrefixes(prefixes []string, separator string) string {
	if len(prefixes) > 0 {
		return strings.Join(prefixes, separator) + separator
	} else {
		return ""
	}
}
