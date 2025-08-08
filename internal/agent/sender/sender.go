package sender

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"net/http"
	"time"

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

	retryIntervals := []time.Duration{0, 1 * time.Second, 3 * time.Second, 5 * time.Second}
	statusCode := 0
	for _, retryInterval := range retryIntervals {
		if retryInterval > 0 {
			s.sugar.Infoln("Повторная отправка данных через", retryInterval)
			time.Sleep(retryInterval)

		}
		statusCode, err := s.Send(compressedBuf)
		if statusCode == http.StatusOK {
			break
		}

		s.sugar.Infoln(
			"unexpected status",
			"json", string(json),
			"status", statusCode,
			"err", err,
		)

	}

	if statusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %v", statusCode)
	}

	s.sugar.Debugln(
		"json", string(json),
	)
	return nil
}

func (s *Sender) Send(compressedBuf bytes.Buffer) (int, error) {
	url := fmt.Sprintf("http://%s/updates/", *s.baseURL)
	req, err := http.NewRequest(http.MethodPost, url, &compressedBuf)

	if err != nil {
		s.sugar.Infoln(
			"failed to create request",
			"uri", url,
			"err", err,
		)
		return 0, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	response, err := s.client.Do(req)
	if err != nil {
		s.sugar.Infoln(
			"failed request",
			"uri", url,
			"err", err,
		)
		return 0, fmt.Errorf("failed request: %v", err)
	}

	defer response.Body.Close()
	return response.StatusCode, nil
}
