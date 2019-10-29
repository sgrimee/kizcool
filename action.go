package kizcool

// ExecIDT is the id of an execution (job)
type ExecIDT string

// ActionGroup is a list of Actions in sequence, with metadata. Think "scenario".
type ActionGroup struct {
	CreationTime          int
	LastUpdateTime        int
	Label                 string
	Shortcut              bool
	NotificationTypeMask  int
	NotificationCondition string
	Actions               []Action
	OID                   string
}

// Action defines a list of commands
type Action struct {
	DeviceURL string
	Commands  []Command
}
