package main

import (
	"fmt"
	"os"
)

const (
	envBotID       = "SLOVNIK_BOT_ID"
	envAPIURL      = "SLOVNIK_API_URL"
	envWebhookHost = "SLOVNIK_WEBHOOK_HOST"
)

// Config represents configuration information
type Config struct {
	BotID      string
	SlovnikURL string
	WebhookURL string
}

// InitConfig initializes bot configuration
func InitConfig() (*Config, error) {
	botID, ok := os.LookupEnv(envBotID)
	if !ok {
		return nil, fmt.Errorf("%s is not set", envBotID)
	}

	slovnikURL, ok := os.LookupEnv(envAPIURL)
	if !ok {
		return nil, fmt.Errorf("%s is not set", envAPIURL)
	}

	webhookURL := ""
	webhookHost, ok := os.LookupEnv(envWebhookHost)
	if ok {
		webhookURL = fmt.Sprintf("%s/bot%s", webhookHost, botID)
	}

	config := Config{
		BotID:      botID,
		SlovnikURL: slovnikURL,
		WebhookURL: webhookURL,
	}

	return &config, nil
}

// IsWebhook returns true if webhook URL is set
func (c *Config) IsWebhook() bool {
	return len(c.WebhookURL) > 0
}
