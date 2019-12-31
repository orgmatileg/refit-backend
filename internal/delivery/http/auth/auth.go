package auth

import (
	"fmt"
	"log"
	"net/http"
	"refit_backend/internal/helpers"
	"refit_backend/internal/services"
	"refit_backend/models"

	"github.com/labstack/echo"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
)

// IAuth interface
type IAuth interface {
	AuthLoginWithEmail(c echo.Context) error
	AuthRegister(c echo.Context) error
	OAuthGoogleLogin(c echo.Context) error
	OAuthGoogleCallback(c echo.Context) error
	OAuthFacebookLogin(c echo.Context) error
	OAuthFacebookCallback(c echo.Context) error
	OAuthTwitterLogin(c echo.Context) error
	OAuthTwitterCallback(c echo.Context) error
}

type auth struct {
	service services.IServices
}

// New auth http handler
func New() IAuth {
	return &auth{
		service: services.New(),
	}
}

func (a auth) AuthLoginWithEmail(c echo.Context) error {
	var ru models.User
	err := c.Bind(&ru)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()

	token, err := a.service.Users().AuthLoginWithEmail(ctx, &ru)
	if err != nil {
		return helpers.MakeDefaultResponse(c, http.StatusBadRequest, err)
	}

	return helpers.MakeDefaultResponse(c, http.StatusOK, struct {
		Token string `json:"token"`
	}{
		Token: token,
	})
}

func (a auth) AuthRegister(c echo.Context) error {
	var ru models.User
	err := c.Bind(&ru)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	_, err = a.service.Users().Create(ctx, &ru)
	if err != nil {
		return helpers.MakeDefaultResponse(c, http.StatusBadRequest, err)
	}

	return helpers.MakeDefaultResponse(c, http.StatusCreated, nil)
}

func (a auth) OAuthGoogleLogin(c echo.Context) error {
	var googleOauthConfig = &oauth2.Config{
		RedirectURL:  "https://refit-api.luqmanul.com/auth/google/callback",
		ClientID:     viper.GetString("google.oauth.client_id"),
		ClientSecret: viper.GetString("google.oauth.secret"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/plus.me",
		},
		Endpoint: google.Endpoint,
	}

	// Create oauthState cookie
	oauthState, cookie := helpers.GenerateStateOauthCookie()
	c.SetCookie(&cookie)
	u := googleOauthConfig.AuthCodeURL(oauthState)

	return c.Redirect(http.StatusTemporaryRedirect, u)
}

func (a auth) OAuthGoogleCallback(c echo.Context) error {

	// Read oauthState from Cookie
	oauthState, _ := c.Cookie("oauthstate")

	if c.FormValue("state") != oauthState.Value {
		log.Println("invalid oauth google state")
		c.Redirect(http.StatusTemporaryRedirect, "luqmanul.com")
		return nil
	}

	data, err := helpers.GetUserDataFromGoogle(c.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		c.Redirect(http.StatusTemporaryRedirect, "luqmanul.com")
		return nil
	}
	fmt.Println(data)

	return c.Redirect(http.StatusTemporaryRedirect, "exp://192.168.43.2:19000/home")

	// GetOrCreate User in your db.
	// Redirect or response with a token.
	// More code .....
	// fmt.Fprintf(w, "UserInfo: %s\n", data)

	// return c.JSON(http.StatusOK, string(data))
}

func (a auth) OAuthFacebookLogin(c echo.Context) error {
	var facebookOauthConfig = &oauth2.Config{
		RedirectURL:  "https://refit-api.luqmanul.com/auth/facebook/callback",
		ClientID:     viper.GetString("facebook.oauth.app_id"),
		ClientSecret: viper.GetString("facebook.oauth.secret"),
		Scopes: []string{
			"user_birthday",
			"email",
		},
		Endpoint: facebook.Endpoint,
	}

	// Create oauthState cookie
	oauthState, cookie := helpers.GenerateStateOauthCookie()
	c.SetCookie(&cookie)
	u := facebookOauthConfig.AuthCodeURL(oauthState)

	return c.Redirect(http.StatusTemporaryRedirect, u)
}

func (a auth) OAuthFacebookCallback(c echo.Context) error {

	// Read oauthState from Cookie
	oauthState, _ := c.Cookie("oauthstate")

	if c.FormValue("state") != oauthState.Value {
		log.Println("invalid oauth facebook state")
		c.Redirect(http.StatusTemporaryRedirect, "luqmanul.com")
		return nil
	}

	data, err := helpers.GetUserDataFromFacebook(c.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		c.Redirect(http.StatusTemporaryRedirect, "luqmanul.com")
		return nil
	}

	// GetOrCreate User in your db.
	// Redirect or response with a token.
	// More code .....
	// fmt.Fprintf(w, "UserInfo: %s\n", data)

	return c.JSON(http.StatusOK, string(data))
}

func (a auth) OAuthTwitterLogin(c echo.Context) error {
	return nil
}

func (a auth) OAuthTwitterCallback(c echo.Context) error {
	return nil
}
