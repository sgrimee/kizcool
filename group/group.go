package group

import "github.com/vapourismo/knx-go/knx/cemi"

// Device maps a device with Url and optional Label tp a list of Command and State
type Device struct {
	Label    string
	URL      string
	Commands []Command
	States   []State
}

// Command maps a command name to a group address
type Command struct {
	Name      string
	GroupAddr cemi.GroupAddr
}

// State maps a state name to a group address
type State struct {
	Name      string
	GroupAddr cemi.GroupAddr
}
