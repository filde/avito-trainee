package models

import "time"

type PullRequest struct {
	PullRequestID     string     `gorm:"primaryKey" json:"pull_request_id"`
	PullRequestName   string     `gorm:"not null" json:"pull_request_name"`
	AuthorID          string     `gorm:"not null" json:"author_id"`
	Status            string     `gorm:"not null;type:status_enum" json:"status"`
	AssignedReviewers []string   `gorm:"not null" json:"assigned_reviewers"`
	CreatedAt         *time.Time `json:"createdAt"`
	MergedAt          *time.Time `json:"mergedAt"`
	User              `gorm:"foreignKey:AuthorID;references:UserID" json:"-"`
}
