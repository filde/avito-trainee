package database

import (
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func autoMigrate(db *gorm.DB) {
	err := db.AutoMigrate()
	if err != nil {
		log.Panic().Msgf("Couldn't auto migrate: %v", err)
	}
}
