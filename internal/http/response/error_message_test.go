package response_test

import (
	"testing"

	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http/response"
	"github.com/stretchr/testify/assert"
)

func TestMessageToCode(t *testing.T) {
	assert.Equal(t, "400-001", response.MessageToCode[response.BadRequest])
	assert.Equal(t, "400-002", response.MessageToCode[response.RequestValidationFailed])
	assert.Equal(t, "400-003", response.MessageToCode[response.EmailIsWrong])
	assert.Equal(t, "400-004", response.MessageToCode[response.PasswordIsWrong])
	assert.Equal(t, "401-001", response.MessageToCode[response.Unauthorized])
	assert.Equal(t, "403-001", response.MessageToCode[response.Forbidden])
	assert.Equal(t, "500-001", response.MessageToCode[response.InternalServerError])
	assert.Equal(t, "503-001", response.MessageToCode[response.ServiceUnavailable])
}
