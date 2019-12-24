package users

import (
	"context"
	"database/sql"
	"errors"
	"refit_backend/internal/helpers"
	"refit_backend/internal/logger"
	"refit_backend/internal/repository"
	"refit_backend/models"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	validation "github.com/go-ozzo/ozzo-validation/v3"

	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

// IUsers interface
type IUsers interface {
	Create(ctx context.Context, mu *models.User) (userID uint, err error)
	AuthLoginWithEmail(ctx context.Context, mu *models.User) (token string, err error)
	FindAll(ctx context.Context)
	Update(ctx context.Context)
	Delete(ctx context.Context)
	Count(ctx context.Context)
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

// Create services users
func (u users) Create(ctx context.Context, mu *models.User) (userID uint, err error) {

	err = mu.Validate()
	if err != nil {
		logger.Infof("could not validate: %s", err.Error())
		return 0, err
	}

	passwordHashed, err := bcrypt.GenerateFromPassword([]byte(mu.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Infof("could not hash password: %s", err.Error())
		return 0, errors.New("could not hash password")
	}

	mu.Password = string(passwordHashed)
	mu.RoleID = 2 // Set role_id to Normal User by default
	mu.CreatedAt = time.Now()
	mu.UpdatedAt = time.Now()

	userID, err = u.repository.Users().Create(ctx, mu)
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

// AuthLoginWithEmail service
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
			return "", errors.New("email or password not exists")
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
func (u users) FindAll(ctx context.Context) {}
func (u users) Update(ctx context.Context)  {}
func (u users) Delete(ctx context.Context)  {}
func (u users) Count(ctx context.Context)   {}
