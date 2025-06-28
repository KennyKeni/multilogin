package client

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type AuthResponse struct {
	Status APIStatus `json:"status"`
	Data   AuthData  `json:"data"`
}

type APIStatus struct {
	ErrorCode string `json:"error_code"`
	HTTPCode  int    `json:"http_code"`
	Message   string `json:"message"`
}

type AuthData struct {
	RefreshToken string `json:"refresh_token"`
	Token        string `json:"token"`
}

func (c *Client) authenticate() error {
	authData := map[string]string{
		"email":    c.email,
		"password": c.passwordHash,
	}

	resp, err := c.makeApiRequest(http.MethodPost, "/user/signin", authData, nil, false)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return err
	}

	if authResp.Status.HTTPCode != 200 {
		return fmt.Errorf("auth failed: %s", authResp.Status.Message)
	}

	tokenExpiration, err := getTokenExpiration(authResp.Data.Token)
	if err != nil {
		return err
	}

	c.accessToken = authResp.Data.Token
	c.refreshToken = authResp.Data.RefreshToken
	c.accessTokenExp = tokenExpiration

	return nil
}

func (c *Client) refreshAccessToken() error {
	payload := map[string]string{
		"email":         c.email,
		"refresh_token": c.refreshToken,
	}

	resp, err := c.makeApiRequest(http.MethodPost, "/user/refresh_token", payload, nil, true)
	defer resp.Body.Close()

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return err
	}

	if authResp.Status.HTTPCode != 200 {
		return fmt.Errorf("auth failed: %s", authResp.Status.Message)
	}

	tokenExpiration, err := getTokenExpiration(authResp.Data.Token)
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

    resp, err := c.makeApiRequest(http.MethodGet, "/automation_token", nil, parameters, true)

    return nil
}

func getTokenExpiration(tokenString string) (time.Time, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return time.Time{}, fmt.Errorf("invalid token format")
	}

	// Add padding if needed
	payload := parts[1]
	if len(payload)%4 != 0 {
		payload += strings.Repeat("=", 4-len(payload)%4)
	}

	decoded, err := base64.URLEncoding.DecodeString(payload)
	if err != nil {
		return time.Time{}, err
	}

	var claims map[string]interface{}
	if err := json.Unmarshal(decoded, &claims); err != nil {
		return time.Time{}, err
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return time.Time{}, fmt.Errorf("exp claim not found or invalid type")
	}

	return time.Unix(int64(exp), 0), nil
}
