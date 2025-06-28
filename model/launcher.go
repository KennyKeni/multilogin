package model

type StartProfileRequest struct {
	AutomationType string `json:"automation_type,omitempty"` // selenium, puppeteer, playwright
	HeadlessMode   bool   `json:"headless_mode,omitempty"`   // true/false
}

// StartProfileResponse represents the response from starting a browser profile
type StartProfileResponse struct {
	Status APIStatus        `json:"status"`
	Data   StartProfileData `json:"data"`
}

// StartProfileData contains the response data from starting a profile
type StartProfileData struct {
	Port      string `json:"port,omitempty"`       // Automation port returned when automation_type is specified
	Message   string `json:"message,omitempty"`    // Response message
	ProfileID string `json:"profile_id,omitempty"` // Profile ID that was started
}

// StopAllProfilesResponse represents the response from stopping all profiles
type StopAllProfilesResponse struct {
	Status APIStatus           `json:"status"`
	Data   StopAllProfilesData `json:"data"`
}

// StopAllProfilesData contains the response data from stopping all profiles
type StopAllProfilesData struct {
	Message         string   `json:"message,omitempty"`
	StoppedCount    int      `json:"stopped_count,omitempty"`
	ProfilesStopped []string `json:"profiles_stopped,omitempty"`
}
