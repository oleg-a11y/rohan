package main

import (
	"fmt"
	"net/http"
	"rohan/internal/config"
	"rohan/internal/logger"
	"rohan/internal/service"
)

func main() {
	cfg := config.LoadConfig()
	log, err := logger.NewLogger("app.log")
	if err != nil {
		fmt.Printf("Error creating logger: %v\n", err)
		return
	}
	defer log.Close()

	log.Info("Configuration loaded successfully")

	notionService := service.NewNotionService(cfg, log)
	telegramService := service.NewTelegramService(cfg, log)

	cronService := service.NewCronService(notionService, telegramService, log)
	cronService.Start()

	log.Info("Starting HTTP server on 0.0.0.0:8080")
	if err := http.ListenAndServe("0.0.0.0:8080", nil); err != nil {
		log.Error("Error starting HTTP server: " + err.Error())
	}

	select {}
}
