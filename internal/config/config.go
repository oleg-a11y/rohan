package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	NotionAPIKey     string
	NotionDatabaseID string
	TelegramBotToken string
	TelegramChatID   string
	TelegramThreadID string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	return &Config{
		NotionAPIKey:     os.Getenv("NOTION_API_KEY"),
		NotionDatabaseID: os.Getenv("NOTION_DATABASE_ID"),
		TelegramBotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		TelegramChatID:   os.Getenv("TELEGRAM_CHAT_ID"),
		TelegramThreadID: os.Getenv("TELEGRAM_THREAD_ID"),
	}
}
