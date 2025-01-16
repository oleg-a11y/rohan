package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"rohan/internal/model"
)

type NotionService struct {
	NotionToken string
	DatabaseId  string
}

func NewNotionService(notionToken, databaseId string) *NotionService {
	return &NotionService{
		NotionToken: notionToken,
		DatabaseId:  databaseId,
	}
}

func (ns *NotionService) GetItemsWithFilter(filter map[string]interface{}) (*model.NotionResponse, error) {
	filterJSON, err := json.Marshal(filter)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.notion.com/v1/databases/%s/query", ns.DatabaseId), bytes.NewBuffer(filterJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+ns.NotionToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Notion-Version", "2022-06-28")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var notionResponse model.NotionResponse
	if err := json.NewDecoder(resp.Body).Decode(&notionResponse); err != nil {
		return nil, err
	}

	return &notionResponse, nil
}
