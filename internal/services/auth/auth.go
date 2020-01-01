package auth

import (
	"context"
	"fmt"
	"refit_backend/internal/helpers"
	"refit_backend/internal/repository"
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
	fmt.Println(m)
	// mu := models.User{
	// 	FullName: m.Name,
	// 	Email:    m.Email,
	// 	RoleID:   2,
	// }

	// fmt.Println(mu)

	// userID, err := a.repository.Users().Create(ctx, &mu)
	// if err != nil {
	// 	return "", err
	// }

	return tokenJWT, nil
}
