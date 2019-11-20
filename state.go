package kizcool

// StateDefinition describes the fields of a State
type StateDefinition struct {
	Type          string
	QualifiedName string
	Values        []string
}

// StateName is the name of a State
type StateName string

// StateType has value 1 (int), 2 (float) or 3 (string)
type StateType int

// States have types determining the value of their type
const (
	StateInt    StateType = 1
	StateFloat  StateType = 2
	StateString StateType = 3
)

// DeviceState encodes a device state
type DeviceState struct {
	Name  StateName
	Type  StateType
	Value interface{}
}
