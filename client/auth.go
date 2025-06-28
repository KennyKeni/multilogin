package client

import (
	"fmt"
	"net/http"
	"time"

	"github.com/KennyKeni/multilogin/model"
	"github.com/KennyKeni/multilogin/util"
)

func (c *Client) signIn() error {
	authData := map[string]string{
		"email":    c.email,
		"password": c.passwordHash,
	}

	var authResp model.AuthResponse
	err := c.makeRequestAndDecode(c.makeUnauthenticatedApiRequest, http.MethodPost, "/user/signin", authData, nil, &authResp)
	if err != nil {
		return err
	}

	if authResp.Status.HTTPCode != 200 {
		return fmt.Errorf("auth failed: %s", authResp.Status.Message)
	}

	tokenExpiration, err := util.GetTokenExpiration(authResp.Data.Token)
	if err != nil {
		return err
	}

	c.accessToken = authResp.Data.Token
	c.refreshToken = authResp.Data.RefreshToken
	c.accessTokenExp = tokenExpiration

	return nil
}

func (c *Client) refreshAccessToken() error {
	if c.email == "" {
		return fmt.Errorf("email cannot be empty")
	}
	payload := map[string]string{
		"email":         c.email,
		"refresh_token": c.refreshToken,
	}

	var authResp model.AuthResponse
	err := c.makeRequestAndDecode(c.makeApiRequest, http.MethodPost, "/user/refresh_token", payload, nil, &authResp)
	if err != nil {
		return err
	}

	if authResp.Status.HTTPCode != 200 {
		return fmt.Errorf("auth failed: %s", authResp.Status.Message)
	}

	tokenExpiration, err := util.GetTokenExpiration(authResp.Data.Token)
	if err != nil {
		return err
	}

	c.accessToken = authResp.Data.Token
	c.refreshToken = authResp.Data.RefreshToken
	c.accessTokenExp = tokenExpiration

	return nil
}

func (c *Client) getAutomationToken() error {
	parameters := map[string]string{
		"expiration_period": "1h",
	}

	var automationResp model.AutomationResponse
	err := c.makeRequestAndDecode(c.makeApiRequest, http.MethodPost, "/user/refresh_token", nil, parameters, &automationResp)
	if err != nil {
		return err
	}

	tokenExpiration, err := util.GetTokenExpiration(automationResp.Data.Token)
	if err != nil {
		return err
	}
	c.automationToken = automationResp.Data.Token
	c.automationTokenExp = tokenExpiration

	return nil
}

func (c *Client) ensureAuth() error {
	if c.automationToken != "" {
		return c.ensureAutomationToken()
	}

	return c.ensureAccessToken()
}

func (c *Client) ensureAutomationToken() error {
	if c.isAutomationTokenExpired() {
		err := c.ensureAccessToken()
		if err != nil {
			return err
		}
		err = c.getAutomationToken()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) ensureAccessToken() error {
	// Case 1: No tokens at all - authenticate from scratch
	if c.accessToken == "" || c.refreshToken == "" {
		return c.signIn()
	}

	// Case 2: Have tokens but access token is expired - refresh it
	if c.isAccessTokenExpired() {
		return c.refreshAccessToken()
	}

	// Case 3: Access token is valid - nothing to do
	return nil
}

func (c *Client) getBestToken() string {
	// Priority: automation token > bearer token
	if c.automationToken != "" && !c.isAutomationTokenExpired() {
		return c.automationToken
	}

	if c.accessToken != "" && !c.isAccessTokenExpired() {
		return c.accessToken
	}

	return ""
}

func (c *Client) isAutomationTokenExpired() bool {
	return time.Now().After(c.automationTokenExp.Add(-10 * time.Minute))
}

func (c *Client) isAccessTokenExpired() bool {
	return time.Now().After(c.accessTokenExp.Add(-5 * time.Minute))
}
