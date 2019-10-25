package kizcool

// StateDefinition describes the fields of a State
type StateDefinition struct {
	Type          string
	QualifiedName string
	Values        []string
}

// State encodes a device state
type State struct {
	Name  string
	Type  int
	Value interface{}
}
