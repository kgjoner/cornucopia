package httpserver

import (
	"bytes"
	"encoding"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/kgjoner/cornucopia/v2/media"
	"github.com/kgjoner/cornucopia/v2/prim"
)

// Binder accumulates the first binding error so callers need check it only once.
// Values are written to destination pointers immediately on each method call.
//
// Usage:
//
//	b := httpserver.Bind(r)
//	b.JSONBody(&input).PathParam("id", &input.ID).Header("auth", &input.Token)
//	b.Pagination(&input.Pagination).Market(&input.Market).IP(&input.IP)
//	if err := b.Err(); err != nil { ... }
type Binder struct {
	r   *http.Request
	err error
}

// Bind creates a Binder for the given request.
func Bind(r *http.Request) *Binder {
	return &Binder{r: r}
}

// Err returns the first error encountered during binding, or nil.
func (b *Binder) Err() error {
	return b.err
}

// JSONBody decodes the request body as JSON into dst.
func (b *Binder) JSONBody(dst any) *Binder {
	if b.err != nil {
		return b
	}
	defer b.r.Body.Close()
	body, err := io.ReadAll(b.r.Body)
	if err != nil {
		b.err = err
		return b
	}
	b.err = json.Unmarshal(body, dst)
	return b
}

// PathParam binds a chi URL path parameter to dst.
// dst must be a *string, a numeric/bool pointer, or implement encoding.TextUnmarshaler.
func (b *Binder) PathParam(name string, dst any) *Binder {
	if b.err != nil {
		return b
	}

	b.err = bindText(chi.URLParam(b.r, name), dst)
	return b
}

// QueryParam binds a single query string value to dst.
// If the param is absent, dst is unchanged.
func (b *Binder) QueryParam(name string, dst any) *Binder {
	if b.err != nil {
		return b
	}
	if val := b.r.URL.Query().Get(name); val != "" {
		b.err = bindText(val, dst)
	}
	return b
}

// QueryParams binds repeated query string values to a *[]string.
// If the param is absent, dst is unchanged.
func (b *Binder) QueryParams(name string, dst *[]string) *Binder {
	if b.err != nil {
		return b
	}
	if vals := b.r.URL.Query()[name]; len(vals) > 0 {
		*dst = vals
	}
	return b
}

// Header binds a request header value to dst.
// If the header is absent, dst is unchanged.
func (b *Binder) Header(name string, dst any) *Binder {
	if b.err != nil {
		return b
	}
	if val := b.r.Header.Get(name); val != "" {
		b.err = bindText(val, dst)
	}
	return b
}

// FromContext binds a context value to dst. dst must be a non-nil pointer to the
// same type stored in the context. If the key is absent, dst is left unchanged.
func (b *Binder) FromContext(key any, dst any) *Binder {
	if b.err != nil {
		return b
	}
	val := b.r.Context().Value(key)
	if val == nil {
		return b
	}
	dv := reflect.ValueOf(dst)
	if dv.Kind() != reflect.Pointer || dv.IsNil() {
		b.err = fmt.Errorf("httpserver: FromContext dst must be a non-nil pointer, got %T", dst)
		return b
	}
	sv := reflect.ValueOf(val)
	elem := dv.Elem()
	if !sv.Type().AssignableTo(elem.Type()) {
		b.err = fmt.Errorf("httpserver: cannot assign context value of type %T to %T", val, elem.Interface())
		return b
	}
	elem.Set(sv)
	return b
}

// ParseMultipartForm parses the request as multipart/form-data.
// Must be called before FormValue or UploadedFile.
// Use 0 for the default 32 MB memory limit.
func (b *Binder) ParseMultipartForm(maxMemory int64) *Binder {
	if b.err != nil {
		return b
	}
	if maxMemory == 0 {
		maxMemory = 32 << 20
	}
	b.err = b.r.ParseMultipartForm(maxMemory)
	return b
}

// FormValue binds a multipart or URL-encoded form text field to dst.
func (b *Binder) FormValue(name string, dst any) *Binder {
	if b.err != nil {
		return b
	}
	if val := b.r.PostFormValue(name); val != "" {
		b.err = bindText(val, dst)
	}
	return b
}

// UploadedFile reads the named file from a parsed multipart form as a *media.Media.
// ParseMultipartForm must be called first. Returns (nil, false) if the file is absent.
func (b *Binder) UploadedFile(name string, svc media.MediaService) (*media.Media, bool) {
	if b.err != nil {
		return nil, false
	}
	file, _, err := b.r.FormFile(name)
	if err == http.ErrMissingFile {
		return nil, false
	}
	if err != nil {
		b.err = err
		return nil, false
	}
	var buf bytes.Buffer
	io.Copy(&buf, file)
	return media.New(&buf, svc), true
}

// Pagination parses limit/page query params into dst.
// Absent or invalid values fall back to defaults.
func (b *Binder) Pagination(dst *prim.Pagination) *Binder {
	if b.err != nil {
		return b
	}
	if dst == nil {
		b.err = fmt.Errorf("httpserver: Pagination dst must be non-nil")
		return b
	}
	*dst = *ParsePagination(b.r)
	return b
}

// ParsePagination parses limit/page query params and returns a Pagination value.
// Absent or invalid values fall back to defaults.
func ParsePagination(r *http.Request) *prim.Pagination {
	q := r.URL.Query()
	var limit, page int64
	if v := q.Get("limit"); v != "" {
		limit, _ = strconv.ParseInt(v, 10, 32)
	}
	if v := q.Get("page"); v != "" {
		page, _ = strconv.ParseInt(v, 10, 32)
	}
	return prim.NewPagination(&prim.PaginationCreationFields{
		Limit: int(limit),
		Page:  int(page),
	})
}

// Market parses the market derived from the x-timezone request header into dst.
func (b *Binder) Market(dst *prim.Market) *Binder {
	if b.err != nil {
		return b
	}
	if dst == nil {
		b.err = fmt.Errorf("httpserver: Market dst must be non-nil")
		return b
	}
	market, err := ParseMarket(b.r)
	if err != nil {
		b.err = err
		return b
	}
	*dst = market
	return b
}

// ParseMarket returns the market derived from the x-timezone request header.
func ParseMarket(r *http.Request) (prim.Market, error) {
	return prim.MarketByTimezone(r.Header.Get("x-timezone"))
}

// Languages parses and quality-sorts values from the Accept-Language header into dst.
func (b *Binder) Languages(dst *[]string) *Binder {
	if b.err != nil {
		return b
	}
	if dst == nil {
		b.err = fmt.Errorf("httpserver: Languages dst must be non-nil")
		return b
	}
	langs, err := ParseLanguages(b.r)
	if err != nil {
		b.err = err
		return b
	}
	*dst = langs
	return b
}

// ParseLanguages returns the parsed and quality-sorted values from the Accept-Language header.
func ParseLanguages(r *http.Request) ([]string, error) {
	type langQ struct {
		lang   string
		weight float64
	}
	acptLang := r.Header.Get("accept-language")
	var lqs []langQ
	for _, part := range strings.Split(acptLang, ",") {
		part = strings.TrimSpace(part)
		segments := strings.Split(part, ";")
		if len(segments) == 1 {
			lqs = append(lqs, langQ{segments[0], 1})
			continue
		}
		qp := strings.Split(segments[1], "=")
		if len(qp) != 2 {
			return nil, fmt.Errorf("malformed Accept-Language quality value: %q", segments[1])
		}
		q, err := strconv.ParseFloat(qp[1], 64)
		if err != nil {
			return nil, err
		}
		lqs = append(lqs, langQ{segments[0], q})
	}
	sort.SliceStable(lqs, func(i, j int) bool {
		return lqs[i].weight > lqs[j].weight
	})
	langs := make([]string, len(lqs))
	for i, lq := range lqs {
		langs[i] = strings.ToLower(lq.lang)
	}
	return langs, nil
}

// IP binds the request's remote address into dst.
func (b *Binder) IP(dst *string) *Binder {
	if b.err != nil {
		return b
	}
	if dst == nil {
		b.err = fmt.Errorf("httpserver: IP dst must be non-nil")
		return b
	}
	*dst = ParseIP(b.r)
	return b
}

// ParseIP returns the request's remote address.
func ParseIP(r *http.Request) string {
	return r.RemoteAddr
}

// ContextValue retrieves a typed value from the request context.
func ContextValue[T any](r *http.Request, key any) (T, bool) {
	val := r.Context().Value(key)
	if val == nil {
		var zero T
		return zero, false
	}
	typed, ok := val.(T)
	return typed, ok
}

// bindText assigns the string s to dst.
// Supported types: *string, all integer/float/bool pointer kinds, and encoding.TextUnmarshaler.
func bindText(s string, dst any) error {
	if s == "" {
		return nil
	}
	if tu, ok := dst.(encoding.TextUnmarshaler); ok {
		return tu.UnmarshalText([]byte(s))
	}
	switch v := dst.(type) {
	case *string:
		*v = s
	case *int:
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		*v = int(n)
	case *int8:
		n, err := strconv.ParseInt(s, 10, 8)
		if err != nil {
			return err
		}
		*v = int8(n)
	case *int16:
		n, err := strconv.ParseInt(s, 10, 16)
		if err != nil {
			return err
		}
		*v = int16(n)
	case *int32:
		n, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return err
		}
		*v = int32(n)
	case *int64:
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		*v = n
	case *uint:
		n, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}
		*v = uint(n)
	case *uint8:
		n, err := strconv.ParseUint(s, 10, 8)
		if err != nil {
			return err
		}
		*v = uint8(n)
	case *uint16:
		n, err := strconv.ParseUint(s, 10, 16)
		if err != nil {
			return err
		}
		*v = uint16(n)
	case *uint32:
		n, err := strconv.ParseUint(s, 10, 32)
		if err != nil {
			return err
		}
		*v = uint32(n)
	case *uint64:
		n, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}
		*v = n
	case *float32:
		f, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return err
		}
		*v = float32(f)
	case *float64:
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		*v = f
	case *bool:
		bv, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}
		*v = bv
	default:
		return fmt.Errorf("httpserver: unsupported bind destination type %T", dst)
	}
	return nil
}
