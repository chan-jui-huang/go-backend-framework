package requestlog

import "github.com/gin-gonic/gin"

// WithPolicy stores an immutable copy of a route's request logging policy.
func WithPolicy(policy Policy) gin.HandlerFunc {
	policy.ErrorLog = append([]string(nil), policy.ErrorLog...)
	policy.OperationRecord = append([]string(nil), policy.OperationRecord...)

	return func(c *gin.Context) {
		c.Set(policyContextKey, policy)
		c.Next()
	}
}
