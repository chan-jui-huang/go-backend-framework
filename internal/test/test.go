package test

import (
	"testing"

	testruntime "github.com/chan-jui-huang/go-backend-framework/v2/internal/test/runtime"
)

type Runtime = testruntime.Runtime
type RuntimeOptions = testruntime.RuntimeOptions

func NewRuntime(tb testing.TB, options RuntimeOptions) *Runtime {
	return testruntime.NewRuntime(tb, options)
}

func NewBaseRuntime(tb testing.TB) *Runtime {
	return testruntime.NewBaseRuntime(tb)
}

func NewRdbmsRuntime(tb testing.TB) *Runtime {
	return testruntime.NewRdbmsRuntime(tb)
}

func NewClickhouseRuntime(tb testing.TB) *Runtime {
	return testruntime.NewClickhouseRuntime(tb)
}

func NewFullRuntime(tb testing.TB) *Runtime {
	return testruntime.NewFullRuntime(tb)
}
