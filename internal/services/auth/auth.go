package auth

import (
	"context"
	"database/sql"
	"io/ioutil"
	"net/http"
	"refit_backend/internal/helpers"
	"refit_backend/internal/logger"
	"refit_backend/internal/repository"
	"refit_backend/models"
	"regexp"
	"strings"
	"time"

	"github.com/dghubble/oauth1"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

var (
	regexNumberOnly = regexp.MustCompile("^[0-9]*$")
)

// IAuth interface
type IAuth interface {
	OAuthGoogleCallback(ctx context.Context, code string) (tokenJWT string, err error)
	OAuthFacebookCallback(ctx context.Context, code string) (tokenJWT string, err error)
	OAuthTwitterCallback(ctx context.Context, oauthToken, oauthVerifier string) (tokenJWT string, err error)
}

type auth struct {
	repository repository.IRepository
}

// New Repository Users
func New() IAuth {
	return &auth{
		repository: repository.New(),
	}
}

// FindOneByID services users
func (a auth) OAuthGoogleCallback(ctx context.Context, code string) (tokenJWT string, err error) {
	m, err := a.repository.Auth().GetUserDataFromGoogle(ctx, helpers.GetOAuthGoogleConfig(), code)
	if err != nil {
		return "", err
	}

	mo := models.OAuth{
		OpenID:  m.OpenID,
		Service: "google",
	}

	exist, _, err := a.repository.Users().IsExistsOAuth(ctx, &mo)
	if err != nil {
		return "", err
	}

	if !exist {

		mu, err := a.repository.Users().FindOneByEmail(ctx, m.Email)
		if err != nil && err == sql.ErrNoRows {
			userID, err := a.repository.Users().Create(ctx, &models.User{
				Gender:    "others",
				RoleID:    2,
				FullName:  m.Name,
				Email:     m.Email,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			})
			if err != nil {
				return "", err
			}
			mo.UserID = userID

		} else if err != nil {
			mo.UserID = mu.ID

		} else {
			logger.Infof("%s", err.Error())
			return "", err
		}

		_, err = a.repository.Users().StoreOAuth(ctx, &mo)
		if err != nil {
			return "", err
		}
	}

	claims := helpers.JWTPayload{
		StandardClaims: &jwt.StandardClaims{
			Audience:  "MOBILE",
			Issuer:    "Luqmanul Hakim API",
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(1440)).Unix(),
		},
	}

	tokenJWT, err = helpers.GetJWTTokenGenerator().GenerateToken(claims)
	if err != nil {
		logger.Infof("could not generate token: %s", err.Error())
		return "", err
	}

	return tokenJWT, nil
}

func (a auth) OAuthFacebookCallback(ctx context.Context, code string) (tokenJWT string, err error) {
	m, err := a.repository.Auth().GetUserDataFromFacebook(ctx, helpers.GetOAuthFacebookConfig(), code)
	if err != nil {
		return "", err
	}

	mo := models.OAuth{
		OpenID:  m.OpenID,
		Service: "facebook",
	}

	exist, _, err := a.repository.Users().IsExistsOAuth(ctx, &mo)
	if err != nil {
		return "", err
	}

	if !exist {
		mu, err := a.repository.Users().FindOneByEmail(ctx, m.Email)
		if err != nil && err == sql.ErrNoRows {
			userID, err := a.repository.Users().Create(ctx, &models.User{
				Gender:    "others",
				RoleID:    2,
				FullName:  m.Name,
				Email:     m.Email,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			})
			if err != nil {
				return "", err
			}
			mo.UserID = userID

		} else if err != nil {
			mo.UserID = mu.ID

		} else {
			logger.Infof("%s", err.Error())
			return "", err
		}

		_, err = a.repository.Users().StoreOAuth(ctx, &mo)
		if err != nil {
			return "", err
		}
	}

	claims := helpers.JWTPayload{
		StandardClaims: &jwt.StandardClaims{
			Audience:  "MOBILE",
			Issuer:    "Luqmanul Hakim API",
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(1440)).Unix(),
		},
	}

	tokenJWT, err = helpers.GetJWTTokenGenerator().GenerateToken(claims)
	if err != nil {
		logger.Infof("could not generate token: %s", err.Error())
		return "", err
	}

	return tokenJWT, nil
}

func (a auth) OAuthTwitterCallback(ctx context.Context, oauthToken, oauthVerifier string) (tokenJWT string, err error) {

	var (
		ConsumerKey    = viper.GetString("twitter.consumer_api_key")
		ConsumerSecret = viper.GetString("twitter.consumer_api_secret")
	)

	req, err := http.NewRequest("POST", "https://api.twitter.com/oauth/access_token", nil)
	if err != nil {
		logger.Infof("%s", err.Error())
	}
	req.Header.Add("Content-Type", "application/json")
	q := req.URL.Query()
	q.Add("oauth_token", oauthToken)
	q.Add("oauth_verifier", oauthVerifier)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Infof("%s", err.Error())
		return "", err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Infof("%s", err.Error())
		return "", err
	}

	resSlice := strings.Split(string(b), "&")
	m := make(map[string]string)
	for _, v := range resSlice {
		datas := strings.Split(v, "=")
		m[datas[0]] = datas[1]
	}

	token := oauth1.NewToken(m["oauth_token"], m["oauth_token_secret"])
	config := oauth1.NewConfig(ConsumerKey, ConsumerSecret)
	httpClient := config.Client(ctx, token)

	mt, err := a.repository.Auth().GetUserDataFromTwitter(httpClient)
	if err != nil {
		logger.Infof("%s", err.Error())
		return "", err
	}

	mo := models.OAuth{
		OpenID:  mt.IDStr,
		Service: "twitter",
	}

	exist, _, err := a.repository.Users().IsExistsOAuth(ctx, &mo)
	if err != nil {
		return "", err
	}

	if !exist {

		mu, err := a.repository.Users().FindOneByEmail(ctx, mt.Email)
		if err != nil && err == sql.ErrNoRows {
			userID, err := a.repository.Users().Create(ctx, &models.User{
				Gender:    "others",
				RoleID:    2,
				FullName:  mt.Name,
				Email:     mt.Email,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			})
			if err != nil {
				return "", err
			}
			mo.UserID = userID

		} else if err != nil {
			mo.UserID = mu.ID

		} else {
			logger.Infof("%s", err.Error())
			return "", err
		}

		_, err = a.repository.Users().StoreOAuth(ctx, &mo)
		if err != nil {
			return "", err
		}

	}

	claims := helpers.JWTPayload{
		StandardClaims: &jwt.StandardClaims{
			Audience:  "MOBILE",
			Issuer:    "Luqmanul Hakim API",
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(1440)).Unix(),
		},
	}

	tokenJWT, err = helpers.GetJWTTokenGenerator().GenerateToken(claims)
	if err != nil {
		logger.Infof("could not generate token: %s", err.Error())
		return "", err
	}

	return tokenJWT, nil
}
