package validator

import (
	"html"
	"regexp"
	"strings"
	"unicode"
)

// SanitizationOptions defines what sanitization rules to apply
type SanitizationOptions struct {
	TrimSpace       bool
	RemoveHTML      bool
	AllowedHTMLTags []string // Only used if RemoveHTML is false
	MaxLength       int      // 0 means no max
	DisallowControl bool
	DisallowScripts bool
}

// DefaultSanitizationOptions returns common sanitization defaults
func DefaultSanitizationOptions() SanitizationOptions {
	return SanitizationOptions{
		TrimSpace:       true,
		RemoveHTML:      true,
		AllowedHTMLTags: []string{},
		MaxLength:       0,
		DisallowControl: true,
		DisallowScripts: true,
	}
}

// RichTextSanitizationOptions returns sanitization options suitable for rich text input
func RichTextSanitizationOptions() SanitizationOptions {
	return SanitizationOptions{
		TrimSpace:       true,
		RemoveHTML:      false,
		AllowedHTMLTags: []string{"p", "br", "b", "i", "strong", "em", "ul", "ol", "li", "a", "h1", "h2", "h3", "h4", "h5", "h6"},
		MaxLength:       0,
		DisallowControl: true,
		DisallowScripts: true,
	}
}

// SanitizeString removes potentially harmful content from a string based on options
func SanitizeString(input string, options SanitizationOptions) string {
	if input == "" {
		return input
	}

	result := input

	// Trim whitespace if required
	if options.TrimSpace {
		result = strings.TrimSpace(result)
	}

	// Apply length limit if specified
	if options.MaxLength > 0 && len(result) > options.MaxLength {
		result = result[:options.MaxLength]
	}

	// Handle HTML content
	if options.RemoveHTML {
		// Completely remove all HTML tags
		result = StripHTMLTags(result)
	} else if len(options.AllowedHTMLTags) > 0 {
		// Remove only disallowed HTML tags
		result = SanitizeHTMLAllowTags(result, options.AllowedHTMLTags)
	}

	// Remove control characters if specified
	if options.DisallowControl {
		result = RemoveControlCharacters(result)
	}

	// Special handling for script-like content
	if options.DisallowScripts {
		result = RemoveScriptContent(result)
	}

	return result
}

// StripHTML removes all HTML tags from the input string
func StripHTMLTags(input string) string {
	// First handle HTML entity decoding
	input = html.UnescapeString(input)

	// Basic HTML tag removal pattern
	re := regexp.MustCompile(`<[^>]*>`)
	return re.ReplaceAllString(input, "")
}

// SanitizeHTMLAllowTags removes all HTML tags except those in the whitelist
func SanitizeHTMLAllowTags(input string, allowedTags []string) string {
	// Construct regex for allowed tags
	if len(allowedTags) == 0 {
		return StripHTMLTags(input)
	}

	// Decode HTML entities
	input = html.UnescapeString(input)

	// Replace all disallowed tags
	result := input
	allowedTagsMap := make(map[string]bool)
	for _, tag := range allowedTags {
		allowedTagsMap[strings.ToLower(tag)] = true
	}

	// This is a simplistic approach - a proper HTML sanitizer library would be better
	tagRegex := regexp.MustCompile(`<\/?([a-zA-Z][a-zA-Z0-9]*)[^>]*>`)
	result = tagRegex.ReplaceAllStringFunc(result, func(match string) string {
		tag := tagRegex.FindStringSubmatch(match)[1]
		if allowedTagsMap[strings.ToLower(tag)] {
			// Keep allowed tags, but sanitize their attributes
			if strings.HasPrefix(match, "</") {
				// Closing tag - keep as is
				return match
			}
			// Opening tag with potential attributes - keep only safe ones
			return sanitizeTagAttributes(match)
		}
		// Remove disallowed tags
		return ""
	})

	return result
}

// sanitizeTagAttributes removes potentially harmful attributes from HTML tags
func sanitizeTagAttributes(tag string) string {
	// This is a simplified version - a real implementation would be more comprehensive

	// Extract tag name
	tagNameRegex := regexp.MustCompile(`<([a-zA-Z][a-zA-Z0-9]*)`)
	matches := tagNameRegex.FindStringSubmatch(tag)
	if len(matches) < 2 {
		return "" // Invalid tag format
	}
	tagName := matches[1]

	// List of allowed attributes (simplified)
	allowedAttrs := map[string]bool{
		"href": true, "title": true, "alt": true,
		"class": true, "id": true, "name": true,
	}

	// Disallow on* attributes (event handlers) and javascript: URLs
	dangerousPatterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)on\w+\s*=`),              // onclick, onload, etc.
		regexp.MustCompile(`(?i)javascript\s*:`),         // javascript: protocol
		regexp.MustCompile(`(?i)data\s*:`),               // data: URLs
		regexp.MustCompile(`(?i)expression\s*\(`),        // CSS expressions
		regexp.MustCompile(`(?i)(document|window)\s*\.`), // Direct JS object access
	}

	// Check for dangerous patterns
	for _, pattern := range dangerousPatterns {
		if pattern.MatchString(tag) {
			// Return just the tag without attributes
			return "<" + tagName + ">"
		}
	}

	// Extract attributes
	attrRegex := regexp.MustCompile(`([a-zA-Z][a-zA-Z0-9]*)\s*=\s*(['"])(.*?)\2`)
	safeTag := "<" + tagName

	attrs := attrRegex.FindAllStringSubmatch(tag, -1)
	for _, attr := range attrs {
		if len(attr) < 4 {
			continue
		}

		attrName := strings.ToLower(attr[1])
		attrValue := attr[3]

		// Only allow specific attributes
		if allowedAttrs[attrName] {
			// Further sanitize href to prevent javascript
			if attrName == "href" {
				if !strings.HasPrefix(strings.TrimSpace(strings.ToLower(attrValue)), "http") &&
					!strings.HasPrefix(strings.TrimSpace(strings.ToLower(attrValue)), "/") &&
					!strings.HasPrefix(strings.TrimSpace(strings.ToLower(attrValue)), "#") {
					continue
				}
			}

			safeTag += " " + attrName + "=\"" + html.EscapeString(attrValue) + "\""
		}
	}

	return safeTag + ">"
}

// RemoveControlCharacters removes Unicode control characters
func RemoveControlCharacters(input string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsControl(r) && r != '\n' && r != '\t' && r != '\r' {
			return -1 // Drop the character
		}
		return r
	}, input)
}

// RemoveScriptContent removes script-like content that might be used for XSS
func RemoveScriptContent(input string) string {
	// Remove <script> tags and their contents
	scriptTagRegex := regexp.MustCompile(`(?is)<script.*?>.*?</script>`)
	result := scriptTagRegex.ReplaceAllString(input, "")

	// Remove javascript: protocol from attributes
	jsProtocolRegex := regexp.MustCompile(`(?i)(href|src|action)\s*=\s*(['"])\s*javascript:.*?\2`)
	result = jsProtocolRegex.ReplaceAllStringFunc(result, func(match string) string {
		// Extract the attribute name
		parts := strings.SplitN(match, "=", 2)
		if len(parts) < 2 {
			return ""
		}
		// Return just the attribute name with empty value
		return parts[0] + "=\"\""
	})

	// Remove event handlers (on*)
	eventHandlerRegex := regexp.MustCompile(`(?i)\s+on\w+\s*=\s*(['"]).*?\1`)
	result = eventHandlerRegex.ReplaceAllString(result, "")

	// Remove other potentially dangerous tags
	dangerousTags := []string{"iframe", "object", "embed", "base", "form", "input", "button", "textarea", "select", "option"}
	for _, tag := range dangerousTags {
		tagRegex := regexp.MustCompile(`(?is)<` + tag + `[^>]*>.*?</` + tag + `>`)
		result = tagRegex.ReplaceAllString(result, "")

		// Also remove self-closing variants
		selfClosingRegex := regexp.MustCompile(`(?is)<` + tag + `[^>]*/>`)
		result = selfClosingRegex.ReplaceAllString(result, "")
	}

	// Remove data: URIs which could contain executable content
	dataURIRegex := regexp.MustCompile(`(?i)(src|href)\s*=\s*(['"])\s*data:.*?\2`)
	result = dataURIRegex.ReplaceAllStringFunc(result, func(match string) string {
		parts := strings.SplitN(match, "=", 2)
		if len(parts) < 2 {
			return ""
		}
		return parts[0] + "=\"\""
	})

	// Remove CSS expressions
	cssExprRegex := regexp.MustCompile(`(?i)expression\s*\(.*?\)`)
	result = cssExprRegex.ReplaceAllString(result, "")

	return result
}

// SanitizeEmail validates and sanitizes email addresses
func SanitizeEmail(email string) string {
	// Trim whitespace
	email = strings.TrimSpace(email)

	// Basic regex for email format
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return "" // Invalid email
	}

	// Convert to lowercase
	email = strings.ToLower(email)

	return email
}

// SanitizeUsername removes potentially dangerous characters from usernames
func SanitizeUsername(username string) string {
	// Trim whitespace
	username = strings.TrimSpace(username)

	// Remove special characters except _ and -
	usernameRegex := regexp.MustCompile(`[^a-zA-Z0-9_\-]`)
	username = usernameRegex.ReplaceAllString(username, "")

	// Ensure it's not empty after sanitization
	if username == "" {
		return ""
	}

	return username
}

// IsValidURL checks if a string is a valid URL
func IsValidURL(url string) bool {
	// Very basic URL validation
	url = strings.TrimSpace(url)

	// Must start with http:// or https://
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return false
	}

	// Basic URL pattern
	urlRegex := regexp.MustCompile(`^(http|https)://[a-zA-Z0-9\-]+(\.[a-zA-Z0-9\-]+)+(:[0-9]+)?(/[a-zA-Z0-9\-._~:/?#[\]@!$&'()*+,;=]*)?$`)
	return urlRegex.MatchString(url)
}

// SanitizeJSON removes potentially dangerous content from JSON input
func SanitizeJSON(jsonInput string) string {
	// This is a simplified approach - consider using a proper JSON parser
	// to validate and sanitize structured data

	// Remove control characters
	jsonInput = RemoveControlCharacters(jsonInput)

	// Remove script-like content
	jsonInput = RemoveScriptContent(jsonInput)

	return jsonInput
}
