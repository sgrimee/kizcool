package client

import (
	"errors"
	"net/http/httptest"
	"testing"

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
