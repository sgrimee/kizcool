package kizcool

import (
	"encoding/json"
	"errors"
	"net/url"
	"regexp"
	"strings"

	"github.com/sgrimee/kizcool/client"
)

// Kiz provides high-level methods and structs to interact with the server.
type Kiz struct {
	clt client.APIClient
}

// New returns an initialized Kiz
// sessionID is optional and used for external caching of sessions
func New(username, password, baseURL, sessionID string) (*Kiz, error) {
	clt, err := client.New(username, password, baseURL, sessionID)
	if err != nil {
		return nil, err
	}
	return NewWithClient(clt)
}

// NewWithClient returns an initialized Kiz
func NewWithClient(c client.APIClient) (*Kiz, error) {
	k := Kiz{
		clt: c,
	}
	return &k, nil
}

// SessionID is the latest known sessionID value
// It can be used for caching sessions externally.
func (k *Kiz) SessionID() string {
	return k.clt.SessionID()
}

// Login to the server
func (k *Kiz) Login() error {
	return k.clt.Login()
}

// GetDevices returns the list of devices
func (k *Kiz) GetDevices() ([]Device, error) {
	resp, err := k.clt.GetWithAuth("/enduserAPI/setup/devices")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result []Device
	json.NewDecoder(resp.Body).Decode(&result)
	return result, nil
}

// GetDevice returns a single device
func (k *Kiz) GetDevice(deviceURL DeviceURL) (Device, error) {
	query := "/enduserAPI/setup/devices/" + url.QueryEscape(string(deviceURL))
	resp, err := k.clt.GetWithAuth(query)
	if err != nil {
		return Device{}, err
	}
	defer resp.Body.Close()
	var result Device
	json.NewDecoder(resp.Body).Decode(&result)
	return result, nil
}

// DeviceFromListByLabel tries to match the given string to the Labels of the given devices
// and returns the found Device. An error is return is zero or more than one devices match.
func DeviceFromListByLabel(label string, devices []Device) (Device, error) {
	var foundDevice Device
	for _, d := range devices {
		if strings.Compare(strings.ToLower(d.Label), strings.ToLower(label)) == 0 {
			if foundDevice.DeviceURL != "" {
				return Device{}, errors.New("More than one device with that label")
			}
			foundDevice = d
		}
	}
	if foundDevice.DeviceURL == "" {
		return Device{}, errors.New("No device with that label")
	}
	return foundDevice, nil
}

// GetDeviceByText returns a Device from a text string
// If first tries to match a DeviceURL. If no match, it tries to match a device Label
func (k *Kiz) GetDeviceByText(text string) (Device, error) {
	validURL := regexp.MustCompile(`^[a-z]+://\d{4}-\d{4}-\d{4}/\d+`)
	if validURL.MatchString(text) {
		// a DeviceURL was given
		device, err := k.GetDevice(DeviceURL(text))
		if err != nil {
			return Device{}, err
		}
		return device, nil
	}
	// try to match a Label from all devices
	devices, err := k.GetDevices()
	if err != nil {
		return Device{}, err
	}
	device, err := DeviceFromListByLabel(text, devices)
	if err != nil {
		return Device{}, err
	}
	return device, nil
}

// GetDeviceState returns the current state with name stateName for the device with URL deviceURL
func (k *Kiz) GetDeviceState(deviceURL DeviceURL, stateName StateName) (State, error) {
	query := "/enduserAPI/setup/devices/" + url.QueryEscape(string(deviceURL)) +
		"/states/" + url.QueryEscape(string(stateName))
	resp, err := k.clt.GetWithAuth(query)
	if err != nil {
		return State{}, err
	}
	defer resp.Body.Close()
	var result State
	json.NewDecoder(resp.Body).Decode(&result)
	return result, nil
}

// RefreshStates tells the server to refresh states.
// But not sure yet what it really means`?
func (k *Kiz) RefreshStates() error {
	return k.clt.RefreshStates()
}

// GetActionGroups returns the list of action groups defined on the box
func (k *Kiz) GetActionGroups() ([]ActionGroup, error) {
	resp, err := k.clt.GetWithAuth("/enduserAPI/actionGroups")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result []ActionGroup
	json.NewDecoder(resp.Body).Decode(&result)
	return result, nil
}

// SupportsCommand returns true if the command is supported by the device.
func SupportsCommand(device Device, command Command) bool {
	for _, supportedCommand := range device.Definition.Commands {
		if command.Name == supportedCommand.CommandName {
			return true
		}
	}
	return false
}

// ActionGroupWithOneCommand returns an action group with a single command for the device
func ActionGroupWithOneCommand(device Device, command Command) (ActionGroup, error) {
	if !SupportsCommand(device, command) {
		return ActionGroup{}, errors.New("Device does not support this command")
	}
	action := Action{
		DeviceURL: device.DeviceURL,
		Commands:  []Command{command},
	}
	actionGroup := ActionGroup{
		Actions: []Action{action},
	}
	return actionGroup, nil
}

// Execute runs an action group and returns a (job) ExecID
func (k *Kiz) Execute(ag ActionGroup) (ExecID, error) {
	jsonStr, err := json.Marshal(ag)
	if err != nil {
		return "", err
	}
	resp, err := k.clt.Execute(jsonStr)
	if err != nil {
		return "", err
	}
	type Result struct {
		ExecID ExecID
	}
	var result Result
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.ExecID, nil
}

// On turns a device on
func (k *Kiz) On(device Device) (ExecID, error) {
	ag, err := ActionGroupWithOneCommand(device, Command{Name: CmdOn})
	if err != nil {
		return "", err
	}
	return k.Execute(ag)
}

// Off turns a device off
func (k *Kiz) Off(device Device) (ExecID, error) {
	ag, err := ActionGroupWithOneCommand(device, Command{Name: CmdOff})
	if err != nil {
		return "", err
	}
	return k.Execute(ag)
}

// Open opens a device
func (k *Kiz) Open(device Device) (ExecID, error) {
	ag, err := ActionGroupWithOneCommand(device, Command{Name: CmdOpen})
	if err != nil {
		return "", err
	}
	return k.Execute(ag)
}

// Close closes a device
func (k *Kiz) Close(device Device) (ExecID, error) {
	ag, err := ActionGroupWithOneCommand(device, Command{Name: CmdClose})
	if err != nil {
		return "", err
	}
	return k.Execute(ag)
}

// Stop interrupts the current activity
func (k *Kiz) Stop(device Device) (ExecID, error) {
	ag, err := ActionGroupWithOneCommand(device, Command{Name: CmdStop})
	if err != nil {
		return "", err
	}
	return k.Execute(ag)
}

// SetIntensity sets the light intensity to given value
func (k *Kiz) SetIntensity(device Device, intensity int) (ExecID, error) {
	ag, err := ActionGroupWithOneCommand(device, Command{
		Name:       CmdSetIntensity,
		Parameters: []int{intensity},
	})
	if err != nil {
		return "", err
	}
	return k.Execute(ag)
}

// SetClosure sets the device closure/position to given value
func (k *Kiz) SetClosure(device Device, position int) (ExecID, error) {
	ag, err := ActionGroupWithOneCommand(device, Command{
		Name:       CmdSetClosure,
		Parameters: []int{position},
	})
	if err != nil {
		return "", err
	}
	return k.Execute(ag)
}
