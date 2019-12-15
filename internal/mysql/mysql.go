package mysql

import (
	"database/sql"
	// mysql need this
	_ "github.com/go-sql-driver/mysql"
)

// CreateConnection open new connection MySQL
func CreateConnection() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:masuk123@tcp(luqmanul.com:3306)/refit?multiStatements=true")
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
