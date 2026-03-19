package httpserver

import (
	"errors"
	"net/http"

	"github.com/kgjoner/cornucopia/v2/apperr"
	log "github.com/sirupsen/logrus"
)

type ctxKey string

const ActorLogKey = ctxKey("actor.log")

func NewLogger(r *http.Request, data interface{}) *log.Entry {
	actor := r.Context().Value(ActorLogKey)

	if err, ok := data.(error); ok {
		var appErr *apperr.AppError
		if errors.As(err, &appErr) {
			return log.WithFields(log.Fields{
				"Method": r.Method,
				"Path":   r.URL.Path,
				"Actor":  actor,
				"Kind":   appErr.Kind,
				"Code":   appErr.Code,
			})
		}
		return log.WithFields(log.Fields{
			"Method": r.Method,
			"Path":   r.URL.Path,
			"Actor":  actor,
			"Kind":   "Unknown",
			"Code":   "Unexpected",
		})
	}

	if data == 201 {
		return log.WithFields(log.Fields{
			"Method": r.Method,
			"Path":   r.URL.Path,
			"Actor":  actor,
			"Kind":   "Creation",
		})
	}

	return log.WithFields(log.Fields{
		"Method": r.Method,
		"Path":   r.URL.Path,
		"Actor":  actor,
	})
}
