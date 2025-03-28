package utils

import (
	"time"
	"log"

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


func ExtractDayFromEvents(dateStr string) (string, error) {
	// Define the format that matches your datetime string
	layout := time.RFC3339 // This corresponds to "2006-01-02T15:04:05Z"
	t, err := time.Parse(layout, dateStr)
	if err != nil {
		log.Println("Error parsing date:", err)
		return "", err
	}
	// Return the full weekday name (e.g., "Monday", "Tuesday", etc.)
	return t.Weekday().String(), nil
}
