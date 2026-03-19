package httpserver_test

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/kgjoner/cornucopia/v2/httpserver"
	"github.com/kgjoner/cornucopia/v2/media"
	"github.com/kgjoner/cornucopia/v2/prim"
)

// Mock media service for testing
type mockMediaService struct {
	storeFunc func(file *bytes.Buffer, kind, id string) (string, error)
}

func (m *mockMediaService) Store(file *bytes.Buffer, kind, id string) (string, error) {
	if m.storeFunc != nil {
		return m.storeFunc(file, kind, id)
	}
	return "http://example.com/test.jpg", nil
}

type TestCtxKey string

func TestBind(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	b := httpserver.Bind(req)
	if b == nil {
		t.Error("Expected Binder to be created, got nil")
	}
	if err := b.Err(); err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestPathParam(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "123")
	rctx.URLParams.Add("name", "john")
	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
	req = req.WithContext(ctx)

	var id string
	var username string
	err := httpserver.Bind(req).
		PathParam("id", &id).
		PathParam("name", &username).
		Err()

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if id != "123" {
		t.Errorf("Expected id '123', got '%s'", id)
	}
	if username != "john" {
		t.Errorf("Expected username 'john', got '%s'", username)
	}
}

func TestQueryParam(t *testing.T) {
	req := httptest.NewRequest("GET", "/test?limit=10&offset=20&search=test", nil)

	var limit, offset, search string
	err := httpserver.Bind(req).
		QueryParam("limit", &limit).
		QueryParam("offset", &offset).
		QueryParam("search", &search).
		Err()

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if limit != "10" {
		t.Errorf("Expected limit '10', got '%s'", limit)
	}
	if offset != "20" {
		t.Errorf("Expected offset '20', got '%s'", offset)
	}
	if search != "test" {
		t.Errorf("Expected search 'test', got '%s'", search)
	}
}

func TestQueryParamInt(t *testing.T) {
	req := httptest.NewRequest("GET", "/test?count=42", nil)

	var count int
	err := httpserver.Bind(req).QueryParam("count", &count).Err()

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if count != 42 {
		t.Errorf("Expected count 42, got %d", count)
	}
}

func TestJSONBody(t *testing.T) {
	type Input struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Age   int    `json:"age"`
	}
	jsonData := `{"name": "John", "email": "john@example.com", "age": 30}`
	req := httptest.NewRequest("POST", "/test", strings.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")

	var input Input
	err := httpserver.Bind(req).JSONBody(&input).Err()

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if input.Name != "John" {
		t.Errorf("Expected name 'John', got '%s'", input.Name)
	}
	if input.Email != "john@example.com" {
		t.Errorf("Expected email 'john@example.com', got '%s'", input.Email)
	}
	if input.Age != 30 {
		t.Errorf("Expected age 30, got %d", input.Age)
	}
}

func TestPagination(t *testing.T) {
	req := httptest.NewRequest("GET", "/test?limit=25&page=2", nil)

	var pag prim.Pagination
	err := httpserver.Bind(req).Pagination(&pag).Err()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if pag.Limit != 25 {
		t.Errorf("Expected limit 25, got %d", pag.Limit)
	}
	if pag.Page != 2 {
		t.Errorf("Expected page 2, got %d", pag.Page)
	}
}

func TestHeader(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer token123")
	req.Header.Set("X-Custom-Header", "custom-value")

	var authorization, customHeader string
	err := httpserver.Bind(req).
		Header("Authorization", &authorization).
		Header("X-Custom-Header", &customHeader).
		Err()

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if authorization != "Bearer token123" {
		t.Errorf("Expected 'Bearer token123', got '%s'", authorization)
	}
	if customHeader != "custom-value" {
		t.Errorf("Expected 'custom-value', got '%s'", customHeader)
	}
}

func TestIP(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:8080"

	var ip string
	err := httpserver.Bind(req).IP(&ip).Err()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if ip != "192.168.1.1:8080" {
		t.Errorf("Expected '192.168.1.1:8080', got '%s'", ip)
	}
}

func TestLanguages(t *testing.T) {
	tests := []struct {
		name           string
		acceptLanguage string
		expected       []string
	}{
		{
			name:           "single language",
			acceptLanguage: "en-US",
			expected:       []string{"en-us"},
		},
		{
			name:           "multiple languages with quality",
			acceptLanguage: "en-US,en;q=0.9,fr;q=0.8",
			expected:       []string{"en-us", "en", "fr"},
		},
		{
			name:           "languages with different quality values",
			acceptLanguage: "fr;q=0.8,en-US;q=0.9,en;q=0.7",
			expected:       []string{"en-us", "fr", "en"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			req.Header.Set("Accept-Language", tt.acceptLanguage)

			var langs []string
			err := httpserver.Bind(req).Languages(&langs).Err()

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !reflect.DeepEqual(langs, tt.expected) {
				t.Errorf("Expected languages %v, got %v", tt.expected, langs)
			}
		})
	}
}

func TestFromContext(t *testing.T) {
	var tokenKey TestCtxKey = "token"
	req := httptest.NewRequest("GET", "/test", nil)
	ctx := context.WithValue(req.Context(), tokenKey, "test-token")
	req = req.WithContext(ctx)

	var token string
	err := httpserver.Bind(req).FromContext(tokenKey, &token).Err()

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if token != "test-token" {
		t.Errorf("Expected token 'test-token', got '%s'", token)
	}
}

func TestContextValue(t *testing.T) {
	type User struct{ ID int }
	var userKey TestCtxKey = "user"
	req := httptest.NewRequest("GET", "/test", nil)
	ctx := context.WithValue(req.Context(), userKey, User{ID: 42})
	req = req.WithContext(ctx)

	user, ok := httpserver.ContextValue[User](req, userKey)

	if !ok {
		t.Error("Expected context value to be found")
	}
	if user.ID != 42 {
		t.Errorf("Expected user ID 42, got %d", user.ID)
	}
}

func TestParseMultipartForm(t *testing.T) {
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)
	writer.WriteField("name", "John Doe")
	writer.WriteField("email", "john@example.com")

	fileWriter, err := writer.CreateFormFile("avatar", "test.txt")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}
	fileWriter.Write([]byte("test file content"))
	writer.Close()

	req := httptest.NewRequest("POST", "/test", &b)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	mockSvc := &mockMediaService{}
	binder := httpserver.Bind(req).ParseMultipartForm(0)

	var name, email string
	binder.FormValue("name", &name).FormValue("email", &email)
	avatar, hasAvatar := binder.UploadedFile("avatar", mockSvc)

	if err := binder.Err(); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if name != "John Doe" {
		t.Errorf("Expected name 'John Doe', got '%s'", name)
	}
	if email != "john@example.com" {
		t.Errorf("Expected email 'john@example.com', got '%s'", email)
	}
	if !hasAvatar || avatar == nil {
		t.Error("Expected avatar to be set")
	}
}

func TestChainedBinding(t *testing.T) {
	req := httptest.NewRequest("GET", "/test?limit=10&page=1", nil)
	req.Header.Set("Authorization", "Bearer token123")
	req.RemoteAddr = "192.168.1.1:8080"

	var tokenKey TestCtxKey = "token"
	ctx := context.WithValue(req.Context(), tokenKey, "test-token")
	req = req.WithContext(ctx)

	var token, authorization string
	var pag prim.Pagination
	var ip string
	b := httpserver.Bind(req)
	b.FromContext(tokenKey, &token).
		Header("Authorization", &authorization).
		Pagination(&pag).
		IP(&ip)

	if err := b.Err(); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if token != "test-token" {
		t.Errorf("Expected token 'test-token', got '%s'", token)
	}
	if authorization != "Bearer token123" {
		t.Errorf("Expected authorization 'Bearer token123', got '%s'", authorization)
	}
	if pag.Limit != 10 || pag.Page != 1 {
		t.Errorf("Unexpected pagination: %+v", pag)
	}
	if ip != "192.168.1.1:8080" {
		t.Errorf("Expected ip '192.168.1.1:8080', got '%s'", ip)
	}
}

func TestParseHelpers(t *testing.T) {
	req := httptest.NewRequest("GET", "/test?limit=15&page=3", nil)
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("x-timezone", "America/Sao_Paulo")
	req.RemoteAddr = "10.0.0.1:9090"

	pag := httpserver.ParsePagination(req)
	if pag.Limit != 15 || pag.Page != 3 {
		t.Errorf("Unexpected pagination: %+v", pag)
	}

	market, err := httpserver.ParseMarket(req)
	if err != nil {
		t.Fatalf("Unexpected ParseMarket error: %v", err)
	}
	if market != prim.MarketBrazil {
		t.Errorf("Expected market %q, got %q", prim.MarketBrazil, market)
	}

	langs, err := httpserver.ParseLanguages(req)
	if err != nil {
		t.Fatalf("Unexpected ParseLanguages error: %v", err)
	}
	if !reflect.DeepEqual(langs, []string{"en-us", "en"}) {
		t.Errorf("Unexpected languages: %v", langs)
	}

	ip := httpserver.ParseIP(req)
	if ip != "10.0.0.1:9090" {
		t.Errorf("Expected IP 10.0.0.1:9090, got %s", ip)
	}
}

func TestBindTextUnmarshaler(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("market", "brazil")
	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
	req = req.WithContext(ctx)

	var market prim.Market
	err := httpserver.Bind(req).PathParam("market", &market).Err()

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if market != prim.MarketBrazil {
		t.Errorf("Expected MarketBrazil, got %q", market)
	}
}

// Ensure unused imports from media package don't slip through.
var _ media.MediaService = (*mockMediaService)(nil)
