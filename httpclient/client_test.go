package httpclient

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kgjoner/cornucopia/v2/apperr"
)

func TestSetOptions(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://example.com/test", nil)
	SetOptions(req, Options{
		Headers: map[string]string{"X-Test": "ok"},
		Params:  map[string]string{"q": "search"},
	})

	if got := req.Header.Get("X-Test"); got != "ok" {
		t.Fatalf("expected header to be set, got %q", got)
	}
	if got := req.URL.Query().Get("q"); got != "search" {
		t.Fatalf("expected query param to be set, got %q", got)
	}
}

func TestDoReqSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"name": "john"})
	}))
	defer srv.Close()

	req, err := http.NewRequest(http.MethodGet, srv.URL, nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}

	var out map[string]any
	res, err := DoReq(srv.Client(), req, &out)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
	if out["name"] != "john" {
		t.Fatalf("unexpected decoded payload: %#v", out)
	}
}

func TestDoReqErrorMapping(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"message":"not allowed"}`))
	}))
	defer srv.Close()

	req, err := http.NewRequest(http.MethodGet, srv.URL, nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}

	_, err = DoReq(srv.Client(), req, nil)
	if err == nil {
		t.Fatal("expected error")
	}

	var appErr *apperr.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got %T", err)
	}
	if appErr.Kind != apperr.External || appErr.Code != apperr.Unauthenticated {
		t.Fatalf("unexpected kind/code: (%s, %s)", appErr.Kind, appErr.Code)
	}
	if !strings.Contains(appErr.Error(), "not allowed") {
		t.Fatalf("expected message to include body message, got %q", appErr.Error())
	}
}

func TestClientRequestWithDefaultOptions(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Default") != "yes" {
			t.Fatalf("default header missing")
		}
		if r.URL.Query().Get("tenant") != "a" {
			t.Fatalf("default query param missing")
		}
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	c := New(srv.URL)
	c.SetDefaultOptions(&Options{
		Headers: map[string]string{"X-Default": "yes"},
		Params:  map[string]string{"tenant": "a"},
	})

	var out map[string]any
	_, err := c.Get("", nil)(&out)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if out["ok"] != true {
		t.Fatalf("unexpected payload: %#v", out)
	}
}
