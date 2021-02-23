package di_test

import (
	"testing"

	"github.com/eklementev/di"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type simpleBean struct {
	value         int
	setupDone     bool
	postSetupDone bool
	shutdownDone  bool
}

func (sb *simpleBean) Setup(di.Container) error { sb.setupDone = true; return nil }
func (sb *simpleBean) PostSetup() error         { sb.postSetupDone = true; return nil }
func (sb *simpleBean) Shutdown()                { sb.shutdownDone = true }

func TestNewConstructsUniqueInstances(t *testing.T) {
	c1 := di.New()
	c2 := di.New()
	assert.False(t, c1 == c2)
}

func TestLookupReturnsSameSingleton(t *testing.T) {
	original := 12
	replacement := 13

	c := di.New()
	c.Define("s", di.ScopeSingleton, di.StaticBuilder(&simpleBean{value: original}))

	bean1, err := c.Lookup("s")
	require.Nil(t, err)
	bean2, err := c.Lookup("s")
	require.Nil(t, err)

	simpleBean1 := bean1.(*simpleBean)
	simpleBean2 := bean2.(*simpleBean)

	assert.NotNil(t, simpleBean1)
	assert.NotNil(t, simpleBean2)

	simpleBean1.value = replacement
	assert.Equal(t, simpleBean1.value, simpleBean2.value)
}

func TestLookupAndBuildReturnsNilForUnknownName(t *testing.T) {
	c := di.New()
	c.Define("s", di.ScopeSingleton, di.StaticBuilder(&simpleBean{}))
	c.Define("p", di.ScopePrototype, func() di.Bean { return &simpleBean{} })

	bean, err := c.Lookup("unknown")
	_, ok := err.(di.ErrUnknownBean)
	assert.True(t, ok)
	assert.Nil(t, bean)

	built, err := c.Build("unknown")
	_, ok = err.(di.ErrUnknownBean)
	assert.True(t, ok)
	assert.Nil(t, built)
}

func TestBuildReturnsDifferentsFromPrototypes(t *testing.T) {
	original := 12
	replacement := 13
	counter := 0
	expectedCounter := 2

	c := di.New()
	c.Define("p", di.ScopePrototype, func() di.Bean { counter++; return &simpleBean{value: original} })

	bean1, err := c.Build("p")
	require.Nil(t, err)
	bean2, err := c.Build("p")
	require.Nil(t, err)

	simpleBean1 := bean1.(*simpleBean)
	simpleBean2 := bean2.(*simpleBean)

	assert.True(t, counter == expectedCounter)
	assert.NotNil(t, simpleBean1)
	assert.NotNil(t, simpleBean2)

	simpleBean1.value = replacement
	assert.NotEqual(t, simpleBean1.value, simpleBean2.value)
}

func TestLifeCycle(t *testing.T) {
	bean0 := &simpleBean{}
	bean1 := &simpleBean{}

	c := di.New()
	c.Define("0", di.ScopeSingleton, di.StaticBuilder(bean0))
	c.Define("1", di.ScopeSingleton, di.StaticBuilder(bean1))
	c.Define("z", di.ScopePrototype, func() di.Bean { return &simpleBean{} })

	err := c.Setup()
	require.Nil(t, err)

	beanzuntyped, err := c.Build("z")
	require.Nil(t, err)
	beanz := beanzuntyped.(*simpleBean)

	assert.True(t, bean0.setupDone)
	assert.True(t, bean0.postSetupDone)
	assert.True(t, bean1.setupDone)
	assert.True(t, bean1.postSetupDone)
	assert.True(t, beanz.setupDone)
	assert.True(t, beanz.postSetupDone)

	c.Shutdown()

	assert.True(t, bean0.shutdownDone)
	assert.True(t, bean1.shutdownDone)
	assert.False(t, beanz.shutdownDone)
}

func TestBeanOverriding(t *testing.T) {
	c := di.New()

	sOriginal := 11
	sReplacement := 22
	pOriginal := 33
	pReplacement := 44

	c.Define("s", di.ScopeSingleton, di.StaticBuilder(&simpleBean{value: sOriginal}))
	c.Define("s", di.ScopeSingleton, di.StaticBuilder(&simpleBean{value: sReplacement}))

	c.Define("p", di.ScopePrototype, func() di.Bean { return &simpleBean{value: pOriginal} })
	c.Define("p", di.ScopePrototype, func() di.Bean { return &simpleBean{value: pReplacement} })

	p, err := c.Build("p")
	require.Nil(t, err)

	bean, err := c.Lookup("s")
	require.Nil(t, err)

	assert.True(t, bean.(*simpleBean).value == sReplacement)
	assert.True(t, p.(*simpleBean).value == pReplacement)
}

func TestDefineWithInvalidScope(t *testing.T) {
	invalidScopeValue := 6

	assert.Panics(t, func() { di.New().Define("fail", di.Scope(invalidScopeValue), di.StaticBuilder(&simpleBean{})) })
}
