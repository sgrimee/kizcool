package client

import "net/http"

// APIClient interface defines a low-level api client
type APIClient interface {
	SessionID() string
	Login() error
	GetWithAuth(query string) (*http.Response, error)
	DoWithAuth(req *http.Request) (*http.Response, error)
	RefreshStates() error
	Execute(json []byte) (*http.Response, error)
}
