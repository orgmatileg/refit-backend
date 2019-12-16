package serve

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"refit_backend/internal/http/router"
	"refit_backend/internal/logger"
	"refit_backend/internal/mysql"

	"github.com/labstack/echo"

	"time"
)

// IAppServe interface
type IAppServe interface {
	// Getter
	GetCtx() context.Context
	GetEchoHTTP() *echo.Echo
	GetLogger() logger.Logger
	GetDBMySQL() *sql.DB

	// Initiator
	InitCtx()
	InitLogger()
	InitMySQL()
	InitHTTP()
}

// AppServe struct
type appServe struct {
	ctx      context.Context
	echoHTTP *echo.Echo
	logger   logger.Logger
	mysql    mysql.IDBMySQL
}

func (a *appServe) GetEchoHTTP() *echo.Echo {
	return a.echoHTTP
}

func (a *appServe) GetCtx() context.Context {
	return a.ctx
}

func (a *appServe) GetLogger() logger.Logger {
	return a.logger
}

func (a *appServe) GetDBMySQL() *sql.DB {
	return a.mysql.GetDB()
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

	loggerz := logger.WithFields(
		logger.Fields{"logger_instance": "zap"},
	)

	a.logger = loggerz
}

func (a *appServe) InitMySQL() {
	a.mysql = mysql.NewDBMySQL()
	err := a.mysql.CreateConnection()
	if err != nil {
		a.GetLogger().Fatalf("could not create to mysql database: %s", err.Error())
	}
}

func (a *appServe) InitHTTP() {
	e := echo.New()
	router.SetHTTPRouter(e)
	a.echoHTTP = e
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

	app.InitCtx()
	app.InitLogger()
	app.InitMySQL()
	app.InitHTTP()

	// Start server
	go func() {
		if err := app.GetEchoHTTP().Start(":1323"); err != nil {
			app.GetLogger().Infof("could not start HTTP Server: %s", err.Error())
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	app.GetLogger().Infof("Shutting down HTTP Server")

	ctx, cancel := context.WithTimeout(app.GetCtx(), 10*time.Second)
	defer cancel()
	if err := app.GetEchoHTTP().Shutdown(ctx); err != nil {
		app.GetLogger().Fatalf("could not shutdown HTTP Server: %s", err.Error())
	}

}
