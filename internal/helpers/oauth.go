package helpers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"refit_backend/internal/constants"
	"time"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
)

// GetOAuthGoogleConfig get config oauth2 google
func GetOAuthGoogleConfig() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  fmt.Sprintf("%s/auth/google/callback", constants.BaseURL),
		ClientID:     viper.GetString("google.oauth.client_id"),
		ClientSecret: viper.GetString("google.oauth.secret"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/plus.me",
		},
		Endpoint: google.Endpoint,
	}
}

// GetOAuthFacebookConfig get config oauth2 facebok
func GetOAuthFacebookConfig() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  fmt.Sprintf("%s/auth/facebook/callback", constants.BaseURL),
		ClientID:     viper.GetString("facebook.oauth.app_id"),
		ClientSecret: viper.GetString("facebook.oauth.secret"),
		Scopes: []string{
			"user_birthday",
			"email",
		},
		Endpoint: facebook.Endpoint,
	}
}

// GenerateStateOauthCookie generate cookie for oauth
func GenerateStateOauthCookie() (string, http.Cookie) {
	var expiration = time.Now().Add(365 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}

	return state, cookie
}

// GetUserDataFromTwitter get data from twitter after get token
func GetUserDataFromTwitter(code string) ([]byte, error) {

	var facebookOAuthConfig = &oauth2.Config{
		RedirectURL:  "https://refit-api.luqmanul.com/auth/twitter/callback",
		ClientID:     viper.GetString("twitter.oauth.app_id"),
		ClientSecret: viper.GetString("twitter.oauth.secret"),
		Scopes: []string{
			"id",
			"name",
			"email",
		},
		Endpoint: facebook.Endpoint,
	}
	// Use code to get token and get user info from Google.

	token, err := facebookOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	response, err := http.Get(fmt.Sprintf("https://graph.facebook.com/v3.2/me?fields=id,name,email&access_token=%s", token.AccessToken))
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}
	return contents, nil
}
