package users

import (
	"github.com/labstack/echo"
)

type IUsers interface {
	Create(c echo.Context) error
	FindOne(c echo.Context) error
	FindAll(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
	Count(c echo.Context) error
}

type users struct{}

// New Repository Users
func New() IUsers {
	return &users{}
}

func (u users) Create(c echo.Context) error {
	return nil
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
