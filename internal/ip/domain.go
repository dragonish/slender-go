package ip

import (
	"regexp"
	"strings"
)

func isIPv4(s string) bool {
	re := regexp.MustCompile(`^\d+\.\d+\.\d+\.\d+$`)
	return re.MatchString(s)
}

// GetDomain returns the domain from the hostname.
func GetDomain(hostname string) string {
	parts := strings.Split(hostname, ".")
	partsLen := len(parts)

	if partsLen <= 2 {
		return hostname
	} else if partsLen == 4 && isIPv4(hostname) {
		return hostname
	}

	return strings.Join(parts[1:], ".")
}
