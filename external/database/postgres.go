package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newPostgresDialector() (gorm.Dialector, error) {
	host := os.Getenv("POSTGRES_HOST")
	if len(host) == 0 {
		host = "localhost"
	}
	port := os.Getenv("POSTGRES_PORT")
	if len(port) == 0 {
		port = "5432"
	}
	user := os.Getenv("POSTGRES_USER")
	if len(user) == 0 {
		user = "postgres"
	}
	dbname := os.Getenv("POSTGRES_DB")
	if len(dbname) == 0 {
		dbname = "postgres"
	}
	password := os.Getenv("POSTGRES_PASSWORD")
	if len(password) == 0 {
		password = "postgres"
	}

	connectStr := fmt.Sprintf("host=%v port=%v user=%v dbname=%v password=%v",
		host, port, user, dbname, password,
	)

	return postgres.Open(connectStr), nil
}
