package client

import "net/http"

// APIClient interface defines a low-level api client to the kiz server
type APIClient interface {
	SessionID() string
	Login() error
	GetWithAuth(query string) (*http.Response, error)
	DoWithAuth(req *http.Request) (*http.Response, error)
	GetDevices() (*http.Response, error)
	GetDevice(deviceURL string) (*http.Response, error)
	GetDeviceState(deviceURL, stateName string) (*http.Response, error)
	RefreshStates() error
	Execute(json []byte) (*http.Response, error)
	RegisterListener() (string, error)
	UnregisterListener(string) error
}
