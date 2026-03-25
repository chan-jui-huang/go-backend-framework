package config

import (
	"net/http"

	"golang.org/x/time/rate"
)

type CsrfConfig struct {
	Cookie struct {
		Name     string
		Path     string
		Domain   string
		MaxAge   int
		Secure   bool
		HttpOnly bool
		SameSite http.SameSite
	}
	Header string
}

type RateLimitConfig struct {
	PutTokenRate rate.Limit
	BurstNumber  int
}
