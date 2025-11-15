package api

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Login2FARequest struct {
	TempToken string `json:"temp_token"`
	Code      string `json:"code"`
}

type UserResponse struct {
	ID            int64  `json:"id"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	TwoFAEnabled  bool   `json:"twofa_enabled"`
}
