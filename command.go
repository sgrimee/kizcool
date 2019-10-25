package kizcool

// CommandDefinition describes the fields of a Command
type CommandDefinition struct {
	CommandName string
	Nparams     int
}

// Command describes a command (duh)
type Command struct {
	Type       int
	Name       string
	Parameters []interface{}
}
