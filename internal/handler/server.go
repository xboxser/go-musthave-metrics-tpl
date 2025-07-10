package handler

import (
	"fmt"
	"html/template"
	models "metrics/internal/model"
	"metrics/internal/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

type serverHandler struct {
	service *service.ServerService
	config  *configSever
}

func newServerHandler(service *service.ServerService, config *configSever) *serverHandler {
	return &serverHandler{
		service: service,
		config:  config,
	}
}

func Run(service *service.ServerService) {
	fmt.Println("Run server")

	congig := newConfigServer()
	h := newServerHandler(service, congig)

	r := chi.NewRouter()
	r.Get("/value/{type}/{name}", h.value)
	r.Post("/update/{type}/{name}/{value}", h.update)
	r.Get("/", h.main)

	err := http.ListenAndServe(h.config.PortSever, r)
	if err != nil {
		panic(err)
	}
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

func valiteValueMetrics(value string) bool {
	_, err := strconv.ParseFloat(value, 64)
	return err == nil
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

	if !valiteValueMetrics(val) {
		http.Error(res, "incorrect value", http.StatusBadRequest)
		return
	}

	err := h.service.Update(typeMetrics, name, val)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}

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
		<h1>List metricks:</h1>
		<ul>
			{{range $key, $value := .}}
			<li><strong>{{$key}}:</strong> {{$value}}</li>
			{{end}}
		</ul>
	</body>
	</html>`

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
