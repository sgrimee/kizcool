// Package kizcool provides access to the overkiz API to control velux devices
// with a tahoma box
package kizcool

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"golang.org/x/net/publicsuffix"
)

// The Kiz api provides methods to call api endpoints.
// It should be created with New, not used directly
type Kiz struct {
	Debug bool

	config Config
	client *http.Client
}

// New Return a new Kiz with the cookie jar and http client set up
func New(config Config) (*Kiz, error) {
	// All users of cookiejar should import "golang.org/x/net/publicsuffix"
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Jar: jar,
	}

	k := Kiz{
		config: config,
		Debug:  config.Debug,
		client: client,
	}
	return &k, nil
}

// checkStatusOk performs simple tests to ensure the request was successful
// if an error occured, try to qualify it then return it. If not return nil.
func checkStatusOk(resp *http.Response) error {
	if resp.StatusCode == 401 {
		err := &AuthenticationError{401, "Authentication error"}
		return err
	}
	if resp.StatusCode != 200 {
		// Decode the body to try to get a meaningful error message
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		msg := result["error"]
		if msg == nil {
			msg = result
		}
		return fmt.Errorf("%v", msg)
	}
	return nil
}

// Login and get a session cookie
func (k *Kiz) Login() error {
	formData := url.Values{"userId": {k.config.Username}, "userPassword": {k.config.Password}}
	resp, err := k.client.PostForm(k.config.BaseURL+"/enduserAPI/login", formData)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if e := checkStatusOk(resp); e != nil {
		return e
	}
	for _, c := range resp.Cookies() {
		if (c.Name == "JSESSIONID") && (c.Value != "") {
			return nil
		}
	}
	return errors.New("JSESSIONID not found in response to /login")
}

// getSetup returns informations about the site setup, including
// devices, device location in rooms, etc
// NOT FULLY IMPLEMEMNTED and may never be. Included for documentation only.
func (k *Kiz) getSetup() ([]interface{}, error) {
	resp, err := k.client.Get(k.config.BaseURL + "/externalAPI/json/getSetup")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if e := checkStatusOk(resp); e != nil {
		return nil, e
	}
	var result []interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	return result, nil
}

// GetDevices returns the list of devices
func (k *Kiz) GetDevices() ([]Device, error) {
	resp, err := k.client.Get(k.config.BaseURL + "/enduserAPI/setup/devices")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if e := checkStatusOk(resp); e != nil {
		return nil, e
	}
	var result []Device
	json.NewDecoder(resp.Body).Decode(&result)
	return result, nil
}

// GetDevice returns a single device
func (k *Kiz) GetDevice(deviceURL DeviceURLT) (Device, error) {
	query := "/enduserAPI/setup/devices/" + url.QueryEscape(string(deviceURL))
	resp, err := k.client.Get(k.config.BaseURL + query)
	if err != nil {
		return Device{}, err
	}
	defer resp.Body.Close()
	if e := checkStatusOk(resp); e != nil {
		return Device{}, e
	}
	var result Device
	json.NewDecoder(resp.Body).Decode(&result)
	return result, nil
}

// DeviceURLByLabel tries to match the given string to the Labels of the given devices
// and returns its DeviceURL. An error is return is zero or more than one devices match.
func DeviceURLByLabel(label string, devices []Device) (DeviceURLT, error) {
	var foundDevice Device
	for _, d := range devices {
		if strings.Compare(strings.ToLower(d.Label), strings.ToLower(label)) == 0 {
			if foundDevice.DeviceURL != "" {
				return "", errors.New("More than one device with that label")
			}
			foundDevice = d
		}
	}
	if foundDevice.DeviceURL == "" {
		return "", errors.New("No device with that label")
	}
	return foundDevice.DeviceURL, nil
}

// GetDeviceState returns the current state with name stateName for the device with URL deviceURL
func (k *Kiz) GetDeviceState(deviceURL DeviceURLT, stateName StateNameT) (State, error) {
	query := "/enduserAPI/setup/devices/" + url.QueryEscape(string(deviceURL)) +
		"/states/" + url.QueryEscape(string(stateName))
	resp, err := k.client.Get(k.config.BaseURL + query)
	if err != nil {
		return State{}, err
	}
	defer resp.Body.Close()
	if e := checkStatusOk(resp); e != nil {
		return State{}, e
	}
	var result State
	json.NewDecoder(resp.Body).Decode(&result)
	return result, nil
}

// GetActionGroups returns the list of action groups defined on the box
func (k *Kiz) GetActionGroups() ([]ActionGroup, error) {
	resp, err := k.client.Get(k.config.BaseURL + "/enduserAPI/actionGroups")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if e := checkStatusOk(resp); e != nil {
		return nil, e
	}
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
	resp, err := k.client.Do(req)
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
	req, err := http.NewRequest("POST", k.config.BaseURL+"/enduserAPI/exec/apply", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return "", err
	}
	resp, err := k.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if err := checkStatusOk(resp); err != nil {
		return "", err
	}
	type Result struct {
		ExecID ExecIDT
	}
	var result Result
	json.NewDecoder(resp.Body).Decode(&result)
	return result.ExecID, nil
}
