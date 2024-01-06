package httputil

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/kgjoner/cornucopia/helpers/normalizederr"
)

func DoReq(client *http.Client, req *http.Request) (*http.Response, error) {
	res, err := client.Do(req)
	if err != nil {
		return res, err
	}

	if res.StatusCode >= 400 {
		var bodyErr map[string]any
		json.NewDecoder(res.Body).Decode(&bodyErr)

		err := normalizederr.NewExternalError(bodyErr["message"].(string), map[string]error{
			"RequestMethod":  errors.New(req.Method),
			"RequestUrl":     errors.New(req.URL.String()),
			"ResponseStatus": errors.New(fmt.Sprintln(res.StatusCode)),
			"ResponseBody":   errors.New(fmt.Sprintln(bodyErr)),
		})

		return res, err
	}

	return res, nil
}
