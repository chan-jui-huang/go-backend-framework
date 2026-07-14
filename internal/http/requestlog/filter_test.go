package requestlog

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestFilter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		body        string
		policy      *Policy
		destination Destination
		expected    string
		expectsErr  bool
	}{
		{
			name:     "default records non-sensitive body",
			body:     "{\"name\":\"John\",\"email\":\"john@test.com\"}",
			expected: "{\"name\":\"John\",\"email\":\"john@test.com\"}",
		},
		{
			name:       "default rejects nested sensitive body",
			body:       "{\"profile\":{\"accessToken\":\"secret\"}}",
			expectsErr: true,
		},
		{
			name:       "default rejects top-level sensitive body",
			body:       "{\"email\":\"john@test.com\",\"password\":\"secret\"}",
			expectsErr: true,
		},
		{
			name:       "default rejects sensitive suffixes",
			body:       "{\"metadata\":{\"new_password\":\"secret\",\"webhookSecret\":\"secret\",\"id-token\":\"secret\"}}",
			expectsErr: true,
		},
		{
			name:       "default rejects credential containers",
			body:       "{\"serviceCredential\":{\"value\":\"secret\"}}",
			expectsErr: true,
		},
		{
			name:   "empty policy denies malformed body without parsing",
			body:   "{\"email\":",
			policy: &Policy{},
		},
		{
			name:     "login policy retains email",
			body:     "{\"email\":\"john@test.com\",\"password\":\"secret\"}",
			policy:   &Policy{ErrorLog: []string{"email"}},
			expected: "{\"email\":\"john@test.com\"}",
		},
		{
			name:     "nested array paths",
			body:     "{\"items\":[{\"id\":1,\"password\":\"secret\"}]}",
			policy:   &Policy{ErrorLog: []string{"items.id"}},
			expected: "{\"items\":[{\"id\":1}]}",
		},
		{
			name:     "nested object paths",
			body:     "{\"settings\":{\"timezone\":\"Asia/Taipei\",\"api_key\":\"secret\"}}",
			policy:   &Policy{ErrorLog: []string{"settings.timezone"}},
			expected: "{\"settings\":{\"timezone\":\"Asia/Taipei\"}}",
		},
		{
			name:       "forbidden allowlist field",
			body:       "{\"Password\":\"secret\"}",
			policy:     &Policy{ErrorLog: []string{"Password"}},
			expectsErr: true,
		},
		{
			name:       "forbidden field below allowed parent",
			body:       "{\"settings\":{\"apiKey\":\"secret\"}}",
			policy:     &Policy{ErrorLog: []string{"settings"}},
			expectsErr: true,
		},
		{
			name:       "invalid allowlist path",
			body:       "{\"settings\":{\"timezone\":\"Asia/Taipei\"}}",
			policy:     &Policy{ErrorLog: []string{"settings..timezone"}},
			expectsErr: true,
		},
		{
			name:       "invalid JSON",
			body:       "{\"email\":",
			policy:     &Policy{ErrorLog: []string{"email"}},
			expectsErr: true,
		},
		{
			name:       "non-object JSON",
			body:       "[1,2,3]",
			policy:     &Policy{ErrorLog: []string{"email"}},
			expectsErr: true,
		},
		{
			name:       "body too large",
			body:       strings.Repeat("a", maxBodySize+1),
			policy:     &Policy{ErrorLog: []string{"email"}},
			expectsErr: true,
		},
		{
			name:        "separate operation record policy",
			body:        "{\"email\":\"john@test.com\",\"name\":\"John\"}",
			policy:      &Policy{ErrorLog: []string{"email"}, OperationRecord: []string{"name"}},
			destination: OperationRecord,
			expected:    "{\"name\":\"John\"}",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Set(gin.BodyBytesKey, []byte(test.body))
			if test.policy != nil {
				WithPolicy(*test.policy)(c)
			}

			originalBody := append([]byte(nil), c.MustGet(gin.BodyBytesKey).([]byte)...)
			result, err := Filter(c, test.destination)
			if test.expectsErr {
				require.Error(t, err)
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				if test.expected == "" {
					require.Empty(t, result)
				} else {
					require.JSONEq(t, test.expected, string(result))
				}
			}
			require.Equal(t, originalBody, c.MustGet(gin.BodyBytesKey))
		})
	}
}
