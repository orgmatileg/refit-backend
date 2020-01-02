package auth

import (
	"context"
	"refit_backend/internal/helpers"
	"refit_backend/internal/logger"
	"refit_backend/internal/repository"
	"refit_backend/models"
	"regexp"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	regexNumberOnly = regexp.MustCompile("^[0-9]*$")
)

// IAuth interface
type IAuth interface {
	OAuthGoogleCallback(ctx context.Context, code string) (tokenJWT string, err error)
	OAuthFacebookCallback(ctx context.Context, code string) (tokenJWT string, err error)
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
