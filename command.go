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

// CommandNames
const (
	CmdClose         = "close"
	CmdDown          = "down"
	CmdIdentify      = "identify"
	CmdOff           = "off"
	CmdOn            = "on"
	CmdOpen          = "open"
	CmdSetIntensity  = "setIntensity"
	CmdSetPosition   = "setPosition"
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
// "setClosure"
// "setCountryCode"
// "setDeployment"
// "setIntensityWithTimer"
// "setIntrusionDetected"
// "setLightingLedPodMode"
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
