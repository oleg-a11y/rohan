package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"rohan/internal/model"
)

type NotionService struct {
	NotionToken string
	DatabaseId  string
}

func NewNotionService(notionToken, databaseId string) *NotionService {
	log.Println("Создание нового NotionService")
	return &NotionService{
		NotionToken: notionToken,
		DatabaseId:  databaseId,
	}
}

func (ns *NotionService) GetItemsWithFilter(filter map[string]interface{}) (*model.NotionResponse, error) {
	log.Printf("Получение элементов с фильтром: %v", filter)
	filterJSON, err := json.Marshal(filter)
	if err != nil {
		log.Printf("Ошибка при маршализации фильтра: %v", err)
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.notion.com/v1/databases/%s/query", ns.DatabaseId), bytes.NewBuffer(filterJSON))
	if err != nil {
		log.Printf("Ошибка при создании запроса: %v", err)
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+ns.NotionToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Notion-Version", "2022-06-28")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Ошибка при выполнении запроса: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	var notionResponse model.NotionResponse
	if err := json.NewDecoder(resp.Body).Decode(&notionResponse); err != nil {
		log.Printf("Ошибка при декодировании ответа: %v", err)
		return nil, err
	}

	log.Printf("Получено %d результатов", len(notionResponse.Results))
	return &notionResponse, nil
}
