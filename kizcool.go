// Package kizcool provides access to the overkiz API to control velux devices
// with a tahoma box
package kizcool

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"golang.org/x/net/publicsuffix"
)

// The Kiz api provides methods to call api endpoints.
// It should be created with NewKiz, not used directly
type Kiz struct {
	username string
	password string
	BaseURL  string // defaults to tahomalink
	Client   *http.Client
}

// NewKiz Return a new Kiz with the cookie jar and http client set up
func NewKiz(username string, password string) Kiz {
	// All users of cookiejar should import "golang.org/x/net/publicsuffix"
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{
		Jar: jar,
	}

	k := Kiz{
		username,
		password,
		"https://tahomalink.com/enduser-mobile-web/enduserAPI",
		client,
	}
	return k
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
func (k Kiz) Login() error {
	formData := url.Values{"userId": {k.username},
		"userPassword": {k.password}}
	resp, err := k.Client.PostForm(k.BaseURL+"/login", formData)
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

// GetDevices returns the list of devices
func (k Kiz) GetDevices() ([]Device, error) {
	resp, err := k.Client.Get(k.BaseURL + "/setup/devices")
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

// GetActionGroups returns the list of action groups defined on the box
func (k Kiz) GetActionGroups() ([]ActionGroup, error) {
	resp, err := k.Client.Get(k.BaseURL + "/actionGroups")
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
