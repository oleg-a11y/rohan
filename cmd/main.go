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
	log.Println("Запуск приложения...")

	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalf("Ошибка при загрузке конфигурации: %v", err)
	}
	log.Println("Конфигурация загружена успешно")

	notionService := service.NewNotionService(cfg.NotionToken, cfg.DatabaseID)
	telegramService := service.NewTelegramService(cfg.BotToken, cfg.BotChatID, cfg.BotThreadID)
	telegramHandler := handler.NewTelegramHandler(telegramService, notionService)

	c := cron.New()

	_, err = c.AddFunc("34 14 * * *", func() {
		log.Println("Запуск задачи SendNotionData")
		if err := telegramHandler.SendNotionData(); err != nil {
			log.Printf("Ошибка при отправке данных: %v", err)
		}
	})

	if err != nil {
		log.Fatalf("Ошибка при добавлении задачи в cron: %v", err)
	}

	_, err = c.AddFunc("* * * * *", func() {
		log.Println("Запуск задачи NotifyUpcomingInterview")
		if err := telegramHandler.NotifyUpcomingInterview(); err != nil {
			log.Printf("Ошибка при отправке уведомления о собеседовании: %v", err)
		}
	})

	if err != nil {
		log.Fatalf("Ошибка при добавлении задачи в cron: %v", err)
	}

	c.Start()

	log.Println("Сервер запущен на 0.0.0.0:8080")
	err = http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}

	select {}
}
