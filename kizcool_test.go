package kizcool

import (
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
		assert.Equal(t, req.URL.String(), "/login")
		rw.WriteHeader(401)
		rw.Write([]byte(`{"errorCode": "AUTHENTICATION_ERROR","error": "Bad credentials"}`))
	}))
	defer server.Close()
	kiz := NewKiz("baduser", "badpass")
	kiz.BaseURL = server.URL
	kiz.Client = server.Client()
	err := kiz.Login()
	assert.EqualError(t, err, "401: Authentication error")
}

func TestGoodLogin(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.URL.String(), "/login")
		cookie := http.Cookie{
			Name:    "JSESSIONID",
			Value:   "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF",
			Expires: time.Now().AddDate(0, 0, 1),
		}
		http.SetCookie(rw, &cookie)
		rw.Write([]byte(`{"success":true,"roles":[{"name":"ENDUSER"}]}`))
	}))
	kiz := NewKiz("gooduser", "goodpass")
	kiz.BaseURL = server.URL
	kiz.Client = server.Client()
	err := kiz.Login()
	assert.Nil(t, err)
}

func helperLoadBytes(t *testing.T, name string) []byte {
	path := filepath.Join("testdata", name) // relative path
	bytes, err := ioutil.ReadFile(path)
	assert.Nil(t, err)
	return bytes
}

func TestGetDevices(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.URL.String(), "/setup/devices")
		rw.Write(helperLoadBytes(t, "getDevices.json"))
	}))
	kiz := NewKiz("gooduser", "goodpass")
	kiz.BaseURL = server.URL
	kiz.Client = server.Client()
	devices, err := kiz.GetDevices()
	assert.Nil(t, err)
	assert.Equal(t, len(devices), 5)
}

func TestGetActionGroups(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.URL.String(), "/actionGroups")
		rw.Write(helperLoadBytes(t, "getActionGroups.json"))
	}))
	kiz := NewKiz("gooduser", "goodpass")
	kiz.BaseURL = server.URL
	kiz.Client = server.Client()
	actionGroups, err := kiz.GetActionGroups()
	assert.Nil(t, err)
	assert.Equal(t, len(actionGroups), 1)
	assert.Equal(t, len(actionGroups[0].Actions), 2)
}
