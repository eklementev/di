package di

// Scope defines beans scope
type Scope uint8

const (
	// ScopeSingleton defines scope for bean with single instance
	ScopeSingleton Scope = iota
	// ScopePrototype defines scope for beans, each time are created with builder
	ScopePrototype Scope = iota
)

// String represents scope as string
func (s Scope) String() string {
	switch s {
	case ScopeSingleton:
		return "singleton"
	case ScopePrototype:
		return "prototype"
	default:
		panic("unreachable")
	}
}
