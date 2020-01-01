package repository

import (
	"refit_backend/internal/repository/auth"
	"refit_backend/internal/repository/bodyweight"
	"refit_backend/internal/repository/todos"
	"refit_backend/internal/repository/users"
)

// IRepository interface
type IRepository interface {
	Auth() auth.IAuth
	Users() users.IUsers
	Todos() todos.ITodos
	BodyWeight() bodyweight.IBodyWeight
}

type repository struct{}

// New Repository
func New() IRepository {
	return &repository{}
}

func (r repository) Auth() auth.IAuth {
	return auth.New()
}

func (r repository) Users() users.IUsers {
	return users.New()
}

func (r repository) Todos() todos.ITodos {
	return todos.New()
}

func (r repository) BodyWeight() bodyweight.IBodyWeight {
	return bodyweight.New()
}
