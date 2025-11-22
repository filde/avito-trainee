package database

import (
	"avito-trainee/common/constants"
	"avito-trainee/domains/models"
	"avito-trainee/external/httpserver"
	"avito-trainee/helpers"
	"errors"
	"gorm.io/gorm"
	"time"
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

func (db *Database) UpdateUserActivity(userID string, isActive bool) error {
	err := db.Model(&models.User{}).Where("user_id = ?", userID).Update("is_active", isActive).Error
	return err
}

func (db *Database) GetUser(userID string) (*models.UserFull, error) {
	var user *models.UserFull
	err := db.Model(&models.User{}).Where("user_id = ?", userID).First(&user).Error
	return user, err
}

func (db *Database) GetUserPR(userID string) (*models.UsersPR, error) {
	var user *models.UsersPR
	err := db.Model(&models.User{}).Preload("PullRequests").First(*user).Error
	return user, err
}

func (db *Database) TeamExists(name string) error {
	var i int
	err := db.Model(&models.Team{}).Where("team_name = ?", name).Select("1").
		First(&i).Error
	return err
}

func (db *Database) GetTeamUsers(name string, author string) ([]*models.User, error) {
	var users []*models.User
	err := db.Model(&models.User{}).Where("team_name = ?", name).
		Where("user_id <> ?", author).Where("is_active = ?", true).
		Order("RANDOM()").Limit(2).Find(&users).Error
	return users, err
}

func (db *Database) CreatePR(pr *models.PullRequest) error {
	err := db.Create(&pr).Error
	return err
}

func (db *Database) GetPR(id string) (*models.PullRequest, error) {
	var pr *models.PullRequest
	err := db.Model(&models.PullRequest{}).Preload("Reviewers").Where("pull_request_id = ?", id).First(&pr).Error
	return pr, err
}

func (db *Database) MergePR(id string, mergeTime *time.Time) error {
	err := db.Model(&models.PullRequest{}).Where("pull_request_id = ?", id).
		Update("status", constants.MERGED_STATUS).Update("merged_at = ?", mergeTime).Error
	return err
}

func (db *Database) GetTeamActiveUser(team string, notAllowed ...string) (*models.User, error) {
	var user *models.User
	err := db.Model(&models.User{}).Where("team_name = ?", team).
		Where("user_id now in ?", notAllowed).Where("is_active = ?", true).
		Order("RANDOM()").First(&user).Error
	return user, err
}

func (db *Database) UpdatePR(pr *models.PullRequest) error {
	err := db.Updates(pr).Error
	return err
}
