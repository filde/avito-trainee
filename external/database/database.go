package database

import (
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

func InitDB() *Database {
	postgresDialector, err := newPostgresDialector()
	if err != nil {
		log.Panic().Msgf("Couldn't init dialector: %v", err)
	}

	var db *gorm.DB
	for i := 0; i < 12; i++ {
		db, err = gorm.Open(postgresDialector, &gorm.Config{})
		if err == nil {
			break
		}
		log.Error().Msgf("Init database error: %v. Will try again in a few seconds...", err)
		time.Sleep(10 * time.Second)
	}

	if err != nil {
		log.Panic().Msgf("Init database error: %v", err)
	}

	autoMigrate(db)

	return &Database{
		db,
	}
}
