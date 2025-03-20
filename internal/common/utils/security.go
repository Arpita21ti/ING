package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/argon2"
)

// SecurityConstants defines security-related constants
const (
	// Token related
	TokenExpiryTime = 24 * time.Hour

	// CSRF related
	CSRFTokenLength = 32
	CSRFCookieName  = "csrf_token"
	CSRFHeaderName  = "X-CSRF-Token"

	// Password hashing related
	ArgonTime    = 1
	ArgonMemory  = 64 * 1024
	ArgonThreads = 4
	ArgonKeyLen  = 32
	ArgonSaltLen = 16
)

// ArgonParams represents the parameters used for Argon2 hashing
type ArgonParams struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

// DefaultArgonParams returns the default parameters for Argon2 hashing
func DefaultArgonParams() *ArgonParams {
	return &ArgonParams{
		Memory:      ArgonMemory,
		Iterations:  ArgonTime,
		Parallelism: ArgonThreads,
		SaltLength:  ArgonSaltLen,
		KeyLength:   ArgonKeyLen,
	}
}

// GenerateRandomBytes generates cryptographically secure random bytes
func GenerateRandomBytes(length int) ([]byte, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

// GenerateToken generates a secure random token as a hex string
func GenerateToken(length int) (string, error) {
	bytes, err := GenerateRandomBytes(length)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

// GenerateCSRFToken generates a new CSRF token
func GenerateCSRFToken() (string, error) {
	return GenerateToken(CSRFTokenLength)
}

// SetCSRFToken sets a CSRF token in both context and cookie
func SetCSRFToken(c *gin.Context) (string, error) {
	token, err := GenerateCSRFToken()
	if err != nil {
		return "", err
	}

	// Store token in context for template rendering
	c.Set("csrf_token", token)

	// Set cookie
	c.SetCookie(
		CSRFCookieName,
		token,
		int(TokenExpiryTime.Seconds()),
		"/",
		"",
		true, // Secure
		true, // HTTP only
	)

	return token, nil
}

// ValidateCSRFToken validates a CSRF token against the one stored in cookie
func ValidateCSRFToken(c *gin.Context) bool {
	// Get token from cookie
	cookieToken, err := c.Cookie(CSRFCookieName)
	if err != nil {
		return false
	}

	// Get token from header
	headerToken := c.GetHeader(CSRFHeaderName)
	if headerToken == "" {
		// If not in header, try from form
		headerToken = c.PostForm("csrf_token")
	}

	if headerToken == "" {
		return false
	}

	// Compare tokens - use constant time comparison to prevent timing attacks
	return subtle.ConstantTimeCompare([]byte(cookieToken), []byte(headerToken)) == 1
}

// CSRFMiddleware is a middleware that validates CSRF tokens
func CSRFMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only validate for state-changing methods
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS" || c.Request.Method == "TRACE" {
			c.Next()
			return
		}

		// Validate CSRF token
		if !ValidateCSRFToken(c) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    "INVALID_CSRF_TOKEN",
				"message": "Invalid or missing CSRF token",
			})
			return
		}

		c.Next()
	}
}

// HashPassword hashes a password using Argon2id
func HashPassword(password string) (string, error) {
	params := DefaultArgonParams()

	// Generate a random salt
	salt, err := GenerateRandomBytes(int(params.SaltLength))
	if err != nil {
		return "", err
	}

	// Hash the password
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		params.Iterations,
		params.Memory,
		params.Parallelism,
		params.KeyLength,
	)

	// Encode the parameters, salt, and hash into a string
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Format: $argon2id$v=19$m=memory,t=time,p=parallel$salt$hash
	encodedHash := fmt.Sprintf(
		"$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		params.Memory,
		params.Iterations,
		params.Parallelism,
		b64Salt,
		b64Hash,
	)

	return encodedHash, nil
}

// VerifyPassword verifies a password against a hash
func VerifyPassword(password, encodedHash string) (bool, error) {
	// Extract the parameters, salt, and hash from the encoded hash
	params, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	// Hash the password with the same parameters and salt
	otherHash := argon2.IDKey(
		[]byte(password),
		salt,
		params.Iterations,
		params.Memory,
		params.Parallelism,
		params.KeyLength,
	)

	// Compare the hashes in constant time
	return subtle.ConstantTimeCompare(hash, otherHash) == 1, nil
}

// decodeHash decodes an Argon2id hash string into its parameters, salt, and hash
func decodeHash(encodedHash string) (*ArgonParams, []byte, []byte, error) {
	// Expected format: $argon2id$v=19$m=memory,t=time,p=parallel$salt$hash
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return nil, nil, nil, errors.New("invalid hash format")
	}

	if parts[1] != "argon2id" {
		return nil, nil, nil, errors.New("unsupported algorithm")
	}

	// Parse the parameters
	var params ArgonParams

	paramParts := strings.Split(parts[3], ",")
	for _, part := range paramParts {
		param := strings.Split(part, "=")
		if len(param) != 2 {
			return nil, nil, nil, errors.New("invalid parameter format")
		}

		switch param[0] {
		case "m":
			val := 0
			fmt.Sscanf(param[1], "%d", &val)
			params.Memory = uint32(val)
		case "t":
			val := 0
			fmt.Sscanf(param[1], "%d", &val)
			params.Iterations = uint32(val)
		case "p":
			val := 0
			fmt.Sscanf(param[1], "%d", &val)
			params.Parallelism = uint8(val)
		}
	}

	// Decode the salt and hash
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, nil, nil, err
	}
	params.SaltLength = uint32(len(salt))

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, nil, nil, err
	}
	params.KeyLength = uint32(len(hash))

	return &params, salt, hash, nil
}

// SanitizeFilename sanitizes a filename to prevent path traversal attacks
func SanitizeFilename(filename string) string {
	// Remove path components
	filename = strings.ReplaceAll(filename, "/", "")
	filename = strings.ReplaceAll(filename, "\\", "")

	// Remove potentially dangerous characters
	safeFilenameRegex := regexp.MustCompile(`[^\w\d\.-]`)
	filename = safeFilenameRegex.ReplaceAllString(filename, "_")

	return filename
}

// SecureHeaders is a middleware to set security headers
func SecureHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Security headers
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self'; object-src 'none'; img-src 'self' data:; style-src 'self' 'unsafe-inline'; font-src 'self'; frame-src 'none'; frame-ancestors 'none'")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "camera=(), microphone=(), geolocation=()")

		c.Next()
	}
}

// RateLimiter is a simple rate limiter based on IP address
type RateLimiter struct {
	requests     map[string][]time.Time
	windowSize   time.Duration
	maxRequests  int
	cleanupTimer *time.Timer
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(windowSize time.Duration, maxRequests int) *RateLimiter {
	rl := &RateLimiter{
		requests:    make(map[string][]time.Time),
		windowSize:  windowSize,
		maxRequests: maxRequests,
	}

	// Start periodic cleanup
	rl.cleanupTimer = time.AfterFunc(windowSize, func() {
		rl.cleanup()
	})

	return rl
}

// Allow checks if a request from the given IP is allowed
func (rl *RateLimiter) Allow(ip string) bool {
	now := time.Now()
	cutoff := now.Add(-rl.windowSize)

	// Filter out requests older than the window
	var recent []time.Time
	for _, t := range rl.requests[ip] {
		if t.After(cutoff) {
			recent = append(recent, t)
		}
	}

	// Update requests
	if len(recent) >= rl.maxRequests {
		rl.requests[ip] = recent
		return false
	}

	rl.requests[ip] = append(recent, now)
	return true
}

// RemainingRequests returns the number of remaining requests for the given IP
func (rl *RateLimiter) RemainingRequests(ip string) int {
	now := time.Now()
	cutoff := now.Add(-rl.windowSize)

	// Filter out requests older than the window
	var count int
	for _, t := range rl.requests[ip] {
		if t.After(cutoff) {
			count++
		}
	}

	return rl.maxRequests - count
}

// cleanup removes old entries from the rate limiter
func (rl *RateLimiter) cleanup() {
	now := time.Now()
	cutoff := now.Add(-rl.windowSize)

	for ip, times := range rl.requests {
		var recent []time.Time
		for _, t := range times {
			if t.After(cutoff) {
				recent = append(recent, t)
			}
		}

		if len(recent) == 0 {
			delete(rl.requests, ip)
		} else {
			rl.requests[ip] = recent
		}
	}

	// Schedule next cleanup
	rl.cleanupTimer.Reset(rl.windowSize)
}

// Close stops the rate limiter's cleanup timer
func (rl *RateLimiter) Close() {
	rl.cleanupTimer.Stop()
}

// RateLimitMiddleware creates a middleware for rate limiting based on IP
func RateLimitMiddleware(windowSize time.Duration, maxRequests int) gin.HandlerFunc {
	limiter := NewRateLimiter(windowSize, maxRequests)

	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !limiter.Allow(ip) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code":    "RATE_LIMIT_EXCEEDED",
				"message": "Rate limit exceeded",
				"details": map[string]interface{}{
					"retryAfter": int(windowSize.Seconds()),
				},
			})
			return
		}

		// Set rate limit headers
		remaining := limiter.RemainingRequests(ip)
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", maxRequests))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(windowSize).Unix()))

		c.Next()
	}
}
