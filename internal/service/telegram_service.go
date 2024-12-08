package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type TelegramService struct {
	BotToken    string
	ChatID      string
	BotThreadID string
}

func NewTelegramService() *TelegramService {
	return &TelegramService{
		BotToken:    os.Getenv("BOT_TOKEN"),
		ChatID:      os.Getenv("BOT_CHAT_ID"),
		BotThreadID: os.Getenv("BOT_THREAD_ID"),
	}
}

func (ts *TelegramService) SendMessage(message string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", ts.BotToken)
	msg := map[string]interface{}{
		"chat_id":           ts.ChatID,
		"text":              message,
		"message_thread_id": ts.BotThreadID,
	}
	msgBytes, _ := json.Marshal(msg)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(msgBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (ts *TelegramService) SendNotionData(notionService *NotionService) error {
	now := time.Now().Format("2006-01-02")
	filter := map[string]interface{}{
		"filter": map[string]interface{}{
			"property": "Date",
			"date": map[string]interface{}{
				"equals": now,
			},
		},
	}

	notionResponse, err := notionService.GetItemsWithFilter(filter)
	if err != nil {
		return err
	}

	var message string
	var hasInterviews bool

	for _, page := range notionResponse.Results {
		dateTime, err := time.Parse(time.RFC3339, page.Properties.Date.Date.Start)
		if err != nil {
			return err
		}

		company := "Нет компании"
		if len(page.Properties.Company.Title) > 0 {
			company = page.Properties.Company.Title[0].PlainText
		}

		stage := "Нет этапа"
		if page.Properties.Stage.Select.Name != "" {
			stage = page.Properties.Stage.Select.Name
		}

		creator := "Нет инициатора"
		if len(page.Properties.Creator.RichText) > 0 {
			creator = page.Properties.Creator.RichText[0].PlainText
		}

		message += fmt.Sprintf("У %s, Компания: %s, Время: %s, Этап: %s\n", creator, company, dateTime.Format("15:04"), stage)
		hasInterviews = true
	}

	if hasInterviews {
		message = "Доброе время суток! На сегодня запланированы собеседования:\n" + message
	} else {
		message = "На сегодня запланированных собеседований нет."
	}

	err = ts.SendMessage(message)
	if err != nil {
		return err
	}
	return nil
}

func (ts *TelegramService) NotifyUpcomingInterview(notionService *NotionService) error {
	now := time.Now()
	currentTime := now.Truncate(time.Minute).Add(10 * time.Minute)

	filter := map[string]interface{}{
		"filter": map[string]interface{}{
			"property": "Date",
			"date": map[string]interface{}{
				"equals": currentTime.Format(time.RFC3339),
			},
		},
	}

	notionResponse, err := notionService.GetItemsWithFilter(filter)
	if err != nil {
		return err
	}

	if len(notionResponse.Results) == 0 {
		return nil
	}

	for _, page := range notionResponse.Results {
		creator := "Нет инициатора"
		if len(page.Properties.Creator.RichText) > 0 {
			creator = page.Properties.Creator.RichText[0].PlainText
		}

		company := "Нет компании"
		if len(page.Properties.Company.Title) > 0 {
			company = page.Properties.Company.Title[0].PlainText
		}

		stage := "Нет этапа"
		if page.Properties.Stage.Select.Name != "" {
			stage = page.Properties.Stage.Select.Name
		}

		message := fmt.Sprintf("Уважаемый %s, у вас через 10 минут начнется собеседование в Компанию: %s, Этап: %s", creator, company, stage)
		err = ts.SendMessage(message)
		if err != nil {
			return err
		}
	}

	return nil
}
