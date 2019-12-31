package mysql

import (
	"database/sql"
	"refit_backend/internal/logger"
	"time"

	// mysql need this
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// GetDB get db
func GetDB() *sql.DB {
	if db != nil {
		return db
	}
	createConnection()
	return db
}

// CreateConnection open new connection MySQL
func createConnection() error {
	dbCon, err := sql.Open("mysql", "root:masuk123@tcp(luqmanul.com:3306)/refit?parseTime=true")
	if err != nil {
		logger.Errorf("could not open mysql database dsn: %s")
		return err
	}
	err = dbCon.Ping()
	if err != nil {
		logger.Errorf("could not ping mysql database: %s", err.Error())
		return err
	}
	logger.Infof("database mysql: Connected!")
	dbCon.SetMaxOpenConns(100)
	dbCon.SetMaxIdleConns(10)
	dbCon.SetConnMaxLifetime(time.Duration(300 * time.Second))
	db = dbCon
	return nil
}
