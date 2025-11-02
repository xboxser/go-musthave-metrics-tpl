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
	"metrics/internal/config/db"
	"metrics/internal/handler/middleware"
	"metrics/internal/hash"
	models "metrics/internal/model"
	"metrics/internal/service"
	"metrics/internal/storage"
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

type serverHandler struct {
	service *service.ServerService
	config  *config.ConfigServer
	m       *middleware.RequestMiddleware
	file    *storage.FileJSON
	db      *db.DB
	hasher  hash.Hasher
	event   *audit.Event
}

func newServerHandler() (*serverHandler, error) {
	config := config.NewConfigServer()
	event := new(audit.Event)
	event.Register(audit.NewFileSubscriber(config.AuditFile))
	event.Register(audit.NewURLSubscriber(config.AuditURL))

	return &serverHandler{
		config: config,
		m:      middleware.NewRequestMiddleware(),
		hasher: nil,
		event:  event,
	}, nil
}

func Run(service *service.ServerService) {

	h, err := newServerHandler()

	if err != nil {
		panic(err)
	}

	h.addService(service)

	file, err := storage.NewFileJSON(h.config.FileStoragePath)
	if err != nil {
		panic(err)
	}
	h.addFile(file)
	h.addHasher(h.config.KEY)

	defer h.file.Close()

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

func (h *serverHandler) startServer() error {

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
	h.file.Save(h.service.GetModels())
	// Пытаемся корректно завершить сервер
	ctxStopServer, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctxStopServer); err != nil {
		log.Printf("Ошибка при завершении сервера: %v\n", err)
		return err
	}
	return nil
}

func (h *serverHandler) addService(service *service.ServerService) {
	h.service = service
}

func (h *serverHandler) addFile(file *storage.FileJSON) {
	h.file = file
}

func (h *serverHandler) addHasher(key string) {
	h.hasher = hash.NewSHA256(key)
}

// регистрируем роуты
func (h *serverHandler) registerRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.GzipMiddleware)
	r.Route("/update", func(r chi.Router) {
		r.Post("/", h.m.WithLogging(h.updateJSON))
		r.Post("/{type}/{name}/{value}", h.m.WithLogging(h.update))
	})
	r.Route("/updates", func(r chi.Router) {
		r.Post("/", h.m.WithLogging(h.updateBatchJSON))
	})
	r.Route("/value", func(r chi.Router) {
		r.Post("/", h.m.WithLogging(h.valueJSON))
		r.Get("/{type}/{name}", h.m.WithLogging(h.value))
	})

	r.Get("/ping", h.m.WithLogging(h.ping))

	r.Get("/", h.m.WithLogging(h.main))
	return r
}

func (h *serverHandler) connectDB(ctx context.Context) error {
	if h.config.DateBaseDSN == "" {
		return nil
	}
	db, err := db.NewDB(ctx, h.config.DateBaseDSN)
	if err != nil {
		return err
	}
	h.db = db
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
func (h *serverHandler) save() {
	if h.saveToDB() {
		return
	}
	h.saveToFile()
}

func (h *serverHandler) saveToDB() bool {
	if h.config.DateBaseDSN == "" || !h.db.Ping() {
		return false
	}
	err := h.db.SaveAll(h.service.GetModels())
	if err != nil {
		log.Printf("Ошибка при записи в БД: %v\n", err)
	} else {
		return true
	}

	return false
}

func (h *serverHandler) saveToFile() {
	err := h.file.Save(h.service.GetModels())
	if err != nil {
		log.Printf("Ошибка при записи в файл: %v\n", err)
	}
}

func (h *serverHandler) read() {
	if !h.config.Restore {
		return
	}

	if h.readFromDB() {
		return
	}

	h.readFromFile()
}

func (h *serverHandler) readFromDB() bool {
	if h.config.DateBaseDSN == "" || !h.db.Ping() {
		return false
	}

	m, err := h.db.ReadAll()
	if err != nil {
		log.Println("Не удалось получить информацию из БД", err)
		return false
	}

	h.service.SetModel(m)
	return true
}

func (h *serverHandler) readFromFile() {
	m, err := h.file.Read()
	if err != nil {
		log.Println("Не удалось прочитать файл", err)
		return
	}

	h.service.SetModel(*m)
}
func (h *serverHandler) updateBatchJSON(res http.ResponseWriter, req *http.Request) {
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
		IP_address: host,
		Metrics:    IDMetrics,
		TS:         int(time.Now().Unix()),
	}
	h.event.Update(audit)

	res.WriteHeader(http.StatusOK)
}

func (h *serverHandler) updateJSON(res http.ResponseWriter, req *http.Request) {
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

	// десериализуем JSON в Visitor
	if err = json.Unmarshal(buf.Bytes(), &metrics); err != nil {
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
		IP_address: host,
		Metrics:    []string{metrics.ID},
		TS:         int(time.Now().Unix()),
	}
	h.event.Update(audit)
	res.WriteHeader(http.StatusOK)
	res.Write(resp)
}

func (h *serverHandler) addMetrics(metrics models.Metrics) error {
	if !validateTypeMetrics(metrics.MType) {
		return fmt.Errorf("incorrect metric type")
	}

	err := h.service.UpdateJSON(&metrics)
	if err != nil {
		return err
	}
	return nil
}

func (h *serverHandler) update(res http.ResponseWriter, req *http.Request) {

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
		IP_address: host,
		Metrics:    []string{name},
		TS:         int(time.Now().Unix()),
	}
	h.event.Update(audit)

	res.WriteHeader(http.StatusOK)
}

func (h *serverHandler) ping(res http.ResponseWriter, req *http.Request) {
	if h.db.Ping() {
		res.WriteHeader(http.StatusOK)
	} else {
		res.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *serverHandler) valueJSON(res http.ResponseWriter, req *http.Request) {
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

func (h *serverHandler) value(res http.ResponseWriter, req *http.Request) {

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

func (h *serverHandler) main(res http.ResponseWriter, req *http.Request) {

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
