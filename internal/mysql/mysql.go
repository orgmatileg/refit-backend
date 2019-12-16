package mysql

import (
	"database/sql"
	// mysql need this
	_ "github.com/go-sql-driver/mysql"
)

// IDBMySQL interface
type IDBMySQL interface {
	CreateConnection() error
	GetDB() *sql.DB
}

// DBMySQL struct
type DBMySQL struct {
	db *sql.DB
}

// NewDBMySQL get new DBMySQL
func NewDBMySQL() IDBMySQL {
	return &DBMySQL{}
}

// GetDB get db
func (d *DBMySQL) GetDB() *sql.DB {
	return d.db
}

// CreateConnection open new connection MySQL
func (d *DBMySQL) CreateConnection() error {
	if d.db != nil {
		return nil
	}
	db, err := sql.Open("mysql", "root:masuk123@tcp(luqmanul.com:3306)/refit")
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}
	d.db = db
	return nil
}
