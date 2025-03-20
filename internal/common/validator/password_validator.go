package validator

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
)

var (
	// Common password errors
	ErrPasswordTooShort        = errors.New("password too short")
	ErrPasswordTooLong         = errors.New("password too long")
	ErrPasswordNoUpper         = errors.New("password must contain at least one uppercase letter")
	ErrPasswordNoLower         = errors.New("password must contain at least one lowercase letter")
	ErrPasswordNoNumber        = errors.New("password must contain at least one number")
	ErrPasswordNoSpecial       = errors.New("password must contain at least one special character")
	ErrPasswordCommonPattern   = errors.New("password contains common patterns")
	ErrPasswordFoundInBreaches = errors.New("password found in known data breaches")
)

// PasswordValidationOptions defines customizable password validation rules
type PasswordValidationOptions struct {
	MinLength      int
	MaxLength      int
	RequireUpper   bool
	RequireLower   bool
	RequireNumber  bool
	RequireSpecial bool
	CheckCommon    bool
	CheckBreaches  bool
}

// DefaultPasswordOptions returns the standard password validation options
func DefaultPasswordOptions() PasswordValidationOptions {
	return PasswordValidationOptions{
		MinLength:      8,
		MaxLength:      128,
		RequireUpper:   true,
		RequireLower:   true,
		RequireNumber:  true,
		RequireSpecial: true,
		CheckCommon:    true,
		CheckBreaches:  true,
	}
}

// ValidatePassword checks if the given password meets the security requirements
func ValidatePassword(password string, options PasswordValidationOptions) error {
	// Check length
	if len(password) < options.MinLength {
		return ErrPasswordTooShort
	}
	if options.MaxLength > 0 && len(password) > options.MaxLength {
		return ErrPasswordTooLong
	}

	// Check character requirements
	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if options.RequireUpper && !hasUpper {
		return ErrPasswordNoUpper
	}
	if options.RequireLower && !hasLower {
		return ErrPasswordNoLower
	}
	if options.RequireNumber && !hasNumber {
		return ErrPasswordNoNumber
	}
	if options.RequireSpecial && !hasSpecial {
		return ErrPasswordNoSpecial
	}

	// Check for common patterns
	if options.CheckCommon && hasCommonPattern(password) {
		return ErrPasswordCommonPattern
	}

	// Check against breach database (simplified - would require external API in production)
	if options.CheckBreaches && isBreachedPassword(password) {
		return ErrPasswordFoundInBreaches
	}

	return nil
}

// hasCommonPattern checks for common patterns like sequences, repetitions, etc.
func hasCommonPattern(password string) bool {
	patterns := []string{
		"123", "1234", "12345", "123456", "654321",
		"qwerty", "asdfgh", "zxcvbn", "password",
		"admin", "user", "login", "welcome",
	}

	lowerPass := strings.ToLower(password)
	for _, pattern := range patterns {
		if strings.Contains(lowerPass, pattern) {
			return true
		}
	}

	// Check for repeating characters (e.g., "aaa", "111")
	repeatingPattern := regexp.MustCompile(`(.)\1{2,}`)
	if repeatingPattern.MatchString(password) {
		return true
	}

	return false
}

// isBreachedPassword checks if a password appears in known breaches
// In a real system, this would connect to a service like HIBP
func isBreachedPassword(password string) bool {
	// This is a placeholder for real breach checking logic
	// In production, you'd use a service like Have I Been Pwned via API
	// or maintain your own breach database

	// For testing/demo - list of top breached passwords
	topBreached := map[string]bool{
		"123456": true, "password": true, "123456789": true,
		"12345678": true, "12345": true, "qwerty": true,
		"abc123": true, "football": true, "1234567": true,
		"monkey": true, "111111": true, "letmein": true,
	}

	return topBreached[strings.ToLower(password)]
}

// PasswordStrengthScore calculates a score from 0-100 representing password strength
func PasswordStrengthScore(password string) int {
	// Base score starts at 0
	score := 0

	// Length contribution (up to 25 points)
	if len(password) >= 8 {
		score += 10
		// Additional points for longer passwords
		extraLength := min(len(password)-8, 16) // Capped at +16 points
		score += extraLength
	}

	// Character variety (up to 40 points)
	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if hasUpper {
		score += 10
	}
	if hasLower {
		score += 10
	}
	if hasDigit {
		score += 10
	}
	if hasSpecial {
		score += 10
	}

	// Pattern deductions (up to -30 points)
	if hasCommonPattern(password) {
		score -= 30
	}

	// Breach database deduction
	if isBreachedPassword(password) {
		score -= 40
	}

	// Ensure score is between 0-100
	return max(0, min(100, score))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}
