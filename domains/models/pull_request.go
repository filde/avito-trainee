package models

import "time"

type PullRequest struct {
	PullRequestID     string     `gorm:"primaryKey" json:"pull_request_id"`
	PullRequestName   string     `gorm:"not null" json:"pull_request_name"`
	AuthorID          string     `gorm:"not null" json:"author_id"`
	Status            string     `gorm:"not null;type:status_enum" json:"status"`
	CreatedAt         *time.Time `json:"createdAt"`
	MergedAt          *time.Time `json:"mergedAt"`
	User              `gorm:"foreignKey:AuthorID;references:UserID" json:"-"`
	AssignedReviewers []*User `gorm:"many2many:pr_reviewers"`
}
