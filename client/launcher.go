package client

import (
	"fmt"
	"github.com/KennyKeni/multilogin/model"
	"net/http"
)

// StartBrowserProfile starts a browser profile with the specified parameters
// folderID and profileID are required UUIDs
// automationType can be "selenium", "puppeteer", or "playwright" (optional)
// headlessMode enables/disables headless mode (optional, defaults to false)
func (c *Client) StartBrowserProfile(folderID, profileID string, automationType string, headlessMode bool) (*model.StartProfileResponse, error) {
	if folderID == "" || profileID == "" {
		return nil, fmt.Errorf("folder_id and profile_id cannot be empty")
	}

	// Build endpoint path
	endpoint := fmt.Sprintf("/api/v2/profile/f/%s/p/%s/start", folderID, profileID)

	// Build query parameters
	params := make(map[string]string)
	if automationType != "" {
		// Validate automation type
		switch automationType {
		case "selenium", "puppeteer", "playwright":
			params["automation_type"] = automationType
		default:
			return nil, fmt.Errorf("invalid automation_type: %s. Must be 'selenium', 'puppeteer', or 'playwright'", automationType)
		}
	}

	// Add headless mode parameter
	if headlessMode {
		params["headless_mode"] = "true"
	} else {
		params["headless_mode"] = "false"
	}

	var response model.StartProfileResponse
	err := c.makeRequestAndDecode(c.makeLauncherRequest, http.MethodGet, endpoint, nil, params, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to start browser profile: %w", err)
	}

	if response.Status.HTTPCode != 200 {
		return nil, fmt.Errorf("failed to start profile: %s", response.Status.Message)
	}

	return &response, nil
}

// StopAllProfiles stops all launched profiles
// profileType can be "all", "regular", or "quick" (optional, defaults to "all")
func (c *Client) StopAllProfiles(profileType string) (*model.StopAllProfilesResponse, error) {
	endpoint := "/api/v1/profile/stop_all"

	// Build query parameters
	params := make(map[string]string)
	if profileType != "" {
		// Validate profile type
		switch profileType {
		case "all", "regular", "quick":
			params["type"] = profileType
		default:
			return nil, fmt.Errorf("invalid profile type: %s. Must be 'all', 'regular', or 'quick'", profileType)
		}
	} else {
		// Default to "all" if not specified
		params["type"] = "all"
	}

	var response model.StopAllProfilesResponse
	err := c.makeRequestAndDecode(c.makeLauncherRequest, http.MethodGet, endpoint, nil, params, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to stop all profiles: %w", err)
	}

	if response.Status.HTTPCode != 200 {
		return nil, fmt.Errorf("failed to stop profiles: %s", response.Status.Message)
	}

	return &response, nil
}
