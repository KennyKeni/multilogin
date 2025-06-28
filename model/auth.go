package model

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Status APIStatus `json:"status"`
	Data   AuthData  `json:"data"`
}

type AutomationResponse struct {
	Status APIStatus      `json:"status"`
	Data   AutomationData `json:"data"`
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

type AutomationData struct {
	Token string `json:"token"`
}
