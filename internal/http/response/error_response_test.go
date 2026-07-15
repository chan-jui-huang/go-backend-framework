package response

import (
	"net/http/httptest"
	"testing"

	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/requestlog"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestErrorResponseMakeLogFieldsFiltersRequestBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("POST", "/api/user/login", nil)
	c.Set(gin.BodyBytesKey, []byte("{\"email\":\"john@test.com\",\"password\":\"secret\"}"))
	requestlog.WithPolicy(requestlog.Policy{ErrorLog: []string{"email"}})(c)
	errResp := NewErrorResponse(BadRequest, errors.New("bad request"), nil, false)

	fields := errResp.MakeLogFields(c)

	var requestBody []byte
	for _, field := range fields {
		if field.Key == "request_body" {
			requestBody = field.Interface.([]byte)
		}
	}
	require.JSONEq(t, "{\"email\":\"john@test.com\"}", string(requestBody))
	require.NotContains(t, string(requestBody), "secret")
}
