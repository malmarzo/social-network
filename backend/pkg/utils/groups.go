package utils

import (
	"time"

)

func ExtractDay(dateTime string) (string, error) {
	// Define custom layout (without seconds and time zone)
	const layout = "2006-01-02T15:04" // This matches the format "YYYY-MM-DDTHH:MM"
	// Parse the date time string to a Go time.Time object
	parsedTime, err := time.Parse(layout, dateTime)
	if err != nil {
		return "", err
	}

	// Extract the weekday name (e.g., "Monday")
	day := parsedTime.Weekday().String()
	return day, nil
}