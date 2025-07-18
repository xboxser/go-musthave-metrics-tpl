package sender

import (
	"bytes"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type Sender struct {
	baseURL *string
	client  *http.Client
	sugar   zap.SugaredLogger
}

func NewSender(baseURL *string) *Sender {
	var sugar zap.SugaredLogger
	logger, err := zap.NewDevelopment()
	if err != nil {
		// вызываем панику, если ошибка
		panic(err)
	}
	defer logger.Sync()
	sugar = *logger.Sugar()
	return &Sender{
		baseURL: baseURL,
		client:  &http.Client{},
		sugar:   sugar,
	}
}

func (s *Sender) SendRequest(json []byte) error {

	url := fmt.Sprintf("http://%s/update/", *s.baseURL)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(json))

	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	response, err := s.client.Do(req)
	if err != nil {
		s.sugar.Debugln(
			"error",
			"uri", url,
			"json", string(json),
		)
		return fmt.Errorf("failed request: %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s, %s", response.Status, url)
	}
	s.sugar.Infoln(
		"uri", url,
		"json", string(json),
	)
	return nil
}
