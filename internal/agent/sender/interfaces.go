package sender

import (
	"net/http"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock_sender.go -package=mocks

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Logger интерфейс для логирования
type Logger interface {
	Infoln(args ...interface{})
	Debugln(args ...interface{})
}

// Hasher интерфейс для хеширования
type Hasher interface {
	StringHash(data []byte) string
}
