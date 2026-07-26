package handler

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpRequest struct {
	Username             string `json:"username"`
	Email                string `json:"email"`
	EmailConfirmation    string `json:"email_confirmation"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}
