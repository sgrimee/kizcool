package config

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vapourismo/knx-go/knx/cemi"
)

var yamlExample = []byte(`
debug: true
password: mysecretpass
username: me@nowhere.com
bla: boum
devices:
  - label: Spots lit parents
    url: io://1111-0000-4444/11111111
    commands:
      - name: setIntensity
        groupaddr: 10/0/3
      - name: setOnOff
        groupaddr: 10/0/4
    states:
      - name: core:OnOffState
        groupaddr: 10/0/5
      - name: core:LightIntensityState
        groupaddr: 10/0/6
`)

func TestMain(m *testing.M) {
	v.SetConfigType("yaml")
	v.ReadConfig(bytes.NewBuffer(yamlExample))

	os.Exit(m.Run())
}

func TestUsername(t *testing.T) {
	assert.Equal(t, "me@nowhere.com", Username())
}

func TestSetUsername(t *testing.T) {
	SetUsername("notme@somewhere.com")
	assert.Equal(t, "notme@somewhere.com", Username())
}

func TestPassword(t *testing.T) {
	assert.Equal(t, "mysecretpass", Password())
}

func TestSetPassword(t *testing.T) {
	SetPassword("newpassword")
	assert.Equal(t, "newpassword", Password())
}

func TestBaseURL(t *testing.T) {
	assert.Equal(t, TahomaBaseURL, BaseURL())
}

func TestSetBaseURL(t *testing.T) {
	SetBaseURL("https://happy.com")
	assert.Equal(t, "https://happy.com", BaseURL())
}

func TestSessionID(t *testing.T) {
	assert.Equal(t, "", SessionID())
}

func TestSetSessionID(t *testing.T) {
	SetSessionID("1234567890")
	assert.Equal(t, "1234567890", SessionID())
}

func TestDevices(t *testing.T) {
	devices, err := Devices()
	require.NoError(t, err)
	require.Equal(t, 1, len(devices))
	assert.Equal(t, "setOnOff", devices[0].Commands[1].Name)
	gaddr, _ := cemi.NewGroupAddrString("10/0/4")
	assert.Equal(t, gaddr, devices[0].Commands[1].GroupAddr)
}

// func TestRead(t *testing.T) {
// }

// func TestWrite(t *testing.T) {
// }

func TestUsingConfigFile(t *testing.T) {
	usingConfigFile = true
	assert.True(t, UsingConfigFile())
	usingConfigFile = false
	assert.False(t, UsingConfigFile())
}

func TestFile(t *testing.T) {
	assert.Equal(t, "", File())
}

func TestDebug(t *testing.T) {
	assert.True(t, Debug())
}
