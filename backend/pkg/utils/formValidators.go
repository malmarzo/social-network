package utils

import (
	"fmt"
	"regexp"
	"time"
)

// ValidateTextOnly checks if the string contains only letters (no numbers or special characters)
func ValidateTextOnly(input string) bool {
	match, _ := regexp.MatchString("^[A-Za-z]+$", input)
	return match
}

// ValidateAlphanumericUnderscore checks if the string contains only letters, numbers, and underscores
func ValidateAlphanumericUnderscore(input string) bool {
	match, _ := regexp.MatchString("^[A-Za-z0-9_]+$", input)
	return match
}

// ValidateEmail checks if the string is a valid email address
func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(pattern, email)
	return match
}

// IsStrongPassword checks if the given password is strong.
// A strong password must contain at least one uppercase letter, one lowercase letter, one digit, and be at least 8 characters long.
func ValidatePassword(password string) bool {
	// Define the regex patterns for each requirement
	uppercasePattern := `[A-Z]`
	lowercasePattern := `[a-z]`
	digitPattern := `\d`
	minLengthPattern := `.{8,}`
	// Compile the regex patterns
	uppercaseRe := regexp.MustCompile(uppercasePattern)
	lowercaseRe := regexp.MustCompile(lowercasePattern)
	digitRe := regexp.MustCompile(digitPattern)
	minLengthRe := regexp.MustCompile(minLengthPattern)
	// Check each requirement
	return minLengthRe.MatchString(password) &&
		uppercaseRe.MatchString(password) &&
		lowercaseRe.MatchString(password) &&
		digitRe.MatchString(password)
}

// ValidateDOB checks if the user is at least 15 years old
func ValidateDOB(dob string) bool {
	// Parse the date string
	birthDate, err := time.Parse("2006-01-02", dob)
	if err != nil {
		fmt.Printf("Error parsing date: %v, Format received: %s\n", err, dob)
		return false
	}

	// Get current time
	now := time.Now()

	// Calculate age
	age := now.Year() - birthDate.Year()

	// Adjust age if birthday hasn't occurred this year
	if now.Month() < birthDate.Month() || (now.Month() == birthDate.Month() && now.Day() < birthDate.Day()) {
		age--
	}

	return age >= 15
}

func SignupValidator(firstName, lastName, email, password, dob, nickname string) (string, bool) {
	if len(firstName) == 0 || len(lastName) == 0 || len(email) == 0 || len(password) == 0 || len(dob) == 0 || len(nickname) == 0 {
		return "Please fill the form!", false
	}

	// Validate first name and last name (text only)
	if !ValidateTextOnly(firstName) || !ValidateTextOnly(lastName) {
		return "Names should include letters only!", false
	}

	// Validate nickname (alphanumeric with underscore)
	if !ValidateAlphanumericUnderscore(nickname) {
		return "Invalid nickname", false
	}

	// Validate email format
	if !ValidateEmail(email) {
		return "Invalid email", false
	}

	if !ValidatePassword(password) {
		return "Use strong password!", false
	}

	return "", true
}
