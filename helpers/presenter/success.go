package presenter

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	Counter200 = promauto.NewCounter(prometheus.CounterOpts{
		Name: "api_200_success_count",
		Help: "The total number of ok responses",
	})
	Counter201 = promauto.NewCounter(prometheus.CounterOpts{
		Name: "api_201_success_count",
		Help: "The total number of created responses",
	})
)

type Success[T any] struct {
	Data T
}

type successResponse struct {
	Data any `json:"data"`
}

func HttpSuccess(data interface{}, w http.ResponseWriter, r *http.Request, status ...int) http.ResponseWriter {
	var statusCode int
	if len(status) == 0 {
		statusCode = http.StatusOK
	} else {
		statusCode = status[0]
	}

	switch statusCode {
	case 201:
		Counter201.Inc()
		NewLogger(r, 201).Info()
	default:
		Counter200.Inc()
	}

	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")

	var res any

	dataV := reflect.Indirect(reflect.ValueOf(data))
	if dataV.Kind() == reflect.Struct {
		field := dataV.FieldByName("Data")
		if field.IsValid() {
			res = data
		} else {
			res = successResponse{data}
		}
	} else {
		res = successResponse{data}
	}

	json.NewEncoder(w).Encode(res)

	return w
}
