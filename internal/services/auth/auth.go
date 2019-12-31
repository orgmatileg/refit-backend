package auth

import (
	"context"
	"database/sql"
	"errors"
	"refit_backend/internal/logger"
	"refit_backend/internal/repository"
	"refit_backend/models"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v3"
)

var (
	regexNumberOnly = regexp.MustCompile("^[0-9]*$")
)

// IAuth interface
type IAuth interface {
	FindOneByID(ctx context.Context, userID string) (mu *models.User, err error)
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
func (u auth) FindOneByID(ctx context.Context, userID string) (mu *models.User, err error) {

	err = validation.Validate(userID, validation.Match(regexNumberOnly))
	if err != nil {
		logger.Infof("could not validate: %s", err.Error())
		return nil, errors.New("invalid userID param, should be number only")
	}
	mu, err = u.repository.Users().FindOneByID(ctx, userID)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			logger.Infof("could not find user by id: %s", err.Error())
			return nil, errors.New("userID not exists")
		default:
			logger.Infof("could not find user by id: %s", err.Error())
			return nil, err
		}
	}

	// Remove sensitive data
	mu.Password = ""

	return mu, nil
}
