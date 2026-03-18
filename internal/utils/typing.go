package utils

import "time"

// CalculateWPM calculates typing speed in words per minute
// Standard definition: 5 characters = 1 word
func CalculateWPM(charsTyped int, duration time.Duration) float64 {
	if duration == 0 {
		return 0
	}
	minutes := duration.Seconds() / 60.0
	words := float64(charsTyped) / 5.0
	return words / minutes
}

// CalculateErrorRate calculates the percentage of errors in typed text
// Compares typed text with target text character by character
func CalculateErrorRate(typed, target string) float64 {
	if len(target) == 0 {
		return 0
	}

	errors := 0
	targetRunes := []rune(target)
	typedRunes := []rune(typed)

	for i := 0; i < len(typedRunes) && i < len(targetRunes); i++ {
		if typedRunes[i] != targetRunes[i] {
			errors++
		}
	}

	if len(typedRunes) > len(targetRunes) {
		errors += len(typedRunes) - len(targetRunes)
	}

	return float64(errors) / float64(len(target)) * 100
}
