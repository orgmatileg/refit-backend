package auth

import (
	"net/http"
	"refit_backend/internal/helpers"
	"refit_backend/internal/services"
	"refit_backend/models"

	"github.com/labstack/echo"
)

// IAuth interface
type IAuth interface {
	AuthLoginWithEmail(c echo.Context) error
	AuthRegister(c echo.Context) error
}

type auth struct {
	service services.IServices
}

// New auth http handler
func New() IAuth {
	return &auth{
		service: services.New(),
	}
}

func (a auth) AuthLoginWithEmail(c echo.Context) error {
	var ru models.User
	err := c.Bind(&ru)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()

	token, err := a.service.Users().AuthLoginWithEmail(ctx, &ru)
	if err != nil {
		return err
	}

	return helpers.MakeDefaultResponse(c, http.StatusOK, struct {
		Token string `json:"token"`
	}{
		Token: token,
	})
}

func (a auth) AuthRegister(c echo.Context) error {
	var ru models.User
	err := c.Bind(&ru)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	_, err = a.service.Users().Create(ctx, &ru)
	if err != nil {
		return helpers.MakeDefaultResponse(c, http.StatusBadRequest, err)
	}

	return helpers.MakeDefaultResponse(c, http.StatusCreated, nil)
}
