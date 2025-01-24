package helper

import (
	"fmt"
	"regexp"
)

func ValidateNIK(nik string) error {
	// Define the regex pattern for a valid NIK
	pattern := `^(1[1-9]|21|[37][1-6]|5[1-3]|6[1-5]|[89][12])\d{2}\d{2}([04][1-9]|[1256][0-9]|[37][01])(0[1-9]|1[0-2])\d{2}\d{4}$`

	// Compile the regex
	re := regexp.MustCompile(pattern)

	if nik == "" {
		return fmt.Errorf("NIK is required")
	}

	// Check if NIK matches the pattern
	if !re.MatchString(nik) {
		return fmt.Errorf("NIK does not match the required format")
	}

	return nil
}
