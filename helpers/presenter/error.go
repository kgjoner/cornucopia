package presenter

import (
	"encoding/json"
	"net/http"

	"github.com/kgjoner/cornucopia/helpers/normalizederr"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
)

var (
	ErrCounters = map[int]prometheus.Counter{
		401: promauto.NewCounter(prometheus.CounterOpts{
			Name: "401_error_count",
			Help: "The total number of unauthorized response errors",
		}),
		403: promauto.NewCounter(prometheus.CounterOpts{
			Name: "403_error_count",
			Help: "The total number of forbidden response errors",
		}),
		400: promauto.NewCounter(prometheus.CounterOpts{
			Name: "400_error_count",
			Help: "The total number of bad request response errors",
		}),
		422: promauto.NewCounter(prometheus.CounterOpts{
			Name: "422_error_count",
			Help: "The total number of validation response errors",
		}),
		500: promauto.NewCounter(prometheus.CounterOpts{
			Name: "500_error_count",
			Help: "The total number of internal server response errors",
		}),
		502: promauto.NewCounter(prometheus.CounterOpts{
			Name: "502_error_count",
			Help: "The total number of bad gateway response errors",
		}),
	}
)

func HttpError(err error, w http.ResponseWriter, r *http.Request) {
	var status int
	var logLevel log.Level
	if e, ok := err.(normalizederr.NormalizedError); ok {
		switch e.Kind {
		case "Unauthorized":
			status = http.StatusUnauthorized
			logLevel = log.WarnLevel
		case "FatalUnauthorized":
			status = http.StatusUnauthorized
			logLevel = log.FatalLevel
		case "Forbidden":
			status = http.StatusForbidden
			logLevel = log.WarnLevel
		case "Request":
			status = http.StatusBadRequest
			logLevel = log.WarnLevel
		case "Validation":
			status = http.StatusUnprocessableEntity
			logLevel = log.WarnLevel
		case "External":
			status = http.StatusBadGateway
			logLevel = log.ErrorLevel
		case "FatalInternal":
			status = http.StatusInternalServerError
			logLevel = log.FatalLevel
		default:
			status = http.StatusInternalServerError
			logLevel = log.ErrorLevel
		}
	} else {
		err = normalizederr.NormalizedError{Message: err.Error(), Kind: "Unexpected", Code: "0000001"}
		status = http.StatusInternalServerError
		logLevel = log.ErrorLevel
	}

	counter := ErrCounters[status]
	counter.Inc()

	NewLogger(r, err).Log(logLevel, err.Error())

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(err)
}
