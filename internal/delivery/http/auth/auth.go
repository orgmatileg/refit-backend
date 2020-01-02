package auth

import (
	"fmt"
	"net/http"
	"refit_backend/internal/constants"
	"refit_backend/internal/helpers"
	"refit_backend/internal/logger"
	"refit_backend/internal/services"
	"refit_backend/models"

	"github.com/labstack/echo"
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
	oauthState, cookie := helpers.GenerateStateOauthCookie()
	c.SetCookie(&cookie)
	u := helpers.GetOAuthGoogleConfig().AuthCodeURL(oauthState)
	return c.Redirect(http.StatusTemporaryRedirect, u)
}

func (a auth) OAuthGoogleCallback(c echo.Context) error {

	var (
		state = c.FormValue("state")
		code  = c.FormValue("code")
		ctx   = c.Request().Context()
	)

	oauthState, err := c.Cookie("oauthstate")
	if err != nil {
		logger.Infof("could not get cookie: %s", err.Error())
		return c.Redirect(http.StatusTemporaryRedirect, constants.RedirectFailOAuth)
	}

	if state != oauthState.Value {
		logger.Infof("oauth google state not equal")
		return c.Redirect(http.StatusTemporaryRedirect, constants.RedirectFailOAuth)
	}

	token, err := a.service.Auth().OAuthGoogleCallback(ctx, code)
	if err != nil {
		logger.Infof("could not handle service callback: %s", err.Error())
		return c.Redirect(http.StatusTemporaryRedirect, constants.RedirectFailOAuth)
	}

	return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("exp://192.168.43.2:19000/--/home?setToken=%s", token))
}

func (a auth) OAuthFacebookLogin(c echo.Context) error {
	oauthState, cookie := helpers.GenerateStateOauthCookie()
	c.SetCookie(&cookie)
	u := helpers.GetOAuthFacebookConfig().AuthCodeURL(oauthState)
	return c.Redirect(http.StatusTemporaryRedirect, u)
}

func (a auth) OAuthFacebookCallback(c echo.Context) error {

	var (
		state = c.FormValue("state")
		code  = c.FormValue("code")
		ctx   = c.Request().Context()
	)

	oauthState, err := c.Cookie("oauthstate")
	if err != nil {
		logger.Infof("could not get cookie: %s", err.Error())
		return c.Redirect(http.StatusTemporaryRedirect, constants.RedirectFailOAuth)
	}

	if state != oauthState.Value {
		logger.Infof("oauth facebook state not equal")
		return c.Redirect(http.StatusTemporaryRedirect, constants.RedirectFailOAuth)
	}

	token, err := a.service.Auth().OAuthFacebookCallback(ctx, code)
	if err != nil {
		logger.Infof("could not handle service callback: %s", err.Error())
		return c.Redirect(http.StatusTemporaryRedirect, constants.RedirectFailOAuth)
	}

	return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("exp://192.168.43.2:19000/--/home?setToken=%s", token))
}

func (a auth) OAuthTwitterLogin(c echo.Context) error {
	return nil
}

func (a auth) OAuthTwitterCallback(c echo.Context) error {
	return nil
}
