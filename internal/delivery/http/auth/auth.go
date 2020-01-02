package auth

import (
	"context"
	"fmt"
	"github.com/dghubble/oauth1"
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
	"strings"
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
		ConsumerKey      = viper.GetString("twitter.consumer_api_key")
		ConsumerSecret   = viper.GetString("twitter.consumer_api_secret")
		verificationCode = c.QueryParam("oauth_verifier")
		tokenKey         = c.QueryParam("oauth_token")
	)

	req, err := http.NewRequest("POST", "https://api.twitter.com/oauth/access_token", nil)
	if err != nil {
		logger.Infof("%s", err.Error())
	}
	req.Header.Add("Content-Type", "application/json")
	q := req.URL.Query()
	q.Add("oauth_token", tokenKey)
	q.Add("oauth_verifier", verificationCode)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Infof("%s", err.Error())
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Infof("%s", err.Error())
	}

	resSlice := strings.Split(string(b), "&")
	m := make(map[string]string)
	for _, v := range resSlice {
		datas := strings.Split(v, "=")
		m[datas[0]] = datas[1]
	}

	oauthWithToken := oauth1.NewToken(m["oauth_token"], m["oauth_token_secret"])
	config := oauth1.NewConfig(ConsumerKey, ConsumerSecret)
	httpClient := config.Client(context.Background(), oauthWithToken)

	path := "https://api.twitter.com/1.1/statuses/home_timeline.json?count=2"
	resp, err = httpClient.Get(path)
	if err != nil {
		logger.Infof("%s", err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Infof("%s", err.Error())
	}
	fmt.Printf("Raw Response Body:\n%v\n", string(body))
	return c.String(200, string(b)+string(body))

	////////////
	// req, err = http.NewRequest("GET", "https://api.twitter.com/1.1/account/verify_credentials.json?include_email=true", nil)
	// if err != nil {
	// 	logger.Infof("%s", err.Error())
	// }
	// req.Header.Add("Content-Type", "application/json")
	// q = req.URL.Query()
	// q.Add("oauth_consumer_key", ConsumerKey)
	// q.Add("oauth_token")
	// req.URL.RawQuery = q.Encode()
	// client = &http.Client{}
	// resp, err = client.Do(req)
	// if err != nil {
	// 	logger.Infof("%s", err.Error())
	// }
	// b, err = ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	logger.Infof("%s", err.Error())
	// }
	// defer resp.Body.Close()

	// return c.String(200, string(b))

	// timelineURL := fmt.Sprintf("http://%s/time", r.Host)
	// return c.Redirect(http.StatusTemporaryRedirect, timelineURL)

	// token, err := a.service.Auth().OAuthFacebookCallback(ctx, code)
	// if err != nil {
	// 	logger.Infof("could not handle service callback: %s", err.Error())
	// 	return c.Redirect(http.StatusTemporaryRedirect, constants.RedirectFailOAuth)
	// }

	// return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("exp://192.168.43.2:19000/--/home?setToken=%s", "token"))

}
