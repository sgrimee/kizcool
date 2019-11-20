package kizcool

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/sgrimee/kizcool/api"
	"github.com/stretchr/testify/assert"
)

// helperLoadBytes loads test data from a file
func helperLoadBytes(t *testing.T, name string) []byte {
	path := filepath.Join("testdata", name) // relative path
	bytes, err := ioutil.ReadFile(path)
	assert.Nil(t, err)
	return bytes
}

func TestNew(t *testing.T) {
	_, err := New("gooduser", "goodpass", "", "")
	assert.NoError(t, err)
}

func TestBadLogin(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/enduserAPI/login", req.URL.String())
		rw.WriteHeader(401)
		rw.Write([]byte(`{"errorCode": "AUTHENTICATION_ERROR","error": "Bad credentials"}`))
	}))
	defer server.Close()
	ac, err := api.NewWithHTTPClient("baduser", "badpass", server.URL, "", server.Client())
	assert.NoError(t, err)
	kiz, _ := NewWithAPIClient(ac)
	err = kiz.Login()
	assert.Equal(t, err, api.NewAuthenticationError("Bad credentials"))
}

func TestGoodLogin(t *testing.T) {
	const sessionID = "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/enduserAPI/login", req.URL.String())
		cookie := http.Cookie{
			Name:    "JSESSIONID",
			Value:   sessionID,
			Expires: time.Now().AddDate(0, 0, 1),
		}
		http.SetCookie(rw, &cookie)
		rw.Write([]byte(`{"success":true,"roles":[{"name":"ENDUSER"}]}`))
	}))
	defer server.Close()
	ac, err := api.NewWithHTTPClient("gooduser", "goodpass", server.URL, "", server.Client())
	assert.NoError(t, err)
	kiz, _ := NewWithAPIClient(ac)
	err = kiz.Login()
	assert.Nil(t, err)
	assert.Equal(t, sessionID, kiz.SessionID())
}

func TestGetDevices(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/enduserAPI/setup/devices", req.URL.String())
		rw.Write(helperLoadBytes(t, "getDevices.json"))
	}))
	defer server.Close()
	ac, err := api.NewWithHTTPClient("gooduser", "goodpass", server.URL, "", server.Client())
	assert.NoError(t, err)
	kiz, _ := NewWithAPIClient(ac)
	devices, err := kiz.GetDevices()
	assert.Nil(t, err)
	assert.Equal(t, len(devices), 5)
}

func TestGetDevice(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/enduserAPI/setup/devices/io%3A%2F%2F1111-0000-4444%2F11784413", req.URL.String())
		rw.Write(helperLoadBytes(t, "getDevice.json"))
	}))
	defer server.Close()
	ac, err := api.NewWithHTTPClient("gooduser", "goodpass", server.URL, "", server.Client())
	assert.NoError(t, err)
	kiz, _ := NewWithAPIClient(ac)
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
	assert.Equal(t, DeviceURL("url1"), d.DeviceURL)

	// case not found
	_, err = DeviceFromListByLabel("bogus", devices)
	assert.NotNil(t, err)

	// case multiple found
	_, err = DeviceFromListByLabel("label2", devices)
	assert.NotNil(t, err)
}

func TestGetDeviceByTextMatchText(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/enduserAPI/setup/devices", req.URL.String())
		rw.Write(helperLoadBytes(t, "getDevices.json"))
	}))
	defer server.Close()
	ac, err := api.NewWithHTTPClient("gooduser", "goodpass", server.URL, "", server.Client())
	assert.NoError(t, err)
	kiz, _ := NewWithAPIClient(ac)
	device, err := kiz.GetDeviceByText("fenetre1")
	assert.Nil(t, err)
	assert.Equal(t, device.Label, "Fenetre1")
}

func TestGetDeviceByTextURI(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/enduserAPI/setup/devices/io%3A%2F%2F1111-0000-4444%2F11784413", req.URL.String())
		rw.Write(helperLoadBytes(t, "getDevice.json"))
	}))
	defer server.Close()
	ac, err := api.NewWithHTTPClient("gooduser", "goodpass", server.URL, "", server.Client())
	assert.NoError(t, err)
	kiz, _ := NewWithAPIClient(ac)
	device, err := kiz.GetDeviceByText("io://1111-0000-4444/11784413")
	assert.Nil(t, err)
	assert.Equal(t, device.Label, "Fenetre1")
}

func TestGetDeviceState(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/enduserAPI/setup/devices/io%3A%2F%2F1111-0000-4444%2F12345678/states/core%3AOnOffState", req.URL.String())
		rw.Write([]byte(`{"name": "core:OnOffState","type": 3,"value": "off"}`))
	}))
	defer server.Close()
	ac, err := api.NewWithHTTPClient("gooduser", "goodpass", server.URL, "", server.Client())
	assert.NoError(t, err)
	kiz, _ := NewWithAPIClient(ac)
	state, err := kiz.GetDeviceState("io://1111-0000-4444/12345678", "core:OnOffState")
	assert.Nil(t, err)
	assert.Equal(t, "off", state.Value)
}

func TestRefreshStates(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/enduserAPI/setup/devices/states/refresh", req.URL.String())
	}))
	defer server.Close()
	ac, err := api.NewWithHTTPClient("gooduser", "goodpass", server.URL, "", server.Client())
	assert.NoError(t, err)
	kiz, _ := NewWithAPIClient(ac)
	err = kiz.RefreshStates()
	assert.Nil(t, err)
}

func TestActionGroupWithOneCommandUnsupported(t *testing.T) {
	command := Command{
		Name: "unsupportedCmd",
	}
	device := Device{
		Definition: DeviceDefinition{
			Commands: []CommandDefinition{},
		},
	}
	_, err := ActionGroupWithOneCommand(device, command)
	assert.Error(t, err)
}

func TestGetActionGroups(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/enduserAPI/actionGroups", req.URL.String())
		rw.Write(helperLoadBytes(t, "getActionGroups.json"))
	}))
	defer server.Close()
	ac, err := api.NewWithHTTPClient("gooduser", "goodpass", server.URL, "", server.Client())
	assert.NoError(t, err)
	kiz, _ := NewWithAPIClient(ac)
	actionGroups, err := kiz.GetActionGroups()
	assert.Nil(t, err)
	assert.Equal(t, len(actionGroups), 1)
	assert.Equal(t, len(actionGroups[0].Actions), 2)
}

func TestSupportsCommand(t *testing.T) {
	goodCmdDef := CommandDefinition{
		CommandName: "goodCmd",
	}
	device := Device{
		Definition: DeviceDefinition{
			Commands: []CommandDefinition{goodCmdDef},
		},
	}
	assert.True(t, SupportsCommand(device, Command{Name: "goodCmd"}))
	assert.False(t, SupportsCommand(device, Command{Name: "badCmd"}))
}

func TestExecute(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/enduserAPI/exec/apply", req.URL.String())
		var ag ActionGroup
		err := json.NewDecoder(req.Body).Decode(&ag)
		assert.Nil(t, err)
		assert.Equal(t, DeviceURL("io://1111-0000-4444/12345678"), ag.Actions[0].DeviceURL)
		rw.Write([]byte(`{"execId": "133a5c55-3655-5455-2355-c33e43535e55"}`))
	}))
	defer server.Close()
	ac, err := api.NewWithHTTPClient("gooduser", "goodpass", server.URL, "", server.Client())
	assert.NoError(t, err)
	kiz, _ := NewWithAPIClient(ac)
	device := Device{
		DeviceURL: "io://1111-0000-4444/12345678",
		Definition: DeviceDefinition{
			Commands: []CommandDefinition{
				{CommandName: CmdOn},
				{CommandName: CmdOff},
				{CommandName: CmdOpen},
				{CommandName: CmdClose},
				{CommandName: CmdStop},
				{CommandName: CmdSetIntensity},
				{CommandName: CmdSetClosure},
			},
		},
	}
	actionGroup, err := ActionGroupWithOneCommand(device, Command{Name: CmdOn})
	assert.NoError(t, err)
	id, err := kiz.Execute(actionGroup)
	assert.NoError(t, err)
	assert.Equal(t, 36, len(string(id)))

	// also test the shortcuts
	_, err = kiz.On(device)
	assert.NoError(t, err)

	_, err = kiz.Off(device)
	assert.NoError(t, err)

	_, err = kiz.Open(device)
	assert.NoError(t, err)

	_, err = kiz.Close(device)
	assert.NoError(t, err)

	_, err = kiz.Stop(device)
	assert.NoError(t, err)

	_, err = kiz.SetIntensity(device, 0)
	assert.NoError(t, err)

	_, err = kiz.SetClosure(device, 0)
	assert.NoError(t, err)
}

func TestPollEventsAllKnownEventTypes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write(helperLoadBytes(t, "pollEvents.json"))
	}))
	defer server.Close()
	ac, err := api.NewWithHTTPClient("gooduser", "goodpass", server.URL, "", server.Client())
	assert.NoError(t, err)
	ac.SetListenerID("not_empty")
	kiz, _ := NewWithAPIClient(ac)
	events, err := kiz.PollEvents()
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(events), 0)
}

func TestPollEventsNoEvent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(`[]`))
	}))
	defer server.Close()
	ac, err := api.NewWithHTTPClient("gooduser", "goodpass", server.URL, "", server.Client())
	assert.NoError(t, err)
	ac.SetListenerID("not_empty")
	kiz, _ := NewWithAPIClient(ac)
	events, err := kiz.PollEvents()
	assert.NoError(t, err)
	assert.Equal(t, 0, len(events))
}
