package audit

import (
	"bytes"
	"encoding/json"
	"metrics/internal/audit/model"
	"net/http"
)

// Реализует интерфейс Observer
type URLSubscriber struct {
	url string
}

func NewURLSubscriber(url string) *URLSubscriber {
	return &URLSubscriber{url: url}
}

func (f *URLSubscriber) Update(audit model.Audit) {
	if f.url == "" {
		return
	}
	resp, err := json.Marshal(audit)
	if err != nil {
		return
	}

	httpResp, err := http.Post(f.url, "application/json", bytes.NewBuffer(resp))
	if err != nil {
		return
	}
	defer httpResp.Body.Close()
}
