package auth

import (
	"context"
	"refit_backend/internal/helpers"
	"refit_backend/internal/repository"
	"refit_backend/models"
	"regexp"
)

var (
	regexNumberOnly = regexp.MustCompile("^[0-9]*$")
)

// IAuth interface
type IAuth interface {
	OAuthGoogleCallback(ctx context.Context, code string) (tokenJWT string, err error)
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
			Gender:   "others",
			RoleID:   2,
			FullName: m.Name,
			Email:    m.Email,
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

	return tokenJWT, nil
}
