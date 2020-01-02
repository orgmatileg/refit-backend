package bodyweight

import (
	"github.com/labstack/echo"
	"net/http"
	"refit_backend/internal/helpers"
	"refit_backend/internal/services"
	"refit_backend/models"
	"strconv"
)

type IBodyWeight interface {
	Create(c echo.Context) error
	FindOneByID(c echo.Context) error
	FindAll(c echo.Context) error
	UpdateByID(c echo.Context) error
	DeleteByID(c echo.Context) error
}

type bodyweight struct {
	service services.IServices
}

// New Repository todos
func New() IBodyWeight {
	return &bodyweight{
		service: services.New(),
	}
}

// Create delivery http users
func (b bodyweight) Create(c echo.Context) error {
	var rm models.BodyWeight
	err := c.Bind(&rm)
	if err != nil {
		return err
	}
	ctx := c.Request().Context()
	_, err = b.service.BodyWeight().Create(ctx, &rm)
	if err != nil {
		return helpers.MakeDefaultResponse(c, http.StatusBadRequest, err)
	}
	return helpers.MakeDefaultResponse(c, http.StatusCreated, nil)
}

// FindOneByID delivery http users
func (b bodyweight) FindOneByID(c echo.Context) error {
	userID := c.Param("id")
	ctx := c.Request().Context()

	mu, err := b.service.Users().FindOneByID(ctx, userID)
	if err != nil {
		return helpers.MakeDefaultResponse(c, http.StatusBadRequest, err)
	}
	return helpers.MakeDefaultResponse(c, http.StatusOK, mu)
}

// FindAll delivery http users
func (b bodyweight) FindAll(c echo.Context) error {
	var (
		ctx    = c.Request().Context()
		limit  = "20"
		offset = "0"
		order  = "desc"
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

	mu, count, err := b.service.Users().FindAll(ctx, limit, offset, order)
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
func (b bodyweight) UpdateByID(c echo.Context) error {
	var ru models.User
	err := c.Bind(&ru)
	if err != nil {
		return err
	}
	userID := c.Param("id")

	ctx := c.Request().Context()
	err = b.service.Users().UpdateByID(ctx, &ru, userID)
	if err != nil {
		return helpers.MakeDefaultResponse(c, http.StatusBadRequest, err)
	}
	return helpers.MakeDefaultResponse(c, http.StatusOK, nil)
}

// DeleteByID http users
func (b bodyweight) DeleteByID(c echo.Context) error {
	userID := c.Param("id")
	ctx := c.Request().Context()
	err := b.service.Users().DeleteByID(ctx, userID)
	if err != nil {
		return helpers.MakeDefaultResponse(c, http.StatusBadRequest, err)
	}
	return helpers.MakeDefaultResponse(c, http.StatusOK, nil)
}
