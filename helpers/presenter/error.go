package presenter

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/kgjoner/cornucopia/helpers/apperr"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
)

var (
	ErrCounters = map[int]prometheus.Counter{
		401: promauto.NewCounter(prometheus.CounterOpts{
			Name: "api_401_error_count",
			Help: "The total number of unauthorized response errors",
		}),
		403: promauto.NewCounter(prometheus.CounterOpts{
			Name: "api_403_error_count",
			Help: "The total number of forbidden response errors",
		}),
		400: promauto.NewCounter(prometheus.CounterOpts{
			Name: "api_400_error_count",
			Help: "The total number of bad request response errors",
		}),
		408: promauto.NewCounter(prometheus.CounterOpts{
			Name: "api_408_error_count",
			Help: "The total number of request timeout response errors",
		}),
		409: promauto.NewCounter(prometheus.CounterOpts{
			Name: "api_409_error_count",
			Help: "The total number of conflict response errors",
		}),
		422: promauto.NewCounter(prometheus.CounterOpts{
			Name: "api_422_error_count",
			Help: "The total number of validation response errors",
		}),
		499: promauto.NewCounter(prometheus.CounterOpts{
			Name: "scapi_499_error_count",
			Help: "The total number of context canceled errors",
		}),
		500: promauto.NewCounter(prometheus.CounterOpts{
			Name: "api_500_error_count",
			Help: "The total number of internal server response errors",
		}),
		502: promauto.NewCounter(prometheus.CounterOpts{
			Name: "api_502_error_count",
			Help: "The total number of bad gateway response errors",
		}),
	}
)

func HTTPError(err error, w http.ResponseWriter, r *http.Request) {
	var status int
	var logLevel log.Level
	var e *apperr.AppError
	if errors.As(err, &e) {
		switch e.Kind {
		case apperr.Unauthorized:
			status = http.StatusUnauthorized
			logLevel = log.WarnLevel
		case apperr.Forbidden:
			status = http.StatusForbidden
			logLevel = log.WarnLevel
		case apperr.Request:
			status = http.StatusBadRequest
			logLevel = log.WarnLevel
		case apperr.Validation:
			status = http.StatusUnprocessableEntity
			logLevel = log.WarnLevel
		case apperr.Conflict:
			status = http.StatusConflict
			logLevel = log.ErrorLevel
		case apperr.External:
			if e.Code == apperr.Unexpected {
				status = http.StatusBadGateway
				logLevel = log.ErrorLevel
			} else {
				status = http.StatusBadRequest
				logLevel = log.WarnLevel
			}
		default:
			status = http.StatusInternalServerError
			logLevel = log.ErrorLevel
		}
	} else {
		if errors.Is(err, context.Canceled) {
			err = apperr.NewRequestError(err.Error(), "CONTEXT_CANCELED")
			status = 499
			logLevel = log.WarnLevel
		} else if errors.Is(err, context.DeadlineExceeded) {
			err = apperr.NewRequestError(err.Error(), "CONTEXT_TIMEOUT")
			status = http.StatusRequestTimeout
			logLevel = log.WarnLevel
		} else {
			err = apperr.NewInternalError(err.Error())
			status = http.StatusInternalServerError
			logLevel = log.ErrorLevel
		}
	}

	if apperr.IsFatal(err) {
		logLevel = log.FatalLevel
	}

	counter := ErrCounters[status]
	counter.Inc()

	NewLogger(r, err).Log(logLevel, err.Error())

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(err)
}
