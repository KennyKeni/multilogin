package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/KennyKeni/multilogin/constants"
	"github.com/KennyKeni/multilogin/util"
)

// Client allows communication with Multilogin API
// Handles all tokens refresh
// TODO: move to automation tokens in the future
type Client struct {
	email              string
	passwordHash       string
	accessToken        string
	accessTokenExp     time.Time
	refreshToken       string
	automationToken    string
	automationTokenExp time.Time
	apiURL             string
	launcherURL        string
	launcherPort       string
	httpClient         *http.Client
}

func New(email string, password string) (*Client, error) {
	if email == "" || password == "" {
		return nil, fmt.Errorf("email and password cannot be empty")
	}
	password = util.GetMD5Hash(password)
	return &Client{
		email:        email,
		passwordHash: password,
		apiURL:       constants.ApiUrl,
		launcherURL:  constants.LauncherURL,
		launcherPort: constants.LauncherPort,
		httpClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}, nil
}

func NewAuthenticated(email string, password string) (*Client, error) {
	client, err := New(email, password)
	if err != nil {
		return nil, err
	}
	err = client.signIn()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func NewAutomation(automationToken string) (*Client, error) {
	if automationToken == "" {
		return nil, fmt.Errorf("automation token cannot be empty")
	}
	automationTokenExp, err := util.GetTokenExpiration(automationToken)
	if err != nil {
		return nil, fmt.Errorf("could not parse token expiration")
	}
	return &Client{
		automationToken:    automationToken,
		automationTokenExp: automationTokenExp,
		apiURL:             constants.ApiUrl,
		launcherURL:        constants.LauncherURL,
		launcherPort:       constants.LauncherPort,
		httpClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}, nil
}

func (c *Client) makeApiRequest(
	method string,
	endpoint string,
	body interface{},
	params map[string]string,
) (*http.Response, error) {
	if err := c.ensureAuth(); err != nil {
		return nil, err
	}
	resp, err := c.makeRequest(method, c.apiURL, endpoint, body, params)
	if err != nil {
		return nil, fmt.Errorf("error creating api request")
	}

	return resp, err
}

func (c *Client) makeLauncherRequest(
	method string,
	endpoint string,
	body interface{},
	params map[string]string,
) (*http.Response, error) {
	if err := c.ensureAuth(); err != nil {
		return nil, err
	}
	launcherURL := fmt.Sprintf("%s:%s", c.launcherURL, c.launcherPort)
	resp, err := c.makeRequest(method, launcherURL, endpoint, body, params)
	if err != nil {
		return nil, fmt.Errorf("error creating api request")
	}

	return resp, err
}

func (c *Client) makeRequest(
	method string,
	baseURL string,
	endpoint string,
	body interface{},
	params map[string]string,
) (*http.Response, error) {
	req, err := c.buildRequest(method, baseURL, endpoint, body, params)
	if err != nil {
		return nil, err
	}

	token := c.getBestToken()
	if token == "" {
		return nil, fmt.Errorf("error finding valid token")
	}
	req.Header.Set("Authorization", "Bearer "+token)

	return c.httpClient.Do(req)
}

func (c *Client) makeUnauthenticatedRequest(
	method string,
	baseURL string,
	endpoint string,
	body interface{},
	params map[string]string,
) (*http.Response, error) {
	req, err := c.buildRequest(method, baseURL, endpoint, body, params)
	if err != nil {
		return nil, err
	}

	return c.httpClient.Do(req)
}

func (c *Client) buildRequest(
	method string,
	baseURL string,
	endpoint string,
	body interface{},
	params map[string]string,
) (*http.Request, error) {
	u, err := url.Parse(baseURL + endpoint)
	if err != nil {
		return nil, err
	}

	if len(params) > 0 {
		q := u.Query()
		for key, value := range params {
			q.Set(key, value)
		}
		u.RawQuery = q.Encode()
	}

	// Handle request body
	var reqBody io.Reader
	if body != nil {
		reqBodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(reqBodyBytes)
	}

	req, err := http.NewRequest(method, u.String(), reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")

	// Only set Content-Type if request body exists
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}
