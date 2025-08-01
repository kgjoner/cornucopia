package httputil

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/kgjoner/cornucopia/helpers/normalizederr"
)

type HTTPUtil struct {
	client         *http.Client
	baseURL        string
	defaultOptions *Options
}

func New(baseURL string) *HTTPUtil {
	return &HTTPUtil{
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
		baseURL: baseURL,
	}
}

type Options struct {
	Params  map[string]string
	Headers map[string]string
}

func (u *HTTPUtil) SetDefaultOptions(opt *Options) {
	u.defaultOptions = opt
}

type Executer func(data any) (*http.Response, error)

func (u HTTPUtil) Get(path string, opt *Options) Executer {
	return u.request("GET", path, nil, opt)
}

func (u HTTPUtil) Delete(path string, opt *Options) Executer {
	return u.request("DELETE", path, nil, opt)
}

func (u HTTPUtil) Post(path string, body map[string]any, opt *Options) Executer {
	return u.request("POST", path, body, opt)
}

func (u HTTPUtil) Put(path string, body map[string]any, opt *Options) Executer {
	return u.request("PUT", path, body, opt)
}

func (u HTTPUtil) Patch(path string, body map[string]any, opt *Options) Executer {
	return u.request("PATCH", path, body, opt)
}

func (u HTTPUtil) request(method string, path string, inputtedBody map[string]any, opt *Options) Executer {
	var body io.Reader = nil
	if inputtedBody != nil {
		jsonBody, err := json.Marshal(inputtedBody)
		if err != nil {
			panic(err)
		}
		body = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, u.baseURL+path, body)
	if err != nil {
		panic(err)
	}

	if u.defaultOptions != nil {
		SetOptions(req, *u.defaultOptions)
	}

	if opt != nil {
		SetOptions(req, *opt)
	}

	if body != nil {
		req.Header.Add("content-type", "application/json")
	}

	return func(data any) (*http.Response, error) {
		return DoReq(u.client, req, data)
	}
}

func SetOptions(req *http.Request, opt Options) {
	if opt.Headers != nil {
		for k, v := range opt.Headers {
			req.Header.Add(k, v)
		}
	}

	if opt.Params != nil {
		q := req.URL.Query()
		for k, v := range opt.Params {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}
}

func DoReq(client *http.Client, req *http.Request, data any) (*http.Response, error) {
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		var bodyErr map[string]any
		json.NewDecoder(res.Body).Decode(&bodyErr)

		msg, _ := bodyErr["message"].(string)
		code, _ := bodyErr["code"].(string)
		err := normalizederr.NewExternalError(msg, map[string]error{
			"RequestMethod":  errors.New(req.Method),
			"RequestURL":     errors.New(req.URL.String()),
			"ResponseStatus": errors.New(fmt.Sprint(res.StatusCode)),
			"ResponseBody":   errors.New(fmt.Sprint(bodyErr)),
		}, code)

		return res, err
	}

	if data != nil {
		json.NewDecoder(res.Body).Decode(&data)
	}

	return res, nil
}
