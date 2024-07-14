package utils

import (
	"fmt"
	"html"
	"log"
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
