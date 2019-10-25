package kizcool

// Device representation of a single device
type Device struct {
	CreationTime     int
	LastUpdateTime   int
	Label            string
	DeviceURL        string
	Shortcut         bool
	ControllableName string
	Definition       DeviceDefinition
	States           []State
	Available        bool
	Enabled          bool
	PlaceOID         string
	Widget           string
	Type             int
	OID              string
	UIClass          string
	//Attributes       []interface{}
}

// DeviceDefinition describes the fields of a Device
type DeviceDefinition struct {
	Commands      []CommandDefinition
	States        []StateDefinition
	WidgetName    string
	UIClass       string
	QualifiedName string
	Type          string
	//DataProperties []
}
