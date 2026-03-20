package integration_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kgjoner/cornucopia/v3/apperr"
	"github.com/kgjoner/cornucopia/v3/httpclient"
	"github.com/kgjoner/cornucopia/v3/httpserver"
	"github.com/kgjoner/cornucopia/v3/prim"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// errHandler returns an http.HandlerFunc that writes an apperr via httpserver.Error.
func errHandler(appErr error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		httpserver.Error(appErr, w, r)
	}
}

// TestHTTPServerErrorStatusMapping verifies that each apperr.Kind is mapped to
// the expected HTTP status code by httpserver.Error.
func TestHTTPServerErrorStatusMapping(t *testing.T) {
	cases := []struct {
		name           string
		err            error
		expectedStatus int
	}{
		{"validation", apperr.NewValidationError("invalid"), http.StatusUnprocessableEntity},
		{"request", apperr.NewRequestError("bad input"), http.StatusBadRequest},
		{"unauthorized", apperr.NewUnauthorizedError("no auth"), http.StatusUnauthorized},
		{"forbidden", apperr.NewForbiddenError("no access"), http.StatusForbidden},
		{"conflict", apperr.NewConflictError("duplicate"), http.StatusConflict},
		{"internal", apperr.NewInternalError("oops"), http.StatusInternalServerError},
		// External with Code=Unexpected (default) → 502 Bad Gateway.
		{"external_unexpected", apperr.NewExternalError("upstream failed"), http.StatusBadGateway},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			srv := httptest.NewServer(errHandler(tc.err))
			defer srv.Close()

			resp, err := http.Get(srv.URL)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)
		})
	}
}

// TestHTTPServerErrorResponseBodyEncoding verifies that the JSON error body
// returned by httpserver.Error contains the expected kind, code, and message.
func TestHTTPServerErrorResponseBodyEncoding(t *testing.T) {
	appErr := apperr.NewValidationError("email is invalid")
	srv := httptest.NewServer(errHandler(appErr))
	defer srv.Close()

	resp, err := http.Get(srv.URL)
	require.NoError(t, err)
	defer resp.Body.Close()

	var body map[string]any
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))

	assert.Equal(t, string(apperr.Validation), body["kind"])
	assert.Equal(t, string(apperr.InvalidData), body["code"])
	assert.Contains(t, body["message"].(string), "email is invalid")
}

// TestHTTPServerContextCanceledMapsTo499 verifies that context.Canceled is
// handled as a 499 status (client closed request).
func TestHTTPServerContextCanceledMapsTo499(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpserver.Error(context.Canceled, w, r)
	}))
	defer srv.Close()

	resp, err := http.Get(srv.URL)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, 499, resp.StatusCode)
}

// TestHTTPServerContextDeadlineExceededMapsTo408 verifies that
// context.DeadlineExceeded is mapped to 408 Request Timeout.
func TestHTTPServerContextDeadlineExceededMapsTo408(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpserver.Error(context.DeadlineExceeded, w, r)
	}))
	defer srv.Close()

	resp, err := http.Get(srv.URL)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusRequestTimeout, resp.StatusCode)
}

// TestHTTPServerUnknownErrorMapsTo500 verifies that a plain error (not an
// apperr.AppError and not a context error) maps to a 500 status.
func TestHTTPServerUnknownErrorMapsTo500(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpserver.Error(errors.New("something went wrong"), w, r)
	}))
	defer srv.Close()

	resp, err := http.Get(srv.URL)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

// TestHTTPServerSuccessWrapsData verifies that httpserver.Success wraps a
// plain struct response in a {"data": ...} envelope.
func TestHTTPServerSuccessWrapsData(t *testing.T) {
	type item struct{ Name string }
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpserver.Success(item{Name: "widget"}, w, r)
	}))
	defer srv.Close()

	resp, err := http.Get(srv.URL)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var body map[string]any
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))

	data, ok := body["data"].(map[string]any)
	require.True(t, ok, "expected 'data' key wrapping the struct")
	assert.Equal(t, "widget", data["Name"])
}

// TestHTTPServerSuccessWithPaginatedData verifies that a prim.PaginatedData
// response (which already has a Data field) is NOT double-wrapped.
func TestHTTPServerSuccessWithPaginatedData(t *testing.T) {
	pagination := prim.Pagination{Page: 0, Limit: 10, HasNext: false}
	pd := prim.NewPaginatedData(pagination, []string{"a", "b"})

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpserver.Success(pd, w, r)
	}))
	defer srv.Close()

	resp, err := http.Get(srv.URL)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var body map[string]any
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))

	// PaginatedData has its own Data field — it should appear directly, not
	// nested under an additional "data" wrapper.
	rawData, hasData := body["data"]
	require.True(t, hasData, "expected 'data' key in paginated response")

	items, ok := rawData.([]any)
	require.True(t, ok, "expected 'data' to be a slice")
	assert.Equal(t, []any{"a", "b"}, items)
}

// TestHTTPServerSuccess201 verifies that passing status 201 to httpserver.Success
// results in a 201 Created response.
func TestHTTPServerSuccess201(t *testing.T) {
	type created struct{ ID int }
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpserver.Success(created{ID: 42}, w, r, http.StatusCreated)
	}))
	defer srv.Close()

	resp, err := http.Get(srv.URL)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

// TestHTTPClientReadsAppErrorFromServer tests the full pipeline where a server
// uses httpserver.Error and a client reads the response via httpclient.DoReq,
// resulting in an apperr.AppError on the client side.
func TestHTTPClientReadsAppErrorFromServer(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpserver.Error(apperr.NewValidationError("name is required"), w, r)
	}))
	defer srv.Close()

	req, err := http.NewRequest(http.MethodGet, srv.URL, nil)
	require.NoError(t, err)

	_, clientErr := httpclient.DoReq(srv.Client(), req, nil)
	require.Error(t, clientErr)

	var appErr *apperr.AppError
	require.True(t, errors.As(clientErr, &appErr))
	// Client side: the error is External since it came from a remote call.
	assert.Equal(t, apperr.External, appErr.Kind)
	// The "message" field in the response body is forwarded into the wrapped error.
	assert.Contains(t, appErr.Error(), "name is required")
}

// TestHTTPClientMapsUnauthorizedFromServer verifies that a 401 response from
// the server is mapped to Code=Unauthenticated on the client side.
func TestHTTPClientMapsUnauthorizedFromServer(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpserver.Error(apperr.NewUnauthorizedError("invalid token"), w, r)
	}))
	defer srv.Close()

	req, err := http.NewRequest(http.MethodGet, srv.URL, nil)
	require.NoError(t, err)

	_, clientErr := httpclient.DoReq(srv.Client(), req, nil)
	require.Error(t, clientErr)

	var appErr *apperr.AppError
	require.True(t, errors.As(clientErr, &appErr))
	assert.Equal(t, apperr.External, appErr.Kind)
	assert.Equal(t, apperr.Unauthenticated, appErr.Code)
}
