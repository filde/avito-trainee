package database

import (
	"avito-trainee/common/constants"
	"avito-trainee/domains/models"
	"avito-trainee/external/httpserver"
	"avito-trainee/helpers"
	"errors"
	"gorm.io/gorm"
)

var _ httpserver.StorageItf = &Database{}

func (db *Database) CreateTeam(team *models.Team) (*models.ErrorType, error) {
	var errorModel *models.ErrorType
	err := db.Transaction(func(tx *gorm.DB) error {
		err := db.Create(&models.Team{TeamName: team.TeamName}).Error
		if err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				errorModel = helpers.GetError(constants.TEAM_EXISTS, team.TeamName)
			}
			return err
		}

		for i := 0; i < len(team.Members); i++ {
			team.Members[i].TeamName = team.TeamName
		}

		err = db.Create(team.Members).Error
		if err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				errorModel = helpers.GetError(constants.USER_EXISTS)
			}
			return err
		}
		return nil
	})
	return errorModel, err
}

func (db *Database) GetTeam(name string) (*models.Team, error) {
	var team *models.Team
	err := db.Preload("Members").Where("team_name = ?", name).
		First(&team).Error
	return team, err
}
