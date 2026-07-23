package response

type LoginResponse struct {
	UUID         string `json:"uuid"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type SignUpResponse struct {
	UUID         string `json:"uuid"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
