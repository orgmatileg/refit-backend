package auth

import (
	"fmt"
	"github.com/kkdai/twitter"
	"github.com/labstack/echo"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"refit_backend/internal/constants"
	"refit_backend/internal/helpers"
	"refit_backend/internal/logger"
	"refit_backend/internal/services"
	"refit_backend/models"
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
	var (
		ConsumerKey    = viper.GetString("twitter.consumer_api_key")
		ConsumerSecret = viper.GetString("twitter.consumer_api_secret")
		CallbackURL    = "https://refit-api.luqmanul.com/auth/twitter/callback"
		twitterClient  *twitter.ServerClient
	)
	twitterClient = twitter.NewServerClient(ConsumerKey, ConsumerSecret)
	u := twitterClient.GetAuthURL(CallbackURL)
	return c.Redirect(http.StatusTemporaryRedirect, u)
}

func (a auth) OAuthTwitterCallback(c echo.Context) error {
	var (
		// ConsumerKey      = viper.GetString("twitter.consumer_api_key")
		// ConsumerSecret   = viper.GetString("twitter.consumer_api_secret")
		verificationCode = c.QueryParam("oauth_verifier")
		tokenKey         = c.QueryParam("oauth_token")
	)

	// authorization: OAuth oauth_consumer_key="CONSUMER_API_KEY", oauth_nonce="OAUTH_NONCE", oauth_signature="OAUTH_SIGNATURE", oauth_signature_method="HMAC-SHA1", oauth_timestamp="OAUTH_TIMESTAMP", oauth_token="ACCESS_TOKEN", oauth_version="1.0"
	req, err := http.NewRequest("GET", "https://api.twitter.com/oauth/access_token", nil)
	if err != nil {
		logger.Infof("%s", err.Error())
	}
	q := req.URL.Query()
	q.Add("oauth_token", tokenKey)
	q.Add("oauth_verifier", verificationCode)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Infof("%s", err.Error())
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Infof("%s", err.Error())
	}
	defer resp.Body.Close()
	fmt.Sprintln(string(b))
	return c.String(200, string(b))
	// http.Get("https://api.twitter.com/1.1/account/verify_credentials.json?include_email=true")

	// timelineURL := fmt.Sprintf("http://%s/time", r.Host)
	// return c.Redirect(http.StatusTemporaryRedirect, timelineURL)

	// token, err := a.service.Auth().OAuthFacebookCallback(ctx, code)
	// if err != nil {
	// 	logger.Infof("could not handle service callback: %s", err.Error())
	// 	return c.Redirect(http.StatusTemporaryRedirect, constants.RedirectFailOAuth)
	// }

	// return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("exp://192.168.43.2:19000/--/home?setToken=%s", token))

}
