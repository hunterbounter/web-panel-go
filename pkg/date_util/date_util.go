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

func DateDiffInDays(t1, t2 time.Time) int {
	return int(t1.Sub(t2).Hours() / 24)
}

const layoutDB = "02.01.2006 15:04:05"      // DB'den gelen tarih formatı
const layoutCurrent = "2006.01.02 15:04:05" // DateYYYYMMDDHH24MISS formatı

// KILL_DOCKER_PROCESS_TIMEOUT süresi (saniye cinsinden)
const KILL_DOCKER_PROCESS_TIMEOUT = 300 // 5 dakika = 300 saniye

// DateDiff iki tarih/zaman stringi arasındaki farkı saniye cinsinden döndürür
func DateDiff(date1, date2 string) (int64, error) {
	t1, err1 := time.Parse(layoutDB, date1)
	if err1 != nil {
		return 0, fmt.Errorf("error parsing date1: %v", err1)
	}

	t2, err2 := time.Parse(layoutCurrent, date2)
	if err2 != nil {
		return 0, fmt.Errorf("error parsing date2: %v", err2)
	}

	diff := t2.Sub(t1)
	return int64(diff.Seconds()), nil
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
