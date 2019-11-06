// Package kizcool provides access to the overkiz API to control velux devices
// with a tahoma box
package kizcool

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// The Kiz api provides methods to call api endpoints.
// It should be created with New, not used directly
type Kiz struct {
	config Config
	api    *API
}

// New returns an initialized Kiz
func New(config Config) (*Kiz, error) {
	api, err := NewAPI(config)
	if err != nil {
		return nil, err
	}
	k := Kiz{
		config: config,
		api:    api,
	}
	return &k, nil
}

// Login and get a session ID
// The new sessionID is stored/updated in the configfile
func (k *Kiz) Login() error {
	_, err := k.api.Login()
	if err != nil {
		return err
	}
	return nil
}

// GetDevices returns the list of devices
func (k *Kiz) GetDevices() ([]Device, error) {
	resp, err := k.api.GetWithAuth("/enduserAPI/setup/devices")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result []Device
	json.NewDecoder(resp.Body).Decode(&result)
	return result, nil
}

// GetDevice returns a single device
func (k *Kiz) GetDevice(deviceURL DeviceURLT) (Device, error) {
	query := "/enduserAPI/setup/devices/" + url.QueryEscape(string(deviceURL))
	resp, err := k.api.GetWithAuth(query)
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
		device, err := k.GetDevice(DeviceURLT(text))
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
func (k *Kiz) GetDeviceState(deviceURL DeviceURLT, stateName StateNameT) (State, error) {
	query := "/enduserAPI/setup/devices/" + url.QueryEscape(string(deviceURL)) +
		"/states/" + url.QueryEscape(string(stateName))
	resp, err := k.api.GetWithAuth(query)
	if err != nil {
		return State{}, err
	}
	defer resp.Body.Close()
	var result State
	json.NewDecoder(resp.Body).Decode(&result)
	return result, nil
}

// GetActionGroups returns the list of action groups defined on the box
func (k *Kiz) GetActionGroups() ([]ActionGroup, error) {
	resp, err := k.api.GetWithAuth("/enduserAPI/actionGroups")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result []ActionGroup
	json.NewDecoder(resp.Body).Decode(&result)
	return result, nil
}

// RefreshStates tells the server to refresh states.
// But not sure yet what it really means`?
func (k *Kiz) RefreshStates() error {
	req, err := http.NewRequest(http.MethodPut, k.config.BaseURL+"/enduserAPI/setup/devices/states/refresh", nil)
	if err != nil {
		return err
	}
	resp, err := k.api.DoWithAuth(req)
	if err != nil {
		return err
	}
	if e := checkStatusOk(resp); e != nil {
		return e
	}
	return nil
}

// Execute initiates the execution of a group of actions
func (k *Kiz) Execute(ag ActionGroup) (ExecIDT, error) {
	jsonStr, err := json.Marshal(ag)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(http.MethodPost, k.config.BaseURL+"/enduserAPI/exec/apply", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return "", err
	}
	resp, err := k.api.DoWithAuth(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	type Result struct {
		ExecID ExecIDT
	}
	var result Result
	json.NewDecoder(resp.Body).Decode(&result)
	return result.ExecID, nil
}

func supportsCommand(device Device, command Command) bool {
	for _, supportedCommand := range device.Definition.Commands {
		if command.Name == supportedCommand.CommandName {
			return true
		}
	}
	return false
}

// ActionGroupWithOneCommand returns an action group with a single command for the device
func ActionGroupWithOneCommand(device Device, command Command) (ActionGroup, error) {
	if !supportsCommand(device, command) {
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

// On turns a device on
func (k *Kiz) On(device Device) (ExecIDT, error) {
	ag, err := ActionGroupWithOneCommand(device, Command{Name: CmdOn})
	if err != nil {
		return "", err
	}
	return k.Execute(ag)
}

// Off turns a device off
func (k *Kiz) Off(device Device) (ExecIDT, error) {
	ag, err := ActionGroupWithOneCommand(device, Command{Name: CmdOff})
	if err != nil {
		return "", err
	}
	return k.Execute(ag)
}

// Open opens a device
func (k *Kiz) Open(device Device) (ExecIDT, error) {
	ag, err := ActionGroupWithOneCommand(device, Command{Name: CmdOpen})
	if err != nil {
		return "", err
	}
	return k.Execute(ag)
}

// Close closes a device
func (k *Kiz) Close(device Device) (ExecIDT, error) {
	ag, err := ActionGroupWithOneCommand(device, Command{Name: CmdClose})
	if err != nil {
		return "", err
	}
	return k.Execute(ag)
}

// Stop interrupts the current activity
func (k *Kiz) Stop(device Device) (ExecIDT, error) {
	ag, err := ActionGroupWithOneCommand(device, Command{Name: CmdStop})
	if err != nil {
		return "", err
	}
	return k.Execute(ag)
}

// SetIntensity sets the light intensity to given value
func (k *Kiz) SetIntensity(device Device, intensity int) (ExecIDT, error) {
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
func (k *Kiz) SetClosure(device Device, position int) (ExecIDT, error) {
	ag, err := ActionGroupWithOneCommand(device, Command{
		Name:       CmdSetClosure,
		Parameters: []int{position},
	})
	if err != nil {
		return "", err
	}
	return k.Execute(ag)
}
