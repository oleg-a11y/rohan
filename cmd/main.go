package main

import (
	"log"
	"net/http"
	"time"

	"github.com/robfig/cron/v3"
	"rohan/internal/handler"
	"rohan/internal/service"
)

func main() {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Fatal(err)
	}
	time.Local = loc

	notionService := service.NewNotionService()
	telegramService := service.NewTelegramService()
	telegramHandler := handler.NewTelegramHandler(telegramService, notionService)

	c := cron.New()

	_, err = c.AddFunc("00 10 * * *", func() {
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
