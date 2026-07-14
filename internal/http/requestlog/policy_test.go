package requestlog

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestWithPolicyDoesNotConsumeRequestBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("POST", "/", strings.NewReader("request body"))
	policy := Policy{
		ErrorLog:        []string{"email"},
		OperationRecord: []string{"name"},
	}
	middleware := WithPolicy(policy)
	policy.ErrorLog[0] = "password"
	policy.OperationRecord[0] = "token"

	middleware(c)

	body, err := io.ReadAll(c.Request.Body)
	require.NoError(t, err)
	require.Equal(t, "request body", string(body))
	errorLogFields, hasPolicy := fieldsFor(c, ErrorLog)
	require.True(t, hasPolicy)
	require.Equal(t, []string{"email"}, errorLogFields)
	operationRecordFields, hasPolicy := fieldsFor(c, OperationRecord)
	require.True(t, hasPolicy)
	require.Equal(t, []string{"name"}, operationRecordFields)
}

func TestFieldsForWithoutPolicy(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	fields, hasPolicy := fieldsFor(c, ErrorLog)

	require.False(t, hasPolicy)
	require.Nil(t, fields)
}
