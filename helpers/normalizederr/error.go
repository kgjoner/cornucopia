package normalizederr

import (
	"encoding/json"
	"fmt"
	"runtime"
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
	if len(code) >= 1 {
		effectiveCode = code[0]
	} else {
		effectiveCode = UnprocessableEntity
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
	if len(code) >= 1 {
		effectiveCode = code[0]
	} else {
		effectiveCode = UnprocessableEntity
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
	if len(code) >= 1 {
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
	if len(code) >= 1 {
		effectiveCode = code[0]
	} else {
		effectiveCode = Forbidden
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
	if len(code) >= 1 {
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
	if len(code) >= 1 {
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

func NewFatalInternalError(message string, code ...string) NormalizedError {
	var effectiveCode string
	if len(code) >= 1 {
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
	if len(code) >= 1 {
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
