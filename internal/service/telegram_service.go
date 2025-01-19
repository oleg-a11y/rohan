package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type TelegramService struct {
	BotToken    string
	ChatID      string
	BotThreadID string
}

func NewTelegramService(botToken, chatID, botThreadID string) *TelegramService {
	return &TelegramService{
		BotToken:    botToken,
		ChatID:      chatID,
		BotThreadID: botThreadID,
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

	message += "Cобеседования на сегодня\n\n"

	for _, page := range notionResponse.Results {
		dateTime, err := time.Parse(time.RFC3339, page.Properties.Date.Date.Start)
		if err != nil {
			return err
		}

		company := "Нет компании"
		if len(page.Properties.Company.RichText) > 0 {
			company = page.Properties.Company.RichText[0].PlainText
		}

		stage := "Нет этапа"
		if page.Properties.Stage.Select.Name != "" {
			stage = page.Properties.Stage.Select.Name
		}

		creator := "Нет инициатора"
		if len(page.Properties.Telegram.RichText) > 0 {
			creator = page.Properties.Telegram.RichText[0].PlainText
		}

		message += fmt.Sprintf(
			"%s\nКуда: %s\nКогда: %s\nКто: %s\n\n",
			stage, company, dateTime.Format("15:04"), creator,
		)

		hasInterviews = true
	}

	if hasInterviews {
		err = ts.SendMessage(message)
		if err != nil {
			return err
		}
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

	message := "Через 10 минут начнется собеседование\n\n"

	for _, page := range notionResponse.Results {
		creator := "Нет инициатора"
		if len(page.Properties.Telegram.RichText) > 0 {
			creator = page.Properties.Telegram.RichText[0].PlainText
		}

		company := "Нет компании"
		if len(page.Properties.Company.RichText) > 0 {
			company = page.Properties.Company.RichText[0].PlainText
		}

		stage := "Нет этапа"
		if page.Properties.Stage.Select.Name != "" {
			stage = page.Properties.Stage.Select.Name
		}

		dateTime, err := time.Parse(time.RFC3339, page.Properties.Date.Date.Start)
		if err != nil {
			return err
		}

		message += fmt.Sprintf(
			"%s\nКуда: %s\nКогда: %s\nКто: %s\n\n",
			stage, company, dateTime.Format("15:04"), creator,
		)
	}

	err = ts.SendMessage(message)
	if err != nil {
		return err
	}

	return nil
}
