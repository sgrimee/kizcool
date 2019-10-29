package kizcool

// ExecIDT is the id of an execution (job)
type ExecIDT string

// ActionGroup is a list of Actions in sequence, with metadata. Think "scenario".
type ActionGroup struct {
	CreationTime          int      `json:"creationTime,omitempty"`
	LastUpdateTime        int      `json:"lastUpdateTime,omitempty"`
	Label                 string   `json:"label,omitempty"`
	Shortcut              bool     `json:"shortcut,omitempty"`
	NotificationTypeMask  int      `json:"notificationTypeMask,omitempty"`
	NotificationCondition string   `json:"notificationCondition,omitempty"`
	Actions               []Action `json:"actions,omitempty"`
	OID                   string   `json:"oid,omitempty"`
}

// Action defines a list of commands
type Action struct {
	DeviceURL DeviceURLT `json:"deviceURL,omitempty"`
	Commands  []Command  `json:"commands,omitempty"`
}
