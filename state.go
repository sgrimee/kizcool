package kizcool

// StateDefinition describes the fields of a State
type StateDefinition struct {
	Type          string
	QualifiedName string
	Values        []string
}

// StateNameT is the name of a State
type StateNameT string

// StateTypeT has value 1 (int), 2 (float) or 3 (string)
type StateTypeT int

// States have types determining the value of their type
const (
	StateInt    StateTypeT = 1
	StateFloat  StateTypeT = 2
	StateString StateTypeT = 3
)

// State encodes a device state
type State struct {
	Name  StateNameT
	Type  StateTypeT
	Value interface{}
}
