package http

import (
	"net/http"
)

func (handler *Handler) AddCsrfToken(req *http.Request) {
	config := handler.csrfConfig
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
