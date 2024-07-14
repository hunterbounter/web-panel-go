package date_util

import (
	"fmt"
	"strings"
	"time"
)

// CurrentDate returns the current date in the format "2006-01-02"
func CurrentDate() string {
	return time.Now().Format("2006-01-02")
}

// CurrentTime returns the current time in the format "15:04:05"
func CurrentTime() string {
	return time.Now().Format("15:04:05")
}

// CurrentDateTime returns the current date and time in the format "2006-01-02 15:04:05"
func CurrentDateTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// CustomFormat returns the current date and time in a custom format
func CustomFormat(format string) string {
	return time.Now().Format(format)
}

// AddDays adds the specified number of days to the current date and returns it in the format "2006-01-02"
func AddDays(days int) string {
	return time.Now().AddDate(0, 0, days).Format("2006-01-02")
}

// AddHours adds the specified number of hours to the current time and returns it in the format "15:04:05"
func AddHours(hours int) string {
	return time.Now().Add(time.Duration(hours) * time.Hour).Format("15:04:05")
}

// AddMinutes adds the specified number of minutes to the current time and returns it in the format "15:04:05"
func AddMinutes(minutes int) string {
	return time.Now().Add(time.Duration(minutes) * time.Minute).Format("15:04:05")
}

// AddSeconds adds the specified number of seconds to the current time and returns it in the format "15:04:05"
func AddSeconds(seconds int) string {
	return time.Now().Add(time.Duration(seconds) * time.Second).Format("15:04:05")
}

func DateYYYYMMDD() string {
	return strings.Replace(time.Now().Format("2006-01-02"), "-", "", -1)
}

// dd.mm.yyyy hh24:mi:ss
func DateYYYYMMDDHH24MISS() string {
	return strings.Replace(time.Now().Format("2006-01-02 15:04:05"), "-", ".", -1)
}

// türkiye için +3 saat ekler
func DateYYYYMMDDHH24MISSWithTRTimezone() string {
	return strings.Replace(time.Now().Add(3*time.Hour).Format("2006-01-02 15:04:05"), "-", ".", -1)

}

func GetUnixTime() int64 {
	return time.Now().UnixNano()
}

// Get Unix Time String
func GetUnixTimeString() string {
	return fmt.Sprint(time.Now().Unix())

}

func DateTomorrowYYYYMMDD() string {
	return strings.Replace(time.Now().AddDate(0, 0, 1).Format("2006-01-02"), "-", "", -1)

}

// StartOfDay returns the start of the current day in the format "2006-01-02 00:00:00"
func StartOfDay() string {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return startOfDay.Format("2006-01-02 00:00:00")
}

// EndOfDay returns the end of the current day in the format "2006-01-02 23:59:59"
func EndOfDay() string {
	now := time.Now()
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
	return endOfDay.Format("2006-01-02 23:59:59")
}

func example() {
	fmt.Println("Current Date:", CurrentDate())
	fmt.Println("Current Time:", CurrentTime())
	fmt.Println("Current DateTime:", CurrentDateTime())
	fmt.Println("Custom Format (Jan 2, 2006):", CustomFormat("Jan 2, 2006"))
	fmt.Println("Add 5 Days:", AddDays(5))
	fmt.Println("Add 3 Hours:", AddHours(3))
	fmt.Println("Add 15 Minutes:", AddMinutes(15))
	fmt.Println("Add 30 Seconds:", AddSeconds(30))
	fmt.Println("Start of Day:", StartOfDay())
	fmt.Println("End of Day:", EndOfDay())
}
