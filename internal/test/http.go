package test

import (
	"net/http"

	pkgHttp "github.com/chan-jui-huang/go-backend-framework/v2/internal/http"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http/middleware"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http/route"
	"github.com/chan-jui-huang/go-backend-package/pkg/booter/config"
	"github.com/gin-gonic/gin"
)

type httpHandler struct {
	engine *gin.Engine
}

var HttpHandler *httpHandler

func newHttpHandler(handlerFuncs ...gin.HandlerFunc) *httpHandler {
	engine, err := pkgHttp.NewEngine()
	if err != nil {
		panic(err)
	}

	handler := &httpHandler{
		engine: engine,
	}
	handler.AttachGlobalMiddleware()
	if len(handlerFuncs) > 0 {
		handler.engine.Use(handlerFuncs...)
	}
	routers := []route.Router{
		route.NewApiRouter(handler.engine),
		route.NewSwaggerRouter(handler.engine),
	}
	for _, router := range routers {
		router.AttachRoutes()
	}

	return handler
}

func NewHttpHandler() *httpHandler {
	return newHttpHandler()
}

func (handler *httpHandler) AttachGlobalMiddleware() {
	csrfConfig := config.Registry.Get("middleware.csrf").(middleware.CsrfConfig)

	handlerFunctions := []gin.HandlerFunc{
		middleware.VerifyCsrfToken(csrfConfig),
	}

	handler.engine.Use(
		handlerFunctions...,
	)
}

func (handler *httpHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handler.engine.ServeHTTP(w, req)
}

func AddCsrfToken(req *http.Request) {
	config := config.Registry.Get("middleware.csrf").(middleware.CsrfConfig)
	cookie := &http.Cookie{
		Name:     config.Cookie.Name,
		Value:    "1234567890",
		Path:     config.Cookie.Path,
		Domain:   config.Cookie.Domain,
		MaxAge:   config.Cookie.MaxAge,
		Secure:   config.Cookie.Secure,
		HttpOnly: config.Cookie.HttpOnly,
		SameSite: config.Cookie.SameSite,
	}
	req.AddCookie(cookie)
	req.Header.Set(config.Header, "1234567890")
}

func AddBearerToken(req *http.Request, token string) {
	req.Header.Set("Authorization", token)
}

func EnableMiddleware(handlerFuncs ...gin.HandlerFunc) func() {
	original := HttpHandler
	handler := newHttpHandler(handlerFuncs...)

	HttpHandler = handler

	return func() {
		HttpHandler = original
	}
}
