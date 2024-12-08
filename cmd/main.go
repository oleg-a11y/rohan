package main

import (
	"github.com/robfig/cron/v3"
	"log"
	"rohan/internal/handler"
	"rohan/internal/service"
)

func main() {
	notionService := service.NewNotionService()
	telegramService := service.NewTelegramService()
	telegramHandler := handler.NewTelegramHandler(telegramService, notionService)

	c := cron.New()

	_, err := c.AddFunc("32 16 * * *", func() {
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

	select {}
}
