// dateTimeUtil is a package providing some convenient functions that is date/time related.
package dateTimeUtil

import (
	"log"
	"time"
)

// DefaultFormat. Highly recommended not to change this. But if you know what you are doing, feel free to change it.
const (
	DefaultFormat = "2006-01-02 15:04:05"
)

// GetCurrentDateTime to get the current date time.
func GetCurrentDateTime() time.Time {
	return time.Now()
}

// Parse to parse the provided datetime string into the time.Time object.
func Parse(datetime string) (*time.Time, error) {
	tt, err := time.Parse(DefaultFormat, datetime)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return &tt, nil
}

// Format to format the provided datetime time.Time into the string object.
func Format(datetime time.Time) string {
	return datetime.Format(DefaultFormat)
}

// ParseCustom is like Parse but take in a custom format parameter.
func ParseCustom(datetime string, format string) (*time.Time, error) {
	tt, err := time.Parse(format, datetime)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return &tt, nil
}

// FormatCustom is like Format but take in a custom format parameter.
func FormatCustom(datetime time.Time, format string) string {
	return datetime.Format(format)
}
