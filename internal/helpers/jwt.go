package helpers

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

// JWTPayload struct
type JWTPayload struct {
	*jwt.StandardClaims
}

// Token JWT Token Model
type Token struct {
	TokenString string
	Key         string
}

// TokenGenerator token generator model
type TokenGenerator struct {
	Key string
	Alg jwt.SigningMethod
}

// ParseJWTClaim func
func ParseJWTClaim(tokenString string) (jwt.MapClaims, error) {

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("jwt.secret")), nil
	})
	if err != nil {
		return nil, err
	}

	return claims, nil
}

// JwtDecodeHelper untuk membantu decode payload data jwt
func JwtDecodeHelper(tokenString string) (jwt.Claims, error) {

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("jwt.secret")), nil
	})

	if err != nil {
		return token.Claims, nil
	}

	return nil, err
}

func parseToken(t *Token) (token *jwt.Token, err error) {
	token, err = jwt.Parse(t.TokenString, func(jt *jwt.Token) (interface{}, error) {
		// Untuk mencegah JWT Signing method NONE attack
		// Maka pastikan untuk memvalidasi juga Algoritma signing nya
		if _, ok := jt.Method.(*jwt.SigningMethodHMAC); ok {
			return []byte(t.Key), nil
		}
		return nil, errors.New("Unknown Signing Method")
	})
	return
}

// GenerateToken generate token
func (gen *TokenGenerator) GenerateToken(claims jwt.Claims) (token string, err error) {
	t := jwt.NewWithClaims(gen.Alg, claims)
	token, err = t.SignedString([]byte(gen.Key))
	return
}

// GetPayload func
func (t *Token) GetPayload() jwt.MapClaims {
	if token, err := parseToken(t); err == nil && token.Valid {
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			return claims
		}
	}
	return nil
}

// IsValidToken validate JWT Token
func (t *Token) IsValidToken() bool {
	if token, err := parseToken(t); err == nil {
		return token.Valid
	}
	return false
}

// GetJWTTokenGenerator ...
func GetJWTTokenGenerator() *TokenGenerator {
	t := &TokenGenerator{
		Key: viper.GetString("jwt.secret"),
		Alg: jwt.SigningMethodHS256,
	}
	return t
}
