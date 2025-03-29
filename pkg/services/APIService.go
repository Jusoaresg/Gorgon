package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type APIService struct {
	Url    string
	Client *http.Client
}

func NewAPIService(url string) (a *APIService) {
	return &APIService{
		Url: url,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (a *APIService) Get(endpoint string, response interface{}) error {
	url := fmt.Sprintf("%s%s", a.Url, endpoint)

	resp, err := a.Client.Get(url)
	if err != nil {
		return fmt.Errorf("Error while get request: %w", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("Error while decoding response body: %w", err)
	}

	return nil
}

func (a *APIService) Post(endpoint string, requestData interface{}, response interface{}) error {
	url := fmt.Sprintf("%s%s", a.Url, endpoint)

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return fmt.Errorf("Error while marshaling request data: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("Error while making POST request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.Client.Do(req)
	if err != nil {
		return fmt.Errorf("Error while making POST request: %w", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("Error while decoding response body: %w", err)
	}

	return nil
}
