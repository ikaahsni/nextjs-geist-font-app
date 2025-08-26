package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	JiraBaseURL    string
	JiraUsername   string
	JiraAPIToken   string
	JiraProjectKey string
	Port           string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	config := &Config{
		JiraBaseURL:    getEnvOrDefault("JIRA_BASE_URL", ""),
		JiraUsername:   getEnvOrDefault("JIRA_USERNAME", ""),
		JiraAPIToken:   getEnvOrDefault("JIRA_API_TOKEN", ""),
		JiraProjectKey: getEnvOrDefault("JIRA_PROJECT_KEY", ""),
		Port:           getEnvOrDefault("PORT", "8080"),
	}

	// Validate required configuration
	if config.JiraBaseURL == "" {
		return nil, fmt.Errorf("JIRA_BASE_URL is required")
	}
	if config.JiraUsername == "" {
		return nil, fmt.Errorf("JIRA_USERNAME is required")
	}
	if config.JiraAPIToken == "" {
		return nil, fmt.Errorf("JIRA_API_TOKEN is required")
	}
	if config.JiraProjectKey == "" {
		return nil, fmt.Errorf("JIRA_PROJECT_KEY is required")
	}

	return config, nil
}

// getEnvOrDefault gets environment variable or returns default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// ValidateConfig validates the configuration and logs warnings for demo credentials
func (c *Config) ValidateConfig() {
	if c.JiraUsername == "demo_user" || c.JiraAPIToken == "demo_token_replace_with_actual" {
		log.Println("⚠️  WARNING: You are using demo credentials!")
		log.Println("⚠️  Please copy .env.sample to .env and update with your actual Jira credentials")
		log.Println("⚠️  The application will not work with real Jira integration until you provide valid credentials")
	}

	log.Printf("✅ Configuration loaded successfully")
	log.Printf("   Jira Base URL: %s", c.JiraBaseURL)
	log.Printf("   Jira Project Key: %s", c.JiraProjectKey)
	log.Printf("   Server Port: %s", c.Port)
}
