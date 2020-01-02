package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"refit_backend/internal/logger"
	"refit_backend/models"

	"golang.org/x/oauth2"
)

// IAuth repository interface
type IAuth interface {
	GetUserDataFromGoogle(ctx context.Context, oauthConfig *oauth2.Config, code string) (m *models.GoogleOAuthUserInfo, err error)
	GetUserDataFromFacebook(ctx context.Context, oauthConfig *oauth2.Config, code string) (m *models.FacebookOAuthUserInfo, err error)
	GetUserDataFromTwitter(c *http.Client) (m *models.TwitterOAuthUserInfo, err error)
}

type auth struct{}

// New Repository Users
func New() IAuth {
	return &auth{}
}

// GetUserDataFromGoogle repository users
func (u auth) GetUserDataFromGoogle(ctx context.Context, oauthConfig *oauth2.Config, code string) (m *models.GoogleOAuthUserInfo, err error) {

	token, err := oauthConfig.Exchange(ctx, code)
	if err != nil {
		logger.Infof("could not exchange authorization code to token google: %s", err.Error())
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	response, err := http.Get(fmt.Sprintf("https://www.googleapis.com/oauth2/v2/userinfo?access_token=%s", token.AccessToken))
	if err != nil {
		logger.Infof("could not getting user info from google: %s", err.Error())
		return nil, err
	}
	defer response.Body.Close()
	resBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Infof("could not read body from response google: %s", err.Error())
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}
	err = json.Unmarshal(resBody, &m)
	if err != nil {
		logger.Infof("could not json unmarshall response from google: %s", err.Error())
		return nil, err
	}
	return m, nil
}

// GetUserDataFromGoogle repository users
func (u auth) GetUserDataFromFacebook(ctx context.Context, oauthConfig *oauth2.Config, code string) (m *models.FacebookOAuthUserInfo, err error) {

	token, err := oauthConfig.Exchange(ctx, code)
	if err != nil {
		logger.Infof("could not exchange authorization code to token facebook: %s", err.Error())
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	response, err := http.Get(fmt.Sprintf("https://graph.facebook.com/v3.2/me?fields=id,name,picture,email,birthday&access_token=%s", token.AccessToken))
	if err != nil {
		logger.Infof("could not getting user info from facebook: %s", err.Error())
		return nil, err
	}
	defer response.Body.Close()
	resBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Infof("could not read body from response facebook: %s", err.Error())
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}
	err = json.Unmarshal(resBody, &m)
	if err != nil {
		logger.Infof("could not json unmarshall response from facebook: %s", err.Error())
		return nil, err
	}
	return m, nil
}

// GetUserDataFromGoogle repository users
func (u auth) GetUserDataFromTwitter(c *http.Client) (m *models.TwitterOAuthUserInfo, err error) {
	resp, err := c.Get("https://api.twitter.com/1.1/account/verify_credentials.json?include_email=true")
	if err != nil {
		logger.Infof("%s", err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Infof("%s", err.Error())
	}
	err = json.Unmarshal(body, &m)
	if err != nil {
		logger.Infof("%s", err.Error())
	}
	return m, nil
}
