package sender

import (
	"fmt"
	"net/http"
)

type Sender struct {
	baseURL *string
}

func NewSender(baseURL *string) *Sender {
	return &Sender{
		baseURL: baseURL,
	}
}

func (s *Sender) SendRequest(metricType string, metricName string, metricValue string) error {

	url := fmt.Sprintf("http://%s/update/%s/%s/%s", *s.baseURL, metricType, metricName, metricValue)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("Content-Length", "0")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed request: %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s, %s", response.Status, url)
	}
	return nil
}
