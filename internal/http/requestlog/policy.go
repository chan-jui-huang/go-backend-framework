package requestlog

import "github.com/gin-gonic/gin"

// Destination identifies a request-body logging consumer.
type Destination uint8

const (
	// ErrorLog selects fields allowed in application error logs.
	ErrorLog Destination = iota
	// OperationRecord selects fields allowed in persisted operation records.
	OperationRecord
)

// Policy defines allowed JSON paths for each logging destination.
type Policy struct {
	ErrorLog        []string
	OperationRecord []string
}

const policyContextKey = "request_log_policy"

func fieldsFor(c *gin.Context, destination Destination) ([]string, bool) {
	value, ok := c.Get(policyContextKey)
	if !ok {
		return nil, false
	}

	policy, ok := value.(Policy)
	if !ok {
		return nil, false
	}

	switch destination {
	case ErrorLog:
		return policy.ErrorLog, true
	case OperationRecord:
		return policy.OperationRecord, true
	default:
		return nil, true
	}
}
