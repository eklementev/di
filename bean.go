package di

// Bean defines interface for bean
type Bean interface {
	// Setup is callback for initial bean setup
	Setup(Container) error
	// PostSetup is callback for additional bean configuration
	PostSetup() error
	// Shutdown is callback for bean shutdown
	Shutdown()
}

// Builder defines function that produces new bean
type Builder func() Bean

// StaticBuilder is Builder with predefined bean instance
func StaticBuilder(bean Bean) Builder {
	return func() Bean { return bean }
}
