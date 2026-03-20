package httpserver_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kgjoner/cornucopia/v3/apperr"
	"github.com/kgjoner/cornucopia/v3/httpserver"
)

func TestNewLogger_withAppError(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	err := apperr.NewValidationError("bad input")

	entry := httpserver.NewLogger(req, err)
	if entry == nil {
		t.Error("Expected a log entry, got nil")
	}
}

func TestNewLogger_withGenericError(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/things", nil)
	err := fmt.Errorf("unexpected failure")

	entry := httpserver.NewLogger(req, err)
	if entry == nil {
		t.Error("Expected a log entry, got nil")
	}
}

func TestNewLogger_withStatus(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/things", nil)

	entry := httpserver.NewLogger(req, 201)
	if entry == nil {
		t.Error("Expected a log entry, got nil")
	}
}
