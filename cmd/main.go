package main

import (
	"github.com/robfig/cron/v3"
	"log"
	"net/http"
	"os"
	"rohan/internal/handler"
	"rohan/internal/service"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	go func() {
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatalf("Ошибка при запуске сервера: %v", err)
		}
	}()

	notionService := service.NewNotionService()
	telegramService := service.NewTelegramService()
	telegramHandler := handler.NewTelegramHandler(telegramService, notionService)

	c := cron.New()

	_, err := c.AddFunc("00 10 * * *", func() {
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
