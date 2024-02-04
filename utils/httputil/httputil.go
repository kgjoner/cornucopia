package httputil

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/kgjoner/cornucopia/helpers/normalizederr"
)

type HttpUtil struct {
	client  *http.Client
	baseUrl string
}

func New(baseUrl string) *HttpUtil {
	return &HttpUtil{
		client:  &http.Client{},
		baseUrl: baseUrl,
	}
}

type Options struct {
	Params  map[string]string
	Headers map[string]string
}

type Executer func(data any) (*http.Response, error)

func (u HttpUtil) Get(path string, opt *Options) Executer {
	req, err := http.NewRequest("GET", u.baseUrl+path, nil)
	if err != nil {
		panic(err)
	}

	if opt != nil {
		SetOptions(req, *opt)
	}

	return func(data any) (*http.Response, error) {
		return DoReq(u.client, req, data)
	}
}

func (u HttpUtil) Delete(path string, opt *Options) Executer {
	req, err := http.NewRequest("DELETE", u.baseUrl+path, nil)
	if err != nil {
		panic(err)
	}

	if opt != nil {
		SetOptions(req, *opt)
	}

	return func(data any) (*http.Response, error) {
		return DoReq(u.client, req, data)
	}
}

func (u HttpUtil) Post(path string, body map[string]string, opt *Options) Executer {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", u.baseUrl+path, bytes.NewBuffer(jsonBody))
	if err != nil {
		panic(err)
	}

	if opt != nil {
		SetOptions(req, *opt)
	}
	
	req.Header.Add("content-type", "application/json")
	return func(data any) (*http.Response, error) {
		return DoReq(u.client, req, data)
	}
}

func (u HttpUtil) Put(path string, body map[string]string, opt *Options) Executer {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("PUT", u.baseUrl+path, bytes.NewBuffer(jsonBody))
	if err != nil {
		panic(err)
	}

	if opt != nil {
		SetOptions(req, *opt)
	}

	req.Header.Add("content-type", "application/json")
	return func(data any) (*http.Response, error) {
		return DoReq(u.client, req, data)
	}
}

func (u HttpUtil) Patch(path string, body map[string]string, opt *Options) Executer {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("PATCH", u.baseUrl+path, bytes.NewBuffer(jsonBody))
	if err != nil {
		panic(err)
	}

	if opt != nil {
		SetOptions(req, *opt)
	}

	req.Header.Add("content-type", "application/json")
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
	defer res.Body.Close()
	if err != nil {
		return res, err
	}

	if res.StatusCode >= 400 {
		var bodyErr map[string]any
		json.NewDecoder(res.Body).Decode(&bodyErr)

		msg, _ := bodyErr["message"].(string)
		err := normalizederr.NewExternalError(msg, map[string]error{
			"RequestMethod":  errors.New(req.Method),
			"RequestUrl":     errors.New(req.URL.String()),
			"ResponseStatus": errors.New(fmt.Sprintln(res.StatusCode)),
			"ResponseBody":   errors.New(fmt.Sprintln(bodyErr)),
		})
		
		return res, err
	}

	json.NewDecoder(res.Body).Decode(&data)
	return res, nil
}