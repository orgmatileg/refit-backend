package serve

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"refit_backend/internal/delivery/http"
	"refit_backend/internal/infrastructures/mysql"
	"refit_backend/internal/infrastructures/s3"

	"github.com/labstack/echo"

	"refit_backend/internal/logger"
	"refit_backend/internal/repository"
	"time"
)

// IAppServe interface
type IAppServe interface {
	// Getter
	GetCtx() context.Context
	GetHTTP() *echo.Echo
	GetDBMySQL() *sql.DB

	// Initiator
	InitCtx()
	InitLogger()
	InitMySQL()
	InitHTTP()
	InitS3()
}

// AppServe struct
type appServe struct {
	ctx        context.Context
	http       http.IServerHTTP
	mysql      *sql.DB
	logger     logger.Logger
	repository repository.IRepository
}

func (a *appServe) GetHTTP() *echo.Echo {
	return a.http.GetHTTP()
}

func (a *appServe) GetCtx() context.Context {
	return a.ctx
}

func (a *appServe) GetDBMySQL() *sql.DB {
	return a.mysql
}

func (a *appServe) InitS3() {
	s3.Init()
}
func (a *appServe) InitLogger() {
	// Init Logger
	config := logger.Configuration{
		EnableConsole:     true,
		ConsoleLevel:      logger.Debug,
		ConsoleJSONFormat: false,
		EnableFile:        false,
		FileLevel:         logger.Info,
		FileJSONFormat:    false,
		FileLocation:      "refit-backend.log",
	}
	err := logger.NewLogger(config, logger.InstanceLogrusLogger)
	if err != nil {
		log.Fatalf("Could not instantiate log %s", err.Error())
	}
}

func (a *appServe) InitMySQL() {
	a.mysql = mysql.GetDB()
}

func (a *appServe) InitHTTP() {
	a.http = http.NewServerHTTP()
	a.http.InitMiddleWare()
	a.http.InitRouter()
}

func (a *appServe) InitCtx() {
	a.ctx = context.Background()
}

func newAppServe() IAppServe {
	return &appServe{}
}

// Start Serve App
func Start() {
	app := newAppServe()

	// Initiator
	app.InitCtx()
	app.InitLogger()
	app.InitMySQL()
	app.InitHTTP()
	app.InitS3()

	// Start server
	go func() {
		if err := app.GetHTTP().Start("0.0.0.0:1323"); err != nil {
			logger.Infof("could not start HTTP Server: %s", err.Error())
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	logger.Infof("Shutting down HTTP Server")

	ctx, cancel := context.WithTimeout(app.GetCtx(), 10*time.Second)
	defer cancel()
	if err := app.GetHTTP().Shutdown(ctx); err != nil {
		logger.Fatalf("could not shutdown HTTP Server: %s", err.Error())
	}
}
