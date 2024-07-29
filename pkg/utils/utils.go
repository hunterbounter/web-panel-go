package utils

import (
	"fmt"
	"html"
	"log"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

// Number to string
func AnyToString(number interface{}) string {
	return fmt.Sprintf("%v", number)
}

func AnyToInt(number interface{}) int {
	// fmt.Sprintf converts string to int
	i, err := strconv.Atoi(AnyToString(number))
	if err != nil {
		return 0 // Hata varsa, 0 döner
	}
	return i // if conversion is successful, return the converted number
}

func TimeToString(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func IsValidDomain(domain string) bool {

	// ilk normal domain mi diye kontrol et

	// Regular expression for validating domain name
	const domainPattern = `^(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z]{2,6}$`
	domainRegex := regexp.MustCompile(domainPattern)

	// Parse the URL
	parsedURL, err := url.Parse(domain)
	if err != nil {
		fmt.Println("Invalid URL:", err)
		return false
	}

	// Check the scheme
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		fmt.Println("Invalid scheme:", parsedURL.Scheme)
		return false
	}

	// Check if host is valid
	host := parsedURL.Hostname()
	if !domainRegex.MatchString(host) {
		fmt.Println("Invalid domain:", host)
		return false
	}

	// Check if port is valid (if provided)
	if parsedURL.Port() != "" {
		port := parsedURL.Port()
		if _, err := net.LookupPort("tcp", port); err != nil {
			fmt.Println("Invalid port:", port)
			return false
		}
	}

	return true
}

func IsValidIP(ip string) bool {
	// Regular expression for validating an IP address
	var ipRegex = regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)
	return ipRegex.MatchString(ip)

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

func RunningDir() string {

	path, err := os.Executable()
	if err != nil {
		return ""
	}

	//log.Println("Executable Path : ", path)

	// the executable directory
	exPath := filepath.Dir(path)

	//log.Println("Executable Path : ", exPath)

	exPath = exPath + "/../../"

	log.Println("RunningDir: ", exPath)

	return exPath
}
