package users

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"refit_backend/internal/logger"
	"refit_backend/internal/repository"
	"refit_backend/models"
	"strings"
	"time"
)

type IUsers interface {
	Create(ctx context.Context, m *models.User) (uint, error)
	AuthLoginWithEmail(ctx context.Context, email string) (*models.User, error)
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

// TODO: add validation
func (u users) Create(ctx context.Context, mu *models.User) (uint, error) {
	if len(mu.FullName) < 2 {
		return 0, errors.New("invalid full_name: should greater than 2")
	}

	passwordHashed, err := bcrypt.GenerateFromPassword([]byte(mu.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, errors.New("could not hash password")
	}
	mu.Password = string(passwordHashed)
	mu.RoleID = 2 // Set role_id to Normal User by default
	mu.CreatedAt = time.Now()
	mu.UpdatedAt = time.Now()

	_, err = u.repository.Users().Create(ctx, mu)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate") {
			return 0, errors.New("invalid email: email already used")
		}
		return 0, err
	}

	return 0, nil
}

// AuthLoginWithEmail service
func (u users) AuthLoginWithEmail(ctx context.Context, email string) (*models.User, error) {
	if email == "" {
		logger.Debugf("could not find user with empty string email")
		return nil, errors.New("invalid email")
	}
	mu, err := u.repository.Users().FindOneByEmail(ctx, email)
	if err != nil {
		logger.Debugf("could not find user by email: %s", err.Error())
		return nil, err
	}
	return mu, nil
}
func (u users) FindAll(ctx context.Context) {}
func (u users) Update(ctx context.Context)  {}
func (u users) Delete(ctx context.Context)  {}
func (u users) Count(ctx context.Context)   {}
