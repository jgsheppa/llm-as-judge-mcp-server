package config

import (
	"errors"
	"os"
	"strings"
)

type Config struct {
	ProviderAPIKeys map[string]string
}

func Load(provider string) (*Config, error) {
	apiKeys := make(map[string]string)

	envKey := strings.ToUpper(provider) + "_API_KEY"
	if apiKey, ok := os.LookupEnv(envKey); ok {
		apiKeys[provider] = apiKey
	}

	if len(apiKeys) == 0 {
		return nil, errors.New("no API keys found. Set provider-specific keys (ANTHROPIC_API_KEY, GEMINI_API_KEY, etc.)")
	}

	return &Config{
		ProviderAPIKeys: apiKeys,
	}, nil
}

func (c *Config) GetAPIKey(provider string) (string, bool) {
	key, ok := c.ProviderAPIKeys[provider]
	return key, ok
}
