package kizcool

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// var kiz *Kiz

func TestBadLogin(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/enduserAPI/login", req.URL.String())
		rw.WriteHeader(401)
		rw.Write([]byte(`{"errorCode": "AUTHENTICATION_ERROR","error": "Bad credentials"}`))
	}))
	defer server.Close()
	kiz, _ := New(Config{"baduser", "badpass", server.URL, false})
	kiz.client = server.Client()
	err := kiz.Login()
	assert.EqualError(t, err, "401: Authentication error")
}

func TestGoodLogin(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/enduserAPI/login", req.URL.String())
		cookie := http.Cookie{
			Name:    "JSESSIONID",
			Value:   "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF",
			Expires: time.Now().AddDate(0, 0, 1),
		}
		http.SetCookie(rw, &cookie)
		rw.Write([]byte(`{"success":true,"roles":[{"name":"ENDUSER"}]}`))
	}))
	kiz, _ := New(Config{"gooduser", "goodpass", server.URL, false})
	kiz.client = server.Client()
	err := kiz.Login()
	assert.Nil(t, err)
}

func helperLoadBytes(t *testing.T, name string) []byte {
	path := filepath.Join("testdata", name) // relative path
	bytes, err := ioutil.ReadFile(path)
	assert.Nil(t, err)
	return bytes
}

func TestGetSetup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/externalAPI/json/getSetup", req.URL.String())
		//rw.Write(helperLoadBytes(t, "getSetup.json"))
	}))
	kiz, _ := New(Config{"gooduser", "goodpass", server.URL, false})
	kiz.client = server.Client()
	_, err := kiz.getSetup()
	assert.Nil(t, err)
}

func TestGetDevices(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/enduserAPI/setup/devices", req.URL.String())
		rw.Write(helperLoadBytes(t, "getDevices.json"))
	}))
	kiz, _ := New(Config{"gooduser", "goodpass", server.URL, false})
	kiz.client = server.Client()
	devices, err := kiz.GetDevices()
	assert.Nil(t, err)
	assert.Equal(t, len(devices), 5)
}

func TestGetDevice(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/enduserAPI/setup/devices/io%3A%2F%2F1111-0000-4444%2F11784413", req.URL.String())
		rw.Write(helperLoadBytes(t, "getDevice.json"))
	}))
	kiz, _ := New(Config{"gooduser", "goodpass", server.URL, false})
	kiz.client = server.Client()
	device, err := kiz.GetDevice("io://1111-0000-4444/11784413")
	assert.Nil(t, err)
	assert.Equal(t, len(device.States), 5)
}

func TestDeviceFromListByLabel(t *testing.T) {
	device1 := Device{
		Label:     "label1",
		DeviceURL: "url1",
	}
	device2 := Device{
		Label:     "label2",
		DeviceURL: "url2",
	}
	devices := []Device{device1, device2, device2}

	// case found
	d, err := DeviceFromListByLabel("label1", devices)
	assert.Nil(t, err)
	assert.Equal(t, DeviceURLT("url1"), d.DeviceURL)

	// case not found
	_, err = DeviceFromListByLabel("bogus", devices)
	assert.NotNil(t, err)

	// case multiple found
	_, err = DeviceFromListByLabel("label2", devices)
	assert.NotNil(t, err)
}

func TestGetDeviceState(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/enduserAPI/setup/devices/io%3A%2F%2F1111-0000-4444%2F12345678/states/core%3AOnOffState", req.URL.String())
		rw.Write([]byte(`{"name": "core:OnOffState","type": 3,"value": "off"}`))
	}))
	kiz, _ := New(Config{"gooduser", "goodpass", server.URL, false})
	kiz.client = server.Client()
	state, err := kiz.GetDeviceState("io://1111-0000-4444/12345678", "core:OnOffState")
	assert.Nil(t, err)
	assert.Equal(t, "off", state.Value)
}

func TestGetActionGroups(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/enduserAPI/actionGroups", req.URL.String())
		rw.Write(helperLoadBytes(t, "getActionGroups.json"))
	}))
	kiz, _ := New(Config{"gooduser", "goodpass", server.URL, false})
	kiz.client = server.Client()
	actionGroups, err := kiz.GetActionGroups()
	assert.Nil(t, err)
	assert.Equal(t, len(actionGroups), 1)
	assert.Equal(t, len(actionGroups[0].Actions), 2)
}

func TestRefreshStates(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/enduserAPI/setup/devices/states/refresh", req.URL.String())
	}))
	kiz, _ := New(Config{"gooduser", "goodpass", server.URL, false})
	kiz.client = server.Client()
	err := kiz.RefreshStates()
	assert.Nil(t, err)
}

func TestExecute(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/enduserAPI/exec/apply", req.URL.String())
		var ag ActionGroup
		err := json.NewDecoder(req.Body).Decode(&ag)
		assert.Nil(t, err)
		assert.Equal(t, DeviceURLT("io://1111-0000-4444/12345678"), ag.Actions[0].DeviceURL)
		assert.Equal(t, "on", ag.Actions[0].Commands[0].Name)
		rw.Write([]byte(`{"execId": "133a5c55-3655-5455-2355-c33e43535e55"}`))
	}))
	kiz, _ := New(Config{"gooduser", "goodpass", server.URL, false})
	kiz.client = server.Client()
	command := Command{
		Name: CmdOn,
	}
	action := Action{
		DeviceURL: "io://1111-0000-4444/12345678",
		Commands:  []Command{command},
	}
	actionGroup := ActionGroup{
		Actions: []Action{action},
	}
	id, err := kiz.Execute(actionGroup)
	assert.Nil(t, err)
	assert.Equal(t, 36, len(string(id)))
}
