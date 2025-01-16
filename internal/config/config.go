package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	NotionToken string
	DatabaseID  string
	BotToken    string
	BotChatID   string
	BotThreadID string
}

func LoadConfig(filePath string) (*Config, error) {
	err := godotenv.Load(filePath)
	if err != nil {
		return nil, err
	}

	return &Config{
		NotionToken: os.Getenv("NOTION_TOKEN"),
		DatabaseID:  os.Getenv("DATABASE_ID"),
		BotToken:    os.Getenv("BOT_TOKEN"),
		BotChatID:   os.Getenv("BOT_CHAT_ID"),
		BotThreadID: os.Getenv("BOT_THREAD_ID"),
	}, nil
}
