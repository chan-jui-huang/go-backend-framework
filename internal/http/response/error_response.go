package response

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/requestlog"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/stacktrace"
	"github.com/gin-gonic/gin"
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

const debugModeKey = "response_debug_mode"

func SetDebugMode(c *gin.Context, debugMode bool) {
	c.Set(debugModeKey, debugMode)
}

func DebugMode(c *gin.Context) bool {
	value, ok := c.Get(debugModeKey)
	if !ok {
		return false
	}

	debugMode, ok := value.(bool)
	return ok && debugMode
}

func NewErrorResponse(message string, err error, context map[string]any, debugMode bool) *ErrorResponse {
	debug := &Debug{
		err:        err,
		Stacktrace: stacktrace.GetStackStrace(err),
	}
	if err != nil {
		debug.Error = err.Error()
	}

	errResp := &ErrorResponse{
		Message: message,
		Code:    MessageToCode[message],
		debug:   debug,
		Context: context,
	}
	if debugMode {
		errResp.Debug = debug
	}

	return errResp
}

func (er *ErrorResponse) StatusCode() int {
	code, err := strconv.ParseInt(
		strings.Split(er.Code, "-")[0],
		10,
		0,
	)
	if err != nil {
		return http.StatusBadRequest
	}

	return int(code)
}

func (er *ErrorResponse) MakeLogFields(c *gin.Context, fields ...zap.Field) []zap.Field {
	req := c.Request

	var internalFields []zap.Field
	requestBody, err := requestlog.Filter(c, requestlog.ErrorLog)
	if err != nil {
		internalFields = append(internalFields, zap.NamedError("request_body_filter_error", err))
	}

	debug := er.debug
	if er.Debug != nil {
		debug = er.Debug
	}

	errorString := ""
	if debug != nil && debug.err != nil {
		errorString = debug.err.Error()
	}

	stacktraceValue := []string(nil)
	if debug != nil {
		stacktraceValue = debug.Stacktrace
	}

	baseFields := []zap.Field{
		zap.String("code", er.Code),
		zap.String("error", errorString),
		zap.Int("status_code", er.StatusCode()),
		zap.String("method", req.Method),
		zap.String("path", req.URL.Path),
		zap.String("query_string", req.URL.Query().Encode()),
		zap.Strings("stacktrace", stacktraceValue),
	}
	if len(requestBody) > 0 {
		baseFields = append(baseFields, zap.ByteString("request_body", requestBody))
	}
	baseFields = append(baseFields, internalFields...)
	baseFields = append(baseFields, fields...)

	return baseFields
}
