package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"rohan/internal/config"
	"rohan/internal/logger"
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
	Telegram string `json:"Telegram"`
	Date     string `json:"Date"`
	Company  string `json:"Company"`
	Stage    string `json:"Stage"`
}

type NotionResponse struct {
	Results []struct {
		Properties struct {
			Telegram struct {
				RichText []struct {
					PlainText string `json:"plain_text"`
				} `json:"rich_text"`
			} `json:"Tелеграм"`
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
		s.log.Error("Error creating request: " + err.Error())
		return "", err
	}

	req.Header.Add("Authorization", "Bearer "+s.cfg.NotionAPIKey)
	req.Header.Add("Notion-Version", "2022-06-28")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		s.log.Error("Error executing request: " + err.Error())
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.log.Error("Error reading response: " + err.Error())
		return "", err
	}

	var pageResp NotionPageResponse
	if err := json.Unmarshal(body, &pageResp); err != nil {
		s.log.Error("Error decoding JSON: " + err.Error())
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
		s.log.Error("Error encoding JSON: " + err.Error())
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		s.log.Error("Error creating request: " + err.Error())
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+s.cfg.NotionAPIKey)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Notion-Version", "2022-06-28")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		s.log.Error("Error executing request: " + err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.log.Error("Error reading response: " + err.Error())
		return nil, err
	}

	var notionResp NotionResponse
	if err := json.Unmarshal(body, &notionResp); err != nil {
		s.log.Error("Error decoding JSON: " + err.Error())
		return nil, err
	}

	var interviews []Interview
	for _, result := range notionResp.Results {
		interview := Interview{}

		if len(result.Properties.Telegram.RichText) > 0 {
			interview.Telegram = result.Properties.Telegram.RichText[0].PlainText
		}

		rawDate := result.Properties.Date.Date.Start
		parsedTime, err := time.Parse(time.RFC3339, rawDate)
		if err != nil {
			s.log.Error("Error parsing date: " + err.Error())
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
				s.log.Error("Error fetching company: " + err.Error())
				companyContent = ""
			}
			interview.Company = companyContent
		}

		interviews = append(interviews, interview)
	}

	return interviews, nil
}
