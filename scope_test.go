package di_test

import (
	"testing"

	"github.com/eklementev/di"
	"github.com/stretchr/testify/assert"
)

func TestScopeString(t *testing.T) {
	assert.Equal(t, di.ScopeSingleton.String(), "singleton")
	assert.Equal(t, di.ScopePrototype.String(), "prototype")
}

func TestInvalidScopeString(t *testing.T) {
	invalidScopeValue := 6

	assert.Panics(t, func() { _ = di.Scope(invalidScopeValue).String() })
}
