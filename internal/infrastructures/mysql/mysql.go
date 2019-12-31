package mysql

import (
	"database/sql"
	"fmt"
	"refit_backend/internal/logger"
	"time"

	// mysql need this
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
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
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		viper.GetString("mysql.username"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetString("mysql.port"),
		viper.GetString("mysql.dbname"),
	)
	fmt.Println(dsn)
	dbCon, err := sql.Open("mysql", dsn)
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
