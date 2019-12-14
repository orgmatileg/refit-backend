package app

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	// golang-migrate need it
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func MigrateRun() {

	db, err := sql.Open("mysql", "root:masuk123@tcp(luqmanul.com:3306)/refit?multiStatements=true")
	if err != nil {
		log.Fatalf("could not connect to the MySQL database... %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("could not ping DB... %v", err)
	}

	fmt.Println("DB Connected")

	// // Run migrations
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatalf("could not start sql migration... %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./internal/app/migrate",
		"mysql", driver)

	if err != nil {
		log.Fatalf("migration failed... %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("An error occurred while syncing the database.. %v", err)
	}

	log.Println("Database migrated")
	// actual logic to start your application
	os.Exit(0)
}
