package constants

import (
	"fmt"
	"github.com/spf13/viper"
)

var (
	// BaseURL App
	BaseURL = viper.GetString("app.baseurl")
	// RedirectFailOAuth URL for redirect when fail oauth
	RedirectFailOAuth = fmt.Sprintf("%s/%s", BaseURL, "redirect-fail-oauth")
)
