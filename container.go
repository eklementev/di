package di

import (
	"sync"
)

// Container defines interface for DI-container
type Container interface {
	Errorer
	// Define add bean definition with specific name, scope and builder
	Define(name string, scope Scope, builder Builder) Container
	// Lookup returns singleton bean
	Lookup(name string) (Bean, error)
	// Build constructs and return new prototype bean; Setup and Bind are invoked for new bean
	Build(name string) (Bean, error)

	// Setup initiates setup & binding procesures for all singleton beans
	Setup() error
	// Shutdown initiates shutdown procesures for all singleton beans
	Shutdown()
}

type container struct {
	*ErrorEmitter
	singletons map[string]Bean
	prototypes map[string]Builder
	mx         *sync.RWMutex
	wg         *sync.WaitGroup
}

// New creates Container instance
func New() Container {
	return &container{
		ErrorEmitter: NewErrorEmitter(),
		singletons:   map[string]Bean{},
		prototypes:   map[string]Builder{},
		mx:           new(sync.RWMutex),
		wg:           new(sync.WaitGroup),
	}
}

func (c *container) Define(name string, scope Scope, builder Builder) Container {
	c.mx.Lock()
	defer c.mx.Unlock()

	switch scope {
	case ScopeSingleton:
		c.singletons[name] = builder()
	case ScopePrototype:
		c.prototypes[name] = builder
	default:
		panic("unreachable")
	}

	return c
}

func (c *container) Lookup(name string) (Bean, error) {
	c.mx.RLock()
	defer c.mx.RUnlock()

	bean, ok := c.singletons[name]
	if !ok {
		return nil, ErrUnknownBean{name}
	}

	return bean, nil
}

func (c *container) Build(name string) (Bean, error) {
	c.mx.RLock()
	defer c.mx.RUnlock()

	builder, ok := c.prototypes[name]
	if !ok {
		return nil, ErrUnknownBean{name}
	}

	bean := builder()
	err := bean.Setup(c)
	if err != nil {
		return nil, err
	}
	err = bean.PostSetup()
	if err != nil {
		return nil, err
	}

	return bean, nil
}

func (c *container) Setup() error {
	c.mx.RLock()
	defer c.mx.RUnlock()

	for _, bean := range c.singletons {
		errorer, ok := bean.(Errorer)
		if !ok {
			continue
		}

		c.wg.Add(1)
		go func() {
			defer c.wg.Done()
			for err := range errorer.ErrorCh() {
				c.FireError(err)
			}
		}()
	}

	for _, bean := range c.singletons {
		err := bean.Setup(c)
		if err != nil {
			return err
		}
	}

	for _, bean := range c.singletons {
		err := bean.PostSetup()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *container) Shutdown() {
	c.mx.RLock()
	defer c.mx.RUnlock()

	c.wg.Add(len(c.singletons))
	for _, bean := range c.singletons {
		go func(bean Bean) {
			if typed, ok := bean.(errorerShutdowner); ok {
				typed.shutdownErrorer()
			}
			bean.Shutdown()
			c.wg.Done()
		}(bean)
	}
	c.wg.Wait()

	c.shutdownErrorer()
}
