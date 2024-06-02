package utils

import (
	"html"
	"regexp"
	"time"
)

func TimeToString(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func IsValidDomain(domain string) bool {
	// Regular expression for validating a domain name
	var domainRegex = regexp.MustCompile(`^(?i)[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?(\.[a-z]{2,})+$`)
	return domainRegex.MatchString(domain)
}

// abc
// SafeEscapeString safely escapes a value if it's a string. Returns an empty string if the value is nil or not a string.
func SafeEscapeString(value interface{}) string {
	if value == nil {
		return ""
	}

	str, ok := value.(string)
	if !ok {
		return ""
	}

	return html.EscapeString(str)
}

func GetString(value interface{}) string {
	if value == nil {
		return ""
	}

	str, ok := value.(string)
	if !ok {
		return ""
	}

	return str
}
