package service

import (
	"fmt"
	"net/http"
	"net/url"
	"rohan/internal/config"
	"rohan/internal/logger"
	"strings"
)

type TelegramService struct {
	cfg *config.Config
	log *logger.Logger
}

func NewTelegramService(cfg *config.Config, log *logger.Logger) *TelegramService {
	return &TelegramService{cfg: cfg, log: log}
}

func (t *TelegramService) SendInterviews(interviews []Interview, isTenMinutes bool) error {
	var message strings.Builder

	if len(interviews) == 0 {
		if !isTenMinutes {
			message.WriteString("Список запланированных собеседований на сегодня на данный момент пуст\n")
		}
	} else {
		if isTenMinutes {
			message.WriteString("⏳ Через 10 минут начнется собеседование\n\n")
		} else {
			message.WriteString("📅 Собеседования на сегодня\n\n")
		}

		for _, interview := range interviews {
			message.WriteString(fmt.Sprintf(
				"🕒 Время: %s\n🏢 Куда: %s\n📌 Этап: %s\n👤 Телеграм: %s\n🔗 Подробнее: [Открыть Notion](%s)\n\n",
				interview.Date, interview.Company, interview.Stage, interview.Telegram, interview.URL,
			))
		}
	}

	return t.sendMessage(message.String())
}

func (t *TelegramService) sendMessage(message string) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.cfg.TelegramBotToken)
	data := url.Values{}
	data.Set("chat_id", t.cfg.TelegramChatID)
	data.Set("text", message)
	data.Set("parse_mode", "Markdown")

	if t.cfg.TelegramThreadID != "" {
		data.Set("message_thread_id", t.cfg.TelegramThreadID)
	}

	resp, err := http.PostForm(apiURL, data)
	if err != nil {
		t.log.Error("Error sending message to Telegram: " + err.Error())
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.log.Error("Error sending message to Telegram, status: " + resp.Status)
	}

	return nil
}
