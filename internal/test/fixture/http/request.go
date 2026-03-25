package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

func (handler *Handler) NewJSONRequest(method string, path string, body any) *http.Request {
	return httptest.NewRequest(method, path, MarshalJSONBody(body))
}

func MarshalJSONBody(body any) *bytes.Reader {
	if body == nil {
		return bytes.NewReader(nil)
	}

	payload, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	return bytes.NewReader(payload)
}
