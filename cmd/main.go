package main

import (
	"github.com/robfig/cron/v3"
	"log"
	"net/http"
	"rohan/internal/handler"
	"rohan/internal/service"
	"time"
)

func main() {
	notionService := service.NewNotionService()
	telegramService := service.NewTelegramService()
	telegramHandler := handler.NewTelegramHandler(telegramService, notionService)

	c := cron.New(cron.WithLocation(time.FixedZone("MSK", 3*60*60)))

	_, err := c.AddFunc("0 10 * * *", func() {
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
