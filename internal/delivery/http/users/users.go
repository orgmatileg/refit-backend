package users

import (
	"net/http"
	"refit_backend/internal/helpers"
	"refit_backend/internal/services"
	"refit_backend/models"

	"github.com/labstack/echo"
)

// IUsers interface delivery http
type IUsers interface {
	Create(c echo.Context) error
	FindOne(c echo.Context) error
	FindAll(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
	Count(c echo.Context) error
}

// users struct delivery http
type users struct {
	service services.IServices
}

// New delivery http Users
func New() IUsers {
	return &users{
		service: services.New(),
	}
}

// Create delivery http users
func (u users) Create(c echo.Context) error {
	var ru models.User
	err := c.Bind(&ru)
	if err != nil {
		return err
	}
	ctx := c.Request().Context()
	_, err = u.service.Users().Create(ctx, &ru)
	if err != nil {
		return helpers.MakeDefaultResponse(c, http.StatusBadRequest, err)
	}
	return helpers.MakeDefaultResponse(c, http.StatusCreated, nil)
}
func (u users) FindOne(c echo.Context) error {
	return nil
}
func (u users) FindAll(c echo.Context) error {
	return nil
}
func (u users) Update(c echo.Context) error {
	return nil
}
func (u users) Delete(c echo.Context) error {
	return nil
}
func (u users) Count(c echo.Context) error {
	return nil
}
