package constants

import (
	"fmt"
)

var (
	// BaseURL App
	BaseURL = "https://refit-api.luqmanul.com"
	// RedirectFailOAuth URL for redirect when fail oauth
	RedirectFailOAuth = fmt.Sprintf("%s/%s", BaseURL, "redirect-fail-oauth")
)
