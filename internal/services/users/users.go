package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"refit_backend/internal/helpers"
	"refit_backend/internal/logger"
	"refit_backend/internal/repository"
	"refit_backend/models"
	"regexp"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	validation "github.com/go-ozzo/ozzo-validation/v3"

	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

var (
	regexNumberOnly = regexp.MustCompile("^[0-9]*$")
)

// IUsers interface
type IUsers interface {
	AuthLoginWithEmail(ctx context.Context, mu *models.User) (token string, err error)
	Create(ctx context.Context, ru *models.User) (userID uint, err error)
	FindOneByID(ctx context.Context, userID string) (mu *models.User, err error)
	FindAll(ctx context.Context, limit, offset, order string) (lmu []*models.User, count uint, err error)
	UpdateByID(ctx context.Context, ru *models.User, userID string) (err error)
	DeleteByID(ctx context.Context, userID string) (err error)
}

type users struct {
	repository repository.IRepository
}

// New Repository Users
func New() IUsers {
	return &users{
		repository: repository.New(),
	}
}

// FindOneByID services users
func (u users) FindOneByID(ctx context.Context, userID string) (mu *models.User, err error) {

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

// Create services users
func (u users) Create(ctx context.Context, ru *models.User) (userID uint, err error) {

	err = ru.ValidateCreate()
	if err != nil {
		logger.Infof("could not validate: %s", err.Error())
		return 0, err
	}

	passwordHashed, err := bcrypt.GenerateFromPassword([]byte(ru.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Infof("could not hash password: %s", err.Error())
		return 0, errors.New("could not hash password")
	}

	ru.Password = string(passwordHashed)
	ru.RoleID = 2 // Set role_id to Normal User by default
	ru.CreatedAt = time.Now()
	ru.UpdatedAt = time.Now()

	userID, err = u.repository.Users().Create(ctx, ru)
	fmt.Print(err)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate") {
			logger.Infof("invalid email: email already used")
			return 0, errors.New("invalid email: email already used")
		}
		logger.Infof("could not create user to db: %s", err.Error())
		return 0, err
	}

	return userID, nil
}

// AuthLoginWithEmail services users
func (u users) AuthLoginWithEmail(ctx context.Context, ru *models.User) (token string, err error) {

	err = validation.ValidateStruct(ru,
		validation.Field(&ru.Email, validation.Required, is.Email),
		validation.Field(&ru.Password, validation.Required),
	)
	if err != nil {
		logger.Infof("could not validate: %s", err.Error())
		return "", err
	}

	mu, err := u.repository.Users().FindOneByEmail(ctx, ru.Email)

	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			logger.Infof("could not find user by email: %s", err.Error())
			return "", errors.New("email or password is wrong")
		default:
			logger.Infof("could not find user by email: %s", err.Error())
			return "", err
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(mu.Password), []byte(ru.Password))
	if err != nil {
		logger.Warnf("could not compare hash password: %s", err.Error())
		return "", errors.New("password yang Anda masukkan salah")
	}

	claims := helpers.JWTPayload{
		StandardClaims: &jwt.StandardClaims{
			Audience:  "MOBILE",
			Issuer:    "Luqmanul Hakim API",
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(1440)).Unix(),
		},
	}

	token, err = helpers.GetJWTTokenGenerator().GenerateToken(claims)
	if err != nil {
		logger.Infof("could not generate token: %s", err.Error())
		return "", err
	}

	return token, nil
}

// FindAll service users
func (u users) FindAll(ctx context.Context, limit, offset, order string) (lmu []*models.User, count uint, err error) {
	err = helpers.ValidationQueryParamFindAll(limit, offset, order)
	if err != nil {
		logger.Infof("could not validate: %s", err.Error())
		return nil, 0, err
	}

	lmu, err = u.repository.Users().FindAll(ctx, limit, offset, order)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			logger.Infof("could not find all user: %s", err.Error())
			return nil, 0, errors.New("no row exists")
		default:
			logger.Infof("could not find all user: %s", err.Error())
			return nil, 0, err
		}
	}

	// remove sensitive information
	for i := range lmu {
		lmu[i].Password = ""
	}

	count, err = u.repository.Users().Count(ctx)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			logger.Infof("could not find all user: %s", err.Error())
			return nil, 0, errors.New("no row exists")
		default:
			logger.Infof("could not find all user: %s", err.Error())
			return nil, 0, err
		}
	}

	return lmu, count, nil
}

//
func (u users) UpdateByID(ctx context.Context, ru *models.User, userID string) (err error) {
	err = ru.ValidateUpdate()
	if err != nil {
		logger.Infof("could not validate: %s", err.Error())
		return err
	}

	mu, err := u.repository.Users().FindOneByID(ctx, userID)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			logger.Infof("could not find user by id: %s", err.Error())
			return errors.New("userID not exists")
		default:
			logger.Infof("could not find user by id: %s", err.Error())
			return err
		}
	}

	ru.Email = mu.Email
	ru.RoleID = mu.RoleID
	ru.Password = mu.Password
	ru.UpdatedAt = time.Now()

	_, err = u.repository.Users().UpdateByID(ctx, ru, userID)
	if err != nil {
		logger.Infof("could not update user by id: %s", err.Error())
		return err
	}

	return nil
}

func (u users) DeleteByID(ctx context.Context, userID string) (err error) {

	err = validation.Validate(userID, validation.Match(regexNumberOnly))
	if err != nil {
		logger.Infof("could not validate: %s", err.Error())
		return errors.New("invalid userID param, should be number only")
	}
	_, err = u.repository.Users().DeleteByID(ctx, userID)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			logger.Infof("could not delete user by id: %s", err.Error())
			return errors.New("userID not exists")
		default:
			logger.Infof("could not delete user by id: %s", err.Error())
			return err
		}
	}

	return nil
}
