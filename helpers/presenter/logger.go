package presenter

import (
	"errors"
	"net/http"
	"reflect"
	"strings"

	"github.com/kgjoner/cornucopia/v2/helpers/apperr"
	"github.com/kgjoner/cornucopia/v2/helpers/controller"
	"github.com/kgjoner/cornucopia/v2/services/media"
	"github.com/kgjoner/cornucopia/v2/utils/structop"
	log "github.com/sirupsen/logrus"
)

type ctxKey string

const ActorLogKey = ctxKey("actor.log")

func NewLogger(r *http.Request, data interface{}) *log.Entry {
	ctx := r.Context()
	actor := ctx.Value(ActorLogKey)

	if err, ok := data.(error); ok {
		input, ok := ctx.Value(controller.InputKey).(map[string]any)
		if ok {
			removePrivateInputs(input)
		}

		var appErr *apperr.AppError
		if errors.As(err, &appErr) {
			return log.WithFields(log.Fields{
				"Method": r.Method,
				"Path":   r.URL.Path,
				"Actor":  actor,
				"Input":  input,
				"Kind":   appErr.Kind,
				"Code":   appErr.Code,
				// "Stack":  appErr.Stack,
			})
		} else {
			return log.WithFields(log.Fields{
				"Method": r.Method,
				"Path":   r.URL.Path,
				"Actor":  actor,
				"Input":  input,
				"Kind":   "Unknown",
				"Code":   "Unexpected",
			})
		}
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

func removePrivateInputs(input map[string]any) {
	privateInput := []string{
		"password",
		"secret",
		"paymentcard",
		"actor",
		"token",
	}

	for key, value := range input {
		normalizedKey := strings.ReplaceAll(strings.ToLower(key), "_", "")
		if containsAny(normalizedKey, privateInput) {
			delete(input, key)
		} else if _, ok := value.(media.Media); ok {
			delete(input, key)
		}

		if v := reflect.ValueOf(value); v.IsValid() && !v.IsZero() && v.Kind() == reflect.Struct {
			vmap := structop.New(value).Map()
			if id, exists := vmap["ID"]; exists {
				input[key] = id
			}
		}
	}
}

func containsAny(key string, slc []string) bool {
	for _, str := range slc {
		if strings.Contains(key, str) {
			return true
		}
	}

	return false
}
