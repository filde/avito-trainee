package database

import (
	"avito-trainee/common/constants"
	"avito-trainee/domains/models"
	"avito-trainee/external/httpserver"
	"avito-trainee/helpers"
	"gorm.io/gorm"
	"time"
)

var _ httpserver.StorageItf = &Database{}

func (db *Database) CreateTeam(team *models.Team) (*models.ErrorType, error) {
	var errorModel *models.ErrorType
	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&models.Team{TeamName: team.TeamName}).Error
		if err != nil {
			if helpers.IsAlreadyExists(err) {
				errorModel = helpers.GetError(constants.TEAM_EXISTS, team.TeamName)
			}
			return err
		}

		for i := 0; i < len(team.Members); i++ {
			team.Members[i].TeamName = team.TeamName
		}

		err = tx.Create(team.Members).Error
		if err != nil {
			if helpers.IsAlreadyExists(err) {
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
	err := db.Where("team_name = ?", name).First(&team).Error
	if err != nil {
		return nil, err
	}
	var users []*models.User
	err = db.Where("team_name = ?", name).Find(&users).Error
	if err != nil {
		return nil, err
	}
	team.Members = users
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
	err := db.Model(&models.User{}).Where("user_id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}

	var prList []*models.PRShort
	err = db.Model(&models.PullRequest{}).Joins("join pull_request_reviewers on pull_request_id=pull_req_id").
		Where("reviewer_id = ?", userID).Find(&prList).Error
	if err != nil {
		return nil, err
	}
	user.PullRequests = prList
	return user, err
}

func (db *Database) GetTeamReviewers(name string, author string) ([]string, error) {
	var users []string
	err := db.Model(&models.User{}).Where("team_name = ?", name).
		Where("user_id <> ?", author).Where("is_active = ?", true).
		Order("RANDOM()").Limit(2).Select("user_id").Find(&users).Error
	return users, err
}

func (db *Database) CreatePR(pr *models.PullRequest) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&pr).Error
		if err != nil {
			return err
		}

		prRevList := make([]*models.PullRequestReviewers, len(pr.AssignedReviewers))
		for i := 0; i < len(prRevList); i++ {
			prRevList[i] = &models.PullRequestReviewers{PullReqID: pr.PullRequestID, ReviewerID: pr.AssignedReviewers[i]}
		}

		return tx.Create(&prRevList).Error
	})

	return err
}

func (db *Database) GetPR(id string) (*models.PullRequest, error) {
	var pr *models.PullRequest
	err := db.Model(&models.PullRequest{}).Where("pull_request_id = ?", id).First(&pr).Error
	if err != nil {
		return nil, err
	}

	err = db.Model(&models.PullRequestReviewers{}).Where("pull_req_id = ?", id).
		Select("reviewer_id").Find(&pr.AssignedReviewers).Error
	return pr, err
}

func (db *Database) MergePR(id string, mergeTime *time.Time) error {
	err := db.Model(&models.PullRequest{}).Where("pull_request_id = ?", id).
		Where("status = ?", constants.OPEN_STATUS).
		Update("status", constants.MERGED_STATUS).Update("merged_at", mergeTime).Error
	return err
}

func (db *Database) GetTeamActiveUser(team string, notAllowed ...string) (string, error) {
	var user string
	err := db.Model(&models.User{}).Where("team_name = ?", team).
		Where("user_id not in ?", notAllowed).Where("is_active = ?", true).
		Order("RANDOM()").Select("user_id").First(&user).Error
	return user, err
}

func (db *Database) ChangeReviewer(oldReviewer *models.NewPRReviewer, newReviewer string) error {
	err := db.Model(&models.PullRequestReviewers{}).Where("pull_req_id = ?", oldReviewer.PullRequestID).
		Where("reviewer_id = ?", oldReviewer.OldReviewerID).Update("reviewer_id", newReviewer).Error
	return err
}
