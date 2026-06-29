package runtime

import "go.uber.org/fx"

// MockServices is the extension point for service mocks supplied by tests.
// Add fields here when the framework introduces an injectable external service.
type MockServices struct{}

func provideMockServices(mockServices MockServices) fx.Option {
	return fx.Options()
}
