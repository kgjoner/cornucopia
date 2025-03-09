package normalizederr

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strconv"
)

type NormalizedError struct {
	Message string `json:"message"`
	Kind    string `json:"kind"`
	Code    string `json:"code"`
	Stack   string `json:"-"`
	msgMap  map[string]string
}

func NewValidationError(message string, code ...string) NormalizedError {
	var effectiveCode string
	if len(code) >= 1 && code[0] != "" {
		effectiveCode = code[0]
	} else {
		effectiveCode = InvalidData
	}

	return NormalizedError{
		message,
		"Validation",
		effectiveCode,
		getStack(),
		nil,
	}
}

func NewValidationErrorFromMap(errorMap map[string]error, code ...string) NormalizedError {
	message := "{ "
	msgMap := map[string]string{}
	for k, v := range errorMap {
		if v != nil {			
			message += fmt.Sprintf("%v: %v \n", k, v)
			msgMap[k] = v.Error()
		}
	}
	message += " }"

	var effectiveCode string
	if len(code) >= 1 && code[0] != "" {
		effectiveCode = code[0]
	} else {
		effectiveCode = InvalidData
	}

	return NormalizedError{
		message,
		"Validation",
		effectiveCode,
		getStack(),
		msgMap,
	}
}

func NewRequestError(message string, code ...string) NormalizedError {
	var effectiveCode string
	if len(code) >= 1 && code[0] != "" {
		effectiveCode = code[0]
	} else {
		effectiveCode = BadRequest
	}

	return NormalizedError{
		message,
		"Request",
		effectiveCode,
		getStack(),
		nil,
	}
}

func NewForbiddenError(message string, code ...string) NormalizedError {
	var effectiveCode string
	if len(code) >= 1 && code[0] != "" {
		effectiveCode = code[0]
	} else {
		effectiveCode = NotAllowed
	}

	return NormalizedError{
		message,
		"Forbidden",
		effectiveCode,
		getStack(),
		nil,
	}
}

func NewUnauthorizedError(message string, code ...string) NormalizedError {
	var effectiveCode string
	if len(code) >= 1 && code[0] != "" {
		effectiveCode = code[0]
	} else {
		effectiveCode = Unauthenticated
	}

	return NormalizedError{
		message,
		"Unauthorized",
		effectiveCode,
		getStack(),
		nil,
	}
}

func NewFatalUnauthorizedError(message string, code ...string) NormalizedError {
	var effectiveCode string
	if len(code) >= 1 && code[0] != "" {
		effectiveCode = code[0]
	} else {
		effectiveCode = Unauthenticated
	}

	return NormalizedError{
		message,
		"FatalUnauthorized",
		effectiveCode,
		getStack(),
		nil,
	}
}

func NewConflictError(message string, code ...string) NormalizedError {
	var effectiveCode string
	if len(code) >= 1 && code[0] != "" {
		effectiveCode = code[0]
	} else {
		effectiveCode = Inconsistency
	}

	return NormalizedError{
		message,
		"Conflict",
		effectiveCode,
		getStack(),
		nil,
	}
}

func NewFatalInternalError(message string, code ...string) NormalizedError {
	var effectiveCode string
	if len(code) >= 1 && code[0] != "" {
		effectiveCode = code[0]
	} else {
		effectiveCode = Unexpected
	}

	return NormalizedError{
		message,
		"FatalInternal",
		effectiveCode,
		getStack(),
		nil,
	}
}

func NewExternalError(mainMsg string, complMsg map[string]error, code ...string) NormalizedError {
	message := fmt.Sprintf("{\nmain: %v \n", mainMsg)
	msgMap := map[string]string{
		"main": mainMsg,
	}

	for k, v := range complMsg {
		if v != nil {
			message += fmt.Sprintf("%v: %v \n", k, v)
			msgMap[k] = v.Error()
		}
	}

	message += " }"

	var effectiveCode string
	if len(code) >= 1 && code[0] != "" {
		effectiveCode = code[0]
	} else {
		effectiveCode = Unexpected
	}

	return NormalizedError{
		message,
		"External",
		effectiveCode,
		getStack(),
		msgMap,
	}
}

func (e NormalizedError) Error() string {
	return e.Message
}

// Make external errors internal. It returns the unchanged error if it is not external.
func (e NormalizedError) MakeItInternal() NormalizedError {
	if e.Kind == "External" {
		statusStr, ok := e.msgMap["ResponseStatus"]
		status, err := strconv.Atoi(statusStr)
		if ok && err == nil {
			switch {
			case status == 401:
				e.Kind = "Unauthorized"
			case status == 403:
				e.Kind = "Forbidden"
			case status == 422:
				e.Kind = "Validation"
			case status == 409:
				e.Kind = "Conflict"
			case status >= 400 && status < 500:
				e.Kind = "Request"
			default:
				e.Kind = "Internal"
			}
		} else {
			e.Kind = "Internal"
		}
	}
	return e
}

func (e NormalizedError) MarshalJSON() ([]byte, error) {
	if e.msgMap != nil && len(e.msgMap) > 0 {
		payload := struct {
			Message map[string]string `json:"message"`
			Kind    string            `json:"kind"`
			Code    string            `json:"code"`
		}{
			e.msgMap,
			e.Kind,
			e.Code,
		}

		return json.Marshal(payload)
	}

	payload := struct {
		Message string `json:"message"`
		Kind    string `json:"kind"`
		Code    string `json:"code"`
	}{
		e.Message,
		e.Kind,
		e.Code,
	}
	return json.Marshal(payload)
}

func getStack() string {
	stackSlice := make([]byte, 512)
	s := runtime.Stack(stackSlice, false)
	return string(stackSlice[0:s])
}
