package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"rohan/internal/config"
	"rohan/internal/logger"
	"sort"
	"time"
)

type NotionService struct {
	cfg *config.Config
	log *logger.Logger
}

func NewNotionService(cfg *config.Config, log *logger.Logger) *NotionService {
	return &NotionService{cfg: cfg, log: log}
}

type Interview struct {
	Telegram  string `json:"Telegram"`
	Date      string `json:"Date"`
	Company   string `json:"Company"`
	Stage     string `json:"Stage"`
	Streaming bool   `json:"Streaming"`
	URL       string `json:"URL"`
}

type NotionResponse struct {
	Results []struct {
		URL        string `json:"url"`
		Properties struct {
			Telegram struct {
				RichText []struct {
					PlainText string `json:"plain_text"`
				} `json:"rich_text"`
			} `json:"ТГ ваш"`
			Date struct {
				Date struct {
					Start string `json:"start"`
				} `json:"date"`
			} `json:"Дата"`
			Company struct {
				Relation []struct {
					ID string `json:"id"`
				} `json:"relation"`
			} `json:"Компания"`
			Stage struct {
				Title []struct {
					PlainText string `json:"plain_text"`
				} `json:"title"`
			} `json:"Этап"`
			Streaming struct {
				Checkbox bool `json:"checkbox"`
			} `json:"Буду стримить"`
		} `json:"properties"`
	} `json:"results"`
}

type NotionPageResponse struct {
	Properties struct {
		Название struct {
			Title []struct {
				PlainText string `json:"plain_text"`
			} `json:"title"`
		} `json:"Название"`
	} `json:"properties"`
}

func (s *NotionService) fetchCompanyContent(companyID string) (string, error) {
	url := fmt.Sprintf("https://api.notion.com/v1/pages/%s", companyID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		s.log.Error("Ошибка при создании запроса: " + err.Error())
		return "", err
	}

	req.Header.Add("Authorization", "Bearer "+s.cfg.NotionAPIKey)
	req.Header.Add("Notion-Version", "2022-06-28")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		s.log.Error("Ошибка при выполнении запроса: " + err.Error())
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.log.Error("Ошибка при чтении ответа: " + err.Error())
		return "", err
	}

	var pageResp NotionPageResponse
	if err := json.Unmarshal(body, &pageResp); err != nil {
		s.log.Error("Ошибка при разборе JSON: " + err.Error())
		return "", err
	}

	if len(pageResp.Properties.Название.Title) > 0 {
		return pageResp.Properties.Название.Title[0].PlainText, nil
	}

	return "Неизвестная компания", nil
}

func (s *NotionService) FetchInterviews(filter map[string]interface{}) ([]Interview, error) {
	url := fmt.Sprintf("https://api.notion.com/v1/databases/%s/query", s.cfg.NotionDatabaseID)

	requestBody := map[string]interface{}{
		"filter": filter,
	}

	reqBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		s.log.Error("Ошибка кодирования JSON: " + err.Error())
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		s.log.Error("Ошибка при создании запроса: " + err.Error())
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+s.cfg.NotionAPIKey)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Notion-Version", "2022-06-28")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		s.log.Error("Ошибка при выполнении запроса: " + err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.log.Error("Ошибка при чтении ответа: " + err.Error())
		return nil, err
	}

	var notionResp NotionResponse
	if err := json.Unmarshal(body, &notionResp); err != nil {
		s.log.Error("Ошибка при разборе JSON: " + err.Error())
		return nil, err
	}

	var interviews []Interview
	for _, result := range notionResp.Results {
		interview := Interview{}

		if len(result.Properties.Telegram.RichText) > 0 {
			telegram := result.Properties.Telegram.RichText[0].PlainText
			if telegram != "" && telegram[0] != '@' {
				telegram = "@" + telegram
			}
			interview.Telegram = telegram
		}

		rawDate := result.Properties.Date.Date.Start
		parsedTime, err := time.Parse(time.RFC3339, rawDate)
		if err != nil {
			s.log.Error("Ошибка при разборе даты: " + err.Error())
			continue
		}
		interview.Date = parsedTime.Format("15:04")

		if len(result.Properties.Stage.Title) > 0 {
			interview.Stage = result.Properties.Stage.Title[0].PlainText
		}

		if len(result.Properties.Company.Relation) > 0 {
			companyID := result.Properties.Company.Relation[0].ID
			companyContent, err := s.fetchCompanyContent(companyID)
			if err != nil {
				s.log.Error("Ошибка при получении компании: " + err.Error())
				companyContent = ""
			}
			interview.Company = companyContent
		}

		if interview.Telegram == "" || interview.Date == "" || interview.Company == "" || interview.Stage == "" {
			continue
		}

		if !result.Properties.Streaming.Checkbox {
			continue
		}

		interview.URL = result.URL

		interview.Streaming = true
		interviews = append(interviews, interview)
	}

	sort.Slice(interviews, func(i, j int) bool {
		return interviews[i].Date < interviews[j].Date
	})

	return interviews, nil
}
