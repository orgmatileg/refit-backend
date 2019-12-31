package services

import (
	"refit_backend/internal/repository"
	"refit_backend/internal/services/auth"
	"refit_backend/internal/services/bodyweight"
	"refit_backend/internal/services/todos"
	"refit_backend/internal/services/users"
)

// IServices interface
type IServices interface {
	Auth() auth.IAuth
	Users() users.IUsers
	Todos() todos.ITodos
	BodyWeight() bodyweight.IBodyWeight
}

type services struct {
	repository repository.IRepository
}

// New Services
func New() IServices {
	return &services{
		repository: repository.New(),
	}
}

func (s services) Auth() auth.IAuth {
	return auth.New()
}

func (s services) Users() users.IUsers {
	return users.New()
}

func (s services) Todos() todos.ITodos {
	return todos.New()
}

func (s services) BodyWeight() bodyweight.IBodyWeight {
	return bodyweight.New()
}
