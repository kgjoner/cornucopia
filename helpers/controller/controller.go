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
	"github.com/kgjoner/cornucopia/helpers/apperr"
	"github.com/kgjoner/cornucopia/helpers/htypes"
	"github.com/kgjoner/cornucopia/services/media"
	"github.com/kgjoner/cornucopia/utils/structop"
)

type oldCtxKey int

// Deprecated: You should declare your own keys.
// Then, use AddFromContext to add them to the controller.
const (
	ApplicationKey oldCtxKey = iota
	ActorKey
	TargetKey
	TokenKey
)

type ctxKey string

const InputKey = ctxKey("input")

type Controller struct {
	req         *http.Request
	fields      map[string]any
	hasJSONBody bool
	err         error
}

func New(req *http.Request) *Controller {
	return &Controller{
		req:    req,
		fields: make(map[string]any),
	}
}

func (c *Controller) AddFromContext(key any, field string) *Controller {
	if c.err != nil {
		return c
	}

	token := c.req.Context().Value(key)
	if token == nil {
		c.err = apperr.NewInternalError("field \"" + field + "\" should have been set")
		return c
	}

	c.fields[field] = token
	return c
}

// Deprecated: Use AddFromContext instead.
func (c *Controller) AddToken() *Controller {
	if c.err != nil {
		return c
	}

	token, ok := c.req.Context().Value(TokenKey).(string)
	if !ok {
		c.err = apperr.NewUnauthorizedError("token required")
		return c
	}

	c.fields["token"] = token
	return c
}

// Deprecated: Use AddFromContext instead.
func (c *Controller) AddActor() *Controller {
	if c.err != nil {
		return c
	}

	actor := c.req.Context().Value(ActorKey)
	if actor == nil {
		c.err = apperr.NewUnauthorizedError("Actor required.")
		return c
	}

	c.fields["actor"] = actor
	return c
}

// Deprecated: Use AddFromContext instead.
func (c *Controller) AddTarget() *Controller {
	if c.err != nil {
		return c
	}

	target := c.req.Context().Value(TargetKey)
	if target == nil {
		c.err = apperr.NewUnauthorizedError("Target required.")
		return c
	}

	c.fields["target"] = target
	return c
}

// Deprecated: Use AddFromContext instead.
func (c *Controller) AddApplication() *Controller {
	if c.err != nil {
		return c
	}

	application := c.req.Context().Value(ApplicationKey)
	if application == nil {
		c.err = apperr.NewUnauthorizedError("Application required.")
		return c
	}

	c.fields["application"] = application
	return c
}

// Deprecated: Should not be used. Delegate the processing of actor to usecase layer. For special cases,
// you may process actor outside controller and use CustomField to add the result.
//
// Receive actor and parse relevant fields. Inputted func must return array of tuples in form [key, value].
func (c *Controller) ParseActorAs(setFields func(actor any, fields map[string]any)) *Controller {
	if c.err != nil {
		return c
	}

	actor := c.req.Context().Value(ActorKey)
	if actor == nil {
		c.err = apperr.NewUnauthorizedError("Actor required.")
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
func (c *Controller) ParseURLParam(param string, field ...string) *Controller {
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
	qstr := query[param]
	if len(qstr) == 1 {
		c.fields[fieldName] = qstr[0]
	} else if len(qstr) > 1 {
		c.fields[fieldName] = qstr
	}

	return c
}

// Mark the controller as having a JSON body. This will unmarshall the body when Write is called.
func (c *Controller) JSONBody() *Controller {
	c.hasJSONBody = true
	return c
}

// Deprecated: Use JSONBody instead.
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

func (c *Controller) CustomField(field string, value any) *Controller {
	if c.err != nil {
		return c
	}

	c.fields[field] = value
	return c
}

func (c *Controller) Write(input any) error {
	if c.err != nil {
		ctx := c.req.Context()
		ctx = context.WithValue(ctx, InputKey, c.fields)
		*(c.req) = *c.req.WithContext(ctx)
		return c.err
	}

	if c.hasJSONBody {
		defer c.req.Body.Close()
		body, err := io.ReadAll(c.req.Body)
		if err != nil {
			return err
		}

		err = json.Unmarshal(body, input)
		if err != nil {
			return err
		}
	}

	err := structop.New(input).UpdateViaMap(c.fields)
	if err != nil {
		return err
	}

	ctx := c.req.Context()
	ctx = context.WithValue(ctx, InputKey, input)
	*(c.req) = *c.req.WithContext(ctx)

	return nil
}
