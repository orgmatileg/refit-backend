package http

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"refit_backend/internal/http/handler"
)

// IServerHTTP interface
type IServerHTTP interface {
	// Getter
	GetHTTP() *echo.Echo

	// Initiator
	InitRouter()
	InitMiddleWare()
}

// ServerHTTP struct
type serverHTTP struct {
	http *echo.Echo
}

// NewServerHTTP ...
func NewServerHTTP() IServerHTTP {
	return &serverHTTP{
		http: echo.New(),
	}
}

// GetHTTP ...
func (s *serverHTTP) GetHTTP() *echo.Echo {
	return s.http
}

// InitMiddleWare ...
func (s *serverHTTP) InitMiddleWare() {
	s.http.Pre(middleware.RemoveTrailingSlash())
}

// InitRouter ...
func (s *serverHTTP) InitRouter() {

	s.http.GET("/health", handler.HealthCheck)

	routerAuth := s.http.Group("auth")
	routerAuth.POST("/login", handler.AuthLoginWithEmail)
	routerAuth.POST("/register", handler.AuthRegister)

	routerUsers := s.http.Group("users")
	routerUsers.POST("", handler.UsersCreate)
	routerUsers.GET("", handler.UsersFindAll)
	routerUsers.GET("/:id", handler.UsersFindOne)
	routerUsers.POST("/:id", handler.UsersUpdate)
	routerUsers.DELETE("/:id", handler.UsersDelete)

	routerTodos := s.http.Group("todos")
	routerTodos.POST("", handler.TodosCreate)
	routerTodos.GET("", handler.TodosFindAll)
	routerTodos.GET("/:id", handler.TodosFindOne)
	routerTodos.POST("/:id", handler.TodosUpdate)
	routerTodos.DELETE("/:id", handler.TodosDelete)

	routerBodyWeight := s.http.Group("bodyweights")
	routerBodyWeight.POST("", handler.BodyWeightCreate)
	routerBodyWeight.GET("", handler.BodyWeightFindAll)
	routerBodyWeight.GET("/:id", handler.BodyWeightFindOne)
	routerBodyWeight.POST("/:id", handler.BodyWeightUpdate)
	routerBodyWeight.DELETE("/:id", handler.BodyWeightDelete)

	routerContinent := s.http.Group("continents")
	routerContinent.GET("", handler.ContinentFindAll)
	routerContinent.GET("/:id", handler.ContinentFindOne)

	routerCountries := s.http.Group("countries")
	routerCountries.GET("", handler.CountriesFindAll)
	routerCountries.GET("/:id", handler.CountriesFindOne)
}
