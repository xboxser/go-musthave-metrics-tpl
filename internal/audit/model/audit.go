package model

type Audit struct {
	TS         int      `json:"ts"`         // /unix timestamp события
	Metrics    []string `json:"metrics"`    // наименование полученных метрик
	IP_address string   `json:"ip_address"` // IP адрес входящего запроса
}
