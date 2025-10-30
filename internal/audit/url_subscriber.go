package audit

import (
	"fmt"
	"metrics/internal/audit/model"
)

// Реализует интерфейс Observer
type UrlSubscriber struct {
	url string
}

func NewUrlSubscriber(id string) *UrlSubscriber {
	return &UrlSubscriber{url: id}
}

func (f *UrlSubscriber) Update(audit model.Audit) {
	fmt.Printf("The %s UrlSubscriber is notified of the %v event\n", f.url, audit)
}
