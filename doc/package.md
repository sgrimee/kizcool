# kizcool
--
    import "github.com/sgrimee/kizcool"

Package kizcool provides a client to the Overkiz IoT API, used by velux, somfy
and other vendors to control velux devices with a Tahoma box.

## Usage

```go
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
)
```
CommandNames

#### func  Output

```go
func Output(w io.Writer, format string, obj interface{}) error
```
Output prints the given object to the writer in the desired format format can be
'text', 'json' or 'yaml'

#### func  SupportsCommand

```go
func SupportsCommand(device Device, command Command) bool
```
SupportsCommand returns true if the command is supported by the device.

#### type Action

```go
type Action struct {
	DeviceURL DeviceURL `json:"deviceURL,omitempty"`
	Commands  []Command `json:"commands,omitempty"`
}
```

Action defines a list of commands

#### type ActionGroup

```go
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
```

ActionGroup is a list of Actions in sequence, with metadata. Think "scenario".

#### func  ActionGroupWithOneCommand

```go
func ActionGroupWithOneCommand(device Device, command Command) (ActionGroup, error)
```
ActionGroupWithOneCommand returns an action group with a single command for the
device

#### type Command

```go
type Command struct {
	Type       int         `json:"type,omitempty"`
	Name       string      `json:"name,omitempty"`
	Parameters interface{} `json:"parameters,omitempty"`
}
```

Command describes a command (duh)

#### type CommandDefinition

```go
type CommandDefinition struct {
	CommandName string `json:"commandName,omitempty"`
	Nparams     int    `json:"nparams,omitempty"`
}
```

CommandDefinition describes the fields of a Command

#### type Device

```go
type Device struct {
	CreationTime     int
	LastUpdateTime   int
	Label            string
	DeviceURL        DeviceURL
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
}
```

Device representation of a single device

#### func  DeviceFromListByLabel

```go
func DeviceFromListByLabel(label string, devices []Device) (Device, error)
```
DeviceFromListByLabel tries to match the given string to the Labels of the given
devices and returns the found Device. An error is return is zero or more than
one devices match.

#### type DeviceDefinition

```go
type DeviceDefinition struct {
	Commands      []CommandDefinition
	States        []StateDefinition
	WidgetName    string
	UIClass       string
	QualifiedName string
	Type          string
}
```

DeviceDefinition describes the fields of a Device

#### type DeviceURL

```go
type DeviceURL string
```

DeviceURL is the full device URL including prefix e.g.
io://1111-0000-4444/12345678

#### type ExecID

```go
type ExecID string
```

ExecID is the id of an execution (job)

#### type Kiz

```go
type Kiz struct {
}
```

Kiz provides high-level methods and structs to interact with the server.

#### func  New

```go
func New(username, password, baseURL, sessionID string) (*Kiz, error)
```
New returns an initialized Kiz sessionID is optional and used for external
caching of sessions

#### func  NewWithClient

```go
func NewWithClient(c client.APIClient) *Kiz
```
NewWithClient returns a Kiz with the given pre-initialised APIClient

#### func (*Kiz) Close

```go
func (k *Kiz) Close(device Device) (ExecID, error)
```
Close closes a device

#### func (*Kiz) Execute

```go
func (k *Kiz) Execute(ag ActionGroup) (ExecID, error)
```
Execute runs an action group and returns a (job) ExecID

#### func (*Kiz) GetActionGroups

```go
func (k *Kiz) GetActionGroups() ([]ActionGroup, error)
```
GetActionGroups returns the list of action groups defined on the box

#### func (*Kiz) GetDevice

```go
func (k *Kiz) GetDevice(deviceURL DeviceURL) (Device, error)
```
GetDevice returns a single device

#### func (*Kiz) GetDeviceByText

```go
func (k *Kiz) GetDeviceByText(text string) (Device, error)
```
GetDeviceByText returns a Device from a text string If first tries to match a
DeviceURL. If no match, it tries to match a device Label

#### func (*Kiz) GetDeviceState

```go
func (k *Kiz) GetDeviceState(deviceURL DeviceURL, stateName StateName) (State, error)
```
GetDeviceState returns the current state with name stateName for the device with
URL deviceURL

#### func (*Kiz) GetDevices

```go
func (k *Kiz) GetDevices() ([]Device, error)
```
GetDevices returns the list of devices

#### func (*Kiz) Login

```go
func (k *Kiz) Login() error
```
Login to the server

#### func (*Kiz) Off

```go
func (k *Kiz) Off(device Device) (ExecID, error)
```
Off turns a device off

#### func (*Kiz) On

```go
func (k *Kiz) On(device Device) (ExecID, error)
```
On turns a device on

#### func (*Kiz) Open

```go
func (k *Kiz) Open(device Device) (ExecID, error)
```
Open opens a device

#### func (*Kiz) RefreshStates

```go
func (k *Kiz) RefreshStates() error
```
RefreshStates tells the server to refresh states. But not sure yet what it
really means`?

#### func (*Kiz) SetClosure

```go
func (k *Kiz) SetClosure(device Device, position int) (ExecID, error)
```
SetClosure sets the device closure/position to given value

#### func (*Kiz) SetIntensity

```go
func (k *Kiz) SetIntensity(device Device, intensity int) (ExecID, error)
```
SetIntensity sets the light intensity to given value

#### func (*Kiz) Stop

```go
func (k *Kiz) Stop(device Device) (ExecID, error)
```
Stop interrupts the current activity

#### type State

```go
type State struct {
	Name  StateName
	Type  StateType
	Value interface{}
}
```

State encodes a device state

#### type StateDefinition

```go
type StateDefinition struct {
	Type          string
	QualifiedName string
	Values        []string
}
```

StateDefinition describes the fields of a State

#### type StateName

```go
type StateName string
```

StateName is the name of a State

#### type StateType

```go
type StateType int
```

StateType has value 1 (int), 2 (float) or 3 (string)

```go
const (
	StateInt    StateType = 1
	StateFloat  StateType = 2
	StateString StateType = 3
)
```
States have types determining the value of their type
