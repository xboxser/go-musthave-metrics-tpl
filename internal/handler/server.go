package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"metrics/internal/audit"
	modelAudit "metrics/internal/audit/model"
	"metrics/internal/config"
	"metrics/internal/handler/middleware"
	"metrics/internal/hash"
	models "metrics/internal/model"
	"metrics/internal/service"
	key_pair "metrics/internal/service/key_pair"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
)

type ServerHandler struct {
	service           *service.ServerService
	config            *config.ConfigServer
	m                 *middleware.RequestMiddleware
	hasher            hash.Hasher
	event             *audit.Event
	storage           *StorageManager
	cryptoCertificate *key_pair.PrivateKey
}

func NewServerHandler(config *config.ConfigServer) (*ServerHandler, error) {
	event := new(audit.Event)
	event.Register(audit.NewFileSubscriber(config.AuditFile))
	event.Register(audit.NewURLSubscriber(config.AuditURL))
	storage := NewStorageManager(config)
	return &ServerHandler{
		config:  config,
		m:       middleware.NewRequestMiddleware(),
		hasher:  nil,
		event:   event,
		storage: storage,
	}, nil
}

// Run - основной метод запуска обработчика сервера
func Run(service *service.ServerService) {
	config := config.NewConfigServer()
	h, err := NewServerHandler(config)

	if err != nil {
		panic(err)
	}

	h.addService(service)
	h.addHasher(h.config.KEY)

	if h.config.CryptoKeyPrivatePath != "" {
		err := h.addCryptoCertificate(h.config.CryptoKeyPrivatePath)
		if err != nil {
			panic(err)
		}
	}

	defer h.storage.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = h.connectDB(ctx)
	if err != nil {
		panic(err)
	}

	h.read()
	if h.config.IntervalSave > 0 {
		saveTicker := time.NewTicker(time.Duration(h.config.IntervalSave) * time.Second)
		defer saveTicker.Stop()
		go func() {
			for range saveTicker.C {
				h.save()
			}
		}()
	}

	h.startServer()

}

func (h *ServerHandler) startServer() error {

	r := h.registerRoutes()

	server := &http.Server{
		Addr:    h.config.Address,
		Handler: r,
	}

	// Канал для перехвата сигналов
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Запускаем сервер в отдельной горутине
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// Ждём сигнал остановки
	<-stop
	// записываем данные в файл
	h.storage.SaveToFile(h.service.GetModels())
	// Пытаемся корректно завершить сервер
	ctxStopServer, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctxStopServer); err != nil {
		log.Printf("Ошибка при завершении сервера: %v\n", err)
		return err
	}
	return nil
}

func (h *ServerHandler) addService(service *service.ServerService) {
	h.service = service
}

func (h *ServerHandler) addHasher(key string) {
	h.hasher = hash.NewSHA256(key)
}

func (h *ServerHandler) addCryptoCertificate(path string) error {
	cert, err := key_pair.NewPrivateKey(path)
	if err != nil {
		return err
	}
	h.cryptoCertificate = cert
	return nil
}

// регистрируем роуты
func (h *ServerHandler) registerRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.GzipMiddleware)

	// Регистрация маршрутов pprof
	// r.Mount("/debug/pprof", http.HandlerFunc(pprof.Index))
	// r.Mount("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	// r.Mount("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	// r.Mount("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	// r.Mount("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))

	r.Route("/update", func(r chi.Router) {
		r.Post("/", h.m.WithLogging(h.updateJSON))
		r.Post("/{type}/{name}/{value}", h.m.WithLogging(h.Update))
	})
	r.Route("/updates", func(r chi.Router) {
		r.Post("/", h.m.WithLogging(h.UpdateBatchJSON))
	})
	r.Route("/value", func(r chi.Router) {
		r.Post("/", h.m.WithLogging(h.ValueJSON))
		r.Get("/{type}/{name}", h.m.WithLogging(h.Value))
	})

	r.Get("/ping", h.m.WithLogging(h.Ping))

	r.Get("/", h.m.WithLogging(h.Main))
	return r
}

func (h *ServerHandler) connectDB(ctx context.Context) error {
	err := h.storage.ConnectDB(ctx)
	if err != nil {
		return err
	}
	return nil
}

func getParamsURL(path string) []string {
	params := strings.Split(path, "/")
	result := []string{}
	for _, v := range params {
		if v == "" {
			continue
		}
		result = append(result, v)
	}
	return result
}

func validateTypeMetrics(typeMode string) bool {
	return typeMode == models.Counter || typeMode == models.Gauge
}

func validateValueMetrics(value string) bool {
	_, err := strconv.ParseFloat(value, 64)
	return err == nil
}

// сохраняем информацию по метрикам
// Если не доступна БД, сохраняем в файл
func (h *ServerHandler) save() {
	if h.saveToDB() {
		return
	}
	h.saveToFile()
}

func (h *ServerHandler) saveToDB() bool {
	return h.storage.SaveToDB(h.service.GetModels())
}

func (h *ServerHandler) saveToFile() {
	h.storage.SaveToFile(h.service.GetModels())
}

func (h *ServerHandler) read() {
	if !h.config.Restore {
		return
	}

	if h.readFromDB() {
		return
	}

	h.readFromFile()
}

func (h *ServerHandler) readFromDB() bool {
	m, ok := h.storage.ReadFromDB()
	if !ok {
		return false
	}

	h.service.SetModel(m)
	return true
}

func (h *ServerHandler) readFromFile() {
	m := h.storage.ReadFromFile()
	h.service.SetModel(m)
}

// UpdateBatchJSON - обновлять несколько метрик одновременно, отправляя массив метрик в формате JSON.
func (h *ServerHandler) UpdateBatchJSON(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	if req.Method != http.MethodPost {
		http.Error(res, "Use method POST", http.StatusMethodNotAllowed)
		return
	}
	// читаем тело запроса
	var buf bytes.Buffer
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		log.Println("error read  body", err)
		return
	}

	data := buf.Bytes()

	if h.cryptoCertificate != nil {
		data, err = h.cryptoCertificate.Decrypt(data)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			log.Println("error decrypt body", err)
			return
		}
	}
	var metrics []models.Metrics
	// десериализуем JSON в Visitor
	if err = json.Unmarshal(data, &metrics); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		log.Println("error read  json", err)
		return
	}

	binaryHash, err := h.hasher.DecodeString(req.Header.Get("HashSHA256"))
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	if len(binaryHash) > 0 && !h.hasher.Compare(data, binaryHash) {
		http.Error(res, "Error HashSHA256", http.StatusBadRequest)
		return
	}

	IDMetrics := make([]string, len(metrics))
	for _, metric := range metrics {
		err := h.addMetrics(metric)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		IDMetrics = append(IDMetrics, metric.ID)
	}

	host, _, _ := net.SplitHostPort(req.RemoteAddr)
	audit := modelAudit.Audit{
		IPAddress: host,
		Metrics:   IDMetrics,
		TS:        int(time.Now().Unix()),
	}
	h.event.Update(audit)

	res.WriteHeader(http.StatusOK)
}

func (h *ServerHandler) updateJSON(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	if req.Method != http.MethodPost {
		http.Error(res, "Use method POST", http.StatusMethodNotAllowed)
		return
	}

	// читаем тело запроса
	var buf bytes.Buffer
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		log.Println("error read  body", err)
		return
	}
	var metrics models.Metrics
	data := buf.Bytes()

	if h.cryptoCertificate != nil {
		data, err = h.cryptoCertificate.Decrypt(data)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			log.Println("error decrypt body", err)
			return
		}
	}

	// десериализуем JSON в Visitor
	if err = json.Unmarshal(data, &metrics); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		log.Println("error read  json", err)
		return
	}

	err = h.addMetrics(metrics)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(metrics)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	host, _, _ := net.SplitHostPort(req.RemoteAddr)
	audit := modelAudit.Audit{
		IPAddress: host,
		Metrics:   []string{metrics.ID},
		TS:        int(time.Now().Unix()),
	}
	h.event.Update(audit)
	res.WriteHeader(http.StatusOK)
	res.Write(resp)
}

func (h *ServerHandler) addMetrics(metrics models.Metrics) error {
	if !validateTypeMetrics(metrics.MType) {
		return fmt.Errorf("incorrect metric type")
	}

	err := h.service.UpdateJSON(&metrics)
	if err != nil {
		return err
	}
	return nil
}

// Update - обновить значение конкретной метрики
func (h *ServerHandler) Update(res http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodPost {
		http.Error(res, "Use method POST", http.StatusMethodNotAllowed)
		return
	}

	typeMetrics := strings.ToLower(chi.URLParam(req, "type"))
	name := strings.ToLower(chi.URLParam(req, "name"))
	val := strings.ToLower(chi.URLParam(req, "value"))

	if !validateTypeMetrics(typeMetrics) {
		http.Error(res, "incorrect metric type", http.StatusBadRequest)
		return
	}

	if !validateValueMetrics(val) {
		http.Error(res, "incorrect value", http.StatusBadRequest)
		return
	}

	err := h.service.Update(typeMetrics, name, val)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}

	host, _, _ := net.SplitHostPort(req.RemoteAddr)
	audit := modelAudit.Audit{
		IPAddress: host,
		Metrics:   []string{name},
		TS:        int(time.Now().Unix()),
	}
	h.event.Update(audit)

	res.WriteHeader(http.StatusOK)
}

// Ping - проверка есть ли подключение к БД
func (h *ServerHandler) Ping(res http.ResponseWriter, req *http.Request) {
	if h.storage.db != nil && h.storage.db.Ping() {
		res.WriteHeader(http.StatusOK)
	} else {
		res.WriteHeader(http.StatusInternalServerError)
	}
}

// ValueJSON - возвращает значение метрики в виде JSON
func (h *ServerHandler) ValueJSON(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Use method POST", http.StatusMethodNotAllowed)
		return
	}

	// читаем тело запроса
	var buf bytes.Buffer
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	var metrics models.Metrics

	// десериализуем JSON в Visitor
	if err = json.Unmarshal(buf.Bytes(), &metrics); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if !validateTypeMetrics(metrics.MType) {
		http.Error(res, "incorrect metric type", http.StatusBadRequest)
		return
	}

	err = h.service.GetValueJSON(&metrics)
	if err != nil {
		http.Error(res, "not value", http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(metrics)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(resp)
}

// Value - получить значение метрики
func (h *ServerHandler) Value(res http.ResponseWriter, req *http.Request) {

	typeMetrics := strings.ToLower(chi.URLParam(req, "type"))
	name := strings.ToLower(chi.URLParam(req, "name"))

	if !validateTypeMetrics(typeMetrics) {
		http.Error(res, "incorrect metric type", http.StatusBadRequest)
		return
	}

	val, err := h.service.GetValue(typeMetrics, name)
	if err != nil {
		http.Error(res, "not value", http.StatusNotFound)
		return
	}

	res.Write([]byte(val))
}

// Main - точка вывода html страницы со всеми метриками
func (h *ServerHandler) Main(res http.ResponseWriter, req *http.Request) {

	data := h.service.GetAll()

	const tpl = `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Metrics server</title>
	</head>
	<body>
		<h1>List metrics:</h1>
		<ul>
			{{range $key, $value := .}}
			<li><strong>{{$key}}:</strong> {{$value}}</li>
			{{end}}
		</ul>
	</body>
	</html>`
	res.Header().Set("Content-Type", "text/html")
	res.WriteHeader(http.StatusOK)
	tmpl, err := template.New("page").Parse(tpl)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(res, data)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

}
