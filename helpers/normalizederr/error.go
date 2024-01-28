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

func NewValidationError(message string) NormalizedError {
	return NormalizedError{
		message,
		"Validation",
		"",
		getStack(),
		nil,
	}
}

func NewValidationErrorFromMap(errorMap map[string]error) NormalizedError {
	message := "{ "
	msgMap := map[string]string{}
	for k, v := range errorMap {
		if v != nil {
			message += fmt.Sprintf("%v: %v \n", k, v)
			msgMap[k] = v.Error()
		}
	}
	message += " }"

	return NormalizedError{
		message,
		"Validation",
		"",
		getStack(),
		msgMap,
	}
}

func NewRequestError(message string, code string) NormalizedError {
	return NormalizedError{
		message,
		"Request",
		code,
		getStack(),
		nil,
	}
}

func NewForbiddenError(message string) NormalizedError {
	return NormalizedError{
		message,
		"Forbidden",
		"0002002",
		getStack(),
		nil,
	}
}

func NewUnauthorizedError(message string) NormalizedError {
	return NormalizedError{
		message,
		"Unauthorized",
		"0002001",
		getStack(),
		nil,
	}
}

func NewFatalUnauthorizedError(message string) NormalizedError {
	return NormalizedError{
		message,
		"FatalUnauthorized",
		"0002001",
		getStack(),
		nil,
	}
}

func NewFatalInternalError(message string) NormalizedError {
	return NormalizedError{
		message,
		"FatalInternal",
		"",
		getStack(),
		nil,
	}
}

func NewExternalError(mainMsg string, complMsg map[string]error) NormalizedError {
	message := fmt.Sprintf("{\nmain: %v \n", mainMsg)
	msgMap := map[string]string{
		"main": mainMsg,
	}

	if complMsg != nil {
		for k, v := range complMsg {
			if v != nil {
				message += fmt.Sprintf("%v: %v \n", k, v)
				msgMap[k] = v.Error()
			}
		}
	}

	message += " }"

	return NormalizedError{
		message,
		"External",
		"",
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
			Message map[string]string
			Kind    string
			Code    string
		}{
			e.msgMap,
			e.Kind,
			e.Code,
		}

		return json.Marshal(payload)
	}

	payload := struct {
		Message string
		Kind    string
		Code    string
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
	return fmt.Sprintf("%s", stackSlice[0:s])
}
