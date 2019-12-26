package users

import (
	"net/http"
	"refit_backend/internal/helpers"
	"refit_backend/internal/services"
	"refit_backend/models"
	"strconv"

	"github.com/labstack/echo"
)

// IUsers interface delivery http
type IUsers interface {
	Create(c echo.Context) error
	FindOneByID(c echo.Context) error
	FindAll(c echo.Context) error
	UpdateByID(c echo.Context) error
	DeleteByID(c echo.Context) error
}

// users struct delivery http
type users struct {
	service services.IServices
}

// New delivery http users
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

// FindOneByID delivery http users
func (u users) FindOneByID(c echo.Context) error {
	userID := c.Param("id")
	ctx := c.Request().Context()

	mu, err := u.service.Users().FindOneByID(ctx, userID)
	if err != nil {
		return helpers.MakeDefaultResponse(c, http.StatusBadRequest, err)
	}
	return helpers.MakeDefaultResponse(c, http.StatusOK, mu)
}

// FindAll delivery http users
func (u users) FindAll(c echo.Context) error {
	var (
		ctx    = c.Request().Context()
		limit  = "20"
		offset = "0"
		order  = "dsc"
	)

	if v := c.QueryParam("limit"); v != "" {
		limit = v
	}
	if v := c.QueryParam("offset"); v != "" {
		offset = c.QueryParam("offset")
	}
	if v := c.QueryParam("order"); v != "" {
		order = v
	}

	mu, count, err := u.service.Users().FindAll(ctx, limit, offset, order)
	if err != nil {
		return helpers.MakeDefaultResponse(c, http.StatusBadRequest, err)
	}

	backURL, nextURL, err := helpers.CreateBackAndNextURL(c, count, limit, offset, order)

	rqp := models.RequestQueryParam{
		"back_url": backURL,
		"next_url": nextURL,
		"count":    strconv.Itoa(int(count)),
		"limit":    limit,
		"offset":   offset,
		"order":    order,
	}

	return helpers.MakeDefaultResponse(c, http.StatusOK, mu, rqp)
}

// UpdateByID delivery http users
func (u users) UpdateByID(c echo.Context) error {
	var ru models.User
	err := c.Bind(&ru)
	if err != nil {
		return err
	}
	userID := c.Param("id")

	ctx := c.Request().Context()
	err = u.service.Users().UpdateByID(ctx, &ru, userID)
	if err != nil {
		return helpers.MakeDefaultResponse(c, http.StatusBadRequest, err)
	}
	return helpers.MakeDefaultResponse(c, http.StatusOK, nil)
}

// DeleteByID http users
func (u users) DeleteByID(c echo.Context) error {
	userID := c.Param("id")
	ctx := c.Request().Context()
	err := u.service.Users().DeleteByID(ctx, userID)
	if err != nil {
		return helpers.MakeDefaultResponse(c, http.StatusBadRequest, err)
	}
	return helpers.MakeDefaultResponse(c, http.StatusOK, nil)
}
