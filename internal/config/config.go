package config

import (
	"errors"
	"os"
	"strings"
)

type Config struct {
	ProviderAPIKey string
}

func Load(provider string) (*Config, error) {

	envKey := strings.ToUpper(provider) + "_API_KEY"
	apiKey, ok := os.LookupEnv(envKey)
	if !ok {
		return nil, errors.New("no API keys found. Set provider-specific keys (ANTHROPIC_API_KEY, GEMINI_API_KEY, etc.)")
	}

	return &Config{
		ProviderAPIKey: apiKey,
	}, nil
}
