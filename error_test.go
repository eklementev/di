package di_test

import (
	"errors"
	"testing"

	"github.com/eklementev/di"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type errorBean struct {
	*di.ErrorEmitter
}

func (bean *errorBean) Setup(di.Container) error { return nil }
func (bean *errorBean) PostSetup() error         { return nil }
func (bean *errorBean) Shutdown()                {}

func TestErrorRaised(t *testing.T) {
	bean := &errorBean{ErrorEmitter: di.NewErrorEmitter()}

	c := di.New()
	c.Define("e", di.ScopeSingleton, di.StaticBuilder(bean))
	err := c.Setup()
	require.Nil(t, err)

	raised := errors.New("raised")
	go func() {
		bean.FireError(raised)
	}()

	err = <-c.ErrorCh()

	assert.Equal(t, raised, err)

	c.Shutdown()
}

func TestRaiseMultipleErrors(t *testing.T) {
	bean := &errorBean{ErrorEmitter: di.NewErrorEmitter()}

	c := di.New()
	c.Define("e", di.ScopeSingleton, di.StaticBuilder(bean))
	err := c.Setup()
	require.Nil(t, err)

	raised := errors.New("raised")
	other := errors.New("other")
	go func() {
		bean.FireError(raised)
		bean.FireError(other)
	}()

	err = <-c.ErrorCh()
	assert.Equal(t, raised, err)

	select {
	case <-c.ErrorCh():
		require.Fail(t, "the other error should not occur")
	default:
	}

	c.Shutdown()
}
