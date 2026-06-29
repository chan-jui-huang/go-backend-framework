package test

import (
	"testing"

	"github.com/chan-jui-huang/go-backend-framework/v3/internal/test/runtime"
)

type Runtime = runtime.Runtime
type RuntimeOptions = runtime.RuntimeOptions
type MockServices = runtime.MockServices

func NewRuntime(tb testing.TB, options RuntimeOptions) *Runtime {
	return runtime.NewRuntime(tb, options)
}

func NewBaseRuntime(tb testing.TB) *Runtime {
	return runtime.NewBaseRuntime(tb)
}

func NewRdbmsRuntime(tb testing.TB) *Runtime {
	return runtime.NewRdbmsRuntime(tb)
}

func NewClickhouseRuntime(tb testing.TB) *Runtime {
	return runtime.NewClickhouseRuntime(tb)
}

func NewFullRuntime(tb testing.TB) *Runtime {
	return runtime.NewFullRuntime(tb)
}
