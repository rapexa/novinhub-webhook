package utils

import (
	"regexp"
	"strings"
)

// IranianPhoneRegex defines patterns for Iranian phone numbers
var IranianPhoneRegex = []*regexp.Regexp{
	// Pattern 1: 09XXXXXXXXX (11 digits starting with 09)
	regexp.MustCompile(`\b09\d{9}\b`),
	// Pattern 2: +989XXXXXXXXX (with country code)
	regexp.MustCompile(`\+989\d{9}\b`),
	// Pattern 3: 00989XXXXXXXXX (with international prefix)
	regexp.MustCompile(`00989\d{9}\b`),
	// Pattern 4: 9XXXXXXXXX (without leading 0)
	regexp.MustCompile(`\b9\d{9}\b`),
}

// ExtractIranianPhoneNumbers extracts all Iranian phone numbers from text
func ExtractIranianPhoneNumbers(text string) []string {
	var phones []string
	seen := make(map[string]bool)

	// Clean the text - remove extra spaces and normalize
	cleanText := strings.TrimSpace(text)

	for _, regex := range IranianPhoneRegex {
		matches := regex.FindAllString(cleanText, -1)
		for _, match := range matches {
			// Normalize the phone number to standard format (09XXXXXXXXX)
			normalized := NormalizeIranianPhone(match)
			if normalized != "" && !seen[normalized] {
				phones = append(phones, normalized)
				seen[normalized] = true
			}
		}
	}

	return phones
}

// NormalizeIranianPhone normalizes Iranian phone numbers to 09XXXXXXXXX format
func NormalizeIranianPhone(phone string) string {
	// Remove all non-digit characters except +
	cleaned := regexp.MustCompile(`[^\d+]`).ReplaceAllString(phone, "")

	// Handle different formats
	if strings.HasPrefix(cleaned, "+989") {
		// +989XXXXXXXXX -> 09XXXXXXXXX
		return "0" + cleaned[3:]
	} else if strings.HasPrefix(cleaned, "00989") {
		// 00989XXXXXXXXX -> 09XXXXXXXXX
		return "0" + cleaned[5:]
	} else if strings.HasPrefix(cleaned, "9") && len(cleaned) == 10 {
		// 9XXXXXXXXX -> 09XXXXXXXXX
		return "0" + cleaned
	} else if strings.HasPrefix(cleaned, "09") && len(cleaned) == 11 {
		// Already in correct format
		return cleaned
	}

	return ""
}

// IsValidIranianPhone validates if a phone number is a valid Iranian mobile number
func IsValidIranianPhone(phone string) bool {
	normalized := NormalizeIranianPhone(phone)
	if normalized == "" {
		return false
	}

	// Check if it's exactly 11 digits and starts with 09
	if len(normalized) != 11 || !strings.HasPrefix(normalized, "09") {
		return false
	}

	// Check if the third digit is valid (Iranian mobile prefixes)
	validPrefixes := []string{"091", "092", "093", "094", "095", "096", "097", "098", "099"}
	prefix := normalized[:3]

	for _, validPrefix := range validPrefixes {
		if prefix == validPrefix {
			return true
		}
	}

	return false
}
