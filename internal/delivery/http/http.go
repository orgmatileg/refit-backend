package http

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
)

// IServerHTTP interface
type IServerHTTP interface {
	// Getter
	GetHTTP() *echo.Echo
	GetHandler() IHandler

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
	e := echo.New()
	// e.Server.Addr = "0.0.0.0"
	// e.Server.
	return &serverHTTP{
		http: e,
	}
}

// GetHTTP ...
func (s *serverHTTP) GetHTTP() *echo.Echo {
	return s.http
}

// InitMiddleWare ...
func (s *serverHTTP) InitMiddleWare() {
	s.http.HTTPErrorHandler = s.GetHandler().Tools().DefaultErrorHandler
	s.http.Pre(middleware.RemoveTrailingSlash())
	s.http.Pre(middleware.CORS())
}

// InitRouter HTTP
func (s *serverHTTP) InitRouter() {

	// ============================== //
	//      Middleware Injection      //
	// ============================== //
	middlewareAuth := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:  []byte(viper.GetString("jwt.secret")),
		TokenLookup: "header:Authorization",
	})

	s.http.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "static",
		Browse: false,
		HTML5:  true,
		Index:  "index.html",
	}))

	// ============================== //
	//        No Need Auth            //
	// ============================== //
	s.http.GET("/health", s.GetHandler().Tools().HealthCheck)

	routerContinent := s.http.Group("continents")
	routerContinent.GET("", s.GetHandler().Continent().FindAll)
	routerContinent.GET("/:id", s.GetHandler().Continent().FindOne)

	routerCountries := s.http.Group("countries")
	routerCountries.GET("", s.GetHandler().Countries().FindAll)
	routerCountries.GET("/:id", s.GetHandler().Countries().FindOne)

	routerAuth := s.http.Group("auth")
	routerAuth.POST("/login", s.GetHandler().Auth().AuthLoginWithEmail)
	routerAuth.POST("/register", s.GetHandler().Auth().AuthRegister)
	routerAuth.GET("/google/login", s.GetHandler().Auth().OAuthGoogleLogin)
	routerAuth.GET("/google/callback", s.GetHandler().Auth().OAuthGoogleCallback)
	routerAuth.GET("/facebook/login", s.GetHandler().Auth().OAuthFacebookLogin)
	routerAuth.GET("/facebook/callback", s.GetHandler().Auth().OAuthFacebookCallback)
	routerAuth.GET("/twitter/login", s.GetHandler().Auth().OAuthTwitterLogin)
	routerAuth.GET("/twitter/callback", s.GetHandler().Auth().OAuthTwitterCallback)

	// ============================== //
	//           Need Auth            //
	// ============================== //

	routerUsers := s.http.Group("users", middlewareAuth)
	routerUsers.POST("", s.GetHandler().Users().Create)
	routerUsers.GET("", s.GetHandler().Users().FindAll)
	routerUsers.GET("/:id", s.GetHandler().Users().FindOneByID)
	routerUsers.PUT("/:id", s.GetHandler().Users().UpdateByID)
	routerUsers.DELETE("/:id", s.GetHandler().Users().DeleteByID)

	routerTodos := s.http.Group("todos", middlewareAuth)
	routerTodos.POST("", s.GetHandler().Todos().Create)
	routerTodos.GET("", s.GetHandler().Todos().FindAll)
	routerTodos.GET("/:id", s.GetHandler().Todos().FindOneByID)
	routerTodos.POST("/:id", s.GetHandler().Todos().UpdateByID)
	routerTodos.DELETE("/:id", s.GetHandler().Todos().DeleteByID)

	routerBodyWeight := s.http.Group("bodyweights", middlewareAuth)
	routerBodyWeight.POST("", s.GetHandler().BodyWeight().Create)
	routerBodyWeight.GET("", s.GetHandler().BodyWeight().FindAll)
	routerBodyWeight.GET("/:id", s.GetHandler().BodyWeight().FindOneByID)
	routerBodyWeight.POST("/:id", s.GetHandler().BodyWeight().UpdateByID)
	routerBodyWeight.DELETE("/:id", s.GetHandler().BodyWeight().DeleteByID)

}
