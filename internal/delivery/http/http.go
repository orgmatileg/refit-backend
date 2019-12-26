package http

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// IServerHTTP interface
type IServerHTTP interface {
	// Getter
	GetHTTP() *echo.Echo

	// Initiator
	InitRouter()
	InitMiddleWare()

	GetHandler() IHandler
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
	s.http.Pre(middleware.CORS())

	// s.http.Pre(middleware.Logger())
	// s.http.Pre(middleware.JWTWithConfig(middleware.JWTConfig{
	// 	SigningKey:  []byte("secret"),
	// 	TokenLookup: "query:token",
	// }))

	s.http.HTTPErrorHandler = s.GetHandler().Tools().DefaultErrorHandler

}

// InitRouter ...
func (s *serverHTTP) InitRouter() {

	s.http.GET("/health", s.GetHandler().Tools().HealthCheck)

	routerAuth := s.http.Group("auth")
	routerAuth.POST("/login", s.GetHandler().Auth().AuthLoginWithEmail)
	routerAuth.POST("/register", s.GetHandler().Auth().AuthRegister)

	routerUsers := s.http.Group("users")
	routerUsers.POST("", s.GetHandler().Users().Create)
	routerUsers.GET("", s.GetHandler().Users().FindAll)
	routerUsers.GET("/:id", s.GetHandler().Users().FindOneByID)
	routerUsers.PUT("/:id", s.GetHandler().Users().UpdateByID)
	routerUsers.DELETE("/:id", s.GetHandler().Users().DeleteByID)

	routerTodos := s.http.Group("todos")
	routerTodos.POST("", s.GetHandler().Todos().Create)
	routerTodos.GET("", s.GetHandler().Todos().FindAll)
	routerTodos.GET("/:id", s.GetHandler().Todos().FindOne)
	routerTodos.POST("/:id", s.GetHandler().Todos().Update)
	routerTodos.DELETE("/:id", s.GetHandler().Todos().Delete)

	routerBodyWeight := s.http.Group("bodyweights")
	routerBodyWeight.POST("", s.GetHandler().BodyWeight().Create)
	routerBodyWeight.GET("", s.GetHandler().BodyWeight().FindAll)
	routerBodyWeight.GET("/:id", s.GetHandler().BodyWeight().FindOne)
	routerBodyWeight.POST("/:id", s.GetHandler().BodyWeight().Update)
	routerBodyWeight.DELETE("/:id", s.GetHandler().BodyWeight().Delete)

	routerContinent := s.http.Group("continents")
	routerContinent.GET("", s.GetHandler().Continent().FindAll)
	routerContinent.GET("/:id", s.GetHandler().Continent().FindOne)

	routerCountries := s.http.Group("countries")
	routerCountries.GET("", s.GetHandler().Countries().FindAll)
	routerCountries.GET("/:id", s.GetHandler().Countries().FindOne)
}
