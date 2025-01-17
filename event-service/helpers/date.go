package helpers

import (
	"log"
	"time"
)

func ParseDate(dateStr string) time.Time {
	// Define the layout based on the date format in req.DateStart and req.DateEnd
	const layout = "2006-01-02" // Example layout for dates like "2025-01-17"
	parsedDate, err := time.Parse(layout, dateStr)
	if err != nil {
		log.Fatalf("Failed to parse date: %v", err)
	}
	return parsedDate
}
