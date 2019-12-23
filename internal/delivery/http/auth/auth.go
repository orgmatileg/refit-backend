package auth

import (
	"errors"
	"net/http"
	"refit_backend/internal/helpers"
	"refit_backend/internal/logger"
	"refit_backend/internal/services"
	"refit_backend/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
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

	mu, err := a.service.Users().AuthLoginWithEmail(ctx, ru.Email)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(mu.Password), []byte(ru.Password))
	if err != nil {
		logger.Warnf("could not compare hash password: %s", err.Error())
		return errors.New("password yang Anda masukkan salah")
	}

	claims := helpers.JWTPayload{
		StandardClaims: &jwt.StandardClaims{
			Audience:  "MOBILE",
			Issuer:    "Luqmanul Hakim API",
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(1440)).Unix(),
		},
	}

	token, err := helpers.GetJWTTokenGenerator().GenerateToken(claims)
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
