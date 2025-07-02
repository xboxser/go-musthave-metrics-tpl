package handler

import (
	"fmt"
	models "metrics/internal/model"
	"metrics/internal/service"
	"net/http"
	"strconv"
	"strings"
)

type serverHandler struct {
	service *service.ServerService
}

func newServerHandler(service *service.ServerService) *serverHandler {
	return &serverHandler{
		service: service,
	}
}

func Run(service *service.ServerService) {
	fmt.Println("Run server")
	h := newServerHandler(service)

	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, h.update)

	err := http.ListenAndServe(`:8080`, mux)
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

func validateCountParams(params []string) bool {
	return len(params) == 4
}

func validate404(params []string) bool {
	return !(len(params) < 4)
}

func validateTypeMetrics(params []string) bool {
	return params[1] == models.Counter || params[1] == models.Gauge
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

	params := getParamsURL(req.URL.Path)

	if !validate404(params) {
		http.Error(res, "Not found", http.StatusNotFound)
		return
	}

	if !validateCountParams(params) {
		http.Error(res, "Incorrect number of parameters", http.StatusMethodNotAllowed)
		return
	}

	if !validateTypeMetrics(params) {
		http.Error(res, "incorrect metric type", http.StatusBadRequest)
		return
	}

	if !valiteValueMetrics(params[3]) {
		http.Error(res, "incorrect value", http.StatusBadRequest)
		return
	}

	err := h.service.Update(params[1], params[2], params[3])
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}

}
