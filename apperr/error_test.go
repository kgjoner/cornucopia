package apperr

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func TestNewAndWrapAppError(t *testing.T) {
	err := New(Request, BadRequest, "invalid input")

	var appErr *AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected *AppError, got %T", err)
	}
	if appErr.Kind != Request || appErr.Code != BadRequest {
		t.Fatalf("unexpected error metadata: %+v", appErr)
	}
	if appErr.Error() != "invalid input" {
		t.Fatalf("unexpected message: %q", appErr.Error())
	}

	wrapped := Wrap(errors.New("root cause"), Validation, InvalidData, "validation failed")
	if !errors.As(wrapped, &appErr) {
		t.Fatalf("expected wrapped *AppError, got %T", wrapped)
	}
	if !strings.Contains(appErr.Error(), "root cause") {
		t.Fatalf("expected wrapped message to include cause, got %q", appErr.Error())
	}
}

func TestMarshalJSONUsesErrorMessage(t *testing.T) {
	err := Wrap(errors.New("db down"), Internal, Unexpected, "operation failed")
	data, marshalErr := json.Marshal(err)
	if marshalErr != nil {
		t.Fatalf("marshal error: %v", marshalErr)
	}

	var got map[string]any
	if unmarshalErr := json.Unmarshal(data, &got); unmarshalErr != nil {
		t.Fatalf("unmarshal error: %v", unmarshalErr)
	}

	if got["kind"] != string(Internal) {
		t.Fatalf("expected kind %q, got %#v", Internal, got["kind"])
	}
	if got["code"] != string(Unexpected) {
		t.Fatalf("expected code %q, got %#v", Unexpected, got["code"])
	}
	if !strings.Contains(got["message"].(string), "db down") {
		t.Fatalf("expected marshaled message to include wrapped cause, got %q", got["message"])
	}
}

func TestSpecificConstructorsAndFatal(t *testing.T) {
	cases := []struct {
		name string
		err  error
		kind Kind
		code Code
	}{
		{"validation", NewValidationError("x"), Validation, InvalidData},
		{"request", NewRequestError("x"), Request, BadRequest},
		{"unauthorized", NewUnauthorizedError("x"), Unauthorized, Unauthenticated},
		{"forbidden", NewForbiddenError("x"), Forbidden, NotAllowed},
		{"conflict", NewConflictError("x"), Conflict, Inconsistency},
		{"internal", NewInternalError("x"), Internal, Unexpected},
		{"external", NewExternalError("x"), External, Unexpected},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var appErr *AppError
			if !errors.As(tc.err, &appErr) {
				t.Fatalf("expected *AppError, got %T", tc.err)
			}
			if appErr.Kind != tc.kind || appErr.Code != tc.code {
				t.Fatalf("unexpected kind/code: got (%s,%s) want (%s,%s)", appErr.Kind, appErr.Code, tc.kind, tc.code)
			}
		})
	}

	fatal := Fatal(errors.New("boom"))
	if !IsFatal(fatal) {
		t.Fatal("expected fatal error to be detected")
	}
	if IsFatal(errors.New("normal")) {
		t.Fatal("did not expect non-fatal error to be detected as fatal")
	}
}

func TestMapError(t *testing.T) {
	err := NewMapError(map[string]string{
		"Email": "invalid",
		"Name":  "required",
	})
	msg := err.Error()
	if !strings.Contains(msg, "Email: invalid") || !strings.Contains(msg, "Name: required") {
		t.Fatalf("unexpected map error message: %q", msg)
	}
}
