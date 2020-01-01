package models

// OAuth struct
type OAuth struct {
	ID          uint   `json:"id"`
	Service     string `json:"service"`
	OpenID      string `json:"open_id"`
	UserID      uint   `json:"user_id"`
	Token       string `json:"token"`
	TokenSecret string `json:"token_secret"`
}

// GoogleOAuthUserInfo struct							|
// endpoint: https://www.googleapis.com/oauth2/v2/userinfo
type GoogleOAuthUserInfo struct {
	OpenID        string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"local"`
}
