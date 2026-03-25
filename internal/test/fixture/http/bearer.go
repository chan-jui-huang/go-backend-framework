package http

import "net/http"

func (handler *Handler) AddBearerToken(req *http.Request, token string) {
	req.Header.Set("Authorization", token)
}
