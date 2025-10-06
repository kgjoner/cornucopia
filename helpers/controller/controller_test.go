package controller_test

import (
	"bytes"
	"context"
	"errors"
	"mime/multipart"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/kgjoner/cornucopia/v2/helpers/apperr"
	"github.com/kgjoner/cornucopia/v2/helpers/controller"
	"github.com/kgjoner/cornucopia/v2/helpers/htypes"
	"github.com/kgjoner/cornucopia/v2/services/media"
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

// Test struct for unmarshaling
type TestStruct struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
	Token string
	Actor interface{}
}

func TestNew(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	ctrl := controller.New(req)

	if ctrl == nil {
		t.Error("Expected controller to be created, got nil")
	}
}

func TestAddToken(t *testing.T) {
	tests := []struct {
		name    string
		token   interface{}
		wantErr bool
		errType string
	}{
		{
			name:    "valid token",
			token:   "valid-token",
			wantErr: false,
		},
		{
			name:    "missing token",
			token:   nil,
			wantErr: true,
			errType: "Token required.",
		},
		{
			name:    "invalid token type",
			token:   123,
			wantErr: true,
			errType: "Token required.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			if tt.token != nil {
				ctx := context.WithValue(req.Context(), controller.TokenKey, tt.token)
				req = req.WithContext(ctx)
			}

			ctrl := controller.New(req)
			result := ctrl.AddToken()

			var testStruct TestStruct
			err := result.Write(&testStruct)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				if !strings.Contains(err.Error(), tt.errType) {
					t.Errorf("Expected error containing %q, got %q", tt.errType, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if testStruct.Token != tt.token {
					t.Errorf("Expected token %q, got %q", tt.token, testStruct.Token)
				}
			}
		})
	}
}

func TestAddActor(t *testing.T) {
	tests := []struct {
		name    string
		actor   interface{}
		wantErr bool
	}{
		{
			name:    "valid actor",
			actor:   "user123",
			wantErr: false,
		},
		{
			name:    "missing actor",
			actor:   nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			if tt.actor != nil {
				ctx := context.WithValue(req.Context(), controller.ActorKey, tt.actor)
				req = req.WithContext(ctx)
			}

			ctrl := controller.New(req)
			result := ctrl.AddActor()

			var testStruct TestStruct
			err := result.Write(&testStruct)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if testStruct.Actor != tt.actor {
					t.Errorf("Expected actor %v, got %v", tt.actor, testStruct.Actor)
				}
			}
		})
	}
}

func TestAddTarget(t *testing.T) {
	target := "target123"
	req := httptest.NewRequest("GET", "/test", nil)
	ctx := context.WithValue(req.Context(), controller.TargetKey, target)
	req = req.WithContext(ctx)

	ctrl := controller.New(req)
	result := ctrl.AddTarget()

	var testStruct struct {
		Target interface{}
	}
	err := result.Write(&testStruct)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if testStruct.Target != target {
		t.Errorf("Expected target %v, got %v", target, testStruct.Target)
	}
}

func TestAddApplication(t *testing.T) {
	app := "app123"
	req := httptest.NewRequest("GET", "/test", nil)
	ctx := context.WithValue(req.Context(), controller.ApplicationKey, app)
	req = req.WithContext(ctx)

	ctrl := controller.New(req)
	result := ctrl.AddApplication()

	var testStruct struct {
		Application interface{}
	}
	err := result.Write(&testStruct)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if testStruct.Application != app {
		t.Errorf("Expected application %v, got %v", app, testStruct.Application)
	}
}

func TestParseActorAs(t *testing.T) {
	actor := map[string]interface{}{
		"id":   "123",
		"name": "John Doe",
	}

	req := httptest.NewRequest("GET", "/test", nil)
	ctx := context.WithValue(req.Context(), controller.ActorKey, actor)
	req = req.WithContext(ctx)

	ctrl := controller.New(req)
	result := ctrl.ParseActorAs(func(actor any, fields map[string]any) {
		if actorMap, ok := actor.(map[string]interface{}); ok {
			fields["userID"] = actorMap["id"]
			fields["userName"] = actorMap["name"]
		}
	})

	var testStruct struct {
		UserID   string
		UserName string
	}
	err := result.Write(&testStruct)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if testStruct.UserID != "123" {
		t.Errorf("Expected userID '123', got '%s'", testStruct.UserID)
	}
	if testStruct.UserName != "John Doe" {
		t.Errorf("Expected userName 'John Doe', got '%s'", testStruct.UserName)
	}
}

func TestParseURLParam(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)

	// Set up chi context with URL params
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "123")
	rctx.URLParams.Add("name", "john")

	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
	req = req.WithContext(ctx)

	ctrl := controller.New(req)
	result := ctrl.ParseURLParam("id").ParseURLParam("name", "username")

	var testStruct struct {
		ID       string
		Username string
	}
	err := result.Write(&testStruct)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if testStruct.ID != "123" {
		t.Errorf("Expected id '123', got '%s'", testStruct.ID)
	}
	if testStruct.Username != "john" {
		t.Errorf("Expected username 'john', got '%s'", testStruct.Username)
	}
}

func TestParseQueryParam(t *testing.T) {
	req := httptest.NewRequest("GET", "/test?limit=10&offset=20&search=test", nil)

	ctrl := controller.New(req)
	result := ctrl.ParseQueryParam("limit").ParseQueryParam("offset", "skip").ParseQueryParam("search")

	var testStruct struct {
		Limit  string
		Skip   string
		Search string
	}
	err := result.Write(&testStruct)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if testStruct.Limit != "10" {
		t.Errorf("Expected limit '10', got '%s'", testStruct.Limit)
	}
	if testStruct.Skip != "20" {
		t.Errorf("Expected skip '20', got '%s'", testStruct.Skip)
	}
	if testStruct.Search != "test" {
		t.Errorf("Expected search 'test', got '%s'", testStruct.Search)
	}
}

func TestJSONBody(t *testing.T) {
	jsonData := `{"name": "John", "email": "john@example.com", "age": 30}`
	req := httptest.NewRequest("POST", "/test", strings.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")

	ctrl := controller.New(req)
	result := ctrl.JSONBody()

	var testStruct TestStruct
	err := result.Write(&testStruct)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if testStruct.Name != "John" {
		t.Errorf("Expected name 'John', got '%s'", testStruct.Name)
	}
	if testStruct.Email != "john@example.com" {
		t.Errorf("Expected email 'john@example.com', got '%s'", testStruct.Email)
	}
	if testStruct.Age != 30 {
		t.Errorf("Expected age 30, got %d", testStruct.Age)
	}
}

func TestParseBody(t *testing.T) {
	jsonData := `{"Name": "John", "Email": "john@example.com", "Age": 30}`
	req := httptest.NewRequest("POST", "/test", strings.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")

	ctrl := controller.New(req)
	result := ctrl.ParseBody("name", "email", "age")

	var testStruct struct {
		Name  string
		Email string
		Age   float64 // JSON numbers are parsed as float64
	}
	err := result.Write(&testStruct)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if testStruct.Name != "John" {
		t.Errorf("Expected name 'John', got '%s'", testStruct.Name)
	}
	if testStruct.Email != "john@example.com" {
		t.Errorf("Expected email 'john@example.com', got '%s'", testStruct.Email)
	}
	if testStruct.Age != 30 {
		t.Errorf("Expected age 30, got %f", testStruct.Age)
	}
}

func TestAddPagination(t *testing.T) {
	req := httptest.NewRequest("GET", "/test?limit=25&page=2", nil)

	ctrl := controller.New(req)
	result := ctrl.AddPagination()

	var testStruct struct {
		Pagination *htypes.Pagination
	}
	err := result.Write(&testStruct)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if testStruct.Pagination == nil {
		t.Error("Expected pagination to be set")
	}
	if testStruct.Pagination.Limit != 25 {
		t.Errorf("Expected limit 25, got %d", testStruct.Pagination.Limit)
	}
	if testStruct.Pagination.Page != 2 {
		t.Errorf("Expected page 2, got %d", testStruct.Pagination.Page)
	}
}

func TestAddHeader(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer token123")
	req.Header.Set("X-Custom-Header", "custom-value")

	ctrl := controller.New(req)
	result := ctrl.AddHeader("Authorization").AddHeader("X-Custom-Header", "customHeader")

	var testStruct struct {
		Authorization string
		CustomHeader  string
	}
	err := result.Write(&testStruct)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if testStruct.Authorization != "Bearer token123" {
		t.Errorf("Expected authorization 'Bearer token123', got '%s'", testStruct.Authorization)
	}
	if testStruct.CustomHeader != "custom-value" {
		t.Errorf("Expected customHeader 'custom-value', got '%s'", testStruct.CustomHeader)
	}
}

func TestAddIp(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:8080"

	ctrl := controller.New(req)
	result := ctrl.AddIp()

	var testStruct struct {
		Ip string
	}
	err := result.Write(&testStruct)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if testStruct.Ip != "192.168.1.1:8080" {
		t.Errorf("Expected ip '192.168.1.1:8080', got '%s'", testStruct.Ip)
	}
}

func TestAddLanguages(t *testing.T) {
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

			ctrl := controller.New(req)
			result := ctrl.AddLanguages()

			var testStruct struct {
				Languages []string
			}
			err := result.Write(&testStruct)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !reflect.DeepEqual(testStruct.Languages, tt.expected) {
				t.Errorf("Expected languages %v, got %v", tt.expected, testStruct.Languages)
			}
		})
	}
}

func TestParseMultipartForm(t *testing.T) {
	// Create a multipart form
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	// Add a text field
	writer.WriteField("name", "John Doe")
	writer.WriteField("email", "john@example.com")

	// Add a file field
	fileWriter, err := writer.CreateFormFile("avatar", "test.txt")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}
	fileWriter.Write([]byte("test file content"))

	writer.Close()

	req := httptest.NewRequest("POST", "/test", &b)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	mockMedia := &mockMediaService{}
	ctrl := controller.New(req)
	result := ctrl.ParseMultipartForm([]string{"avatar"}, []string{"name", "email"}, mockMedia)

	var testStruct struct {
		Name   string
		Email  string
		Avatar *media.Media
	}
	err = result.Write(&testStruct)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if testStruct.Name != "John Doe" {
		t.Errorf("Expected name 'John Doe', got '%s'", testStruct.Name)
	}
	if testStruct.Email != "john@example.com" {
		t.Errorf("Expected email 'john@example.com', got '%s'", testStruct.Email)
	}
	if testStruct.Avatar == nil {
		t.Error("Expected avatar to be set")
	}
}

func TestErrorHandling(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)

	ctrl := controller.New(req)
	// This should cause an error since no token is provided
	result := ctrl.AddToken()

	var testStruct TestStruct
	err := result.Write(&testStruct)

	if err == nil {
		t.Error("Expected error but got none")
	}

	// Check if it's a normalized error
	var apperr *apperr.AppError
	if errors.As(err, &apperr) {
		if strings.Contains(apperr.Message, "token") {
			t.Errorf("Expected error message containing 'token', got '%s'", apperr.Message)
		}
	} else {
		t.Errorf("Expected AppError, got %T", err)
	}
}

func TestChainedOperations(t *testing.T) {
	req := httptest.NewRequest("GET", "/test?limit=10&page=1", nil)
	req.Header.Set("Authorization", "Bearer token123")
	req.RemoteAddr = "192.168.1.1:8080"

	ctx := context.WithValue(req.Context(), controller.TokenKey, "test-token")
	ctx = context.WithValue(ctx, controller.ActorKey, "user123")
	req = req.WithContext(ctx)

	ctrl := controller.New(req)
	result := ctrl.AddToken().
		AddActor().
		AddPagination().
		AddHeader("Authorization").
		AddIp()

	var testStruct struct {
		Token         string
		Actor         interface{}
		Pagination    *htypes.Pagination
		Authorization string
		Ip            string
	}
	err := result.Write(&testStruct)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if testStruct.Token != "test-token" {
		t.Errorf("Expected token 'test-token', got '%s'", testStruct.Token)
	}
	if testStruct.Actor != "user123" {
		t.Errorf("Expected actor 'user123', got '%v'", testStruct.Actor)
	}
	if testStruct.Pagination == nil {
		t.Error("Expected pagination to be set")
	}
	if testStruct.Authorization != "Bearer token123" {
		t.Errorf("Expected authorization 'Bearer token123', got '%s'", testStruct.Authorization)
	}
	if testStruct.Ip != "192.168.1.1:8080" {
		t.Errorf("Expected ip '192.168.1.1:8080', got '%s'", testStruct.Ip)
	}
}

func TestErrorPropagation(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)

	ctrl := controller.New(req)
	// First operation fails
	result := ctrl.AddToken().AddActor() // Both should fail, but only first error should be returned

	var testStruct TestStruct
	err := result.Write(&testStruct)

	if err == nil {
		t.Error("Expected error but got none")
	}

	// Should be the first error (Token required)
	if !strings.Contains(err.Error(), "Token required") {
		t.Errorf("Expected error containing 'Token required', got '%s'", err.Error())
	}
}
