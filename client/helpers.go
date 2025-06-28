package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type RequestMaker func(method, endpoint string, body interface{}, params map[string]string) (*http.Response, error)

func (c *Client) makeRequestAndDecode(
	requestMaker RequestMaker,
	method, endpoint string,
	body interface{},
	params map[string]string,
	result interface{},
) error {
	resp, err := requestMaker(method, endpoint, body, params)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(result)
}

func (c *Client) makeUnauthenticatedApiRequest(method, endpoint string, body interface{}, params map[string]string) (*http.Response, error) {
	return c.makeUnauthenticatedRequest(method, c.apiURL, endpoint, body, params)
}

func (c *Client) makeUnauthenticatedLauncherRequest(method, endpoint string, body interface{}, params map[string]string) (*http.Response, error) {
	launcherURL := fmt.Sprintf("%s:%s", c.launcherURL, c.launcherPort)
	return c.makeUnauthenticatedRequest(method, launcherURL, endpoint, body, params)
}
