package bodyweight

import (
	"github.com/labstack/echo"
	"net/http"
	"refit_backend/internal/helpers"
	"refit_backend/internal/logger"
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

	fh, err := c.FormFile("image")
	if err != nil && err.Error() != "http: no such file" {
		logger.Infof("could not read form file from request: %s", err.Error())
		return helpers.MakeDefaultResponse(c, http.StatusBadRequest, err)
	}

	var (
		ctx    = c.Request().Context()
		weight = c.FormValue("weight")
		date   = c.FormValue("date")
		userID = c.FormValue("user_id")
	)

	_, err = b.service.BodyWeight().Create(ctx, weight, date, userID, fh)
	if err != nil {
		return helpers.MakeDefaultResponse(c, http.StatusBadRequest, err)
	}
	return helpers.MakeDefaultResponse(c, http.StatusCreated, nil)
}

// FindOneByID delivery http users
func (b bodyweight) FindOneByID(c echo.Context) error {
	bodyweightID := c.Param("id")
	ctx := c.Request().Context()
	m, err := b.service.BodyWeight().FindOneByID(ctx, bodyweightID)
	if err != nil {
		return helpers.MakeDefaultResponse(c, http.StatusBadRequest, err)
	}
	return helpers.MakeDefaultResponse(c, http.StatusOK, m)
}

// FindAll delivery http users
func (b bodyweight) FindAll(c echo.Context) error {
	var (
		ctx    = c.Request().Context()
		limit  = "20"
		offset = "0"
		order  = "desc"
		userID = ""
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
	if v := c.QueryParam("user_id"); v != "" {
		userID = v
	}

	mu, count, err := b.service.BodyWeight().FindAll(ctx, limit, offset, order, userID)
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
	var rm models.BodyWeight
	err := c.Bind(&rm)
	if err != nil {
		return err
	}
	bodyWeightID := c.Param("id")
	ctx := c.Request().Context()
	err = b.service.BodyWeight().UpdateByID(ctx, &rm, bodyWeightID)
	if err != nil {
		return helpers.MakeDefaultResponse(c, http.StatusBadRequest, err)
	}
	return helpers.MakeDefaultResponse(c, http.StatusOK, nil)
}

// DeleteByID http users
func (b bodyweight) DeleteByID(c echo.Context) error {
	bodyweightID := c.Param("id")
	ctx := c.Request().Context()
	err := b.service.BodyWeight().DeleteByID(ctx, bodyweightID)
	if err != nil {
		return helpers.MakeDefaultResponse(c, http.StatusBadRequest, err)
	}
	return helpers.MakeDefaultResponse(c, http.StatusOK, nil)
}
