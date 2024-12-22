package presenter

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/kgjoner/cornucopia/helpers/controller"
	"github.com/kgjoner/cornucopia/helpers/normalizederr"
	"github.com/kgjoner/cornucopia/services/media"
	"github.com/kgjoner/cornucopia/utils/structop"
	log "github.com/sirupsen/logrus"
)

func NewLogger(r *http.Request, data interface{}) *log.Entry {
	ctx := r.Context()

	actor := ctx.Value(controller.ActorKey)
	actorV := reflect.ValueOf(actor)
	actorMap := map[string]interface{}{}
	if actorV.IsValid() && !actorV.IsZero() {
		actorMap = structop.New(actor).Map()
	}

	if err, ok := data.(error); ok {
		input, ok := ctx.Value(controller.InputKey).(map[string]any)
		if ok {
			removePrivateInputs(input)
		}

		if len(actorMap) == 0 && input["Application"] != nil {
			actorMap["Id"] = fmt.Sprintf("Application[%v]", input["Application"])
		}

		if normErr, ok := err.(normalizederr.NormalizedError); ok {
			return log.WithFields(log.Fields{
				"Method": r.Method,
				"Path":   r.URL.Path,
				"Actor":  actorMap["Id"],
				"Input":  input,
				"Kind":   normErr.Kind,
				"Code":   normErr.Code,
				"Stack":  normErr.Stack,
			})
		} else {
			return log.WithFields(log.Fields{
				"Method": r.Method,
				"Path":   r.URL.Path,
				"Actor":  actorMap["Id"],
				"Input":  input,
				"Kind":   "Unexpected",
				"Code":   "0000001",
			})
		}
	}

	if data == 201 {
		return log.WithFields(log.Fields{
			"Method": r.Method,
			"Path":   r.URL.Path,
			"Actor":  actorMap["Id"],
			"Kind":   "Creation",
		})
	}

	return log.WithFields(log.Fields{
		"Method": r.Method,
		"Path":   r.URL.Path,
		"Actor":  actorMap["Id"],
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
			if id, exists := vmap["Id"]; exists {
				input[key] = id
			}
		}
	}
}

func containsAny(key string, slc []string,) bool {
	for _, str := range slc {
		if strings.Contains(key, str) {
			return true
		}
	}

	return false
}
