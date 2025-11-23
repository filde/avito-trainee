package database

import (
	"avito-trainee/common/constants"
	"avito-trainee/domains/models"
	"fmt"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func autoMigrate(db *gorm.DB) {
	err := db.Exec(fmt.Sprintf(`
			DO
			$$
			BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status_enum') THEN
    			create type status_enum AS ENUM ('%v', '%v');
  			END IF;
			END
			$$;
	`, constants.OPEN_STATUS, constants.MERGED_STATUS)).Error
	if err != nil {
		log.Panic().Msgf("Couldn't create enum type: %v", err)
	}
	err = db.AutoMigrate(
		&models.Team{},
		&models.User{},
		&models.PullRequest{},
		&models.PullRequestReviewers{},
	)
	if err != nil {
		log.Panic().Msgf("Couldn't auto migrate: %v", err)
	}
}
