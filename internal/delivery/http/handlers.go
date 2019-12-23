package http

import (
	"refit_backend/internal/delivery/http/auth"
	"refit_backend/internal/delivery/http/bodyweight"
	"refit_backend/internal/delivery/http/continent"
	"refit_backend/internal/delivery/http/countries"
	"refit_backend/internal/delivery/http/todos"
	"refit_backend/internal/delivery/http/tools"
	"refit_backend/internal/delivery/http/users"
)

// IHandler interface
type IHandler interface {
	Auth() auth.IAuth
	Users() users.IUsers
	Todos() todos.ITodos
	BodyWeight() bodyweight.IBodyWeight
	Tools() tools.ITools
	Continent() continent.IContinent
	Countries() countries.ICountries
}

type handler struct{}

// GetHandler ...
func (s *serverHTTP) GetHandler() IHandler {
	return &handler{}
}

func (h handler) Continent() continent.IContinent {
	return continent.New()
}

func (h handler) Countries() countries.ICountries {
	return countries.New()
}

func (h handler) Users() users.IUsers {
	return users.New()
}

func (h handler) Auth() auth.IAuth {
	return auth.New()
}

func (h handler) Todos() todos.ITodos {
	return todos.New()
}

func (h handler) BodyWeight() bodyweight.IBodyWeight {
	return bodyweight.New()
}

func (h handler) Tools() tools.ITools {
	return tools.New()
}
