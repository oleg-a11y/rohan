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
	return &TelegramHandler{
		TelegramService: telegramService,
		NotionService:   notionService,
	}
}

func (h *TelegramHandler) SendNotionData() error {
	err := h.TelegramService.SendNotionData(h.NotionService)
	if err != nil {
		log.Printf("Ошибка при отправке данных: %v", err)
		return err
	}
	return nil
}

func (h *TelegramHandler) NotifyUpcomingInterview() error {
	err := h.TelegramService.NotifyUpcomingInterview(h.NotionService)
	if err != nil {
		log.Printf("Ошибка при отправке уведомления: %v", err)
		return err
	}
	return nil
}
