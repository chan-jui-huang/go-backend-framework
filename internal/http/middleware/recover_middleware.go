package middleware

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/response"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type RecoverMiddleware struct {
	logger *zap.Logger
}

func NewRecoverMiddleware(logger *zap.Logger) *RecoverMiddleware {
	return &RecoverMiddleware{
		logger: logger,
	}
}

func (m *RecoverMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			err := recover()
			if err == nil {
				return
			}
			// condition that warrants a panic stack trace.
			var isBrokenPipe bool
			if ne, ok := err.(*net.OpError); ok {
				var se *os.SyscallError
				if errors.As(ne, &se) {
					if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
						isBrokenPipe = true
					}
				}
			}

			errResp := response.NewErrorResponse(response.InternalServerError, errors.New(fmt.Sprintf("%v", err)), nil, response.DebugMode(c))
			m.logger.Error(response.InternalServerError, errResp.MakeLogFields(c)...)
			if isBrokenPipe {
				c.Abort()
			} else {
				c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
			}
		}()
		c.Next()
	}
}
