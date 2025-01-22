package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	NotionToken string
	DatabaseID  string
	BotToken    string
	BotChatID   string
	BotThreadID string
}

func LoadConfig(filePath string) (*Config, error) {
	log.Printf("Загрузка конфигурации из файла: %s", filePath)
	err := godotenv.Load(filePath)
	if err != nil {
		log.Printf("Ошибка при загрузке .env файла: %v", err)
		return nil, err
	}
	log.Println("Файл .env загружен успешно")

	return &Config{
		NotionToken: os.Getenv("NOTION_TOKEN"),
		DatabaseID:  os.Getenv("DATABASE_ID"),
		BotToken:    os.Getenv("BOT_TOKEN"),
		BotChatID:   os.Getenv("BOT_CHAT_ID"),
		BotThreadID: os.Getenv("BOT_THREAD_ID"),
	}, nil
}
