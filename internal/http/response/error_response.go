package response

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/chan-jui-huang/go-backend-framework/v2/internal/deps"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/stacktrace"
	"go.uber.org/zap"
)

type Debug struct {
	Error      string `json:"error" example:"error message" validate:"required"`
	err        error
	Stacktrace []string `json:"stacktrace" validate:"required"`
}

type ErrorResponse struct {
	Message string `json:"message" validate:"required"`
	Code    string `json:"code" validate:"required"`
	debug   *Debug
	Debug   *Debug         `json:"debug,omitempty"`
	Context map[string]any `json:"context,omitempty"`
}

func NewErrorResponse(message string, err error, context map[string]any) *ErrorResponse {
	debug := &Debug{
		err:        err,
		Stacktrace: stacktrace.GetStackStrace(err),
	}
	if err != nil {
		debug.Error = err.Error()
	}

	booterConfig := deps.BooterConfig()
	if booterConfig.Debug {
		return &ErrorResponse{
			Message: message,
			Code:    MessageToCode[message],
			Debug:   debug,
			Context: context,
		}
	}

	return &ErrorResponse{
		Message: message,
		Code:    MessageToCode[message],
		debug:   debug,
		Context: context,
	}
}

func (er *ErrorResponse) StatusCode() int {
	code, err := strconv.ParseInt(
		strings.Split(er.Code, "-")[0],
		10,
		0,
	)
	if err != nil {
		logger := deps.Logger()
		logger.Error(err.Error())
		code = http.StatusBadRequest
	}

	return int(code)
}

func (er *ErrorResponse) MakeLogFields(req *http.Request, fields ...zap.Field) []zap.Field {
	requestBody, err := io.ReadAll(req.Body)
	if err != nil {
		logger := deps.Logger()
		logger.Error(err.Error())
		requestBody = nil
	}
	if requestBody != nil && json.Valid(requestBody) {
		buffer := bytes.NewBuffer(make([]byte, 0, len(requestBody)))
		err = json.Compact(buffer, requestBody)
		if err != nil {
			logger := deps.Logger()
			logger.Error(err.Error())
			requestBody = nil
		} else {
			requestBody = buffer.Bytes()
		}
	}

	var debug *Debug
	booterConfig := deps.BooterConfig()
	if booterConfig.Debug {
		debug = er.Debug
	} else {
		debug = er.debug
	}

	errorString := ""
	if debug.err != nil {
		errorString = debug.err.Error()
	}

	fields = append(
		[]zap.Field{
			zap.String("code", er.Code),
			zap.String("error", errorString),
			zap.Int("status_code", er.StatusCode()),
			zap.String("method", req.Method),
			zap.String("path", req.URL.Path),
			zap.String("query_string", req.URL.Query().Encode()),
			zap.ByteString("request_body", requestBody),
			zap.Strings("stacktrace", debug.Stacktrace),
		},
		fields...,
	)

	return fields
}
