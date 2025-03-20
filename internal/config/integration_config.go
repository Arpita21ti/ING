package config

import (
	"fmt"
	"time"
)

// IntegrationConfig contains configuration for external service integrations
type IntegrationConfig struct {
	Email       EmailConfig
	SMS         SMSConfig
	Storage     StorageConfig
	Monitoring  MonitoringConfig
	ExternalAPI ExternalAPIConfig
}

// EmailConfig contains email service configuration
type EmailConfig struct {
	Provider          string
	SenderEmail       string
	SenderName        string
	APIKey            string
	TemplateDirectory string
	MaxRetries        int
	RetryInterval     time.Duration
	Enabled           bool
}

// SMSConfig contains SMS service configuration
type SMSConfig struct {
	Provider      string
	AccountSID    string
	AuthToken     string
	PhoneNumber   string
	MaxRetries    int
	RetryInterval time.Duration
	Enabled       bool
}

// StorageConfig contains file storage configuration
type StorageConfig struct {
	Provider   string // "s3", "local", etc.
	BucketName string
	Region     string
	BasePath   string
	Enabled    bool
}

// MonitoringConfig contains monitoring and logging configuration
type MonitoringConfig struct {
	Provider       string // "cloudwatch", "datadog", etc.
	APIKey         string
	FlushInterval  time.Duration
	SamplingRate   float64
	EnabledMetrics []string
	Enabled        bool
}

// ExternalAPIConfig contains configuration for external API integrations
type ExternalAPIConfig struct {
	BaseURL      string
	APIKey       string
	Timeout      time.Duration
	MaxRetries   int
	RetryBackoff time.Duration
	Enabled      bool
}

// loadIntegrationConfig initializes integration configurations
func loadIntegrationConfig() (*IntegrationConfig, error) {
	// Load environment to check if in production
	env, err := loadEnvironment()
	if err != nil {
		return nil, fmt.Errorf("failed to load environment for integration config: %w", err)
	}

	// Get credentials for API keys if needed
	creds, err := loadCredentials()
	if err != nil {
		return nil, fmt.Errorf("failed to load credentials for integration config: %w", err)
	}

	// Email configuration
	emailConfig := EmailConfig{
		Provider:          getEnv("EMAIL_PROVIDER", "ses"),
		SenderEmail:       getEnv("EMAIL_SENDER", "no-reply@tnprgpv.com"),
		SenderName:        getEnv("EMAIL_SENDER_NAME", "TNP RGPV"),
		APIKey:            getAPIKey(creds, "email_provider", getEnv("EMAIL_API_KEY", "")),
		TemplateDirectory: getEnv("EMAIL_TEMPLATE_DIR", "./templates/email"),
		MaxRetries:        getEnvAsInt("EMAIL_MAX_RETRIES", 3),
		RetryInterval:     time.Duration(getEnvAsInt("EMAIL_RETRY_INTERVAL", 5)) * time.Second,
		Enabled:           getEnvAsBool("EMAIL_ENABLED", true),
	}

	// SMS configuration
	smsConfig := SMSConfig{
		Provider:      getEnv("SMS_PROVIDER", "twilio"),
		AccountSID:    getEnv("SMS_ACCOUNT_SID", ""),
		AuthToken:     getAPIKey(creds, "sms_provider", getEnv("SMS_AUTH_TOKEN", "")),
		PhoneNumber:   getEnv("SMS_PHONE_NUMBER", ""),
		MaxRetries:    getEnvAsInt("SMS_MAX_RETRIES", 3),
		RetryInterval: time.Duration(getEnvAsInt("SMS_RETRY_INTERVAL", 5)) * time.Second,
		Enabled:       getEnvAsBool("SMS_ENABLED", env.Production),
	}

	// Storage configuration
	storageConfig := StorageConfig{
		Provider:   getEnv("STORAGE_PROVIDER", "s3"),
		BucketName: getEnv("STORAGE_BUCKET_NAME", "tnp-rgpv-files"),
		Region:     getEnv("AWS_REGION", "us-east-1"),
		BasePath:   getEnv("STORAGE_BASE_PATH", "uploads"),
		Enabled:    getEnvAsBool("STORAGE_ENABLED", true),
	}

	// Monitoring configuration
	monitoringConfig := MonitoringConfig{
		Provider:      getEnv("MONITORING_PROVIDER", "cloudwatch"),
		APIKey:        getAPIKey(creds, "monitoring", getEnv("MONITORING_API_KEY", "")),
		FlushInterval: time.Duration(getEnvAsInt("MONITORING_FLUSH_INTERVAL", 10)) * time.Second,
		SamplingRate:  float64(getEnvAsInt("MONITORING_SAMPLING_RATE", 100)) / 100.0,
		EnabledMetrics: getEnvAsSlice(
			"MONITORING_ENABLED_METRICS",
			[]string{"api.requests", "db.queries", "errors"},
			",",
		),
		Enabled: getEnvAsBool("MONITORING_ENABLED", env.Production),
	}

	// External API configuration
	externalAPIConfig := ExternalAPIConfig{
		BaseURL:      getEnv("EXTERNAL_API_BASE_URL", "https://api.example.com"),
		APIKey:       getAPIKey(creds, "external_service", getEnv("EXTERNAL_API_KEY", "")),
		Timeout:      time.Duration(getEnvAsInt("EXTERNAL_API_TIMEOUT", 30)) * time.Second,
		MaxRetries:   getEnvAsInt("EXTERNAL_API_MAX_RETRIES", 3),
		RetryBackoff: time.Duration(getEnvAsInt("EXTERNAL_API_RETRY_BACKOFF", 5)) * time.Second,
		Enabled:      getEnvAsBool("EXTERNAL_API_ENABLED", false),
	}

	return &IntegrationConfig{
		Email:       emailConfig,
		SMS:         smsConfig,
		Storage:     storageConfig,
		Monitoring:  monitoringConfig,
		ExternalAPI: externalAPIConfig,
	}, nil
}

// getAPIKey retrieves an API key from the credentials or falls back to a provided default
func getAPIKey(creds *Credentials, keyName string, defaultValue string) string {
	if creds != nil && creds.APIKeys != nil {
		if key, exists := creds.APIKeys[keyName]; exists && key != "" {
			return key
		}
	}
	return defaultValue
}