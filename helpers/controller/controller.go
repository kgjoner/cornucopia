package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/kgjoner/cornucopia/helpers/htypes"
	"github.com/kgjoner/cornucopia/helpers/normalizederr"
	"github.com/kgjoner/cornucopia/services/media"
	"github.com/kgjoner/cornucopia/utils/structop"
)

type Controller struct {
	req    *http.Request
	fields map[string]any
	err    error
}

func New(req *http.Request) *Controller {
	return &Controller{
		req:    req,
		fields: make(map[string]any),
	}
}

func (c *Controller) AddToken() *Controller {
	if c.err != nil {
		return c
	}

	token, ok := c.req.Context().Value("token").(string)
	if !ok {
		c.err = normalizederr.NewUnauthorizedError("Token required.")
		return c
	}

	c.fields["token"] = token
	return c
}

func (c *Controller) AddActor() *Controller {
	if c.err != nil {
		return c
	}

	actor := c.req.Context().Value("actor")
	if actor == nil {
		c.err = normalizederr.NewUnauthorizedError("Actor required.")
		return c
	}

	c.fields["actor"] = actor
	return c
}

func (c *Controller) AddApplication() *Controller {
	if c.err != nil {
		return c
	}

	application := c.req.Context().Value("application")
	if application == nil {
		c.err = normalizederr.NewUnauthorizedError("Application required.")
		return c
	}

	c.fields["application"] = application
	return c
}

// Receive actor and parse relevant fields. Inputted func must return array of tuples in form [key, value].
func (c *Controller) ParseActorAs(setFields func(actor any, fields map[string]any)) *Controller {
	if c.err != nil {
		return c
	}

	actor := c.req.Context().Value("actor")
	if actor == nil {
		c.err = normalizederr.NewUnauthorizedError("Actor required.")
		return c
	}

	setFields(actor, c.fields)
	return c
}

func (c *Controller) ParseMultipartForm(files, values []string, mediaService media.MediaService) *Controller {
	if c.err != nil {
		return c
	}

	err := c.req.ParseMultipartForm(32 << 20) // maxMemory 32MB
	if err != nil {
		c.err = err
		return c
	}

	for _, valueInputName := range values {
		c.fields[valueInputName] = c.req.PostFormValue(valueInputName)
	}

	for _, fileInputName := range files {
		file, _, err := c.req.FormFile(fileInputName)
		if err == http.ErrMissingFile {
			continue
		} else if err != nil {
			c.err = err
			return c
		}

		var bufPic bytes.Buffer
		io.Copy(&bufPic, file)

		c.fields[fileInputName] = media.New(&bufPic, mediaService)
	}

	return c
}

// Get param from URL string. Field is the name used in input; if omitted, it is used the same
// param name.
func (c *Controller) ParseUrlParam(param string, field ...string) *Controller {
	var fieldName string
	if len(field) == 0 {
		fieldName = param
	} else {
		fieldName = field[0]
	}

	c.fields[fieldName] = chi.URLParam(c.req, param)

	return c
}

// Get param from query string. Field is the name used in input; if omitted, it is used the same
// param name.
func (c *Controller) ParseQueryParam(param string, field ...string) *Controller {
	var fieldName string
	if len(field) == 0 {
		fieldName = param
	} else {
		fieldName = field[0]
	}

	query := c.req.URL.Query()
	qstr := query[fieldName]
	if len(qstr) != 0 {
		c.fields[fieldName] = qstr[0]
	}

	return c
}

func (c *Controller) ParseBody(fields ...string) *Controller {
	var bodyMap map[string]any
	json.NewDecoder(c.req.Body).Decode(&bodyMap)
	defer c.req.Body.Close()

	normalizedBodyMap := map[string]any{}
	for key, value := range bodyMap {
		normalizedBodyMap[strings.ToLower(key)] = value
	}

	for _, field := range fields {
		normalizedField := strings.ToLower(field)
		c.fields[normalizedField] = normalizedBodyMap[normalizedField]
	}

	return c
}

func (c *Controller) AddPagination() *Controller {
	query := c.req.URL.Query()

	limit := query["limit"]
	var parsedLimit int64
	if len(limit) != 0 {
		parsedLimit, _ = strconv.ParseInt(limit[0], 10, 32)
	}

	page := query["page"]
	var parsedPage int64
	if len(page) != 0 {
		parsedPage, _ = strconv.ParseInt(page[0], 10, 32)
	}

	c.fields["pagination"] = htypes.NewPagination(&htypes.PaginationCreationFields{
		Limit: int(parsedLimit),
		Page:  int(parsedPage),
	})

	return c
}

// Get header. Field is the name used in input; if omitted, it is used the same
// header name.
func (c *Controller) AddHeader(header string, field ...string) *Controller {
	var fieldName string
	if len(field) == 0 {
		fieldName = header
	} else {
		fieldName = field[0]
	}

	headerValue := c.req.Header.Get(header)

	c.fields[fieldName] = headerValue
	return c
}

func (c *Controller) AddIp() *Controller {
	c.fields["ip"] = c.req.RemoteAddr
	return c
}

func (c *Controller) AddMarket() *Controller {
	timezone := c.req.Header.Get("x-timezone")
	market, err := htypes.MarketByTimezone(timezone)
	if err != nil {
		c.err = err
		return c
	}

	c.fields["market"] = market
	return c
}

func (c *Controller) AddLanguages() *Controller {
	acptLang := c.req.Header.Get("accept-language")

	type LangQ struct {
		lang   string
		weight float64
	}

	var lqs []LangQ

	langQStrs := strings.Split(acptLang, ",")
	for _, langQStr := range langQStrs {
		trimedLangQStr := strings.Trim(langQStr, " ")

		langQ := strings.Split(trimedLangQStr, ";")
		if len(langQ) == 1 {
			lq := LangQ{langQ[0], 1}
			lqs = append(lqs, lq)
		} else {
			qp := strings.Split(langQ[1], "=")
			q, err := strconv.ParseFloat(qp[1], 64)
			if err != nil {
				c.err = err
			}

			lq := LangQ{langQ[0], q}
			lqs = append(lqs, lq)
		}
	}

	sort.SliceStable(lqs, func(i, j int) bool {
		return lqs[i].weight > lqs[j].weight
	})

	langs := make([]string, len(lqs))
	for i, lq := range lqs {
		langs[i] = strings.ToLower(lq.lang)
	}

	c.fields["languages"] = langs

	return c
}

func (c *Controller) Write(input any) error {
	ctx := c.req.Context()
	ctx = context.WithValue(ctx, "input", c.fields)
	*(c.req) = *c.req.WithContext(ctx)

	if c.err != nil {
		return c.err
	}

	err := structop.New(input).UpdateViaMap(c.fields)
	if err != nil {
		return err
	}

	return nil
}
