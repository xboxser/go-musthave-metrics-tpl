package model

type Audit struct {
	TS        int      `json:"ts"`         // /unix timestamp события
	Metrics   []string `json:"metrics"`    // наименование полученных метрик
	IPAddress string   `json:"ip_address"` // IP адрес входящего запроса
}
