package audit

import (
	"bytes"
	"encoding/json"
	"metrics/internal/audit/model"
	"net/http"
)

// Реализует интерфейс Observer
type UrlSubscriber struct {
	url string
}

func NewUrlSubscriber(url string) *UrlSubscriber {
	return &UrlSubscriber{url: url}
}

func (f *UrlSubscriber) Update(audit model.Audit) {
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
