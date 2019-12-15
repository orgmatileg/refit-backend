package serve

import (
	"context"
	"database/sql"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"os"
	"os/signal"
	"refit_backend/internal/logger"
	"refit_backend/internal/mysql"
	"time"
)

// IAppServe interface
type IAppServe interface {
	// Getter
	GetCtx() context.Context
	GetEchoHTTP() *echo.Echo
	GetLogger() logger.Logger
	GetMySQL() *sql.DB

	// Initiator
	InitLogger()
	InitHTTP()
	InitCtx()
}

// AppServe struct
type appServe struct {
	ctx      context.Context
	echoHTTP *echo.Echo
	logger   logger.Logger
	mysql    *sql.DB
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

func (a *appServe) GetMySQL() *sql.DB {
	return a.mysql
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
	db, err := mysql.CreateConnection()
	if err != nil {
		a.logger.Fatalf("could not create mysql connectioon: %s", err.Error())
	}
	a.mysql = db
}

func (a *appServe) InitHTTP() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		time.Sleep(5 * time.Second)
		return c.JSON(http.StatusOK, "OK")
	})

	a.echoHTTP = e
}

func (a *appServe) InitCtx() {
	a.ctx = context.Background()
}

func getAppServe() IAppServe {
	return &appServe{}
}

// Start Serve App
func Start() {
	app := getAppServe()

	app.InitCtx()
	app.InitLogger()
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
