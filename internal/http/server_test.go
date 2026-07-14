package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chan-jui-huang/go-backend-framework/v3/internal/config"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/middleware"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/response"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/booter"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestRouterGroupsInheritGlobalMiddlewares(t *testing.T) {
	accessCore, accessLogs := observer.New(zapcore.DebugLevel)
	globalMiddlewares := newGlobalMiddlewares(zap.New(accessCore), zap.NewNop(), false)

	engine, err := NewEngine(globalMiddlewares)
	if err != nil {
		t.Fatalf("create engine: %v", err)
	}
	apiRouter := engine.Group("api")
	apiRouter.GET("ping", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/ping", nil)
	engine.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, responseRecorder.Code)
	}
	if accessLogs.Len() != 1 {
		t.Fatalf("expected one access log, got %d", accessLogs.Len())
	}
	if message := accessLogs.All()[0].Message; message != "GET /api/ping" {
		t.Fatalf("expected access log message %q, got %q", "GET /api/ping", message)
	}
}

func TestGlobalAccessLogHandlesMatchedAndUnmatchedClientErrors(t *testing.T) {
	accessCore, accessLogs := observer.New(zapcore.DebugLevel)
	globalMiddlewares := newGlobalMiddlewares(zap.New(accessCore), zap.NewNop(), false)

	engine, err := NewEngine(globalMiddlewares)
	if err != nil {
		t.Fatalf("create engine: %v", err)
	}
	apiRouter := engine.Group("api")
	apiRouter.GET("unauthorized", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusUnauthorized)
	})

	testCases := []struct {
		name    string
		path    string
		status  int
		message string
	}{
		{
			name:    "matched route",
			path:    "/api/unauthorized",
			status:  http.StatusUnauthorized,
			message: "GET /api/unauthorized",
		},
		{
			name:    "unmatched route",
			path:    "/missing",
			status:  http.StatusNotFound,
			message: "GET /missing",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			responseRecorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, testCase.path, nil)
			engine.ServeHTTP(responseRecorder, request)

			if responseRecorder.Code != testCase.status {
				t.Fatalf("expected status %d, got %d", testCase.status, responseRecorder.Code)
			}
			entries := accessLogs.FilterMessage(testCase.message).All()
			if len(entries) != 1 {
				t.Fatalf("expected one access log for %s, got %d", testCase.path, len(entries))
			}
			if entries[0].Level != zapcore.WarnLevel {
				t.Fatalf("expected warn-level access log, got %s", entries[0].Level)
			}
		})
	}
}

func TestRouterGroupsUseGlobalRecoverMiddleware(t *testing.T) {
	accessCore, accessLogs := observer.New(zapcore.DebugLevel)
	applicationCore, applicationLogs := observer.New(zapcore.DebugLevel)
	globalMiddlewares := newGlobalMiddlewares(zap.New(accessCore), zap.New(applicationCore), false)

	engine, err := NewEngine(globalMiddlewares)
	if err != nil {
		t.Fatalf("create engine: %v", err)
	}
	apiRouter := engine.Group("api")
	apiRouter.GET("panic", func(_ *gin.Context) {
		panic("test panic")
	})

	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/panic", nil)
	engine.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, responseRecorder.Code)
	}
	errorResponse := struct {
		Message string `json:"message"`
	}{}
	if err := json.Unmarshal(responseRecorder.Body.Bytes(), &errorResponse); err != nil {
		t.Fatalf("decode error response: %v", err)
	}
	if errorResponse.Message != response.InternalServerError {
		t.Fatalf("expected error message %q, got %q", response.InternalServerError, errorResponse.Message)
	}
	if applicationLogs.FilterMessage(response.InternalServerError).Len() != 1 {
		t.Fatalf("expected one recovery log, got %d", applicationLogs.Len())
	}
	if accessLogs.Len() != 1 {
		t.Fatalf("expected one access log, got %d", accessLogs.Len())
	}
	if level := accessLogs.All()[0].Level; level != zapcore.ErrorLevel {
		t.Fatalf("expected error-level access log, got %s", level)
	}
}

func TestRouterGroupsUseGlobalResponseContextMiddleware(t *testing.T) {
	globalMiddlewares := newGlobalMiddlewares(zap.NewNop(), zap.NewNop(), true)

	engine, err := NewEngine(globalMiddlewares)
	if err != nil {
		t.Fatalf("create engine: %v", err)
	}
	apiRouter := engine.Group("api")
	apiRouter.GET("debug", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"debug": response.DebugMode(c)})
	})

	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/debug", nil)
	engine.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, responseRecorder.Code)
	}
	responseBody := struct {
		Debug bool `json:"debug"`
	}{}
	if err := json.Unmarshal(responseRecorder.Body.Bytes(), &responseBody); err != nil {
		t.Fatalf("decode response body: %v", err)
	}
	if !responseBody.Debug {
		t.Fatal("expected response debug mode to be enabled")
	}
}

func TestRouterGroupsUseGlobalCsrfMiddleware(t *testing.T) {
	globalMiddlewares := newGlobalMiddlewares(zap.NewNop(), zap.NewNop(), false)

	engine, err := NewEngine(globalMiddlewares)
	if err != nil {
		t.Fatalf("create engine: %v", err)
	}
	apiRouter := engine.Group("api")
	apiRouter.POST("protected", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/protected", nil)
	engine.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, responseRecorder.Code)
	}
}

func newGlobalMiddlewares(accessLogger *zap.Logger, applicationLogger *zap.Logger, debugMode bool) *middleware.GlobalMiddlewares {
	csrfConfig := &config.CsrfConfig{Header: "X-CSRF-Token"}
	csrfConfig.Cookie.Name = "csrf-token"

	return middleware.NewGlobalMiddlewares(
		middleware.NewAccessLogMiddleware(accessLogger),
		middleware.NewRecoverMiddleware(applicationLogger),
		middleware.NewCsrfMiddleware(applicationLogger, csrfConfig),
		middleware.NewResponseContextMiddleware(booter.NewConfig("", "", debugMode)),
	)
}
