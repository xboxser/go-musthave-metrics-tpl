package sender

import (
	"bytes"
	"compress/gzip"
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

	var compressedBuf bytes.Buffer
	gz := gzip.NewWriter(&compressedBuf)

	// Сжатие данных
	if _, err := gz.Write(json); err != nil {
		return fmt.Errorf("ошибка сжатия данных: %w", err)
	}
	if err := gz.Close(); err != nil {
		return fmt.Errorf("ошибка закрытия gzip writer: %w", err)
	}

	url := fmt.Sprintf("http://%s/updates/", *s.baseURL)
	req, err := http.NewRequest(http.MethodPost, url, &compressedBuf)

	if err != nil {
		s.sugar.Infoln(
			"failed to create request",
			"uri", url,
			"json", string(json),
			"err", err,
		)
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	response, err := s.client.Do(req)
	if err != nil {
		s.sugar.Infoln(
			"failed request",
			"uri", url,
			"json", string(json),
			"err", err,
		)
		return fmt.Errorf("failed request: %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		s.sugar.Infoln(
			"unexpected status",
			"uri", url,
			"json", string(json),
			"status", response.Status,
			"err", err,
		)
		return fmt.Errorf("unexpected status: %s, %s", response.Status, url)
	}
	s.sugar.Debugln(
		"uri", url,
		"json", string(json),
	)
	return nil
}
