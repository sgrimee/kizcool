package client

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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

func TestRegisterListener(t *testing.T) {
	t.Skip("Need to write test")
}

func TestUnregisterListener(t *testing.T) {
	t.Skip("Need to write test")
}

func TestFetchEvents(t *testing.T) {
	t.Skip("Need to write test")
}
