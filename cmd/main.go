package main

import (
	"log"
	"net/http"
	"rohan/internal/config"
	"rohan/internal/handler"
	"rohan/internal/service"

	"github.com/robfig/cron/v3"
)

func main() {
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalf("Ошибка при загрузке конфигурации: %v", err)
	}

	notionService := service.NewNotionService(cfg.NotionToken, cfg.DatabaseID)
	telegramService := service.NewTelegramService(cfg.BotToken, cfg.BotChatID, cfg.BotThreadID)
	telegramHandler := handler.NewTelegramHandler(telegramService, notionService)

	c := cron.New()

	_, err = c.AddFunc("34 14 * * *", func() {
		if err := telegramHandler.SendNotionData(); err != nil {
			log.Printf("Ошибка при отправке данных: %v", err)
		}
	})

	if err != nil {
		log.Fatalf("Ошибка при добавлении задачи в cron: %v", err)
	}

	_, err = c.AddFunc("* * * * *", func() {
		if err := telegramHandler.NotifyUpcomingInterview(); err != nil {
			log.Printf("Ошибка при отправке уведомления о собеседовании: %v", err)
		}
	})

	if err != nil {
		log.Fatalf("Ошибка при добавлении задачи в cron: %v", err)
	}

	c.Start()

	err = http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}

	select {}
}
