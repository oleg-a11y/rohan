package handler

import (
	"log"
	"rohan/internal/service"
)

type TelegramHandler struct {
	TelegramService *service.TelegramService
	NotionService   *service.NotionService
}

func NewTelegramHandler(telegramService *service.TelegramService, notionService *service.NotionService) *TelegramHandler {
	log.Println("Создание нового TelegramHandler")
	return &TelegramHandler{
		TelegramService: telegramService,
		NotionService:   notionService,
	}
}

func (h *TelegramHandler) SendNotionData() error {
	log.Println("Отправка данных из Notion")
	err := h.TelegramService.SendNotionData(h.NotionService)
	if err != nil {
		log.Printf("Ошибка при отправке данных: %v", err)
		return err
	}
	log.Println("Данные успешно отправлены")
	return nil
}

func (h *TelegramHandler) NotifyUpcomingInterview() error {
	log.Println("Уведомление о предстоящем собеседовании")
	err := h.TelegramService.NotifyUpcomingInterview(h.NotionService)
	if err != nil {
		log.Printf("Ошибка при отправке уведомления: %v", err)
		return err
	}
	log.Println("Уведомление успешно отправлено")
	return nil
}
