package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	const sID = "myTestSessionID"
	c, err := New("user", "pass", "http://dummy.org", sID)
	assert.NoError(t, err)
	assert.Equal(t, sID, c.SessionID())
}

func TestNewWithHTTPClientGetsSessionCookie(t *testing.T) {
	const sID = "myTestSessionID"
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/enduserAPI/login", req.URL.String())
		cookie := http.Cookie{
			Name:     "JSESSIONID",
			Value:    sID,
			Expires:  time.Now().AddDate(0, 0, 1),
			HttpOnly: true,
			Secure:   false,
		}
		http.SetCookie(rw, &cookie)
		rw.Write([]byte(`{"success":true,"roles":[{"name":"ENDUSER"}]}`))
	}))
	defer server.Close()
	c, err := NewWithHTTPClient("user", "pass", server.URL, "", server.Client())
	assert.NoError(t, err)
	assert.Equal(t, "", c.SessionID())
	err = c.Login()
	assert.NoError(t, err)
	assert.Equal(t, sID, c.SessionID())
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
	ac, err := NewWithHTTPClient("gooduser", "goodpass", server.URL, "", server.Client())
	assert.NoError(t, err)
	err = ac.Login()
	assert.Nil(t, err)
	assert.Equal(t, sessionID, ac.SessionID())
}

func TestBadLogin(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/enduserAPI/login", req.URL.String())
		rw.WriteHeader(401)
		rw.Write([]byte(`{"errorCode": "AUTHENTICATION_ERROR","error": "Bad credentials"}`))
	}))
	defer server.Close()
	ac, err := NewWithHTTPClient("baduser", "badpass", server.URL, "", server.Client())
	assert.NoError(t, err)
	err = ac.Login()
	assert.Equal(t, err, NewAuthenticationError("Bad credentials"))
}

func TestRefreshStates(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/enduserAPI/setup/devices/states/refresh", req.URL.String())
	}))
	defer server.Close()
	ac, _ := NewWithHTTPClient("gooduser", "goodpass", server.URL, "", server.Client())
	err := ac.RefreshStates()
	assert.Nil(t, err)
}

func TestRegisterListener(t *testing.T) {
	const lid = "77777777-3333-5555-2222-cccccccccccc"
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/events/register", req.URL.String())
		rw.Write([]byte(`{"id":"` + lid + `"}`))
	}))
	defer server.Close()
	c, err := NewWithHTTPClient("user", "pass", server.URL, "", server.Client())
	assert.NoError(t, err)
	assert.Equal(t, "", c.ListenerID())
	assert.NoError(t, c.registerListener())
	assert.Equal(t, lid, c.ListenerID())
}

func TestUnregisterListener(t *testing.T) {
	const lid = "77777777-3333-5555-2222-cccccccccccc"
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/events/"+lid+"/unregister", req.URL.String())
	}))
	defer server.Close()
	c, err := NewWithHTTPClient("user", "pass", server.URL, "", server.Client())
	assert.NoError(t, err)
	err = c.unregisterListener()
	assert.NoError(t, err)
	assert.Equal(t, "", c.ListenerID())
}

func TestPollEventsWithIDGood(t *testing.T) {
	const lid = "77777777-3333-5555-2222-cccccccccccc"
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/events/"+lid+"/fetch", req.URL.String())
	}))
	defer server.Close()
	c, err := NewWithHTTPClient("user", "pass", server.URL, "", server.Client())
	assert.NoError(t, err)
	_, err = c.pollEventsWithID(lid)
	assert.NoError(t, err)
}

func TestPollEventsWithIDEmpty(t *testing.T) {
	c, err := New("user", "pass", "http://bla", "")
	assert.NoError(t, err)
	_, err = c.pollEventsWithID("")
	assert.Error(t, err)
	_, ok := err.(*NoRegisteredEventListenerError)
	assert.True(t, ok)
}

func TestPollEventsWithIDBad(t *testing.T) {
	const lid = "77777777-3333-5555-2222-cccccccccccc"
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/events/"+lid+"/fetch", req.URL.String())
		rw.WriteHeader(400)
		rw.Write([]byte(`{
			"errorCode": "UNSPECIFIED_ERROR",
			"error": "No registered event listener"
		  }`))
	}))
	defer server.Close()
	c, err := NewWithHTTPClient("user", "pass", server.URL, "", server.Client())
	assert.NoError(t, err)
	_, err = c.pollEventsWithID(lid)
	assert.Error(t, err)
	_, ok := err.(*NoRegisteredEventListenerError)
	assert.True(t, ok)
}

func TestPollEventsDoesRegisterListener(t *testing.T) {
	const validLID = "77777777-3333-5555-2222-cccccccccccc"
	const expiredLID = "expired_lid"
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		switch query := req.URL.String(); query {
		case "/events/register":
			rw.Write([]byte(`{"id":"` + validLID + `"}`))
		case "/events/" + validLID + "/fetch":
			break
		case "/events/" + expiredLID + "/fetch":
			rw.WriteHeader(400)
			rw.Write([]byte(`{
				"errorCode": "UNSPECIFIED_ERROR",
				"error": "No registered event listener"
			}`))
		default:
			panic("Unexpected query")
		}
	}))
	defer server.Close()
	c, err := NewWithHTTPClient("gooduser", "goodpass", server.URL, "", server.Client())
	assert.NoError(t, err)

	c.SetListenerID("")
	_, err = c.PollEvents()
	assert.NoError(t, err)
	assert.Equal(t, validLID, c.ListenerID())

	c.SetListenerID(expiredLID)
	_, err = c.PollEvents()
	assert.NoError(t, err)
	assert.Equal(t, validLID, c.ListenerID())

	c.SetListenerID(validLID)
	_, err = c.PollEvents()
	assert.NoError(t, err)
	assert.Equal(t, validLID, c.ListenerID())
}

func TestCheckStatusOk(t *testing.T) {
	var tests = []struct {
		name     string
		code     int
		bodyText string
		e        error
	}{
		{"200", 200, "", nil},
		{"401-auth", 401, `{"errorCode":"AUTHENTICATION_ERROR","error":"Bad credentials"}`,
			NewAuthenticationError("Bad credentials")},
		{"401-toomany", 401, `{"errorCode":"AUTHENTICATION_ERROR","error":"Too many requests, try again later : login with user@domain.com"}`,
			NewTooManyRequestsError("Too many requests, try again later : login with user@domain.com")},
		{"500", 500, `{"errorCode":"WEIRD_ERROR","error":"Unexpected"}`, errors.New("{WEIRD_ERROR Unexpected}")},

		{"bad-json", 999, "Not json", errors.New("json decode: invalid character 'N' looking for beginning of value")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			w.Write([]byte(tt.bodyText))
			w.Code = tt.code
			assert.Equal(t, tt.e, checkStatusOk(w.Result()))
		})
	}
}
