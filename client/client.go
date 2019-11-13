// Package client provides a low-level client to the overkiz api
package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"golang.org/x/net/publicsuffix"
)

// checkStatusOk performs simple tests to ensure the request was successful
// if an error occured, try to qualify it then return it. In this case the Body of the
// response will not be usable later on.
func checkStatusOk(resp *http.Response) error {
	if resp.StatusCode == 200 {
		return nil
	}
	// Decode the body to try to get a meaningful error message
	type errResult struct {
		ErrorCode string `json:"errorCode"`
		ErrorMsg  string `json:"error"`
	}
	var result errResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("json decode: %v", err)
	}
	if resp.StatusCode == 401 {
		if strings.Contains(result.ErrorMsg, "Too many requests") {
			return NewTooManyRequestsError(result.ErrorMsg)
		}
		return NewAuthenticationError(result.ErrorMsg)
	}
	return fmt.Errorf("%v", result)
}

// Client provides methods to make http requests to the api server while making the
// authentification and session ID renewal transparent.
type Client struct {
	username  string
	password  string
	baseURL   string
	sessionID string
	hc        *http.Client
}

// New returns a new Client
// sessionID is optional and used when caching sessions externally
func New(username, password, baseURL, sessionID string) (APIClient, error) {
	hc := http.Client{}
	return NewWithHC(username, password, baseURL, sessionID, &hc)
}

// NewWithHC returns a new Client, injecting the HTTP client to use. See New.
func NewWithHC(username, password, baseURL, sessionID string, hc *http.Client) (APIClient, error) {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}
	if sessionID != "" {
		c := http.Cookie{
			Name:  "JSESSIONID",
			Value: sessionID,
		}
		url, err := url.Parse(baseURL)
		if err != nil {
			return nil, err
		}
		jar.SetCookies(url, []*http.Cookie{&c})
	}
	hc.Jar = jar
	client := Client{
		username:  username,
		password:  password,
		baseURL:   baseURL,
		sessionID: sessionID,
		hc:        hc,
	}
	return &client, nil
}

// SessionID is the latest known sessionID value
// It can be used for caching sessions externally.
func (c *Client) SessionID() string {
	return c.sessionID
}

// Login to the api server
func (c *Client) Login() error {
	formData := url.Values{"userId": {c.username}, "userPassword": {c.password}}

	resp, err := c.hc.PostForm(c.baseURL+"/enduserAPI/login", formData)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if err := checkStatusOk(resp); err != nil {
		return err
	}
	for _, cookie := range resp.Cookies() {
		if (cookie.Name == "JSESSIONID") && (cookie.Value != "") {
			c.sessionID = cookie.Value
			return nil
		}
	}
	return errors.New("JSESSIONID not found in response to /login")
}

// GetWithAuth performs a GET request for the given query which is appended
// to the baseURL. It tries to renew the session ID if needed.
func (c *Client) GetWithAuth(query string) (*http.Response, error) {
	req, err := http.NewRequest("GET", c.baseURL+query, nil)
	if err != nil {
		return nil, err
	}
	return c.DoWithAuth(req)
}

// DoWithAuth performs the given request. If an authentication error occurs,
// it tries to login to renew the sessionID, then tries the request again.
func (c *Client) DoWithAuth(req *http.Request) (*http.Response, error) {
	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, err
	}
	if err := checkStatusOk(resp); err != nil {
		switch err.(type) {
		case *AuthenticationError:
			fmt.Println("Need to re-authenticate")
			if err := c.Login(); err != nil {
				return nil, err
			}
			resp, err := c.hc.Do(req)
			if err != nil {
				return nil, err
			}
			return resp, nil
		default:
			return nil, err
		}
	}
	fmt.Println("Using cached sessionID")
	return resp, nil
}

// RefreshStates tells the server to refresh states.
// But not sure yet what it really means`?
func (c *Client) RefreshStates() error {
	req, err := http.NewRequest(http.MethodPut, c.baseURL+"/enduserAPI/setup/devices/states/refresh", nil)
	if err != nil {
		return err
	}
	_, err = c.DoWithAuth(req)
	if err != nil {
		return err
	}
	return nil
}

// Execute initiates the execution of a group of actions
// json needs to be marshalled from an ActionGroup
func (c *Client) Execute(json []byte) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, c.baseURL+"/enduserAPI/exec/apply", bytes.NewBuffer(json))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.DoWithAuth(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
