package handler

import (
	"fmt"
	models "metrics/internal/model"
	"net/http"
	"strconv"
	"strings"
)

var mStorage = models.NewMemStorage()

func Run() {
	fmt.Println("Run server")
	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, update)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}

func getParamsUrl(path string) []string {
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

func validateTypeMetrics(params []string) bool {
	return params[1] != models.Counter && params[1] != models.Gauge
}

func valiteValueMetrics(value string) bool {
	_, err := strconv.Atoi(value)
	return err != nil
}

func update(res http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodPost {
		http.Error(res, "Use method POST", http.StatusMethodNotAllowed)
		return
	}

	params := getParamsUrl(req.URL.Path)

	if !validateCountParams(params) {
		http.Error(res, "Incorrect number of parameters", http.StatusMethodNotAllowed)
		return
	}

	if validateTypeMetrics(params) {
		http.Error(res, "incorrect metric type", http.StatusMethodNotAllowed)
		return
	}

	if valiteValueMetrics(params[3]) {
		http.Error(res, "incorrect value", http.StatusMethodNotAllowed)
		return
	}

	err := mStorage.Update(params[1], params[2], params[3])
	if err != nil {
		http.Error(res, err.Error(), http.StatusMethodNotAllowed)
	}

}
