package kizcool

// CommandDefinition describes the fields of a Command
type CommandDefinition struct {
	CommandName string `json:"commandName,omitempty"`
	Nparams     int    `json:"nparams,omitempty"`
}

// Command describes a command (duh)
type Command struct {
	Type       int         `json:"type,omitempty"`
	Name       string      `json:"name,omitempty"`
	Parameters interface{} `json:"parameters,omitempty"`
}

// CommandNames
const (
	CmdClose         = "close"
	CmdDown          = "down"
	CmdIdentify      = "identify"
	CmdOff           = "off"
	CmdOn            = "on"
	CmdOpen          = "open"
	CmdSetIntensity  = "setIntensity"
	CmdSetClosure    = "setClosure"
	CmdStartIdentify = "startIdentify"
	CmdStop          = "stop"
	CmdStopIdentify  = "stopIdentify"
	CmdUp            = "up"

// "activateCalendar"
// "alarmOff"
// "alarmOn"
// "alarmPartial1"
// "alarmPartial2"
// "deactivateCalendar"
// "delayedStopIdentify"
// "getName"
// "my"
// "onWithTimer"
// "refreshAlarmDelay"
// "refreshBatteryStatus"
// "refreshCurrentAlarmMode"
// "refreshIntrusionDetected"
// "refreshMemorized1Position"
// "refreshPodMode"
// "refreshUpdateStatus"
// "setAlarmDelay"
// "setCalendar"
// "setCountryCode"
// "setDeployment"
// "setIntensityWithTimer"
// "setIntrusionDetected"
// "setLightingLedPodMode"
// "SetPosition"
// "setMemorized1Position"
// "setName"
// "setOnOff"
// "setPodLedOff"
// "setPodLedOn"
// "setSecuredPosition"
// "setTargetAlarmMode"
// "update"
// "wink"
)
