package handler

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"refit_backend/internal/helpers"
	"refit_backend/models"
	"time"
)

type IHandlers interface {
	GetDBMySQL()
}

// HealthCheck handler
func HealthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

func AuthLoginWithEmail(c echo.Context) error {
	var u models.User
	err := c.Bind(&u)
	if err != nil {
		return err
	}

	mu, err := models.FindUserByEmail(u.Email)
	if err != nil {
		return err
	}

	// Jika user berhasil ditemukan, maka langkah selanjutnya adalah memvalidasi password
	// hash yang ada di dalam database dengan inputan yang dikirimkan
	err = bcrypt.CompareHashAndPassword([]byte(mu.Password), []byte(u.Password))
	if err != nil {
		fmt.Println(err.Error())
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

func AuthRegister(c echo.Context) error {
	return c.String(http.StatusOK, "Register")
}
