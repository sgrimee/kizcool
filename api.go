package kizcool

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"golang.org/x/net/publicsuffix"
)

// API contains an http client and methods to manage authentication
type API struct {
	config Config
	client *http.Client
}

// NewAPI returns a new API
func NewAPI(config Config) (*API, error) {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}
	if config.SessionID != "" {
		c := http.Cookie{
			Name:  "JSESSIONID",
			Value: config.SessionID,
		}
		url, err := url.Parse(config.BaseURL)
		if err != nil {
			return nil, err
		}
		jar.SetCookies(url, []*http.Cookie{&c})
	}
	api := API{
		config: config,
		client: &http.Client{
			Jar: jar,
		},
	}
	return &api, nil
}

// Login and return a session ID
// The new sessionID is stored/updated in the configfile
func (api *API) Login() (string, error) {
	formData := url.Values{"userId": {api.config.Username}, "userPassword": {api.config.Password}}
	resp, err := api.client.PostForm(api.config.BaseURL+"/enduserAPI/login", formData)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if e := checkStatusOk(resp); e != nil {
		return "", e
	}
	for _, c := range resp.Cookies() {
		if (c.Name == "JSESSIONID") && (c.Value != "") {
			SaveSessionID(c.Value)
			return c.Value, nil
		}
	}
	return "", errors.New("JSESSIONID not found in response to /login")
}

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

// GetWithAuth performs a GET request for the given query which is appended
// to the BaseURL. It tries to renew the session ID if needed.
func (api *API) GetWithAuth(query string) (*http.Response, error) {
	req, err := http.NewRequest("GET", api.config.BaseURL+query, nil)
	if err != nil {
		return nil, err
	}
	return api.DoWithAuth(req)
}

// DoWithAuth performs the given request. If an authentication error occurs,
// it tries to login to renew the sessionID, then tries the request again.
func (api *API) DoWithAuth(req *http.Request) (*http.Response, error) {
	resp, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := checkStatusOk(resp); err != nil {
		switch err.(type) {
		case *AuthenticationError:
			fmt.Println("Auth failure, attempting login then re-try.")
			api.Login()
			resp, err := api.client.Do(req)
			if err != nil {
				return nil, err
			}
			return resp, nil
		default:
			return nil, err
		}
	}
	return resp, nil
}

// getSetup returns informations about the site setup, including
// devices, device location in rooms, etc
// NOT FULLY IMPLEMEMNTED and may never be. Included for documentation
// of the API endpoint only.
//
// func (api *API) getSetup() ([]interface{}, error) {
// 	resp, err := api.client.Get(api.config.BaseURL + "/externalAPI/json/getSetup")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()
// 	if err := checkStatusOk(resp); err != nil {
// 		return nil, err
// 	}
// 	var result []interface{}
// 	json.NewDecoder(resp.Body).Decode(&result)
// 	return result, nil
// }
