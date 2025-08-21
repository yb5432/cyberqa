// Package openai provides an OpenAI-compatible API interface.
package openai

import (
	"os"
)

// Config holds the configuration for the OpenAI client.
// It includes the API key, base URL, and system prompt.
type Config struct {
	// APIKey is the API key for the OpenAI API.
	// It is used to authenticate requests to the API.
	APIKey string

	// APIBase is the base URL for the OpenAI API.
	// If not set, it defaults to "https://api.openai.com/v1".
	APIBase string

	// SystemPrompt is the system prompt to use for chat completions.
	// If set, it will be prepended to the messages list as a system message.
	SystemPrompt string
}

// NewConfig creates a new Config with default values.
// It reads API_KEY and API_BASE from environment variables if available.
// The environment variables are:
// - OPENAI_API_KEY for the API key
// - OPENAI_API_BASE for the base URL
func NewConfig() *Config {
	return &Config{
		APIKey:  os.Getenv("OPENAI_API_KEY"),
		APIBase: os.Getenv("OPENAI_API_BASE"),
		// SystemPrompt can be set later or via environment variable if needed.
	}
}

// WithAPIKey sets the API key for the Config.
// This is the key used to authenticate requests to the OpenAI API.
func (c *Config) WithAPIKey(apiKey string) *Config {
	c.APIKey = apiKey
	return c
}

// WithAPIBase sets the API base URL for the Config.
// This is the base URL for the OpenAI API.
// If not set, it defaults to "https://api.openai.com/v1".
func (c *Config) WithAPIBase(apiBase string) *Config {
	c.APIBase = apiBase
	return c
}

// WithSystemPrompt sets the system prompt for the Config.
// This prompt will be prepended to the messages list as a system message.
func (c *Config) WithSystemPrompt(systemPrompt string) *Config {
	c.SystemPrompt = systemPrompt
	return c
}
